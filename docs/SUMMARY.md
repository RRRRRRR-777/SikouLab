# ドキュメント一覧
* docs/requirements_definition.md
  * 概要: シコウラボプラットフォームの要件定義書。背景・目的、ユーザーロール、プラン・課金、機能要件（記事・ニュース、個別銘柄ページ、検索、アンケート）を定義。
    ROM専ユーザーを前提とした「読むだけで価値がある体験」を最優先する設計思想を記載。
* docs/system_design.md
  * 概要: シコウラボプラットフォームの基本設計書。全体アーキテクチャ（Next.js + Go）、認証・認可（OAuth）、プラン・権限制御の考え方、画面構成、
    機能別基本設計（記事・ニュース・銘柄ページ・アンケート・管理画面）、データ取得方式、アナリティクス設計を記載。
* docs/er_diagram.md
  * 概要: データベース設計（ER図）。users, plans, articles, news, likes, stocks, stock_news, polls, poll_options, poll_votes の各テーブル定義と
    それらのリレーション（Mermaid形式のER図）を記載。
