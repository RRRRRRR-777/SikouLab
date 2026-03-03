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
	"github.com/RRRRRRR-777/SicouLab/backend/internal/infrastructure/univapay"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/repository"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/router"
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
	univapayClient := univapay.NewClient(cfg.UnivaPayStoreID, cfg.UnivaPayStoreSecret)
	subscriptionUsecase := usecase.NewSubscriptionUsecase(subscriptionRepo, univapayClient, logger)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionUsecase, logger, cfg.UnivaPayWebhookSecret)

	// ルーティング設定
	routerBuilder := router.NewBuilder(
		router.Handlers{
			Health:       &handler.HealthHandler{},
			Auth:         authHandler,
			Subscription: subscriptionHandler,
		},
		router.Middlewares{
			AuthUsecase: authUsecase,
		},
		cfg,
	)
	h := routerBuilder.Build(logger)

	addr := ":" + cfg.Port
	logger.Info().Str("addr", addr).Msg("サーバー起動")
	if err := http.ListenAndServe(addr, h); err != nil {
		fmt.Fprintf(os.Stderr, "サーバー起動失敗: %v\n", err)
		os.Exit(1)
	}
}
