# ドキュメント一覧

## プロジェクト管理ドキュメント
* docs/development_guidelines.md
  * 概要: 開発に関する技術的なガイドライン。技術スタック（Next.js + Go）、リポジトリ構成（モノレポ）、インフラ（Google Cloud）、外部サービス連携（OAuth, Stripe, Metabase）を定義。
    テスト方針、CI/CD、ブランチ戦略、バージョニング戦略を記載。使用ライブラリ（shadcn/ui, Tailwind CSS, ORM等）の検討事項を含む。
* docs/documentation_guidelines.md
  * 概要: ドキュメント作成・管理・運用の方針を定義。ディレクトリ構成、ファイル命名規則、更新ルール、ドキュメント種別を記載。
    SUMMARY.md更新ルールやAIエージェント向けガイドラインを含む。

## バージョン別ドキュメント（versions/）
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
