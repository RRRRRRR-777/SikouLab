# v1.0.0 開発計画

## 概要

設計完了済み（要件定義・基本設計・詳細設計18機能分）の状態から、実装フェーズに移行する。
バックエンド → フロントエンドの順で各機能を開発する。

### 開発順序の原則

```
各機能ごとに:
  1. TBD項目の決定（該当機能に未決定事項がある場合）
  2. OpenAPI追記（backend/api/openapi.yaml）
  3. DBマイグレーション作成（backend/migrations/）
  4. バックエンド実装（TDD）
  5. フロントエンド実装（TDD）
  6. レビュー → コミット
```

### 凡例

- **MUST**: 実装前に必ず決定が必要
- **SHOULD**: 実装中に決定可能だが、事前決定が望ましい
- **RISK**: 工数変動の可能性があるリスク要因

---

## Phase 0: 共通基盤

全機能の前提となるプロジェクト骨格を構築する。

### 横断的な決定事項

全機能に共通するTBD項目。Phase 0の作業開始前に一括で決定する。

| # | 項目 | 選択肢 | 備考 | 決定 |
|---|------|--------|------|------|
| D-01 | ログ出力方針 | 構造化ログ（JSON）+ slog / zerolog / zap | インフラ設計書 §6.1に準拠 | **zerolog**（JSON特化・最速・学習コスト低） |
| D-02 | ログレベル方針 | ERROR/WARN/INFO/DEBUG の使い分け基準 | 全機能共通で適用 | **標準4レベル**（ERROR=即時対応/WARN=注意/INFO=主要イベント/DEBUG=開発時） |
| D-03 | スケルトン表示 | あり / なし | 読み込み中のUI方針（全画面共通） | **あり**（shadcn/ui Skeleton使用） |
| D-04 | エラー表示 | トースト / インライン / 両方 | エラーメッセージのUI方針 | **トースト**（操作フローを邪魔しない） |
| D-05 | Pencil未定義画面の対応 | 実装時にPencilでデザイン作成 / テキストベースで進める | 対象: S-04, S-05, S-07, S-08, S-09, S-11, S-12, S-15, S-16, S-18〜S-25 | **テキストベース**（設計書ベースで実装、後からデザイン調整） |

### 0-1. バックエンド プロジェクト構造

| 作業 | 内容 |
|------|------|
| ディレクトリ構成 | Clean Architecture（cmd/api, internal/handler, internal/usecase, internal/repository, internal/domain） |
| ルーター | 標準ライブラリ `net/http` または軽量ルーター |
| ミドルウェア | CORS、リクエストログ、エラーハンドリング、リカバリ |
| 設定管理 | 環境変数読み込み（.env / godotenv） |
| ヘルスチェック | GET /health（既存を移行） |

### 0-2. DB接続 + マイグレーション基盤

| 作業 | 内容 |
|------|------|
| docker-compose | PostgreSQL 16コンテナ定義（既存docker-compose.ymlに追記） |
| DB接続 | sqlxによる接続プール設定 |
| マイグレーション | golang-migrate CLI + Makefile統合 |
| 初期マイグレーション | plansテーブル（マスタデータ） |

### 0-3. OpenAPI初期ファイル

| 作業 | 内容 |
|------|------|
| 雛形作成 | `backend/api/openapi.yaml`（info, servers, components/schemas共通定義） |
| 共通スキーマ | ErrorResponse, PaginationParams, CursorPaginationResponse |
| ヘルスチェック | GET /health の定義 |

### 0-4. フロントエンド プロジェクト構造

| 作業 | 内容 |
|------|------|
| shadcn/ui | 初期化 + 基本コンポーネント導入 |
| 共通レイアウト | サイドバー（要件定義 §8.1）、ヘッダー（検索バー） |
| テーマ | ダークモード対応（CSS変数、requirements §8.2） |
| API連携基盤 | axios インスタンス + TanStack Query Provider |
| 認証基盤 | AuthProvider（後続F-01で実装） |

---

## Phase 1: F-01 ログイン（認証基盤）

