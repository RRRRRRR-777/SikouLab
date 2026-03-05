#!/usr/bin/env bash
set -euo pipefail

# E2Eテスト用DBリセットスクリプト
# マイグレーション再適用 + seed投入を行う

CONTAINER_NAME="sicoulab-postgres-1"
DB_NAME="sicoulab_database_development"
DB_USER="postgres"

echo "=== E2E DB リセット開始 ==="

# マイグレーション再適用（down → up）
echo "--- マイグレーション再適用 ---"
DATABASE_URL="postgres://${DB_USER}:${DB_USER}@localhost:5432/${DB_NAME}?sslmode=disable"
migrate -path backend/migrations -database "${DATABASE_URL}" down -all || true
migrate -path backend/migrations -database "${DATABASE_URL}" up

# シードデータ投入
echo "--- シードデータ投入 ---"
for file in backend/seeds/*.sql; do
  [ -f "$file" ] || continue
  echo "Applying ${file}..."
  docker exec -i "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" < "${file}"
done

echo "=== E2E DB リセット完了 ==="
