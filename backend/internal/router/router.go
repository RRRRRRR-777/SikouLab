// Package router はHTTPルーティングの構築を提供する。
//
// ハンドラーとミドルウェアを組み合わせ、ルーティング設定済みのhttp.Handlerを生成する。
package router

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/config"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/handler"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// Handlers はルーティングに必要な全ハンドラーを集約した構造体。
type Handlers struct {
	Health       *handler.HealthHandler
	Auth         *handler.AuthHandler
	Subscription *handler.SubscriptionHandler
}

// Middlewares はルーティングに必要なミドルウェアを集約した構造体。
// AuthUsecase は middleware.RequireAuth が受け取る authSessionVerifier を満たす型を設定する。
// テスト時はインターフェースを満たすモックを直接設定可能。
type Middlewares struct {
	AuthUsecase *usecase.AuthUsecase
}

// Builder はルーターを構築するビルダー。
type Builder struct {
	handlers    Handlers
	middlewares Middlewares
	cfg         *config.Config
}

// NewBuilder はルータービルダーを作成する。
func NewBuilder(handlers Handlers, middlewares Middlewares, cfg *config.Config) *Builder {
	return &Builder{
		handlers:    handlers,
		middlewares: middlewares,
		cfg:         cfg,
	}
}

// Build はルーティング設定済みのhttp.Handlerを返す。
//
// ミドルウェアチェーン: Recovery -> Logger -> CORS -> Handler
func (b *Builder) Build(logger zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health: 認証なし
	mux.Handle("GET /health", b.handlers.Health)

	// Auth: 認証なし
	mux.HandleFunc("POST /api/v1/auth/login", b.handlers.Auth.ServeLogin)
	mux.HandleFunc("GET /api/v1/auth/me", b.handlers.Auth.ServeMe)
	mux.HandleFunc("POST /api/v1/auth/logout", b.handlers.Auth.ServeLogout)

	// Subscription: 公開エンドポイント
	mux.HandleFunc("GET /api/v1/plans", b.handlers.Subscription.ServeGetPlans)

	// Subscription: 認証必須エンドポイント
	mux.Handle("POST /api/v1/univapay/checkout",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Subscription.ServeCheckout),
		),
	)

	// Webhook: 認証なし（独自の署名検証）
	mux.HandleFunc("POST /api/v1/univapay/webhook", b.handlers.Subscription.ServeWebhook)

	// ミドルウェアチェーン適用
	var h http.Handler = mux
	h = middleware.CORS(b.cfg.AllowedOrigins)(h)
	h = middleware.Logger(logger)(h)
	h = middleware.Recovery(logger)(h)

	return h
}
