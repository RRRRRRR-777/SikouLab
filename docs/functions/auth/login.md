# ログイン機能

## 機能概要

Firebase Authenticationを使用し、Google / Apple / X によるOAuth認証を実装して、ユーザーがサービスに安全にログインできるようにする。認証プロセスはOAuthプロバイダ側に完全に委任し、初回ログイン時は取得したユーザー情報を基に本システム側で自動登録を行い、既存ユーザーはダッシュボードへ遷移する。

## 目的

- 未ログインユーザーを認証し、サービスへアクセス可能にする
- 認証プロセスはOAuthプロバイダに完全に委任し、初回ログイン時は取得したユーザー情報を基に本システム側で自動登録を行い、サブスクリプション登録へ誘導する
- セッション管理によりログイン状態を維持する

## 機能条件

### アクセス可否

| 状態 | 操作可否 |
|--------|----------|
| 未ログイン | ○ |
| ログイン済み | - |

※ ログイン画面は未ログイン状態の場合のみアクセス可能

### 制約事項
✅ **決定済み**

- セッション有効期限
  - **決定: 7日固定**（ADR-017関連決定）
  - リメンバーme機能はなし（実装コスト削減）

- CSRF対策方式
  - **決定: SameSite=Lax**（ADR-017承認済み）
  - 費用対効果の観点から採用（実装15分 vs Double Submitの6-8時間）

- OAuthプロバイダごとの取得項目
  - **決定: デフォルト案**
  - Google: name, email, picture
  - Apple: name, email（非公開の場合はemailのみ）
  - X: name, profile_image_url

## 画面設計図
🟡 **中程度**

Pencil未定義（実装のみ）

### レイアウト構成（暫定）

```
┌─────────────────────────────────────────────────────────┐
│                                                           │
│                    シコウラボ                             │
│                                                           │
│                  専門家の思考プロセスを                    │
│                   のぞけるプラットフォーム                  │
│                                                           │
│  ┌───────────────────────────────────────────────────┐  │
│  │                                                   │  │
│  │            OAuthプロバイダでログイン               │  │
│  │                                                   │  │
│  │  ┌──────────────────────────────────────────┐    │  │
│  │  │ [G]  Google でログイン                      │    │  │
│  │  └──────────────────────────────────────────┘    │  │
│  │                                                   │  │
│  │  ┌──────────────────────────────────────────┐    │  │
│  │  │   Apple でログイン                          │    │  │
│  │  └──────────────────────────────────────────┘    │  │
│  │                                                   │  │
│  │  ┌──────────────────────────────────────────┐    │  │
│  │  │   X でログイン                               │    │  │
│  │  └──────────────────────────────────────────┘    │  │
│  │                                                   │  │
│  └───────────────────────────────────────────────────┘  │
│                                                           │
│  ログインすることで利用規約に同意したものとみなされます      │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

## 関連テーブル

```mermaid
erDiagram
    %% 正: docs/versions/1_0_0/data-model.md
    users ||--o| user_settings : "has"
    users ||--o| plans : "subscribes"

    users {
        bigint id PK
        string oauth_provider
        string oauth_user_id
        string name
        string display_name
        string avatar_url
        string role "admin/writer/user"
        bigint plan_id FK
        string univapay_customer_id
        string subscription_status "active/canceled/past_due"
        datetime created_at
        datetime updated_at
    }

    user_settings {
        bigint id PK
        bigint user_id FK
        boolean sidebar_article_expanded "記事タブの折りたたみ状態"
        boolean sidebar_admin_expanded "管理タブの折りたたみ状態"
        datetime created_at
        datetime updated_at
    }

    plans {
        bigint id PK
        string name
        string description
        integer amount "月額料金（最小通貨単位）"
        string currency "通貨コード ISO-4217（例: JPY）"
        boolean is_active
        datetime created_at
        datetime updated_at
    }
```

## フロー図

```mermaid
flowchart TD
    A[ユーザー情報] --> B{ログイン済み?}
    B -->|Yes| C[ダッシュボード表示]
    B -->|No| D[ログイン画面表示]

    D --> E{OAuthプロバイダ選択}
    E --> F[Google]
    E --> G[Apple]
    E --> H[X]

    F --> I[OAuth認証]
    G --> I
    H --> I

    I --> J{初回ログイン?}
    J -->|Yes| K[ユーザー作成]
    K --> L[UnivaPayカスタマー作成]
    L --> M[サブスクリプション登録画面へ]
    J -->|No| N[ダッシュボードへ]
