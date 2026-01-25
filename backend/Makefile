.PHONY: dev run build test lint fmt migrate-up migrate-down docker-build help

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
	@echo "  docker-build - Build Docker image"

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
	migrate -path db/migrations -database "$$DATABASE_URL" up

# マイグレーション戻し
migrate-down:
	migrate -path db/migrations -database "$$DATABASE_URL" down

# Dockerイメージビルド
docker-build:
	docker build -t sikoulab-api .
