# SicouLab バックエンド

Clean Architecture で実装された Go 1.25 のバックエンドアプリケーション。

## 技術スタック

| 項目 | 技術 |
|------|------|
| 言語 | Go 1.25.3 |
| アーキテクチャ | Clean Architecture |
| データベース | PostgreSQL |
| ロガー | zerolog |
| 設定管理 | godotenv |

## 起動

```bash
# 開発サーバー
make run

# バイナリビルド
make build
```

## Docker

```bash
# イメージビルド
make docker-build

# または docker compose（リポジトリルートから）
docker compose up backend
```

## 開発コマンド

```bash
make lint      # golangci-lint
make fmt       # gofmt + goimports
make test      # 単体テスト
make check-doc # ドキュメントコメントチェック
```

## マイグレーション

```bash
make migrate-up   # マイグレーション適用
make migrate-down # マイグレーション戻し
```

## CI

- **backend-test.yml**: Go 1.25で単体テスト実行（push to main / PR）
- **build-check.yml**: ビルド検証（PR）
- **doc-check.yml**: ドキュメントコメントチェック（PR）

詳細は [環境構築ガイド](../docs/setup.md) を参照。
