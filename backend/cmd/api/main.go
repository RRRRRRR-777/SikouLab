// Package main はAPIサーバーのエントリポイントを提供する。
//
// 設定読み込み、DB接続、ミドルウェア設定、ルーティングを行い、HTTPサーバーを起動する。
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/bokuyamada/SicouLab/backend/internal/config"
	"github.com/bokuyamada/SicouLab/backend/internal/handler"
	"github.com/bokuyamada/SicouLab/backend/internal/middleware"
	"github.com/bokuyamada/SicouLab/backend/internal/repository"
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

	// ルーティング設定
	mux := http.NewServeMux()
	mux.Handle("GET /health", &handler.HealthHandler{})

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
