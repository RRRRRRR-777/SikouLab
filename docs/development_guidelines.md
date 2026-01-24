# 開発ガイドライン

## 1. 技術スタック

### 1.1 フロントエンド
- **言語・フレームワーク**: Next.js
- **開発ツール**: Claude Code / Cursor
- **開発フロー**:
  - Figmaでデザインモック作成
  - Claude Code / Cursorでフロントエンドを実装
  - モック確認用環境:
    - 第一候補: Vercel（個人利用無料、自動デプロイ）
    - 第二候補: ラズパイ + Cloudflare（コスト最小化）
    - 第三候補: 無理に工数をかける必要もないので、その場合は採用しないものとする
- **使用ライブラリ**: → [ADR-001](./adr/001-ui-library.md), [ADR-002](./adr/002-form-management.md), [ADR-003](./adr/003-data-fetching.md)
  - UIコンポーネント: shadcn/ui + Figma Kit
  - CSSフレームワーク: Tailwind CSS
  - フォーム管理: React Hook Form + Zod
  - データフェッチ: TanStack Query
    - 必要機能: 無限スクロール、楽観的更新、リアルタイム更新
- **記事エディタ**: Markdown形式

### 1.2 バックエンド
- **言語**: Go
- **アーキテクチャ**: Clean Architecture
- **開発フロー**:
  - docs配下の設計書を元にOpenAPIを作成
  - OpenAPIを元にGoでバックエンドを実装
- **使用ライブラリ**: → [ADR-004](./adr/004-backend-database.md)
  - DBアクセス: sqlx
  - マイグレーション: golang-migrate（SQLファイル管理）

### 1.3 リポジトリ構成
- **構成管理**: モノレポ
  ```
  /frontend   - Next.jsアプリケーション
  /backend    - Goアプリケーション
  /docs       - 要件定義書、設計書、ガイドライン
  ```

### 1.4 インフラ
- **クラウドプラットフォーム**: Google Cloud
- **デプロイ先**: Cloud Run vs GKE（費用対効果計算後に決定）
  - 想定ユーザー数:
    - 現在: 1200人（アクティブユーザー800人）
    - 3年後: 2400人
- **ファイルストレージ**: Cloud Storage
  - 用途: 本番環境用画像、VRTベースライン画像
  - バケット/ディレクトリで用途を分離
- **CDN**: Cloud CDN（画像配信用）
- **ログ監視**: → [ADR-007](./adr/007-log-monitoring.md)
  - ログ収集: Cloud Logging
  - 監視・アラート: Cloud Monitoring
  - 通知先: Cloud Consoleアプリ（iOSプッシュ通知）
  - 検知対象: 500エラー、サービスダウン、外部API通信失敗

### 1.5 環境構成
→ [ADR-005](./adr/005-staging-environment.md)
- **prod（本番環境）**: 必須
- **dev（開発環境）**: 用意できたら望ましい
- **staging**: 初期は無し、必要に応じて後から追加

### 1.6 外部サービス連携
- **認証**: OAuth（Google / Apple 等）
- **課金**: Stripe
- **アナリティクス**: Metabase

### 1.7 API連携
→ [ADR-008](./adr/008-api-integration.md)

| 項目 | 決定 |
|------|------|
| 通信ライブラリ | axios |
| タイムアウト | 10秒 |
| リトライ | TanStack Queryに任せる |
| 環境変数 | `NEXT_PUBLIC_API_BASE_URL` |

- **認証（Web）**: JWT + httpOnly Cookie（SameSite=Lax）
- **認証（Mobile）**: 将来対応時にGo APIで両対応（Cookie/Header）
- **型定義**: OpenAPIから手動でTypeScript型作成（camelCase統一）
- **エラー処理**: 共通ラッパーで一元化（401→ログイン、500→トースト）

### 1.8 フロントエンドディレクトリ構成
→ [ADR-009](./adr/009-frontend-directory-structure.md)

```
frontend/
├── app/                      # ルーティング専用
├── components/
│   ├── ui/                  # shadcn/ui
│   └── features/            # 機能別コンポーネント
├── hooks/                   # カスタムフック
├── lib/                     # ライブラリ初期化
├── utils/                   # 純粋関数
├── stores/                  # グローバルステート管理
├── types/                   # 型定義
└── public/                  # 静的ファイル
```

