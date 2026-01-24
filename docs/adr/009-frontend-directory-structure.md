# ADR-009: フロントエンドディレクトリ構成

## ステータス

採用

## 背景

Next.js 15（App Router）のベストプラクティスに沿ったディレクトリ構成を定義する必要がある。

## 決定事項

### ディレクトリ構成

```
frontend/
├── app/                      # ルーティング専用
│   ├── (auth)/              # Route Group（認証系）
│   │   ├── login/
│   │   └── register/
│   ├── (main)/              # Route Group（メイン）
│   │   ├── articles/
│   │   ├── news/
│   │   └── settings/
│   ├── api/                 # API Routes（必要に応じて）
│   ├── layout.tsx
│   ├── page.tsx
│   └── globals.css
├── components/
│   ├── ui/                  # shadcn/ui（自動生成）
│   └── features/            # 機能別コンポーネント
│       ├── article/
│       └── user/
├── hooks/                   # カスタムフック
├── lib/                     # ライブラリ初期化
│   ├── api-client.ts
│   └── db.ts
├── utils/                   # 純粋関数
│   ├── format.ts
│   └── validation.ts
├── stores/                  # グローバルステート管理
├── types/                   # 型定義
│   └── api/
├── public/                  # 静的ファイル
└── .env.local
```

### 各ディレクトリの役割

| ディレクトリ | 役割 | 例 |
|--------------|------|-----|
| **app/** | ルーティング専用。ビジネスロジックを置かない | ページ、レイアウト、API Routes |
| **components/ui/** | shadcn/uiの自動生成先 | Button, Input, Card |
| **components/features/** | 機能別コンポーネント | ArticleCard, UserProfile |
| **hooks/** | カスタムフック | useAuth, useInfiniteScroll |
| **lib/** | ライブラリ初期化・設定 | api-client.ts |
| **utils/** | 純粋関数・ユーティリティ | formatDate, validateEmail |
| **stores/** | グローバルステート管理（Zustand等） | useAuthStore |
| **types/** | TypeScript型定義 | Article, User |

### 設計方針

1. **app/はルーティング専用**
   - Server Componentをデフォルト
   - ビジネスロジックはcomponentsやlibに配置
   - Route Groupsでレイアウトを共有

2. **Featureベースのコンポーネント**
   - `components/features/`で機能ごとに整理
   - 単一責任の原則に従う

3. **libとutilsの分離**
   - lib: ライブラリ初期化（副作用あり）
   - utils: 純粋関数（副作用なし）

4. **グローバルステート**
   - TanStack Query: サーバー状態
   - stores: クライアント状態（UI、フォーム入力）

## 採用理由

- **2025年のベストプラクティスに準拠**: app/はルーティング専用とする構成が推奨されている
- **保守性**: 機能ベースの分割によりコードの所在地が明確
- **スケーラビリティ**: 大規模化しても構造が崩れにくい

## 参考資料

- [Next.js App Router構成の実案件設計ガイド【2025年版】](https://qiita.com/mukai3/items/62e07582294630345902)
- [【2025版】Next.js 最適フォルダ(ディレクトリ)構成・設計](https://zenn.dev/yamu_official/articles/70f59488e8415d)
- [【保存版】Next.js(App Router)のベストプラクティス](https://qiita.com/k_morimori/items/32bc97e524a26a183b30)
- [【Next.js】App RouterとRSCによるアーキテクチャ設計](https://zenn.dev/kiwichan101kg/articles/b44305e3049bac)

## 関連

- [ADR-001: UIライブラリ](./001-ui-library.md)
- [ADR-008: API連携](./008-api-integration.md)
- [開発ガイドライン](../development_guidelines.md)
