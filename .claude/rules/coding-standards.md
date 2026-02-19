# コーディング規約

## 基本ルール

- **日本語で書く**: コメント・ドキュメントは日本語
- **WhatよりWhy**: コードから自明な「何をするか」より、意図を説明

## コメント例

```dockerfile
# ベースイメージ
FROM node:20-alpine

# 依存関係をインストール
RUN npm ci

# ポートを公開
EXPOSE 3000
```

```typescript
/**
 * 記事一覧を取得する
 *
 * @param page - ページ番号
 * @returns 記事リスト
 */
```

```go
// ArticleList は記事一覧を取得する
func ArticleList(page, limit int) ([]Article, error) {
```

## テストコード

```typescript
test("記事一覧ページが表示される", async () => {
  // 準備: 記事データをモック
  const mockData = [{ id: 1, title: "テスト記事" }];
  mockApi.getArticles.mockResolvedValue(mockData);

  // 実行: ページをレンダリング
  const result = render(<ArticleList />);

  // 検証: 記事タイトルが表示されている
  expect(result.getByText("テスト記事")).toBeInTheDocument();
});

test("無効なページ番号ではエラー", async () => {
  // 実行: page=0でAPI呼び出し
  await expect(getArticles(0, 20)).rejects.toThrow("無効なページ番号");
});
```
