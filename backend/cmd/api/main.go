// Package main はAPIサーバーのエントリポイントを提供する。
//
// 設定読み込み、DB接続、ミドルウェア設定、ルーティングを行い、HTTPサーバーを起動する。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/config"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/firebase"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/handler"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/repository"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "設定読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// zerologの設定
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "不正なLOG_LEVEL: %v\n", err)
		os.Exit(1)
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)

	// DB接続
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB接続失敗")
	}
	defer db.Close()
	logger.Info().Msg("DB接続成功")

	// Firebase初期化
	ctx := context.Background()
	firebaseClient, err := firebase.NewClient(ctx, cfg.FirebaseProjectID, cfg.FirebaseServiceAccountJSON)
	if err != nil {
		logger.Fatal().Err(err).Msg("Firebase初期化失敗")
	}
	logger.Info().Msg("Firebase初期化成功")

	// リポジトリ・ユースケース・ハンドラーの初期化
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(firebaseClient, userRepo)
	authHandler := handler.NewAuthHandler(authUsecase, logger, strings.EqualFold(cfg.AppEnv, "production"))

	// サブスクリプション
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionUsecase := usecase.NewSubscriptionUsecase(subscriptionRepo, nil) // UnivaPayClientはTBD
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionUsecase, logger, cfg.UnivaPayWebhookSecret)

	// ルーティング設定
	mux := http.NewServeMux()
	mux.Handle("GET /health", &handler.HealthHandler{})

	// Go 1.22+: メソッドをHandleFuncに渡す場合、明示的にレシーバーをバインドする必要がある
	mux.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		authHandler.ServeLogin(w, r)
	})
	mux.HandleFunc("GET /api/v1/auth/me", func(w http.ResponseWriter, r *http.Request) {
		authHandler.ServeMe(w, r)
	})
	mux.HandleFunc("POST /api/v1/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		authHandler.ServeLogout(w, r)
	})
	mux.HandleFunc("GET /api/v1/plans", func(w http.ResponseWriter, r *http.Request) {
		subscriptionHandler.ServeGetPlans(w, r)
	})
	mux.HandleFunc("POST /api/v1/univapay/checkout", func(w http.ResponseWriter, r *http.Request) {
		middleware.RequireAuth(authUsecase)(http.HandlerFunc(subscriptionHandler.ServeCheckout)).ServeHTTP(w, r)
	})
	mux.HandleFunc("POST /api/v1/univapay/webhook", func(w http.ResponseWriter, r *http.Request) {
		subscriptionHandler.ServeWebhook(w, r)
	})

	// ミドルウェアチェーン: Recovery → Logger → CORS → Handler
	var h http.Handler = mux
	h = middleware.CORS(cfg.AllowedOrigins)(h)
	h = middleware.Logger(logger)(h)
	h = middleware.Recovery(logger)(h)

	addr := ":" + cfg.Port
	logger.Info().Str("addr", addr).Msg("サーバー起動")
	if err := http.ListenAndServe(addr, h); err != nil {
		fmt.Fprintf(os.Stderr, "サーバー起動失敗: %v\n", err)
		os.Exit(1)
	}
}
