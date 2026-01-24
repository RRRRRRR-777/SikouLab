# 環境構築ガイド

## 1. 必須ツール

| ツール | バージョン | 備考 |
|--------|-----------|------|
| Node.js | v24.13.0 | LTS推奨（v22.22.0, v20.20.0 も可） |
| Go | 1.25.6 | |
| Docker Compose | v5.0.1+ | Docker Desktop/Engine に含まれるプラグイン |
| pnpm | 10.27.0 | パッケージマネージャ |
| golang-migrate | latest | DBマイグレーション |
| golangci-lint | latest | Go静的解析 |

### ツールインストール

```bash
# golang-migrate（macOS）
brew install golang-migrate

# golang-migrate（Go install）
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# golangci-lint（macOS）
brew install golangci-lint

# golangci-lint（Go install）
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

→ 詳細: [golang-migrate](https://github.com/golang-migrate/migrate), [golangci-lint](https://golangci-lint.run/welcome/install/)

## 2. 前提ファイル

以下のファイルが必要（実装時に作成）:

| ファイル | 用途 | 必須スクリプト/ターゲット |
|----------|------|--------------------------|
| `docker-compose.yml` | ローカルDB定義 | postgres サービス |
| `frontend/Makefile` | FEビルド/開発 | dev, build, lint, fmt, test, test-e2e, test-vrt, docker-build |
| `backend/Makefile` | BEビルド/開発 | run, build, lint, fmt, test, migrate-up, migrate-down, docker-build |
| `frontend/package.json` | FE依存管理 | scripts: dev, build, lint, format, test, test:e2e, test:vrt |
| `frontend/.env.sample` | FE環境変数テンプレート | - |
| `backend/.env.sample` | BE環境変数テンプレート | - |

## 3. リポジトリ初期化

```bash
# クローン
git clone <repository-url>
cd SikouLab

# フロントエンド依存インストール
cd frontend && pnpm install

# バックエンド依存インストール
cd ../backend && go mod download
```

## 4. 環境変数

| 変数名 | 用途 | 必須 | 例 |
|--------|------|------|-----|
| `NEXT_PUBLIC_API_BASE_URL` | API基底URL | ○ | `http://localhost:8080` |
| `DATABASE_URL` | PostgreSQL接続 | ○ | `postgres://postgres:postgres@localhost:5432/sikoulab` |
| TBD | TBD | - | - |

**`.env.sample`**: `frontend/.env.sample`, `backend/.env.sample`

各ディレクトリの `.env.sample` をコピーして `.env` を作成し、値を設定する。

## 5. ローカルサービス起動

### DB起動

```bash
docker compose up -d
```

| 項目 | 値 |
|------|-----|
| 設定ファイル | `docker-compose.yml`（ルート） |
| DB種別 | PostgreSQL |
| サービス名 | `postgres` |
| ポート | 5432 |
| ユーザー/パスワード | `postgres` / `postgres` |
| DB名 | `sikoulab` |

**注意**: `docker-compose.yml` の設定と `backend/.env` の `DATABASE_URL` が一致していること。

### マイグレーション

```bash
cd backend && make migrate-up
```

Makefile: `backend/Makefile`

### シードデータ

TBD（必要に応じて追記）

## 6. アプリケーション起動

| 対象 | コマンド | ポート | ヘルスチェック |
|------|---------|--------|---------------|
| frontend | `cd frontend && make dev` | 3000 | `localhost:3000` |
| backend | `cd backend && make run` | 8080 | `localhost:8080/health` |

## 7. 開発コマンド

| 対象 | lint | format | test |
|------|------|--------|------|
| frontend | `make lint` | `make fmt` | `make test` |
| backend | `make lint` | `make fmt` | `make test` |

→ 詳細は [開発ガイドライン](./development_guidelines.md) を参照

### Makefile ターゲット（frontend/Makefile）

| ターゲット | コマンド | 説明 |
|-----------|---------|------|
| `dev` | `pnpm dev` | 開発サーバー起動 |
| `build` | `pnpm build` | 本番ビルド |
| `start` | `pnpm start` | 本番サーバー起動 |
| `test` | `pnpm test` | 単体テスト実行 |
| `lint` | `pnpm lint` | ESLint実行 |
| `fmt` | `pnpm format` | Prettierフォーマット |
| `test-e2e` | `pnpm test:e2e` | E2Eテスト実行 |
| `test-vrt` | `pnpm test:vrt` | VRT実行 |
| `docker-build` | `docker build -t sikoulab-web .` | Dockerイメージビルド |

### Makefile ターゲット（backend/Makefile）

| ターゲット | コマンド | 説明 |
|-----------|---------|------|
| `run` | `go run ./cmd/api` | 開発サーバー起動 |
| `build` | `go build -o bin/api ./cmd/api` | バイナリビルド |
| `test` | `go test ./...` | 単体テスト実行 |
| `lint` | `golangci-lint run` | 静的解析 |
| `fmt` | `gofmt -w .` | コードフォーマット |
| `migrate-up` | `migrate -path db/migrations -database "$(DATABASE_URL)" up` | マイグレーション適用 |
| `migrate-down` | `migrate -path db/migrations -database "$(DATABASE_URL)" down` | マイグレーション戻し |
| `docker-build` | `docker build -t sikoulab-api .` | Dockerイメージビルド |

## 8. E2E / VRT

| 項目 | コマンド | 前提条件 |
|------|---------|---------|
| E2Eテスト | `cd frontend && make test-e2e` | frontend/backend 起動済み |
| VRT | `cd frontend && make test-vrt` | frontend/backend 起動済み |

→ VRTベースライン画像は Cloud Storage に保存（[開発ガイドライン](./development_guidelines.md) 参照）

## 9. トラブルシューティング

| 問題 | 解決策 |
|------|--------|
| TBD | TBD |

---

→ 技術スタック・CI/CD詳細: [開発ガイドライン](./development_guidelines.md)
→ ADR一覧: [_sidebar.md](./_sidebar.md)