```

## シーケンス図

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Front as フロントエンド
    participant Firebase as Firebase Auth
    participant API as バックエンドAPI
    participant DB as データベース
    participant UnivaPay as UnivaPay

    User->>Front: ログイン画面アクセス
    Front->>Front: 未ログイン確認

    User->>Front: Googleボタンクリック
    Front->>Firebase: OAuth認証開始
    Firebase->>User: 認証画面表示
    User->>Firebase: 認証実行
    Firebase-->>Front: ID Token

    Front->>API: POST /api/v1/auth/login {id_token}
    API->>Firebase: ID Token検証
    Firebase-->>API: ユーザー情報（uid, email, name, picture）

    alt 初回ログイン
        API->>DB: ユーザー作成
        API->>DB: user_settings作成
        API->>UnivaPay: カスタマー作成
        UnivaPay-->>API: customer_id
        API->>DB: univapay_customer_id保存
        API-->>Front: is_first_login=true
        Front->>User: サブスクリプション登録画面へ遷移
    else 既存ユーザー
        API->>DB: ユーザー取得
        API->>DB: サブスクリプション状態確認
        API-->>Front: user情報
        Front->>User: ダッシュボードへ遷移
    end
```

## 機能要件
🟡 **中程度**

### 機能要件1: OAuth認証（F-01）
- 機能仕様1: Firebase Authenticationを使用したGoogle / Apple / X によるOAuth認証

### 機能要件2: ユーザー自動登録（F-01）
- 機能仕様1: 初回ログイン時のユーザー自動登録（users、user_settings）
- 機能仕様2: 初回ログイン時のUnivaPayカスタマー作成

### 機能要件3: セッション管理（F-01）
- 機能仕様1: セッション管理（httpOnly Cookie、SameSite=Lax）
- 機能仕様2: ログアウト機能

### 機能要件4: アクセス制御（F-01）
- 機能仕様1: 未ログイン時のアクセス制限（保護されたページへリダイレクト）

### 機能要件5: セキュリティ詳細（F-01）
- 機能仕様1: CSRF対策: **SameSite=Lax**（HttpOnly + Secure + SameSite=LaxのCookie設定）。追加のCSRFトークンは不要（ADR-017参照、費用対効果で決定）
- 機能仕様2: セッション有効期限: **7日**（CookieのMaxAgeで設定）
- 機能仕様3: OAuthステートの検証: **Firebase SDK任せ**（Firebase Authenticationのデフォルト実装に依存）
- 機能仕様4: リメンバーme機能: **なし**（7日固定）

## 非機能要件
🟢 **後回し可**

### 非機能要件1: パフォーマンス
- 非機能仕様1: OAuthリダイレクト: 3秒以内
- 非機能仕様2: セッション確認: 500ms以内

### 非機能要件2: セキュリティ
- 非機能仕様1: Firebase Authenticationのセキュリティ機能を活用（OAuth 2.0 / OpenID Connect準拠）
- 非機能仕様2: ID Token検証によるユーザー認証（Firebase Admin SDK）
- 非機能仕様3: httpOnly Cookie（HttpOnly + Secure + SameSite=Lax）
- 非機能仕様4: セッションハイジャック対策（SameSite=Lax + HTTPS強制）
- 非機能仕様5: CSRF対策（SameSite=Lax、追加トークン不要・ADR-017参照）

### 非機能要件3: 可用性
- 非機能仕様1: OAuthプロバイダ障害時: エラーメッセージ表示、リトライ誘導

## ログ
🟢 **後回し可**

### 出力タイミング
- 案1: 全認証操作時に出力（Firebase認証・ID Token検証・ユーザー作成・UnivaPay連携） → 追跡しやすいがログ量増加
- 案2: エラー時のみ出力 → ログ量削減だが正常系追跡困難
- 案3: 重要操作のみ出力（初回ログイン時のユーザー作成・UnivaPay連携・認証エラー） → バランス型
- **決定: TBD**

### ログレベル方針
- 案1: INFO中心（認証開始・成功・ユーザー作成をINFO） → 詳細追跡可能
- 案2: WARN/ERROR中心（認証エラー・UnivaPay連携エラーのみ） → 異常検知に特化
- 案3: INFO（認証成功・ユーザー作成）+ WARN（認証失敗）+ ERROR（システムエラー） → バランス型
- **決定: TBD**

