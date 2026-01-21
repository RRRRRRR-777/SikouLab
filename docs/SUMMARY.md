# ドキュメント一覧

## プロジェクト管理ドキュメント
* docs/development_guidelines.md
  * 概要: 開発に関する技術的なガイドライン。技術スタック（Next.js + Go）、リポジトリ構成（モノレポ）、インフラ（Google Cloud）、外部サービス連携（OAuth, Stripe, Metabase）を定義。
    テスト方針、CI/CD、ブランチ戦略、バージョニング戦略、デプロイフローを記載。
* docs/documentation_guidelines.md
  * 概要: ドキュメント作成・管理・運用の方針を定義。ディレクトリ構成、ファイル命名規則、更新ルール、ドキュメント種別を記載。
    SUMMARY.md更新ルールやAIエージェント向けガイドラインを含む。

## ADR（Architecture Decision Records）
技術選定の意思決定記録。選定理由と検討した選択肢を記載。

* docs/adr/001-ui-library.md
  * 概要: UIライブラリの選定。shadcn/ui + Tailwind CSS + Figma Kitを採用。
    カスタマイズ自由度とFigma連携を重視した選定理由を記載。
* docs/adr/002-form-management.md
  * 概要: フォーム管理ライブラリの選定。React Hook Form + Zodを採用。
    プロジェクトで想定されるフォーム一覧と選定理由を記載。
* docs/adr/003-data-fetching.md
  * 概要: データフェッチライブラリの選定。TanStack Queryを採用。
    無限スクロール、楽観的更新、リアルタイム更新の要件を満たす選定理由を記載。
* docs/adr/004-backend-database.md
  * 概要: バックエンドDBアクセス・マイグレーションの選定。sqlx + golang-migrateを採用。
    SQL重視のアプローチとSQLファイル管理の選定理由を記載。
* docs/adr/005-staging-environment.md
  * 概要: staging環境の方針。初期は構築せず、必要に応じて後から追加。
    dev → prod直接のフローと、staging追加を検討するタイミングを記載。
* docs/adr/006-deploy-flow.md
  * 概要: デプロイフローの選定。GitHub Environment承認を採用。
    トリガー、承認プロセス、ロールバック方針を記載。

## バージョン別ドキュメント（versions/）
* docs/versions/1_0_0/requirements.md
  * 概要: v1.0.0の要件定義書（テンプレート準拠）。「記事購読型コミュニティ」としての核心価値を明記。成功基準（記事閲覧率・継続率含む5指標）、
    ペルソナ（ROM専・投資初心者・アクティブ投資家）、ユーザーストーリー（US-01〜09）、機能一覧（F-01〜F-14、子機能含む計30項目）を整理。
* docs/versions/1_0_0/requirements_definition_document.md
  * 概要: シコウラボプラットフォームの要件定義書（v1.0.0）。背景・目的、ユーザーロール、プラン・課金、機能要件（記事・ニュース、個別銘柄ページ、検索、アンケート）を定義。
    ROM専ユーザーを前提とした「読むだけで価値がある体験」を最優先する設計思想を記載。
* docs/versions/1_0_0/system_design.md
  * 概要: シコウラボプラットフォームの基本設計書（v1.0.0）。全体アーキテクチャ（Next.js + Go）、認証・認可（OAuth）、プラン・権限制御の考え方、画面構成、
    機能別基本設計（記事・ニュース・銘柄ページ・アンケート・管理画面）、データ取得方式、アナリティクス設計を記載。
* docs/versions/1_0_0/er_diagram.md
  * 概要: データベース設計（ER図、v1.0.0）。users, plans, articles, news, likes, stocks, stock_news, polls, poll_options, poll_votes の各テーブル定義と
    それらのリレーション（Mermaid形式のER図）を記載。

## 機能別詳細仕様（functions/）
* TBD: 今後機能実装時に追加予定
