# ADR-006: デプロイフローの選定

## ステータス
承認済み

## コンテキスト
本番デプロイのトリガー方式、承認プロセス、ロールバック方針の選定が必要。
自動デプロイは避け、承認プロセスを設けたい。

## 決定
- **トリガー**: GitHub Environment承認
- **承認**: Environment Protection Rulesで承認者指定
- **ロールバック**: 前バージョンのコンテナイメージに戻す

## 検討した選択肢

| 方式 | 厳格さ | 特徴 |
|------|--------|------|
| ボタンクリック | ★☆☆ | シンプルだが緩い |
| PRコメント方式 | ★★☆ | `/deploy prod` でトリガー |
| 承認者レビュー + コメント | ★★★ | Approve後に `/deploy` |
| **GitHub Environment承認** | ★★★ | GitHub標準機能、承認UI |
| 専用デプロイPR方式 | ★★★★ | 最も厳格 |

## 理由
1. **GitHub標準機能**: 設定のみで実現、カスタム実装不要
2. **監査ログ**: GitHub Audit Logで承認履歴が残る
3. **複数承認者対応**: N人中M人承認を設定可能
4. **タイムアウト設定**: 待機時間を設定可能
5. **Secrets/Variables分離**: 環境ごとに設定を分離できる

## デプロイフロー

```
1. ticket → dev マージ（CI実行: lint, format, 単体テスト, docker build）
2. dev → main のPRを作成
3. PRマージ（CI実行: 上記 + E2E, VRT）
4. GitHub Actionsが `environment: production` に到達
5. 承認者に通知 → GitHub UI上で承認ボタンをクリック
6. デプロイ実行
```

## ロールバック手順

```
1. 問題発生を検知
2. Cloud Run / GKE で前バージョンのコンテナイメージにリビジョン切り替え
3. 原因調査 → 修正 → 通常のデプロイフローで再デプロイ
```

## 影響
- GitHub Repositoryの Settings > Environments で `production` 環境を作成
- Required reviewersで承認者を指定
- GitHub Actionsのワークフローに `environment: production` を追加
