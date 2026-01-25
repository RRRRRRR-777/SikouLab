# SikouLab

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

環境構築の手順は[環境構築ガイド](./docs/setup.md)を参照してください。

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

Copyright (c) 2025 SikouLab
