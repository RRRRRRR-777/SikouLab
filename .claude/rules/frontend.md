---
paths:
  - "frontend/**/*"
---

# frontendディレクトリルール

## モバイルファースト原則

### ブレークポイント定義（プロジェクト共通）

| 名称 | 範囲 | Tailwindプレフィックス |
|------|------|----------------------|
| モバイル | < 768px | （なし・ベース） |
| タブレット | 768px〜1023px | `md:` |
| デスクトップ | ≥ 1024px | `lg:` |

### 記述原則

- **ベースクラス = モバイル**として書き、`md:` / `lg:` で上書きする
- 情報を「非表示」にするのではなく「形を変える」を優先する

### 禁止パターン

| アンチパターン | 問題 | 代替 |
|---|---|---|
| `hidden md:flex` | スマホで消える | `flex flex-col md:flex-row`（形を変える） |
| `hidden lg:block` | スマホ・タブレットで消える | レイアウトを工夫して表示維持 |

### 例外（許容する非表示）

同じ情報を**別の形で提供している場合のみ** `hidden` を許容する。
その際は必ずコメントで「モバイルでは〇〇で代替」と明記すること。

```tsx
{/* デスクトップのサイドナビ。モバイルではドロワーナビ（Sidebar コンポーネント）で代替 */}
<nav className="hidden lg:flex ...">...</nav>
```

## モバイルUI実装ルール

### タッチターゲット
- ボタン・リンクは最小 **44×44px** を確保する（WCAG 2.1 / Apple HIG 準拠）
- Tailwindでは `min-h-[44px] min-w-[44px]` または `h-11 w-11`（44px）で担保

### フォントサイズ
- 本文は最小 `text-base`（16px）以上にする
- 16px未満だとiOSが自動ズームし、UXが壊れる

### フォーム最適化
- `type` 属性でモバイルキーボードを最適化する

| 入力内容 | type属性 |
|---|---|
| メールアドレス | `type="email"` |
| 電話番号 | `type="tel"` |
| 数値 | `type="number"` または `inputMode="numeric"` |
| 検索 | `type="search"` |

### 横スクロール禁止
- ルートレイアウトに `overflow-x: hidden` を設定する
- テーブル・コードブロックなど固定幅コンテンツは `overflow-x: auto` のラッパーで囲む

### 画像
- Next.js `<Image>` に `sizes` 属性を必ず指定する

```tsx
<Image
  src="..."
  alt="..."
  sizes="(max-width: 768px) 100vw, (max-width: 1024px) 50vw, 33vw"
/>
```

### Safe Area（ノッチ・ホームバー対応）
- iPhoneのノッチ・ホームバーに被る要素（フローティングボタン等）は `env(safe-area-inset-*)` で回避する
- `tailwind.config` に `safe-area` プラグインを追加するか、インラインで指定する

```tsx
{/* ホームバーに被らないよう safe-area-inset-bottom を加算 */}
<div className="fixed bottom-6 right-6 pb-[env(safe-area-inset-bottom)]">
```

### hover: スタイル
- タッチデバイスでは `hover:` は基本的に発火しないため、全デバイスで許容する
- ただし hover に**重要な情報**を載せてはならない（タッチで見えないため）

### z-index 管理（プロジェクト共通定義）

| 用途 | z-index | Tailwindクラス |
|---|---|---|
| ヘッダー（sticky） | 30 | `z-30` |
| サイドバー・ドロワー | 40 | `z-40` |
| フローティングボタン | 40 | `z-40` |
| モーダル・ダイアログ | 50 | `z-50` |
| トースト通知 | 60 | `z-[60]` |

- この定義外の値を使う場合はコメントで理由を明記する

## ダークモード原則

### 基本方針

- **OS設定に追従**する（`prefers-color-scheme`）。手動切り替えは将来拡張
- Tailwindの `dark:` プレフィックスを使用する（`darkMode: "class"` 設定済み前提）
- **ハードコードした色は使わない** — 必ずCSSカスタムプロパティまたはTailwindセマンティックトークンで定義する

