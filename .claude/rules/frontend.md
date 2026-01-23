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

### 参考情報
- [JSDoc Reference - TypeScript](https://www.typescriptlang.org/docs/handbook/jsdoc-supported-types.html)
- [JSDoc & TypeDoc Guide](https://dev.to/mirzaleka/learn-how-to-document-javascripttypescript-code-using-jsdoc-typedoc-359h)