## ユースケース
🟡 **中程度**

### シナリオ1: 初回ログイン（早期決定）
1. ユーザーがサービスにアクセス
2. ログイン画面が表示される
3. Googleボタンをクリック
4. Google認証画面が表示され、ユーザー認証はGoogle側に完全に委任される
5. 認証を完了
6. 取得したユーザー情報を基に本システム側で自動登録される
7. サブスクリプション登録画面へ遷移

### シナリオ2: 既存ユーザーログイン（早期決定）
1. ユーザーがサービスにアクセス
2. ログイン画面が表示される
3. Appleボタンをクリック
4. Apple認証画面が表示され、ユーザー認証はApple側に完全に委任される
5. 認証を完了
6. ダッシュボードへ遷移

### シナリオ3: ログアウト（TBD可）
1. ユーザーがログアウトボタンをクリック
2. セッションが削除される
3. ログイン画面へ遷移

## テストケース
🟢 **実装済み**

### 単体テスト（バックエンド）

#### `TestAuthUsecase_Login`（`backend/internal/usecase/`）

| ケース名 | 期待値 |
|---------|--------|
| 有効なIDトークンで既存ユーザーがログインする | ユーザー情報が返される、isFirstLogin=false |
| 無効なIDトークンでエラーが返される | ErrInvalidTokenが返される |
| 初回ログインでユーザーが新規作成される | isFirstLogin=true、Createが呼ばれる |

#### `TestAuthUsecase_GetCurrentUser`（`backend/internal/usecase/`）

| ケース名 | 期待値 |
|---------|--------|
| 有効なセッショントークンでユーザー情報が返される | ユーザー情報が返される |
| 無効なセッショントークンでエラーが返される | ErrInvalidTokenが返される |

#### `TestAuthMiddleware_RequireAuth`（`backend/internal/middleware/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| Cookieあり・有効なトークンの場合はnextが呼ばれる | 200、next handler実行 |
| Cookieなしの場合は401が返される | 401 |
| Cookieあり・無効なトークン（検証失敗）の場合は401が返される | 401 |
| Cookieあり・トークンは有効だがDBにユーザーが存在しない場合は401 | 401 |
| contextにUserが設定されること | next handlerからuser情報を取得できる |

#### `TestAuthMiddleware_RequireRole`（`backend/internal/middleware/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| admin roleでadminエンドポイントにアクセスできる | 200 |
| user roleでadminエンドポイントにアクセスすると403 | 403 |
| writer roleでadminエンドポイントにアクセスすると403 | 403 |
| writer roleでwriterエンドポイントにアクセスできる | 200 |
| user roleでwriterエンドポイントにアクセスすると403 | 403 |
| admin roleはwriter権限エンドポイントにもアクセスできる（admin > writer > user） | 200 |
| contextにUserがない場合は401 | 401 |
| 不明なroleの場合は403 | 403 |

#### `TestAuthMiddleware_RequireSubscription`（`backend/internal/middleware/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| subscription_status=activeのユーザーはアクセスできる | 200 |
| subscription_status=canceledのユーザーは403 | 403 |
| subscription_status=past_dueのユーザーは403 | 403 |
| contextにUserがない場合は401 | 401 |

#### `TestAuthHandler_ServeLogin`（`backend/internal/handler/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| 本番環境ではSecure=trueのCookieが返される | 200、Cookie Secure=true |
| 開発環境ではSecure=falseのCookieが返される | 200、Cookie Secure=false |
| id_tokenが空の場合は400が返される | 400 |
| 不正なJSONボディの場合は400が返される | 400 |
| ErrInvalidTokenの場合は401が返される | 401 |
| その他のエラーの場合は500が返される | 500 |

#### `TestAuthHandler_ServeMe`（`backend/internal/handler/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| 有効なCookieで200とユーザー情報が返される | 200、ユーザー情報 |
| Cookieがない場合は401が返される | 401 |
| ErrInvalidTokenの場合は401が返される | 401 |
| その他のエラーの場合は500が返される | 500 |

#### `TestAuthHandler_ServeLogout`（`backend/internal/handler/auth_test.go`）