### カラートークン定義（プロジェクト共通）

CLAUDE.mdのカラーコードに基づく定義。`globals.css` の `:root` / `.dark` で管理する。

| トークン名 | ライトモード | ダークモード | 用途 |
|-----------|------------|------------|------|
| `--color-bg` | `#FFFFFF` | `#000000` | メイン背景 |
| `--color-text` | `#000000` | `#FFFFFF` | サブタイトル・見出し |
| `--color-primary` | `#E86D00` | `#E86D00` | ボタン・ジャンル名・ラベル（共通） |
| `--color-ticker` | `#63B7E2` | `#63B7E2` | 銘柄コード（共通） |
| `--color-muted` | `#D5D5D5` | `#EBEBEB` | 日付・時間などサブテキスト |

```css
/* globals.css */
:root {
  --color-bg: #FFFFFF;
  --color-text: #000000;
  --color-primary: #E86D00;
  --color-ticker: #63B7E2;
  --color-muted: #D5D5D5;
}

.dark {
  --color-bg: #000000;
  --color-text: #FFFFFF;
  --color-muted: #EBEBEB;
  /* primary / ticker はライト・ダーク共通のため省略 */
}
```

### 記述原則

- ベースクラス = ライトモード として書き、`dark:` で上書きする
- プライマリカラー（`#E86D00`）・ティッカー色（`#63B7E2`）はモード間で変わらないため `dark:` 不要

### 禁止パターン

| アンチパターン | 問題 | 代替 |
|---|---|---|
| `bg-white dark:bg-black` を直書き | トークン管理外になる | `bg-[var(--color-bg)]` または Tailwindカスタムトークン |
| `text-gray-500` のみ（dark:なし） | ダークで読めない色になる | `text-[var(--color-muted)]` |
| インラインstyleに色を直書き | ダークモード切り替え不可 | CSSカスタムプロパティ経由 |

### 実装パターン

```tsx
{/* 背景・テキストはトークンで指定 */}
<div className="bg-[var(--color-bg)] text-[var(--color-text)]">

{/* プライマリカラーはモード不問で同色 */}
<button className="bg-[#E86D00] text-white">

{/* サブテキスト（日付等） */}
<span className="text-[var(--color-muted)]">2026-02-26</span>

{/* 境界線など微妙な色はTailwindのdark:で対応 */}
<div className="border border-gray-200 dark:border-gray-800">
```

### shadcn/ui との連携

- shadcn/uiのデフォルトCSSトークン（`--background`, `--foreground` 等）は上記プロジェクトトークンと**別管理**
- shadcn/uiコンポーネントを使う場合はshadcnのトークン体系に乗る。プロジェクトトークンを混在させない

## UI/UX原則
* **受動的体験を最優先** - ユーザーは「読む」ことが主目的
* **チャット前提UIは完全排除**
* **能動操作は最小限** - 投票・いいね程度に留める
* **見た目より「読める・迷わない」を優先**

## コンポーネント設計
* 「今日一番読まれている」を軸にレイアウト構成
* 過度な抽象化を避ける - 必要になってから抽象化
* 3回同じパターンが現れたら抽象化を検討

## Pencil定義
* **配置場所**: `docs/versions/1_0_0/SicouLab.pen`
* **作成タイミング**: 詳細設計（/feature-spec）フェーズで作成・更新
* **フロー**:
  1. feature-specでUI仕様を設計
  2. **feature-specを参照して**Pencilに該当画面/コンポーネントを追記
  3. 実装時は**Pencilを正**として実装
* **運用ルール**:
  * 新規画面/コンポーネント追加時は必ずPencilを先に更新
  * 実装とPencilに差異が出た場合は**Pencilを正**とする
  * バージョニングはGitで管理
  * **Pencil記載時にわからないことがあれば必ずユーザーに確認する**
  * **Pencil記載時に決定した内容は詳細設計書（feature-spec）に反映する**

## OpenAPI参照
* **配置場所**: `backend/api/openapi.yaml`
* **フロー**:
  1. API連携実装前にbackend/api/openapi.yamlを確認
  2. **openapi.yamlを正**として型定義・リクエストを実装
