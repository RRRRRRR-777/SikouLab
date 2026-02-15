# SikouLab ドキュメント

## サービス定義
- [サービス定義書](docs/service.md)

## プロジェクト管理
- [環境構築ガイド](docs/setup.md)
- [開発ガイドライン](docs/development_guidelines.md)
- [ドキュメントガイドライン](docs/documentation_guidelines.md)

## バージョン別ドキュメント
- **v1.0.0**
  - [要件定義書](docs/versions/1_0_0/requirements.md)
  - [基本設計書](docs/versions/1_0_0/system-design.md)
  - [データ設計書](docs/versions/1_0_0/data-model.md)

## ADR（技術選定記録）
- [001: UIライブラリ](docs/adr/001-ui-library.md)
- [002: フォーム管理](docs/adr/002-form-management.md)
- [003: データフェッチ](docs/adr/003-data-fetching.md)
- [004: バックエンドDB](docs/adr/004-backend-database.md)
- [005: staging環境](docs/adr/005-staging-environment.md)
- [006: デプロイフロー](docs/adr/006-deploy-flow.md)
- [007: ログ監視](docs/adr/007-log-monitoring.md)
- [008: API連携](docs/adr/008-api-integration.md)
- [009: フロントエンドディレクトリ構成](docs/adr/009-frontend-directory-structure.md)
- [010: 品質計測](docs/adr/010-quality-metrics.md)
- [011: 権限モデル](docs/adr/011-rbac-vs-abac.md)
- [012: Repo Map](docs/adr/012-repo-map.md)
- [013: Markdownライブラリ](docs/adr/013-markdown-library.md)
- [014: 検索エンジン](docs/adr/014-search-engine.md)

## 機能別詳細仕様
- **認証機能（F-01）**
  - [ログイン](docs/functions/auth/login.md)
- **ダッシュボード機能（F-02）**
  - [ダッシュボード](docs/functions/dashboard/home.md)
- **銘柄機能（F-03）**
  - [個別銘柄ページ](docs/functions/stock/home.md)
- **記事機能（F-04）**
  - [記事ホーム](docs/functions/article/home.md)
  - [記事詳細](docs/functions/article/detail.md)
  - [記事作成・編集](docs/functions/article/create-edit.md)
  - [予約投稿](docs/functions/article/schedule.md)
- **ニュース機能（F-05）**
  - [ニュースホーム](docs/functions/news/home.md)
  - [ニュース自動取得・翻訳](docs/functions/news/fetch.md)
  - [ジャンル詳細ページ](docs/functions/news/genre-detail.md)
- **アンケート機能（F-06）**
  - [アンケート](docs/functions/poll/home.md)
- **検索機能（F-07）**
  - [検索](docs/functions/search/home.md)
- **ブックマーク機能（F-08）**
  - [ブックマーク](docs/functions/bookmark/home.md)
- **設定機能（F-10）**
  - [設定画面](docs/functions/settings/home.md)
- **ニュースレター機能（F-11）**
  - [ニュースレター機能](docs/functions/newsletter/home.md)
- **ポートフォリオ機能（F-09）**
  - [ポートフォリオ](docs/functions/portfolio/home.md)
- **管理者機能（F-12）**
  - [管理者ページ](docs/functions/admin/home.md)

---

## 開発ルール（.claude/）

### プロジェクト指針
- [CLAUDE.md](.claude/CLAUDE.md)

### rules
- [workflow.md](.claude/rules/workflow.md)
- [docs.md](.claude/rules/docs.md)
- [frontend.md](.claude/rules/frontend.md)
- [backend.md](.claude/rules/backend.md)

### skills
- [adr](.claude/skills/adr/SKILL.md)
- [feature-spec](.claude/skills/feature-spec/SKILL.md)
- [feature-spec-estimate-cost](.claude/skills/feature-spec-estimate-cost/SKILL.md)
- [frontend-design](.claude/skills/frontend-design/SKILL.md)
- [go-standards](.claude/skills/go-standards/SKILL.md)
- [nextjs-15](.claude/skills/nextjs-15/SKILL.md)
- [requirements](.claude/skills/requirements/SKILL.md)
- [review-spec](.claude/skills/review-spec/SKILL.md)
- [system-design](.claude/skills/system-design/SKILL.md)
- [workflow](.claude/skills/workflow/SKILL.md)

---

[ドキュメント一覧](docs/SUMMARY.md)
