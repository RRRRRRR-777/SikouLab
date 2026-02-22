# SicouLab フロントエンド

Next.js 16 (App Router) で実装されたフロントエンドアプリケーション。

## 技術スタック

| 項目             | 技術                        |
| ---------------- | --------------------------- |
| フレームワーク   | Next.js 16.1.6 (App Router) |
| 言語             | TypeScript 5                |
| スタイリング     | Tailwind CSS 4              |
| UIコンポーネント | Radix UI                    |
| フォーム         | React Hook Form + Zod       |
| データフェッチ   | @tanstack/react-query       |
| テスト           | Vitest + Testing Library    |
| ドキュメント     | TypeDoc                     |

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

## CI

- **frontend-test.yml**: Node.js 24で単体テスト実行（push to main / PR）
- **build-check.yml**: ビルド検証（PR）
- **doc-check.yml**: ESLint + TypeDoc生成（PR: チェック、push to main: デプロイ）

詳細は [環境構築ガイド](../docs/setup.md) を参照。