- **app/**: ルーティング専用。ビジネスロジックを置かない
- **components/**: UIコンポーネント。features/で機能別に整理
- **lib/utils**: libはライブラリ初期化、utilsは純粋関数
- **stores/**: グローバルステート管理（Zustand等）

## 2. テスト方針

### 2.1 バックエンド
- **単体テスト**: 実施

### 2.2 フロントエンド
- **単体テスト**: 実施
- **E2Eテスト**: Playwright
- **VRT（Visual Regression Testing）**: Playwright標準機能（toHaveScreenshot）
  - ベースライン画像はCloud Storageに保存
  - CI環境はDockerで統一（フォント差異対策）

### 2.3 TDD（テスト駆動開発）
- **進め方**: Red → Green → Refactor サイクルで開発
- **対象**: バックエンド・フロントエンドのロジック
- **方針**: 基本的に必須（例外は状況に応じて判断）

## 3. CI/CD

### 3.1 Pre-commit（ローカル）

コミット前にlint/formatを自動実行し、CIでの失敗を未然に防ぐ。

| 対象 | 項目 | ツール |
|------|------|--------|
| Frontend | lint | ESLint |
| Frontend | format | Prettier |
| Backend | lint | golangci-lint（未導入時は go vet） |
| Backend | format | gofmt |

**セットアップ**
```bash
./scripts/setup-hooks.sh
```

**スキップ方法**（緊急時のみ）
```bash
git commit --no-verify
```

### 3.2 CI（Continuous Integration）

| 項目 | ticket | dev | main |
|------|--------|-----|------|
| lint | ✓ | ✓ | ✓ |
| format | ✓ | ✓ | ✓ |
| 単体テスト | ✓ | ✓ | ✓ |
| docker build（BE） | ✓ | ✓ | ✓ |
| E2Eテスト | - | - | ✓ |
| VRT | - | - | ✓ |

### 3.3 CD（Continuous Deployment）
→ [ADR-006](./adr/006-deploy-flow.md)
- **トリガー**: GitHub Environment承認
- **承認**: Environment Protection Rulesで承認者指定
- **ロールバック**: 前バージョンのコンテナイメージに戻す

## 4. コーディング規約

- **フロントエンド**: `.claude/skills/nextjs-15/SKILL.md`
- **バックエンド**: `.claude/skills/go-standards/SKILL.md`

## 5. コード内ドキュメント

関数・型の仕様をコード内に記述し、TypeDoc/doc2goでGitHub Pagesに公開する。

### 5.1 フロントエンド（TypeScript）

| 項目 | 内容 |
|------|------|
| 形式 | **JSDoc** |
| ツール | **TypeDoc** |
| 公開先 | GitHub Pages |

**参考**: [TypeDoc公式](https://typedoc.org), [GitHub Pages使用例](https://github.com/TypeStrong/typedoc-site)

### 5.2 バックエンド（Go）

| 項目 | 内容 |
|------|------|
| 形式 | **godoc** |
| ツール | **doc2go** |
| 公開先 | GitHub Pages |

**参考**: [doc2go GitHub](https://github.com/abhinav/doc2go)

### 5.3 CI/CD

mainブランチへのマージ時に自動生成・デプロイする。

## 6. ブランチ戦略

### 6.1 ブランチ種別
- **main**: 本番ブランチ
- **dev-x.x.x.x**: 開発ブランチ（例: dev-1.2.5.0）
- **ticket-xxx**: チケットブランチ（例: ticket-155）

### 6.2 ブランチフロー
- チケットブランチで開発
- 開発ブランチにマージ
- 本番ブランチ（main）にマージしてリリース

## 7. バージョニング戦略

### 7.1 バージョン表記
- **形式**: a.b.c.d（例: 1.2.3.4）

### 7.2 各桁の意味
- **a（メジャーバージョン）**: アプリの大きな仕様変更やアップデートで変更
- **b（マイナーバージョン）**: 仕様変更や機能追加で変更（基本的にこのバージョンが上がる）
- **c（パッチバージョン）**: パッチリリース
- **d（マイナーパッチ）**: マイナーなパッチリリース（基本的に変更されない）

## 8. デプロイフロー
→ [ADR-006](./adr/006-deploy-flow.md)

```
1. ticket → dev マージ（CI実行: lint, format, 単体テスト, docker build）
2. dev → main のPRを作成
3. PRマージ（CI実行: 上記 + E2E, VRT）
4. GitHub Actionsが `environment: production` に到達
5. 承認者に通知 → GitHub UI上で承認
6. デプロイ実行
```
