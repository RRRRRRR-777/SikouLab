# .env が存在する場合は自動で読み込み、環境変数としてエクスポートする
-include .env
export

.PHONY: dev run build test lint fmt migrate-up migrate-down seed docker-build check-doc swagger validate-api tunnel help

# デフォルトターゲット
help:
	@echo "Available targets:"
	@echo "  dev          - Start development server with hot reload"
	@echo "  run          - Start development server"
	@echo "  build        - Build binary"
	@echo "  test         - Run unit tests"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code with gofmt"
	@echo "  migrate-up   - Apply migrations"
	@echo "  migrate-down - Rollback migrations"
	@echo "  seed         - Seed database with test data"
	@echo "  docker-build - Build Docker image"
	@echo "  check-doc    - Check documentation comments"
	@echo "  swagger      - Start Swagger UI (http://localhost:8081)"
	@echo "  validate-api - Validate OpenAPI specification"
	@echo "  tunnel      - Start Cloudflare tunnel (HTTPS for localhost:8080)"

# 開発サーバー起動（ホットリロード）
dev:
	air

# 開発サーバー起動
run:
	go run ./cmd/api

# バイナリビルド
build:
	go build -o bin/api ./cmd/api

# 単体テスト実行
test:
	go test ./...

# 静的解析
lint:
	golangci-lint run

# コードフォーマット
fmt:
	gofmt -w .

# マイグレーション適用
migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

# マイグレーション戻し
migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down 1

# シードデータ適用
seed:
	@for file in seeds/*.sql; do \
		echo "Applying $$file..."; \
		docker exec -i sicoulab-postgres-1 psql -U postgres -d sicoulab_database_development < "$$file"; \
	done

# Dockerイメージビルド
docker-build:
	docker build -t sicoulab-api .

# ドキュメントコメントチェック
check-doc:
	go run check_doc.go

# Swagger UI 起動
swagger:
	docker compose -f ../docker-compose.yml up -d swagger-ui

# OpenAPI仕様バリデーション
validate-api:
	npx @apidevtools/swagger-cli validate api/openapi.yaml

# Cloudflare Tunnel 起動（ローカルHTTPS）
tunnel:
	docker compose -f ../docker-compose.yml --profile tunnel up -d --no-deps cloudflared