| ケース名 | 期待値 |
|---------|--------|
| 本番環境のログアウトCookieはSecure=true | 204、Cookie value=""、MaxAge=-1、Secure=true |
| 開発環境のログアウトCookieはSecure=false | 204、Cookie value=""、MaxAge=-1、Secure=false |

---

### 単体テスト（フロントエンド）

#### ミドルウェア パスマッチング（`frontend/middleware.test.ts`）

| ケース名 | 期待値 |
|---------|--------|
| `'/'` は完全一致のみ | `/login` や `/articles` にはマッチしない |
| 非ルートパスはサブパスにもマッチする | `/articles/123` は `/articles` にマッチ、`/articles-legacy` はマッチしない |
| 公開パス・保護パスの判定が正しい | `/login` は公開、`/` と `/articles/2026` は保護 |

#### ミドルウェア ルートガード（`frontend/middleware.test.ts`）

| ケース名 | 期待値 |
|---------|--------|
| セッションなしで保護パスへのリクエストは/loginにリダイレクト | 307、Location: /login |
| セッションありで/loginへのリクエストは/にリダイレクト | 307、Location: / |
| セッションなしで/subscriptionへのアクセスは通過する | 200（NextResponse.next()） |
| セッションありで保護パスへのアクセスは通過する | 200（NextResponse.next()） |

#### Providers（`frontend/app/__tests__/providers.test.tsx`）

| ケース名 | 期待値 |
|---------|--------|
| AuthProviderが含まれており、useAuthが使える | isLoading=false（AuthProvider内でonAuthStateChangedが呼ばれる） |
| QueryClientProviderが含まれている | useQueryClientが値を返す |

#### 認証APIクライアント（`frontend/lib/auth/__tests__/auth-api.test.ts`）

| グループ | ケース名 | 期待値 |
|---------|---------|--------|
| `authApi.login` | 成功時はユーザー情報と初回ログインフラグを返す | user情報、is_first_login=false |
| `authApi.login` | 初回ログイン時はis_first_loginがtrue | is_first_login=true |
| `authApi.login` | 無効なID Token時は401エラーを投げる | 401エラー |
| `authApi.login` | ネットワークエラー時はエラーを投げる | Network Error |
| `authApi.getMe` | 成功時はユーザー情報を返す | user情報 |
| `authApi.getMe` | 未認証時は401エラーを投げる | 401エラー |
| `authApi.logout` | 成功時は204 No Contentを返す | 正常完了 |
| `authApi.logout` | 未認証時は401エラーを投げる | 401エラー |

#### 認証コンテキスト（`frontend/lib/auth/__tests__/auth-context.test.tsx`）

| グループ | ケース名 | 期待値 |
|---------|---------|--------|
| 初期状態 | 認証状態は未ログイン・ローディング完了 | user=null、isLoading=false、isAuthenticated=false |
| `loginWithGoogle` | 成功時はユーザー情報をセットする | user情報がセット、isAuthenticated=true |
| `loginWithGoogle` | 初回ログイン時はサブスクリプション画面へ遷移する | `/subscription` へpush |
| `loginWithApple` | 成功時はユーザー情報をセットする | isAuthenticated=true |
| `loginWithX` | 成功時はユーザー情報をセットする | isAuthenticated=true |
| `logout` | 成功時はユーザー情報をクリアする | user=null、isAuthenticated=false |
| ログインページでの動作 | Firebaseユーザーがいても/auth/meも/auth/loginも呼ばない | getMe未呼び出し、authApi.login未呼び出し、user=null |
| `refresh` | 成功時は最新のユーザー情報を取得する | getMeが呼ばれる |
| `onAuthStateChanged リフレッシュ` | Firebaseユーザーがいる場合、getIdTokenが呼ばれてからPOST /auth/loginが呼ばれる | getIdToken(forceRefresh=true) が呼ばれ、authApi.login が呼ばれる |
| `onAuthStateChanged リフレッシュ` | トークンリフレッシュ成功後、loginレスポンスからユーザー情報がセットされる | user情報がセット、isAuthenticated=true |
| `onAuthStateChanged リフレッシュ` | getIdToken失敗時はユーザーをnullにする | user=null、isAuthenticated=false |
| `onAuthStateChanged リフレッシュ` | login失敗（401等）時はユーザーをnullにする | user=null、isAuthenticated=false |

#### ログイン画面（`frontend/components/auth/__tests__/LoginPage.test.tsx`）