> 詳細設計: `docs/functions/auth/login.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| セッション有効期限 | 7日 / 24時間 / リメンバーme機能 | **MUST** |
| CSRF対策方式 | Double Submit Cookie / Synchronizer Token / SameSite属性のみ | **MUST** |
| OAuthステート検証 | Firebase SDK任せ / 独自検証 | **MUST** |
| OAuthプロバイダごとの取得項目 | Google(name,email,avatar) / Apple(name,email) / X(name,avatar) | **SHOULD** |

### リスク要因

- Appleプロバイダの設定手順が多い（Apple Developer Program設定）
- Stripe連携：テスト環境での挙動確認が必要

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 1-1 | OpenAPI追記 | POST /auth/login, GET /auth/me, POST /auth/logout, POST /auth/refresh |
| 1-2 | DBマイグレーション | users, plans, user_settings テーブル |
| 1-3 | BE: Firebase Admin SDK統合 | IDトークン検証、ユーザー初回作成 |
| 1-4 | BE: 認証ミドルウェア | JWT検証、ロール判定、サブスク状態チェック |
| 1-5 | BE: Stripe連携（初回） | カスタマー作成、Checkout Session |
| 1-6 | FE: ログイン画面（S-01） | Firebase JS SDK、OAuth3プロバイダ |
| 1-7 | FE: 認証状態管理 | AuthProvider、ルートガード |
| 1-8 | FE: サブスク登録フロー | 初回ログイン後のStripe Checkout遷移 |

---

## Phase 2: F-10-2 サブスクリプション管理（課金基盤）

> 詳細設計: `docs/functions/settings/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| メールアドレス暗号化方式 | AES-256-GCM / bcrypt / 平文保存（初期） | **SHOULD** |

### リスク要因

- Stripeカスタマーポータルの初回連携手順（+1-2sp）
- メールアドレス暗号化保存の実装（+3sp）

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 2-1 | OpenAPI追記 | POST /stripe/checkout, POST /stripe/webhook, GET /settings/subscription |
| 2-2 | DBマイグレーション | newsletter_subscriptions テーブル |
| 2-3 | BE: Stripe Webhook処理 | subscription.created/updated/deleted イベント |
| 2-4 | BE: サブスク状態管理 | subscription_status更新、アクセス制御 |
| 2-5 | FE: 設定画面（S-17） | プロフィール、サブスク管理、メール登録、FAQ |
| 2-6 | FE: メール登録ウィジェット | 各画面に配置可能な共通コンポーネント |

---

## Phase 3: F-04 記事機能（核心価値）

> 詳細設計: `docs/functions/article/home.md`, `create-edit.md`, `detail.md`, `schedule.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| Markdownエディタライブラリ | @uiw/react-md-editor / Milkdown / Novel / Tiptap | **MUST** |
| 画像保存先 | Google Cloud Storage / Cloudinary | **MUST** |
| 自動保存間隔 | 30秒 / 60秒 / 変更検知時のみ | **SHOULD** |
| 目次の表示位置 | 左サイド固定 / 記事上部 | **SHOULD** |
| 関連銘柄のリアルタイム価格表示 | 表示 / 非表示（初期） | **SHOULD** |
| バッチ処理方式（予約投稿） | Cloud Scheduler + Cloud Run Jobs / cron | **MUST** |

### リスク要因

- Markdownエディタの選定により実装難易度が変動（5-8sp幅）
- 画像保存先により実装アプローチが異なる
- 指数減衰スコアのクエリパフォーマンス（記事数増加時）

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 3-1 | OpenAPI追記 | 記事CRUD、ジャンル一覧、人気記事、いいね、閲覧記録 |
| 3-2 | DBマイグレーション | articles, genres, article_genres, article_views, article_likes, article_summaries, posting_users |
| 3-3 | BE: 記事CRUD API | 作成・取得・更新・削除、下書き/公開/予約 |
| 3-4 | BE: 人気記事スコア | 指数減衰スコア算出ロジック（Hacker News方式） |
| 3-5 | BE: 予約投稿バッチ | 毎分実行、scheduled_at <= NOW の記事を公開 |
| 3-6 | BE: 画像アップロードAPI | Cloud Storage連携 |
| 3-7 | FE: 記事一覧（S-03） | ジャンルタブ、人気記事、無限スクロール |
| 3-8 | FE: 記事詳細（S-04） | Markdown表示、目次、いいね、ブックマーク |
| 3-9 | FE: 記事作成・編集（S-05） | Markdownエディタ、投稿ユーザー選択、予約設定 |

---

## Phase 4: F-05 ニュース機能

> 詳細設計: `docs/functions/news/home.md`, `fetch.md`, `genre-detail.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| ニュースAPI | NewsAPI.org / Alpha Vantage / Finnhub | **MUST** |
| 翻訳API | Google Cloud Translation / DeepL / OpenAI GPT-4 | **MUST** |
| ニュース取得頻度 | 1分 / 3分 / 10分 | **MUST** |
| ジャンル・銘柄紐付けロジック | キーワードマッチング(3sp) / LLM分類(8sp) | **MUST** |

