---
paths:
  - "frontend/**/*"
---

# frontendディレクトリルール

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
* **配置場所**: `docs/versions/1_0_0/SikouLab.pen`
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
cd frontend && make dev

# 本番ビルド
cd frontend && make build

# 本番サーバー起動（ビルド後）
cd frontend && make start
```

### コード品質

```bash
# ESLint 実行
cd frontend && make lint

# Prettier フォーマット
cd frontend && make fmt

# 単体テスト
cd frontend && make test
```

### E2E / VRT

```bash
# E2Eテスト（frontend/backend 起動済み前提）
cd frontend && make test-e2e

# VRT（frontend/backend 起動済み前提）
cd frontend && make test-vrt
```

### Docker

```bash
# Dockerイメージビルド
cd frontend && make docker-build
```

### 参考情報
- [JSDoc Reference - TypeScript](https://www.typescriptlang.org/docs/handbook/jsdoc-supported-types.html)
- [JSDoc & TypeDoc Guide](https://dev.to/mirzaleka/learn-how-to-document-javascripttypescript-code-using-jsdoc-typedoc-359h)
