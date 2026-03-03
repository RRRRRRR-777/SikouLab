# Seeds

開発環境用のシードデータ（初期データ）を管理するディレクトリ。

## ファイル命名規則

`数字_テーブル名.sql` の形式で命名し、数字は依存順序に従う。

- `01_plans.sql` — plansテーブル（他テーブルから参照されるため最初）
- `02_test_users.sql` — テスト用ユーザー

## 使用方法

```bash
# psqlで直接実行
psql -U postgres -d sicoulab_database_development -f seeds/01_plans.sql
psql -U postgres -d sicoulab_database_development -f seeds/02_test_users.sql

# Docker経由で実行
docker exec -i sicoulab-postgres-1 psql -U postgres -d sicoulab_database_development < seeds/01_plans.sql
docker exec -i sicoulab-postgres-1 psql -U postgres -d sicoulab_database_development < seeds/02_test_users.sql

# 全て実行（Makefileターゲット使用推奨）
make seed
```

## 注意事項

- **開発環境専用**: 本番環境では絶対に使用しないこと
- **重複回避**: `ON CONFLICT DO NOTHING` で再実行時のエラーを防止
- **外部キー順序**: 依存関係のあるテーブルは、参照先のテーブルを先に実行すること