### リスク要因

- 外部API選定により実装工数が大幅に変動
- 紐付けロジック: キーワードマッチング vs LLMで5sp差
- APIレート制限への対応

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 4-1 | OpenAPI追記 | ニュース一覧、トレンド、ジャンル別、ジャンル詳細、手動作成 |
| 4-2 | DBマイグレーション | news, news_genres, news_tickers, news_views |
| 4-3 | BE: ニュース取得バッチ | 外部API取得 → 翻訳 → DB保存 |
| 4-4 | BE: ニュースCRUD API | 一覧・トレンド・ジャンル別・フィルター |
| 4-5 | BE: ニュース手動作成API | admin/writer向け、予約投稿対応 |
| 4-6 | BE: ピン固定管理API | is_pinned, pin_order の管理 |
| 4-7 | FE: ニュース一覧（S-06） | トレンドニュース、ジャンル別表示 |
| 4-8 | FE: ジャンル詳細（S-07） | フィルター、無限スクロール |
| 4-9 | FE: ニュース作成・編集（S-24） | admin/writer向け管理画面 |

---

## Phase 5: F-06 アンケート + F-08 ブックマーク

> 詳細設計: `docs/functions/poll/home.md`, `docs/functions/bookmark/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| 投票率の端数処理 | 四捨五入（合計≠100%許容） / 最大剰余法 | **SHOULD** |

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 5-1 | OpenAPI追記 | アンケートCRUD、投票、ブックマークCRUD |
| 5-2 | DBマイグレーション | polls, poll_options, poll_votes, bookmarks |
| 5-3 | BE: アンケートAPI | 作成・投票・結果取得（UNIQUE制約で重複防止） |
| 5-4 | BE: ブックマークAPI | 追加・削除・一覧（記事/ニュース/アンケート） |
| 5-5 | FE: アンケート一覧（S-14） | カテゴリ絞り込み |
| 5-6 | FE: アンケート詳細（S-15） | 投票UI、結果プログレスバー |
| 5-7 | FE: ブックマーク一覧（S-16） | タブ形式、解除操作 |

---

## Phase 6: F-03 個別銘柄ページ

> 詳細設計: `docs/functions/stock/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| 外部株価API | Alpha Vantage / Polygon.io / Finnhub / Yahoo Finance | **MUST** |
| リアルタイム価格表示 | リアルタイム / 15分遅延 / ページロード時のみ | **MUST** |
| 移動平均線の表示 | あり / なし（初期） | **SHOULD** |
| チャートライブラリ | lightweight-charts（ADR確定済み） | 確定 |

### リスク要因

- 外部API依存が最も多い機能（株価・ファンダ・インサイダー・UOA・レーティング・Massive API）
- 並列データ取得のパフォーマンス最適化

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 6-1 | OpenAPI追記 | 銘柄詳細、チャート、ファンダ、インサイダー、UOA、レーティング |
| 6-2 | DBマイグレーション | stocks, stock_prices, fundamentals, insider_trades, uoa, ratings, analyst_ratings |
| 6-3 | BE: 株価データ取得バッチ | 外部API → DB保存 |
| 6-4 | BE: 銘柄詳細API | チャート・ファンダ・インサイダー・UOA・レーティング統合 |
| 6-5 | FE: 銘柄詳細（S-10） | lightweight-charts、各セクション表示 |

---

## Phase 7: F-07 検索機能

