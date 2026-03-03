# ADR-017: 認証システムにおけるCSRF対策

## ステータス
承認済み

## コンテキスト

### 背景
本プロジェクトではFirebase AuthenticationによるOAuth認証（Google/Apple/X）を採用し、セッション管理にはFirebase IDトークンをHTTP-only Cookieに保存する方式を検討している。

### システム構成
- **認証方式**: Firebase Authentication（OAuth 2.0 / OpenID Connect準拠）
- **セッション管理**: Firebase IDトークンをHTTP-only Cookieに保存
- **セッション有効期限**: 7日（ADR-017関連決定）
- **OAuthステート検証**: Firebase SDK任せ

### 課題
Firebase IDトークンをCookieで保存する場合、CSRF（Cross-Site Request Forgery）攻撃への対策が必要かどうか、必要な場合はどの方式を採用すべきか検討が必要。

### 制約
- 認証プロセスはOAuthプロバイダに完全に委任
- 実装コストとセキュリティのバランス
- Firebase Authenticationの既存セキュリティ機能との相性

## 決定

**SameSite属性（Lax）+ Firebase Authenticationの標準機能** を採用する。

### 詳細構成

| 項目 | 設定値 | 説明 |
|------|--------|------|
| Cookie属性 | HttpOnly + Secure + SameSite=Lax | XSS対策+HTTPS限定+CSRF対策 |
| CSRF追加対策 | なし | SameSite=Laxで十分と判断 |
| トークン検証 | Firebase Admin SDK | IDトークンの署名検証 |
| 有効期限チェック | Firebase Admin SDK任せ | トークンの有効期限を自動検証 |

```go
// Cookie設定例（Go）
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    idToken,
    Path:     "/",
    MaxAge:   60 * 60 * 24 * 7, // 7日
    HttpOnly: true,
    Secure:   true,              // HTTPSのみ
    SameSite: http.SameSiteLax,  // CSRF対策
})
```

## 検討した選択肢

### 比較表

| 選択肢 | セキュリティ | 実装コスト | Firebase連携 | 総合評価 |
|--------|-------------|-----------|--------------|----------|
| **SameSite=Laxのみ** | ○ | ◎ | ◎ | **◎** |
| Double Submit Cookie | ◎ | △ | ○ | ○ |
| Synchronizer Token | ◎ | △ | △ | △ |

### 詳細評価

#### 選択肢A: SameSite=Laxのみ（採用）
- メリット:
  - 実装コストが最も低い
  - Firebaseの標準機能と競合しない
  - モダンブラウザで広くサポート
  - OAuthリダイレクト時にCookie送信される
- デメリット:
  - 古いブラウザではSameSiteが無効化される可能性
  - GETリクエストによるCSRFは防げない（状態変更系はPOSTで対応）

#### 選択肢B: Double Submit Cookie
- メリット:
  - 古いブラウザでも有効
  - 二重送信防御の仕組みが確立
- デメリット:
  - 実装コストが高い（トークン生成・検証ロジック）
  - Firebase Authenticationとの整合性に課題
  - Cookieとリクエストヘッダーの両方にトークンが必要

#### 選択肢C: Synchronizer Token
- メリット:
  - 最強のCSRF対策
  - セッションごとに一意なトークン
- デメリット:
  - 実装コストが高い
  - サーバーサイドの状態管理が必要
  - Firebaseのステートレス認証との相性が悪い

## 理由

### SameSite=Laxのみが採用された理由

1. **Firebase Authenticationの特性**:
   - FirebaseはOAuth 2.0 / OpenID Connect準拠
   - IDトークンは署名付きJWTで、改ざん検証が可能
   - CSRF攻撃でIDトークンを盗難しても、有効なトークンでないと検証で弾かれる
   - 攻撃者がトークンを取得できても、署名付きJWTは偽造不可能

2. **SameSite=Laxの十分性**:
   - `SameSite=Lax`はGETリクエスト（トップレベルナビゲーション）ではCookieを送信
   - POSTリクエスト（クロスサイト）ではCookieを送信しない
   - OAuthリダイレクトはGETリクエストであり、正常に動作
   - 状態変更系APIはPOSTで実装し、CSRFを防御

3. **費用対効果（Cost-Benefit Analysis）**:
   - **SameSite=Lax**: 実装工数15分（Cookie属性設定のみ）
   - **Double Submit Cookie**: 実装工数6-8時間（トークン生成・検証・ミドルウェア・フロント連携）
   - コスト差分: **24〜32倍**
   - 得られる利益: 古いブラウザ対策のみ（IE11等、モダンブラウザで不要）
   - 結論: 費用に見合う利益が得られないため、SameSite=Laxを採用

4. **古いブラウザへの対応**:
   - IE11等のレガシーブrowserはサポート対象外
   - モダンブラウザ（Chrome 51+, Firefox 60+, Safari 12+）でSameSiteは安定動作

### 補足防御策

SameSite=Laxに加え、以下の対策で多層防御を構築:

| 対策 | 内容 |
|------|------|
| XSS対策 | HttpOnly Cookieでトークン保護 |
| HTTPS強制 | Secure属性で平文通信を禁止 |
| トークン検証 | Firebase Admin SDKで署名・有効期限検証 |
| CORS | 適切なオリジン設定 |
| 状態変更はPOST | GETリクエストで状態変更しない設計 |

## 影響

### 機能への影響
- **F-01 ログイン機能**: SameSite=Lax設定でセッションCookieを発行
- **全APIエンドポイント**: Cookieが自動送信され、追加のCSRFトークン不要

### 依存パッケージ
- 追加なし（標準ライブラリで完結）

### 参考ADR
- ADR-016: Phase0横断的決定（セッション有効期限7日）

## 参考情報

### RFC・仕様
- [RFC 6265bis: SameSite Cookies](https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-rfc6264bis)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

### Firebase公式ドキュメント
- [Firebase Authentication: Session Management](https://firebase.google.com/docs/auth/web/manage-users)
- [Firebase Security Rules](https://firebase.google.com/docs/rules)

### 記事
- [web.dev: SameSite cookies explained](https://web.dev/samesite-cookies-explained/)
- [Firebase Blog: Protecting users from CSRF attacks](https://firebase.googleblog.com/2016/09/protecting-users-from-csrf-attacks.html)
