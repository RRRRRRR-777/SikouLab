# SicouLab

シコウラボ式会社のWebアプリケーション。

## プロダクト原則

- **受動的体験を最優先** - 読む体験が価値の中心
- **作りすぎない** - 初期は最小限、将来拡張を壊さない設計

## スコープ（v1.0.0）

**対象**: 記事、ニュース、アンケート、米国株銘柄、サブスク、いいね、アナリティクス
**対象外**: チャット、コメント、メール認証、Discord連携

## 技術スタック

| 項目 | 技術 |
|------|------|
| フロントエンド | Next.js, TypeScript, Tailwind CSS |
| バックエンド | Go (Clean Architecture) |
| データベース | PostgreSQL |
| インフラ | Google Cloud (Cloud Run/GKE) |

## クイックスタート

### 全サービスをDockerで起動

`backend/.env` を用意した上で（`backend/.env.sample` 参照）:

```bash
docker compose -f docker-compose.yml up -d
```

| サービス | URL |
|---------|-----|
| フロントエンド | http://localhost:3000 |
| バックエンドAPI | http://localhost:8080 |
| PostgreSQL | localhost:5432 |

### ホストで直接起動

```bash
cd backend && make run      # バックエンド（:8080）
cd frontend && make dev     # フロントエンド（:3000）
```

詳細な手順は [環境構築ガイド](./docs/setup.md) を参照してください。

## Docker Compose ファイル

| ファイル | 用途 |
|---------|------|
| `docker-compose.yml` | ローカル開発環境（frontend / backend / postgres） |

> 本番デプロイは Cloud Run + GitHub Actions で行います。

## CI/CD

| ワークフロー | トリガー | 実行内容 |
|-------------|---------|---------|
| **repomap.yml** | push（main/featureブランチ）、PR | REPO_MAP.md自動更新 |
| **backend-test.yml** | push（main）、PR（backend配下） | Go 1.25でバックエンド単体テスト実行 |
| **frontend-test.yml** | push（main）、PR（frontend配下） | Node.js 24でフロントエンド単体テスト実行 |
| **build-check.yml** | PR | バックエンド（Go 1.25）・フロントエンド（Node.js 24）のビルド検証 |
| **doc-check.yml** | PR: ドキュメントコメントチェック<br>push（main）: Cloudflare Pagesへデプロイ | - バックエンド: Go 1.22でdocコメント検証<br>- フロントエンド: Node.js 20でlint実行<br>- デプロイ: TypeDoc生成・docsify・Swagger UI |
| **docker-build.yml** | push（main）、PR（Docker関連） | Docker Compose V2でコンテナビルド検証 |

- **バージョン**: Go 1.25、Node.js 24がメイン（doc-check.ymlのみNode.js 20・Go 1.22使用）
- **デプロイ先**: Cloudflare Pages（`sicoulab-docs` プロジェクト）

## ドキュメント

- [環境構築ガイド](./docs/setup.md) - 必須ツール、初期化手順、環境変数、ローカル起動方法
- [開発ガイドライン](./docs/development_guidelines.md) - 技術スタック、CI/CD、コーディング規約
- [ADR一覧](./docs/adr/) - 重要な設計決定記録

## ディレクトリ構成

```
.
├── frontend/     # Next.jsアプリケーション
├── backend/      # Goアプリケーション
├── docs/         # 要件定義書、設計書、ガイドライン
└── scripts/      # セットアップスクリプト
```

## ライセンス

Copyright (c) 2025 SicouLab