* **運用ルール**:
  * 型定義はopenapi.yamlから自動生成または手動で同期
  * リクエスト/レスポンスの形式はopenapi.yamlに従う
  * **OpenAPIに記載のないエンドポイントは使用しない**

## セキュリティ
* XSS対策を必須とする - ユーザー入力は常にサニタイズ
* 認証状態はバックエンドで必ず検証（フロントは表示制御のみ）
* APIキーやシークレットをフロントエンドに含めない

## パフォーマンス
* 初期表示速度を最優先
* 必要になるまで機能を読み込まない（Lazy Loading）
* 画像は適切なサイズ・フォーマットで配信

## 実装時の注意
* エラーハンドリングは境界（ユーザー入力、外部API）でのみ
* 内部のコード・フレームワークは信頼する
* 使われていないコードは完全削除（後方互換性ハック不要）

## Next.js 15 必須ルール

### Breaking Changes（必ず守る）
* `params`, `searchParams`, `cookies()`, `headers()` は **async** で取得
  ```typescript
  // ✅ 正しい
  export default async function Page({ params }) {
    const { id } = await params;
  }
  ```
* `useFormState` → `useActionState`（`react` からimport）
* `fetch()` はデフォルトでキャッシュされない（明示的に設定が必要）

### 基本原則
* Server Components をデフォルトとし、Client Components は必要時のみ
* データ取得はサーバーで行い、クライアントでの `useEffect` + `fetch` は避ける
* Server Actions でミューテーションを実装

## コード内ドキュメンテーション

### 基本方針
- **設計に関するドキュメントはコード内に完結**
- feature-specには実装仕様・テストケースのみ記載
- エクスポートされる関数・コンポーネント・型には必ずドキュメント

### インラインコメント
- **タイミング**: 実装中に記載
- **対象**: 関数・条件分岐・ループ・例外処理など、処理の単位で記載
- **ルール**: WhatよりもWhy（コードから自明なことは書かない）

### JSDoc必須ルール

#### 対象
- エクスポートされる全ての関数・コンポーネント
- エクスポートされる型・インターフェース
- 複雑なロジックを含む関数

#### 記述内容
```typescript
/**
 * 関数の概要（1行で「何をするか」）
 *
 * 補足説明が必要な場合に記述。why（なぜそうするか）を中心に。
 *
 * @param paramName - パラメータの説明
 * @returns 戻り値の説明
 * @throws {Error} 例外の条件
 *
 * @example
 * ```ts
 * const result = functionName(arg);
 * ```
 */
```

#### ルール
- **WhatよりもWhy** - コードから自明な「何をするか」より、意図や制約を説明
- 完全な文で書く（文末にピリオド）
- パラメータ・戻り値は必ず記述
- 使い方が自明でない場合は `@example` を追加

### 詳細リファレンス
コード例やパターンの詳細が必要な場合は `/nextjs-15` スキルを呼び出してください。

## Makefile コマンド

`frontend/Makefile` で定義されたターゲットを使用する。

### 開発

```bash
# 開発サーバー起動
make -C frontend dev

# 本番ビルド
make -C frontend build

# 本番サーバー起動（ビルド後）
make -C frontend start
```

### コード品質

```bash
# ESLint 実行
make -C frontend lint

# Prettier フォーマット
make -C frontend fmt

# 単体テスト
make -C frontend test
```

### E2E / VRT

```bash
# E2Eテスト（frontend/backend 起動済み前提）
make -C frontend test-e2e

# VRT（frontend/backend 起動済み前提）
make -C frontend test-vrt
```

### Docker

```bash
# Dockerイメージビルド
make -C frontend docker-build
```

### 参考情報
- [JSDoc Reference - TypeScript](https://www.typescriptlang.org/docs/handbook/jsdoc-supported-types.html)
- [JSDoc & TypeDoc Guide](https://dev.to/mirzaleka/learn-how-to-document-javascripttypescript-code-using-jsdoc-typedoc-359h)
