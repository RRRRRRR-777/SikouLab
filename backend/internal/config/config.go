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
	// AppEnv はアプリケーションの実行環境。
	// "production" の場合はSecure Cookieを有効化するなど、本番向けの設定が適用される。
	AppEnv string
	// FirebaseServiceAccountJSON はFirebaseサービスアカウントのJSON。
	// 空の場合はApplication Default Credentials (ADC) を使用する。
	FirebaseServiceAccountJSON string
	// FirebaseProjectID はFirebaseプロジェクトID。
	// IDトークン検証時のaudience確認に必要。
	FirebaseProjectID string
	// UnivaPayWebhookSecret は UnivaPay Webhook 署名検証用シークレット。
	UnivaPayWebhookSecret string
	// UnivaPayStoreID は UnivaPay ストアID（APIトークンのJWT部分）。
	UnivaPayStoreID string
	// UnivaPayStoreSecret は UnivaPay ストアシークレット（APIトークンのシークレット部分）。
	UnivaPayStoreSecret string
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

	// FirebaseプロジェクトID（IDトークン検証のaudience確認に必要）
	cfg.FirebaseProjectID = os.Getenv("FIREBASE_PROJECT_ID")
	if cfg.FirebaseProjectID == "" {
		missing = append(missing, "FIREBASE_PROJECT_ID")
	}

	if len(missing) > 0 {
		return nil, errors.New("必須の環境変数が未設定です: " + strings.Join(missing, ", "))
	}

	// 任意設定
	cfg.AppEnv = os.Getenv("APP_ENV")
	cfg.FirebaseServiceAccountJSON = os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	cfg.UnivaPayWebhookSecret = os.Getenv("UNIVAPAY_WEBHOOK_SECRET")
	cfg.UnivaPayStoreID = os.Getenv("UNIVAPAY_STORE_ID")
	cfg.UnivaPayStoreSecret = os.Getenv("UNIVAPAY_STORE_SECRET")

	return cfg, nil
}
