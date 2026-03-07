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
	User         *handler.UserHandler
	Newsletter   *handler.NewsletterHandler
}

// Middlewares はミドルウェアの依存をBuild時に注入するための構造体。
//
// AuthUsecase は middleware.RequireAuth が受け取る authSessionVerifier を満たす。
// テスト時はインターフェースを満たすモックを直接差し替え可能。
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

// Build はハンドラーとミドルウェアを組み合わせてルーティング設定済みのhttp.Handlerを構築する。
//
// ミドルウェアチェーン: Recovery -> Logger -> CORS -> Handler。
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

	// Subscription: 認証必須エンドポイント（サブスク状態取得・ポータル）
	mux.Handle("GET /api/v1/subscriptions/me",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Subscription.ServeGetMySubscription),
		),
	)
	mux.Handle("POST /api/v1/subscriptions/portal",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Subscription.ServeGeneratePortalURL),
		),
	)

	// Webhook: 認証なし（独自の署名検証）
	mux.HandleFunc("POST /api/v1/univapay/webhook", b.handlers.Subscription.ServeWebhook)

	// User Profile: 認証必須エンドポイント
	mux.Handle("PATCH /api/v1/users/me",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.User.ServeUpdateProfile),
		),
	)
	mux.Handle("POST /api/v1/users/avatar",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.User.ServeUploadAvatar),
		),
	)
	mux.Handle("DELETE /api/v1/users/avatar",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.User.ServeDeleteAvatar),
		),
	)

	// Newsletter: 認証必須エンドポイント
	mux.Handle("GET /api/v1/newsletter/subscription",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Newsletter.ServeGetSubscription),
		),
	)
	mux.Handle("POST /api/v1/newsletter/subscribe",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Newsletter.ServeSubscribe),
		),
	)
	mux.Handle("POST /api/v1/newsletter/unsubscribe",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Newsletter.ServeUnsubscribe),
		),
	)
	mux.Handle("PUT /api/v1/newsletter/subscription",
		middleware.RequireAuth(b.middlewares.AuthUsecase)(
			http.HandlerFunc(b.handlers.Newsletter.ServeUpdateEmail),
		),
	)

	// ミドルウェアチェーン適用
	var h http.Handler = mux
	h = middleware.CORS(b.cfg.AllowedOrigins)(h)
	h = middleware.Logger(logger)(h)
	h = middleware.Recovery(logger)(h)

	return h
}
