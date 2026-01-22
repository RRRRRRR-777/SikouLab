# SikouLab プロジェクト指針

## プロダクト原則
- **受動的体験を最優先** - 読む体験が価値の中心
- **作りすぎない** - 初期は最小限、将来拡張を壊さない設計

## スコープ（ver1.0.0）
**対象**: 記事、ニュース、アンケート、米国株銘柄、サブスク、いいね、アナリティクス
**対象外**: チャット、コメント、メール認証、Discord連携

## 禁止コマンド
以下のコマンドは直接実行せず、ユーザーにコマンドを提示して実行を委ねる。

```bash
# Git破壊系
git push --force
git push -f
git reset --hard
git clean -fd

# ファイル削除
rm -rf

# 本番環境操作
kubectl delete
docker system prune -f
```

## タスク管理
- Jira MCPツールを使用
- 新規タスクは適切なプロジェクト・課題タイプを指定

## 参照情報
- **Figma**: https://www.figma.com/design/KLDR5porwKHiMz5EWp3ABs
- **File Key**: KLDR5porwKHiMz5EWp3ABs

## 詳細ルール
@.claude/rules/workflow.md - 開発フロー（設計→実装→テスト）
@.claude/rules/docs.md - ドキュメント規約
@.claude/rules/frontend.md - フロントエンド規約
@.claude/rules/backend.md - バックエンド規約
