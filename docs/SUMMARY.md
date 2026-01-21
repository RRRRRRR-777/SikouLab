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
* docs/versions/1_0_0/basic-design.md
  * 概要: v1.0.0の基本設計書。機能一覧・相関図・依存関係、画面一覧・遷移図、データ設計（ER図）、外部インターフェース、権限マトリクスを定義。
    要件定義書と機能設計書（feature-spec）の間に位置し、「何を作るか」を機能・画面・データの観点で整理。

## 機能別詳細仕様（functions/）
* TBD: 今後機能実装時に追加予定
