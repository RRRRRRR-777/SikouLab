## API リファレンス
- [Frontend API](https://sicoulab-docs.pages.dev/frontend/)
- [Backend API](https://sicoulab-docs.pages.dev/backend/)

---

## サービス定義
- [サービス定義書](service.md)
- [環境構築ガイド](setup.md)
- [開発ガイドライン](development_guidelines.md)
- [ドキュメントガイドライン](documentation_guidelines.md)

## バージョン毎
- **v1.0.0**
  - [要件定義書](versions/1_0_0/requirements.md)
  - [基本設計書](versions/1_0_0/system-design.md)
  - [データ設計書](versions/1_0_0/data-model.md)
  - [インフラ設計書](versions/1_0_0/infrastructure.md)
  - [開発計画書](versions/1_0_0/development-plan.md)

## ADR（技術選定記録）
- [001: UIライブラリ](adr/001-ui-library.md)
- [002: フォーム管理](adr/002-form-management.md)
- [003: データフェッチ](adr/003-data-fetching.md)
- [004: バックエンドDB](adr/004-backend-database.md)
- [005: staging環境](adr/005-staging-environment.md)
- [006: デプロイフロー](adr/006-deploy-flow.md)
- [007: ログ監視](adr/007-log-monitoring.md)
- [008: API連携](adr/008-api-integration.md)
- [009: フロントエンドディレクトリ構成](adr/009-frontend-directory-structure.md)
- [010: 品質計測](adr/010-quality-metrics.md)
- [011: 権限モデル](adr/011-rbac-vs-abac.md)
- [012: Repo Map](adr/012-repo-map.md)
- [013: Markdownライブラリ](adr/013-markdown-library.md)
- [014: 検索エンジン](adr/014-search-engine.md)
- [015: インフラ環境構築](adr/015-deploy-platform.md)
- [016: Phase 0横断的決定事項](adr/016-phase0-cross-cutting-decisions.md)
- [017: CSRF対策](adr/017-csrf-protection.md)

## 機能仕様書
- **認証機能（F-01）**
  - [ログイン](functions/auth/login.md)
- **ダッシュボード機能（F-02）**
  - [ダッシュボード](functions/dashboard/home.md)
- **銘柄機能（F-03）**
  - [個別銘柄ページ](functions/stock/home.md)
- **記事機能（F-04）**
  - [記事ホーム](functions/article/home.md)
  - [記事詳細](functions/article/detail.md)
  - [記事作成・編集](functions/article/create-edit.md)
  - [予約投稿](functions/article/schedule.md)
- **ニュース機能（F-05）**
  - [ニュースホーム](functions/news/home.md)
  - [ニュース自動取得・翻訳](functions/news/fetch.md)
  - [ジャンル詳細ページ](functions/news/genre-detail.md)
- **アンケート機能（F-06）**
  - [アンケート](functions/poll/home.md)
- **検索機能（F-07）**
  - [検索](functions/search/home.md)
- **ブックマーク機能（F-08）**
  - [ブックマーク](functions/bookmark/home.md)
- **設定機能（F-10）**
  - [設定画面](functions/settings/home.md)
- **ニュースレター機能（F-11）**
  - [ニュースレター機能](functions/newsletter/home.md)
- **ポートフォリオ機能（F-09）**
  - [ポートフォリオ](functions/portfolio/home.md)
- **管理者機能（F-12）**
  - [管理者ページ](functions/admin/home.md)

---

## 開発ルール（.claude/）

### プロジェクト指針
- [CLAUDE.md](../.claude/CLAUDE.md)

### rules
- [workflow.md](../.claude/rules/workflow.md)
- [docs.md](../.claude/rules/docs.md)
- [frontend.md](../.claude/rules/frontend.md)
- [backend.md](../.claude/rules/backend.md)

### skills
- [adr](../.claude/skills/adr/SKILL.md)
- [feature-spec](../.claude/skills/feature-spec/SKILL.md)
- [feature-spec-estimate-cost](../.claude/skills/feature-spec-estimate-cost/SKILL.md)
- [frontend-design](../.claude/skills/frontend-design/SKILL.md)
- [go-standards](../.claude/skills/go-standards/SKILL.md)
- [nextjs-15](../.claude/skills/nextjs-15/SKILL.md)
- [requirements](../.claude/skills/requirements/SKILL.md)
- [review-spec](../.claude/skills/review-spec/SKILL.md)
- [system-design](../.claude/skills/system-design/SKILL.md)
- [workflow](../.claude/skills/workflow/SKILL.md)

---

[ドキュメント一覧](SUMMARY.md)
