# ADR-018: E2Eテスト認証方式の選定

## ステータス
承認済み

## コンテキスト

本プロジェクトではFirebase OAuthによる認証を採用している。E2Eテスト（Playwright）を導入するにあたり、テスト実行時の認証をどう扱うかが課題となった。

Google OAuthはボット検知（reCAPTCHA）や2FAにより自動化が困難であり、CI環境ではほぼ動作しない。テスト用の認証方式を別途確立する必要がある。

### 制約

- Firebase Authentication（Google / Apple OAuth）を使用
- CI環境（GitHub Actions）での自動実行を将来的に想定
- 本番コードへのセキュリティリスクを持ち込まない
- セットアップが簡素であること

## 決定

**Firebase Auth Emulator** を採用する。

### 詳細構成

| 用途 | 技術 | 備考 |
|------|------|------|
| 認証エミュレーション | Firebase Auth Emulator | `localhost:9099` |
| テスト用ログイン | メール/パスワード | OAuth不要 |
| 環境分岐 | `NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST` | `true` でEmulatorに接続 |

```bash
# Emulator起動
firebase emulators:start --only auth

# または docker-compose で起動（将来）
```

## 検討した選択肢

### 比較表

| 選択肢 | CI対応 | 実装難度 | セキュリティ | メンテナンス | 認証フロー検証 |
|--------|--------|---------|------------|------------|--------------|
| **Firebase Auth Emulator** | ◎ | ○ | ◎ | ○ | ○ |
| Admin SDK カスタムトークン | ◎ | ○ | ○ | ○ | △ |
| storageState + 手動OAuth | △ | ◎ | ○ | △ | ◎ |
| バックエンド認証バイパス | ◎ | ◎ | ✕ | ◎ | ✕ |

### 詳細評価

#### A. storageState + 手動OAuth
- メリット: Playwright公式パターン、追加ライブラリ不要
- デメリット: Google OAuthのボット検知でCI環境ではほぼ動かない。セッション期限切れ時に手動再実行が必要

#### B. Firebase Auth Emulator（採用）
- メリット: Firebase公式推奨。本番Firebaseにトラフィックを送らない。CI完全対応。テストデータが本番を汚染しない
- デメリット: firebase-toolsのインストール・設定が必要。アプリコードにEmulator接続の分岐が必要

#### C. Admin SDK カスタムトークン
- メリット: 本物のFirebase認証状態を生成。`nearform/playwright-firebase`プラグインあり
- デメリット: サービスアカウント秘密鍵の管理が必要。実装がやや複雑

#### D. バックエンド認証バイパス
- メリット: 実装が最も簡単
- デメリット: 本番コードにセキュリティホールを埋め込む構造。CVE-2025-29927のような脆弱性リスク。非推奨

## 理由

### Firebase Auth Emulatorが採用された理由
1. **Firebase公式推奨**: Firebase Developers公式ドキュメント・ブログで標準手法として紹介されている
2. **セキュリティ**: 本番コードへの変更が`connectAuthEmulator()`の条件分岐のみ。サービスアカウント秘密鍵の管理不要
3. **CI完全対応**: `firebase emulators:start --only auth`で完結。外部依存なし
4. **認証フロー検証**: メール/パスワード方式だがFirebase AuthのAPIを実際に通すため、認証ミドルウェアの動作検証が可能

## 影響

- **フロントエンド**: Firebase初期化コードに`connectAuthEmulator()`の分岐追加
- **開発環境**: `firebase-tools`のインストールが必要（`npm install -g firebase-tools` または devDependencies）
- **設定ファイル**: `firebase.json`にEmulator設定を追加
- **E2Eテスト**: `auth.setup.ts`をEmulator対応のメール/パスワードログインに変更

## 参考情報

- [Connect your app to the Authentication Emulator | Firebase](https://firebase.google.com/docs/emulator-suite/connect_auth)
- [Run CI tests using the Firebase Emulator Suite | Firebase Developers](https://medium.com/firebase-developers/run-continuous-integration-tests-using-the-firebase-emulator-suite-9090cefefd69)
- [Authentication | Playwright](https://playwright.dev/docs/auth)