> 詳細設計: `docs/functions/search/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| PGroonga + MeCabのセットアップ | docker-compose内 / Cloud SQL拡張 | **MUST** |
| 入力補完（サジェスト） | あり / なし（初期） | **SHOULD** |

### リスク要因

- PGroonga導入のトラブルシューティング
- 4テーブル横断スコアリングの複雑さ

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 7-1 | PGroonga + MeCabセットアップ | docker-compose、インデックス作成 |
| 7-2 | OpenAPI追記 | GET /search（クエリ、タブ、ページネーション） |
| 7-3 | DBマイグレーション | PGroongaインデックス追加 |
| 7-4 | BE: 検索API | 4テーブル横断検索、スコアリング |
| 7-5 | FE: 検索（S-09） | グローバル検索バー、タブ形式結果表示 |

---

## Phase 8: F-09 ポートフォリオ + F-02 ダッシュボード

> 詳細設計: `docs/functions/portfolio/home.md`, `docs/functions/dashboard/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| リアルタイム株価更新戦略 | リアルタイム(30秒) / 15分キャッシュ / 1時間キャッシュ | **MUST** |
| 相場インデックス更新頻度 | リアルタイム(30秒〜2分) / ページロード時のみ(5分TTL) | **MUST** |

### リスク要因

- 異常検出バッチのパフォーマンス（全市場銘柄スキャン）
- 外部株価APIのレート制限

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 8-1 | OpenAPI追記 | ポートフォリオCRUD、ダッシュボード各ウィジェット |
| 8-2 | DBマイグレーション | portfolios, portfolio_stocks, market_anomalies |
| 8-3 | BE: ポートフォリオAPI | CRUD、パフォーマンス計算、関連ニュース |
| 8-4 | BE: 異常検出バッチ | price_surge/volume_surge/volatility_surge |
| 8-5 | BE: ダッシュボードAPI | 各ウィジェットデータ取得 |
| 8-6 | FE: ポートフォリオ詳細（S-12） | 銘柄管理、パフォーマンス表示 |
| 8-7 | FE: ダッシュボード（S-02） | 5ウィジェット統合 |

---

## Phase 9: F-11 ニュースレター + F-10 設定画面（残り）

> 詳細設計: `docs/functions/newsletter/home.md`, `manage.md`, `settings/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| LLM API（記事要約） | OpenAI GPT-4 / Claude / Google Gemini | **MUST** |
| メール配信API | SendGrid / Amazon SES / Resend | **MUST** |

### リスク要因

- LLMプロンプト設計：要約品質の安定化に調整工数（5→8sp幅）
- メール配信APIのセットアップ

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 9-1 | OpenAPI追記 | ニュースレター配信、管理API |
| 9-2 | DBマイグレーション | newsletter_articles, newsletter_logs |
| 9-3 | BE: 記事要約生成 | LLM API連携、要約キャッシュ |
| 9-4 | BE: メール配信バッチ | 毎日07:30 JST実行 |
| 9-5 | BE: ニュースレター管理API | 配信リスト管理、配信履歴 |
| 9-6 | FE: ニュースレター管理（S-23） | 配信リスト、ドラッグ&ドロップ |

---

## Phase 10: F-12 管理者ページ

> 詳細設計: `docs/functions/admin/home.md`

### 実装前の決定事項

| 項目 | 選択肢 | 優先度 |
|------|--------|--------|
| アカウント停止方式 | subscription_statusにsuspended追加 / is_suspendedカラム追加 | **MUST** |
| お知らせ通知配信方式 | アプリ内通知 / メール通知 / 両方 | **MUST** |
| Metabase連携方式（WANT） | Embedding SDK + JWT署名 / 外部リンク | **SHOULD** |

### リスク要因

- Pencilデザイン未定義（S-18〜S-25）：実装時にUI設計が必要
- 依存機能（F-01/F-04/F-05/F-06）の完了が前提
- admin/writerの権限分岐が複雑

### 作業項目

| # | 作業 | 対象 |
|---|------|------|
| 10-1 | OpenAPI追記 | ユーザー管理、設定管理、お知らせ、投稿ユーザー管理 |
| 10-2 | DBマイグレーション | notifications, system_settings |
| 10-3 | BE: ユーザー管理API | 一覧・権限付与・アカウント停止 |
| 10-4 | BE: 設定管理API | ジャンル設定、表示件数、トレンド係数 |
| 10-5 | BE: お知らせ通知API | 通知作成・配信 |
| 10-6 | BE: 投稿ユーザー管理API | CRUD |
| 10-7 | FE: 管理者ダッシュボード（S-18） | ナビゲーション拠点 |
| 10-8 | FE: 各管理画面 | S-19〜S-25 |
