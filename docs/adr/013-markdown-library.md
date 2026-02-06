# ADR-013: Markdownライブラリの選定

## ステータス
承認済み

## コンテキスト
記事機能（F-04）において、以下2つのMarkdown処理が必要:
- **記事詳細ページ（F-04-5）**: MarkdownをHTMLにレンダリング、目次自動生成、見出しスクロール連携
- **記事作成・編集（F-04-1）**: Markdownエディタ、ツールバー（見出し・水平線等のボタン）

実装コストを最小化するライブラリ選定が必要。

## 決定
**react-markdown（表示）+ @uiw/react-md-editor（エディタ）** を採用する。

### 詳細構成

| 用途 | ライブラリ | プラグイン |
|------|------------|------------|
| **記事詳細表示** | react-markdown | rehype-slug, rehype-autolink-headings, @jsdevtools/rehype-toc |
| **記事作成・編集** | @uiw/react-md-editor | 標準ツールバー |

```bash
# 表示用
npm install react-markdown rehype-slug rehype-autolink-headings @jsdevtools/rehype-toc

# エディタ用
npm install @uiw/react-md-editor
```

## 検討した選択肢

### 表示用ライブラリ

| ライブラリ | 目次対応 | 実装コスト | 導入容易さ |
|------------|----------|------------|------------|
| **react-markdown** | ◎ プラグイン豊富 | ◎ 低 | ◎ |
| marked-react | △ 自前実装必要 | ○ 中 | ○ |
| @next/mdx | ○ プラグイン対応 | △ 高 | △ 設定複雑 |
| markdown-to-jsx | △ 自前実装必要 | ○ 中 | ○ |

### エディタ用ライブラリ

| ライブラリ | ツールバー | 実装コスト | カスタマイズ |
|------------|------------|------------|-------------|
| **@uiw/react-md-editor** | ◎ 標準装備 | ◎ 低 | ○ |
| TipTap | ○ カスタム可能 | △ 中 | ◎ 高い |
| Lexical | ○ カスタム可能 | △ 高 | ◎ 高い |
| BlockNote | ○ 標準装備 | ○ 中 | ○ |
| CKEditor 5 | ○ 標準装備 | ○ 中 | ○ |

## 理由

### react-markdown（表示用）
1. **実装コスト最小**: インストールして即使用可能
2. **目次生成**: `rehype-toc`で自動生成、スクロール連携に必要なIDを`rehype-slug`で付与
3. **セキュリティ**: デフォルトでXSS対策済み
4. **エコシステム**: remark/rehypeプラグインが豊富

### @uiw/react-md-editor（エディタ用）
1. **ツールバー標準装備**: 見出し、水平線、太字、斜体等のボタンが最初から利用可能
2. **導入容易**: コンポーネント1つで使用可能
3. **ライブプレビュー**: リアルタイムプレビュー機能標準装備
4. **軽量**: 必要最低限の機能に絞られている

## 実装例

### 記事詳細表示（F-04-5）

```tsx
import ReactMarkdown from 'react-markdown'
import rehypeSlug from 'rehype-slug'
import rehypeAutolinkHeadings from 'rehype-autolink-headings'

<ReactMarkdown
  rehypePlugins={[rehypeSlug, rehype-autolink-headings]}
>{markdown}</ReactMarkdown>
```

### 記事作成・編集（F-04-1）

```tsx
import MDEditor from '@uiw/react-md-editor'

<MDEditor
  value={markdown}
  onChange={(value) => setMarkdown(value)}
  toolbarButtons={['bold', 'italic', 'heading', 'hr', 'list']}
/>
```

## 影響
- **F-04-1（記事作成・編集）**: @uiw/react-md-editorを採用
- **F-04-5（記事詳細）**: react-markdown + rehypeプラグインを採用
- **依存パッケージ**: 5つのnpmパッケージ追加

## 参考情報
- [react-markdown](https://github.com/remarkjs/react-markdown)
- [@uiw/react-md-editor](https://github.com/uiwjs/react-md-editor)
- [rehype-toc](https://github.com/Microflash/rehype-toc)
- [Next.js MDX Guide](https://nextjs.org/docs/app/guides/mdx)