| グループ | ケース名 | 期待値 |
|---------|---------|--------|
| 画面表示 | ロゴとタイトルが表示される | "SikouLab"、利用規約文言が表示される |
| 画面表示 | 3つのOAuthボタンが表示される | Google / Apple / X のボタンが表示される |
| OAuthボタン操作 | Googleボタンクリック時はGoogleログインを実行する | `loginWithGoogle` が1回呼ばれる |
| OAuthボタン操作 | Appleボタンクリック時はAppleログインを実行する | `loginWithApple` が1回呼ばれる |
| OAuthボタン操作 | Xボタンクリック時はXログインを実行する | `loginWithX` が1回呼ばれる |
| アクセシビリティ | 各OAuthボタンには適切なtype属性が設定されている | `type="button"` |
| ダークモード対応 | ダークモードでも正しく表示される | 基本要素が表示される |

---

### E2Eテスト（実装完了後に記載）

| テストシナリオ | 観点 | 期待値 |
|----------------|------|--------|
| 初回ログインフロー | 未ログイン→Firebase認証→ユーザー作成→サブスク画面遷移 | TBD（実装完了後に記載） |
| 既存ユーザーログインフロー | 未ログイン→Firebase認証→ダッシュボード遷移 | TBD（実装完了後に記載） |
| ログアウトフロー | ログイン済み→Firebaseサインアウト→ログイン画面遷移 | TBD（実装完了後に記載） |
| 保護されたページへのアクセス | 未ログイン状態で保護ページアクセス→ログイン画面へリダイレクト | TBD（実装完了後に記載） |

## 影響範囲一覧

### 機能影響範囲

| 関連機能 | 影響内容 |
|----------|----------|
| F-10-2 | 初回ログイン時にサブスクリプション登録へ遷移 |
| F-12-2 | ユーザー情報が自動登録される |
| 全機能 | 未ログイン時はアクセス制限 |

### コード影響範囲
🟢 **後回し可**

- フロントエンド: Firebase Authentication SDK、認証画面、セッション管理
- バックエンド: Firebase Admin SDK（ID Token検証）、セッション管理、ユーザー登録
- 外部サービス: Firebase Authentication、UnivaPay
- **決定: TBD**（実装時に確定）

## API仕様（参考）

### ログイン
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "id_token": "string"
}
```

### セッション確認
```http
GET /api/v1/auth/me
```

### ログアウト
```http
POST /api/v1/auth/logout
```

## 作業見積もり

### 見積もりサマリー

| 項目 | ストーリーポイント | 目安時間 |
|------|------------------|----------|
| **合計** | 22-24sp | 5.5-6時間 |

**目安**: 4sp = 1時間（実装＋単体テスト＋レビューを含む、あくまで参考値）

### タスク一覧

| タスク | ストーリーポイント | 備考 |
|--------|------------------|------|
| **バックエンド** | | |
| Firebaseプロジェクト設定 | 2 | コンソール設定・プロバイダ有効化・環境変数 |
| Firebase Admin SDK導入 | 2 | SDKインストール・初期化・ID Token検証 |
| ユーザー登録ロジック | 3 | users/user_settings作成・既存ユーザー判定 |
| UnivaPay連携 | 2-3 | カスタマー作成・エラーハンドリング |
| ログアウトAPI | 1 | Firebaseサインアウト・セッション削除 |
| セッション確認API | 1 | ユーザー情報返却 |
| **フロントエンド** | | |
| Firebase SDK導入 | 2 | SDKインストール・初期化・認証関数実装 |
| 認証画面実装 | 2 | UI実装・OAuthボタン・ローディング状態 |
| セッション管理 | 2 | Cookie管理・未ログイン時リダイレクト |
| **テスト** | | |
| 単体テスト | 3sp | Firebase/UnivaPayモック・認証フロー |
| E2Eテスト | 2-3sp | OAuthログイン/ログアウト/セッション確認の主要フロー |

### リスク要因

- **UnivaPay連携**: テスト環境での挙動確認が必要
- **Firebase設定**: プロバイダごとの設定差異（Appleは特に手順が多い）
- **セッション管理**: httpOnly Cookieの設定・ドメイン跨ぎ対応

### 依存関係

- Firebaseプロジェクト作成・設定完了後、実装可能
- UnivaPayテスト環境の事前準備が必要
