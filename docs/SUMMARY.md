# ドキュメント一覧

## サービス定義
* docs/service.md
  * 概要: サービス定義書。コンセプト、ターゲット、ミッション・ビジョン・バリュー、提供価値、サービス品質（パフォーマンス・可用性・セキュリティ・UX原則）、ブランドガイドラインを定義。

## プロジェクト管理ドキュメント
* docs/setup.md
  * 概要: 環境構築の手順書。必須ツール・バージョン、リポジトリ初期化、環境変数、ローカルサービス起動、開発コマンドを記載。
* docs/development_guidelines.md
  * 概要: 開発に関する技術的なガイドライン。技術スタック（Next.js + Go）、リポジトリ構成（モノレポ）、インフラ（Google Cloud）、外部サービス連携（OAuth, Stripe, Metabase）、API連携方針、フロントエンドディレクトリ構成を定義。
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
* docs/adr/007-log-monitoring.md
  * 概要: ログ監視方式の選定。Cloud Logging + Cloud Monitoring（GCP完結）を採用。
    障害発生時のフロー、通知先（Cloud Consoleアプリ）、将来拡張（LLM分析）を記載。
* docs/adr/008-api-integration.md
  * 概要: フロントエンドとバックエンドのAPI連携方針。axios + httpOnly Cookie（SameSite=Lax）を採用。
    認証方式、型定義管理、エラーハンドリング、将来のモバイル対応方針を記載。
* docs/adr/009-frontend-directory-structure.md
  * 概要: Next.js 15（App Router）のベストプラクティスに沿ったフロントエンドディレクトリ構成。
    app/はルーティング専用、components/features/で機能別整理、lib/utils/storesの役割を定義。
* docs/adr/010-quality-metrics.md
  * 概要: 品質計測の方針。パフォーマンス（Lighthouse CI + Cloud Monitoring）、可用性（SLA 99.5%、Uptime Check）、
    リグレッション防止（市場障害・ロールバック、Jira + Metabase）の計測フロー・アラート設定を記載。

## バージョン別ドキュメント（versions/）
* docs/versions/1_0_0/requirements.md
  * 概要: v1.0.0の要件定義書。「記事購読型コミュニティ」としての核心価値を明記。成功基準（記事閲覧率・継続率含む5指標）、
    ペルソナ（ROM専・投資初心者・アクティブ投資家）、ユーザーストーリー（US-01〜12）、機能一覧（F-01〜F-12、子機能含む計51項目）。
    機能ID体系: F-01ログイン、F-02ダッシュボード、F-03個別銘柄、F-04記事、F-05ニュース、F-06アンケート、F-07検索、F-08ブックマーク、
    F-09ポートフォリオ、F-10設定、F-11ニュースレター、F-12管理者ページ。
    UI/UX共通要件: サイドバー（記事タブ折りたたみ・状態永続化）、ダークモード、レスポンシブ、カラーコード。
* docs/versions/1_0_0/system-design.md
  * 概要: v1.0.0の基本設計書。機能一覧・相関図・依存関係（2.機能一覧）、画面一覧・遷移図（3.画面設計、S-01〜S-23、S-02はS-02-1〜S-02-5のサブ画面を含む、グローバル検索は補足記載、FigmaノードID付き）、
    共通コンポーネント（サイドバー）、外部インターフェース（5.外部API連携）、権限マトリクス（6.ロール×機能、サブスク状態×アクセス）を定義。
    データ設計（4.データ設計）は別ドキュメント（system_datas.md）に切り出し。
    要件定義書のID体系（F-01〜F-12）に準拠。管理者ページ（F-12、画面S-18〜S-23）、インサイダー取引（F-03-4）、ニュースレター関連テーブル、user_settingsテーブル追加。
* docs/versions/1_0_0/system_datas.md
  * 概要: v1.0.0のデータ設計書。ER図（全34テーブルのリレーション）、テーブル一覧、カラム定義（全テーブルの詳細なカラム仕様）を定義。

## 機能別詳細仕様（functions/）
* docs/functions/dashboard.md
  * 概要: ダッシュボード（F-02）の詳細設計書。ログイン後のホーム画面として、相場インデックス表示（S&P500等）、記事・ニュースハイライト、
    ポートフォリオウィジェット、関連ニュースウィジェット、相場異変レーダー（WANT）の仕様を定義。
    画面設計図・レイアウト構成・関連テーブル・シーケンス図・機能要件・非機能要件・テストケースを含む。
