.PHONY: dev build start test lint fmt test-e2e test-vrt docker-build help

# デフォルトターゲット
help:
	@echo "Available targets:"
	@echo "  dev         - Start development server"
	@echo "  build       - Build for production"
	@echo "  start       - Start production server"
	@echo "  test        - Run unit tests"
	@echo "  lint        - Run ESLint"
	@echo "  fmt         - Run Prettier format"
	@echo "  test-e2e    - Run E2E tests"
	@echo "  test-vrt    - Run VRT"
	@echo "  docker-build - Build Docker image"

# 開発サーバー起動
dev:
	npm run dev

# 本番ビルド
build:
	npm run build

# 本番サーバー起動
start:
	npm run start

# 単体テスト実行
test:
	npm run test

# ESLint実行
lint:
	npm run lint

# Prettierフォーマット
fmt:
	npm run format

# E2Eテスト実行
test-e2e:
	npm run test:e2e

# VRT実行
test-vrt:
	npm run test:vrt

# Dockerイメージビルド
docker-build:
	docker build -t sikoulab-web .
