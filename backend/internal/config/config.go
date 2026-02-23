// Package config はアプリケーション設定の読み込みを担当する。
//
// 全ての設定値は環境変数（.envファイル）から読み込む。
// 未設定の場合はエラーとして起動を停止する。
package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定を保持する。
type Config struct {
	// Port はHTTPサーバーのリッスンポート。
	Port string
	// AllowedOrigins はCORSで許可するオリジンのリスト。
	AllowedOrigins []string
	// LogLevel はzerologのログレベル。
	LogLevel string
	// DatabaseURL はPostgreSQLの接続文字列。
	DatabaseURL string
	// FirebaseServiceAccountJSON はFirebaseサービスアカウントのJSON。
	// 空の場合はApplication Default Credentials (ADC) を使用する。
	FirebaseServiceAccountJSON string
}

// Load は.envファイルおよび環境変数から設定を読み込む。
// 必須の環境変数が未設定の場合はエラーを返す。
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New(".envファイルの読み込みに失敗しました。.env.sampleを参考に.envを作成してください")
	}

	cfg := &Config{}
	var missing []string

	cfg.Port = os.Getenv("SERVER_PORT")
	if cfg.Port == "" {
		missing = append(missing, "SERVER_PORT")
	}

	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	if cfg.LogLevel == "" {
		missing = append(missing, "LOG_LEVEL")
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		missing = append(missing, "ALLOWED_ORIGINS")
	} else {
		for _, o := range strings.Split(origins, ",") {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, strings.TrimSpace(o))
		}
	}

	if len(missing) > 0 {
		return nil, errors.New("必須の環境変数が未設定です: " + strings.Join(missing, ", "))
	}

	// 任意: Firebase認証情報（空の場合はADCを使用）
	cfg.FirebaseServiceAccountJSON = os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")

	return cfg, nil
}
