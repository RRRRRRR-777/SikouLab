.PHONY: dev build start test lint fmt test-e2e test-vrt docker-build docs e2e-up e2e-down e2e-reset-db e2e-auth help

ROOT_DIR := $(shell dirname $(CURDIR))

# デフォルトターゲット
help:
	@echo "Available targets:"
	@echo "  dev          - Start development server"
	@echo "  build        - Build for production"
	@echo "  start        - Start production server"
	@echo "  test         - Run unit tests"
	@echo "  lint         - Run ESLint"
	@echo "  fmt          - Run Prettier format"
	@echo "  test-e2e     - Run E2E tests"
	@echo "  test-vrt     - Run VRT"
	@echo "  e2e-up       - Start E2E environment (services + DB)"
	@echo "  e2e-down     - Stop E2E environment"
	@echo "  e2e-reset-db - Reset DB (migrations + seeds)"
	@echo "  e2e-auth     - Auth setup (Firebase Emulator + auto login)"
	@echo "  docs         - Generate TypeDoc documentation"
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

# E2Eテスト実行（Emulator起動 → Next.js自動起動(:3100) → テスト → 停止）
# バックエンドは事前に起動済みであること
test-e2e:
	@echo "Starting Firebase Auth Emulator..."
	@firebase emulators:start --only auth --project demo-sicoulab &
	@sleep 3
	@npx playwright test; \
		EXIT_CODE=$$?; \
		echo "Stopping Firebase Auth Emulator..."; \
		kill %1 2>/dev/null || true; \
		exit $$EXIT_CODE

# VRT実行
test-vrt:
	npm run test:vrt

# TypeDocドキュメント生成
docs:
	npm run docs

# E2E環境起動（サービス起動 + DB準備）
e2e-up:
	docker compose -f $(ROOT_DIR)/docker-compose.yml up -d
	@echo "Waiting for backend..."
	@until docker compose -f $(ROOT_DIR)/docker-compose.yml exec backend wget -qO- http://localhost:8080/health > /dev/null 2>&1; do sleep 2; done
	@echo "Services are ready."
	$(ROOT_DIR)/scripts/e2e-reset-db.sh

# E2E環境停止（ボリューム削除）
e2e-down:
	docker compose -f $(ROOT_DIR)/docker-compose.yml down -v

# DB リセット（マイグレーション + シード）
e2e-reset-db:
	$(ROOT_DIR)/scripts/e2e-reset-db.sh

# 認証セットアップのみ（Emulator起動 → ログイン → storageState保存）
e2e-auth:
	@echo "Starting Firebase Auth Emulator..."
	@firebase emulators:start --only auth --project demo-sicoulab &
	@sleep 3
	@npx playwright test --project=setup; \
		EXIT_CODE=$$?; \
		echo "Stopping Firebase Auth Emulator..."; \
		kill %1 2>/dev/null || true; \
		exit $$EXIT_CODE

# Dockerイメージビルド
docker-build:
	docker build -t sicoulab-web .
