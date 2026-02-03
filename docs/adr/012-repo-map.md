# ADR: Repo Map（AI向けコード索引）の導入

## ステータス
採用

## コンテキスト
AI（LLM）がコードを探索・修正・実装する際、ファイル数が多くなると「どこに何があるか」の探索コストが高くなる。Aiderの Repo Map 思想に基づき、AIが効率的にコードベースを理解できる自動生成索引を導入する。

---

## 決定事項
**Go統一ツールでフロントエンド・バックエンド両方のRepo Mapを自動生成する**

| 項目 | 内容 |
|------|------|
| **ファイル配置** | `backend/REPO_MAP.md`, `frontend/REPO_MAP.md`（分離） |
| **生成ツール** | `tools/repomap/main.go`（Go統一） |
| **CI連携** | GitHub Actionsで自動生成・自動コミット |

---

## Repo Mapの定義（Aider準拠）

Repo Map = **実装内容を説明しない、構造だけの地図**

### 含めるもの
- Entry points（`package main`, `app/`ルート）
- Packages / ディレクトリ
- Exported symbols（公開API）

### 含めないもの
- ロジック説明
- 詳細コメント
- フロー解説（別ドキュメントで管理）

---

## 選択肢の比較

### Go統一

**仕組み**: Goで両言語を解析（Go: `go/ast`, TS: 正規表現）

| メリット | デメリット |
|----------|-----------|
| 単一バイナリ、CI設定がシンプル | TypeScript解析は `ts-morph` ほど精密でない |
| 依存なし（Node.js不要） | JSXの解析が複雑になる可能性 |
| 処理速度が最速 | - |

---

### 各言語で実装

**仕組み**: Go版はGo、Next.js版はTypeScriptで実装

| メリット | デメリット |
|----------|-----------|
| 最も正確（ネイティブAST） | 2つの生成ツールを保守 |
| TypeScriptは `ts-morph` で型情報も取得可 | CIで両方のランタイム必要 |

---

### 軽量LLM生成

**仕組み**: `haiku` や `gpt-4o-mini` にソースを渡して索引生成

| メリット | デメリット |
|----------|-----------|
| ツール実装不要 | API呼び出しコスト |
| 形式変更がプロンプト修正だけで可能 | 出力の揺れ（再現性低い） |
| 命名規則の意味理解も可能 | ファイル数増加で遅くなる |

---

## 決定の理由

1. **「作りすぎない」原則**: 初期は最小限（Go統一）で開始し、将来的な拡張を壊さない設計
2. **コストパフォーマンス**: LLM APIコスト不要、CI設定がシンプル
3. **処理速度**: AST解析は最速
4. **段階的移行**: TypeScript解析が複雑になったら `ts-morph` に分離可能

---

## 出力フォーマット

### backend/REPO_MAP.md
```markdown
# Backend Repo Map
> Auto-generated. Do not edit manually.

## Entry Points
- cmd/api/main.go

## Packages
### internal/handler
- type ArticleHandler
- func NewArticleHandler
- func (h *ArticleHandler) List
```

### frontend/REPO_MAP.md
```markdown
# Frontend Repo Map
> Auto-generated. Do not edit manually.

## Routes (app/)
- app/page.tsx → HomePage
- app/articles/page.tsx → ArticlesPage

## Components
### components/ui/
- Button
- Card

## Hooks
- useAuth

## Utilities (lib/)
- cn
```

---

## CI設計

```yaml
name: Update Repo Map
on:
  push:
    paths:
      - 'backend/**/*.go'
      - 'frontend/**/*.ts'
      - 'frontend/**/*.tsx'
jobs:
  update:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go run ./tools/repomap --target=backend
      - run: go run ./tools/repomap --target=frontend
      - uses: stefanzweifel/git-auto-commit-action@v5
        with:
          commit_message: "chore: update REPO_MAP.md"
          file_pattern: "*/REPO_MAP.md"
```

---

## 将来の拡張（想定）

- TypeScript解析の精度向上が必要になったら `ts-morph` に分離
- 大規模化したら部分 Repo Map 分割
- Flow ドキュメント（Repo Map とは分離）の追加

---

## 影響

- `CLAUDE.md`: Repo Map参照ルールを追記
- `tools/repomap/`: 生成ツールを新規作成
- `.github/workflows/`: CI設定を追加
- `backend/REPO_MAP.md`, `frontend/REPO_MAP.md`: 自動生成ファイルを追加

---

## 参考資料

- [AiderのRepository mapを理解する](https://zenn.dev/mopemope/articles/4fada49bd26eea)
