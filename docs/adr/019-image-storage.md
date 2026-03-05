# ADR-019: 画像ストレージの選定

## ステータス
承認済み

## コンテキスト

SicouLabでは複数の機能で画像を扱う必要がある。

| 機能 | 用途 | 状態 |
|------|------|------|
| アバター画像 | ユーザープロフィール | 実装済み（LocalStorage） |
| 記事本文画像 | Markdownエディタ内に埋め込む画像 | 設計済み・保存先TBD |
| 記事サムネイル | ダッシュボード・一覧のカード表示 | 設計のみ |
| ニュースサムネイル | ダッシュボード・一覧のカード表示 | 設計のみ |
| アンケートサムネイル | 一覧のカード表示 | 設計のみ |
| 投稿ユーザーアバター | 記事・ニュースの著者表示 | 設計のみ |

### 要件

- **本番環境**: GCPエコシステム統合（Cloud Run + Cloud SQL構成）
- **ローカル開発環境**: 本番との差異を最小化しつつ、セットアップを簡易に
- **拡張性**: 画像種別の追加に対応できる設計
- **既存設計**: インフラ設計書でCloud Storageバケット構成が定義済み（[infrastructure.md](../versions/1_0_0/infrastructure.md)）

### 既存実装の課題

アバター画像は `LocalStorage`（`os.WriteFile` でローカルディスクに保存）で実装済みだが、本番環境（Cloud Run）ではコンテナ再起動時にファイルが消失する。また、記事画像など大量データを扱う機能では、ローカルディスク保存はスケーラビリティに欠ける。

## 決定

### 本番環境: **Cloud Storage（GCP）** を採用する

### ローカル開発環境: **fake-gcs-server** を採用する

### 詳細構成

| 環境 | 技術 | 備考 |
|------|------|------|
| 本番 | Cloud Storage（GCP） | `sicoulab-assets-prod` バケット |
| ローカル開発 | [fake-gcs-server](https://github.com/fsouza/fake-gcs-server) | Docker Compose で起動 |
| Go SDK | `cloud.google.com/go/storage` | 本番・ローカル共通 |

### ディレクトリ構成（バケット内）

```
sicoulab-assets-prod/
  articles/          - 記事本文画像
  users/             - ユーザーアバター
  news/              - ニュースサムネイル
```

### 環境切替方式

環境変数 `STORAGE_EMULATOR_HOST` でエンドポイントを切り替える。Go SDK（`cloud.google.com/go/storage`）がこの環境変数を自動認識するため、アプリケーションコードの分岐は不要。

```bash
# ローカル開発（.env）
STORAGE_EMULATOR_HOST=localhost:4443

# 本番（環境変数なし → 自動的にGCP本番に接続）
```

### docker-compose への追加

```yaml
fake-gcs:
  image: fsouza/fake-gcs-server
  ports:
    - "4443:4443"
  command: ["-scheme", "http", "-port", "4443"]
  volumes:
    - fake_gcs_data:/storage
  networks:
    - sicoulab
```

## 検討した選択肢

### 比較表

| 選択肢 | 本番との差異 | セットアップ | SDK統一 | 管理UI |
|--------|-------------|-------------|---------|--------|
| **fake-gcs-server** | ◎ 最小 | ○ Docker 1コンテナ | ◎ 同一SDK | × なし |
| MinIO | △ S3互換（GCSと微妙に異なる） | ○ Docker 1コンテナ | × 別SDK必要 | ◎ WebUI付 |
| LocalStorage拡張 | × API全く別 | ◎ 追加なし | × 別実装 | × なし |

### 詳細評価

#### fake-gcs-server
- メリット: 本番と同じGo SDK（`cloud.google.com/go/storage`）をそのまま使える。`STORAGE_EMULATOR_HOST` で自動切替
- デメリット: 管理UIがない（ファイル確認はCLIまたはボリューム直接参照）

#### MinIO
- メリット: 管理UIでブラウザから画像を確認できる。S3互換で広く使われている
- デメリット: GCSとAPIが異なるため、SDKを2つ持つか変換層が必要。本番切替時にコード変更が発生する

#### LocalStorage拡張
- メリット: 追加依存なし。最もシンプル
- デメリット: 本番とAPI・SDKが完全に異なる。`AvatarStorage` インターフェース実装を環境ごとに作り分ける必要がある

## 理由

### fake-gcs-serverが採用された理由

1. **SDK統一**: 本番と同じ `cloud.google.com/go/storage` を使えるため、環境ごとのコード分岐が不要
2. **切替の容易さ**: `STORAGE_EMULATOR_HOST` 環境変数の有無だけで本番/ローカルが切り替わる（Go SDKの公式サポート）
3. **GCPエコシステム適合**: 本番インフラがCloud Run + Cloud SQL + Cloud StorageのGCP構成であり、ローカルもGCS互換にすることで一貫性が保たれる

## 影響

- **既存のLocalStorage実装**: `cloud.google.com/go/storage` を使ったGCS実装に置き換える。`AvatarStorage` インターフェースはそのまま活用可能
- **docker-compose.yml**: `fake-gcs` サービスを追加
- **環境変数**: `STORAGE_EMULATOR_HOST` を `.env` / `.env.sample` に追加
- **依存パッケージ**: `cloud.google.com/go/storage` を `go.mod` に追加
- **記事画像のTBD**: `docs/functions/article/create-edit.md` の画像保存先TBDが解消される

## 参考情報
- [fake-gcs-server](https://github.com/fsouza/fake-gcs-server)
- [Cloud Storage Client Libraries（Go）](https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go)
- [STORAGE_EMULATOR_HOST](https://cloud.google.com/storage/docs/emulator)
