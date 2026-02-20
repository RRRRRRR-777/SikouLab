# SicouLab フロントエンド

Next.js 16 (App Router) で実装されたフロントエンドアプリケーション。

## 起動

```bash
# 開発サーバー
make dev

# 本番ビルド
make build

# 本番サーバー起動（ビルド後）
make start
```

## Docker

```bash
# イメージビルド
make docker-build

# または docker compose（リポジトリルートから）
docker compose up frontend
```

## 開発コマンド

```bash
make lint     # ESLint
make fmt      # Prettier
make test     # 単体テスト
make test-e2e # E2Eテスト
```

詳細は [環境構築ガイド](../docs/setup.md) を参照。
