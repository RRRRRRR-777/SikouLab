package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// allowedOrigins は環境変数から読み取った許可オリジンリスト
var allowedOrigins []string

func init() {
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// 環境変数から許可オリジンを読み取り
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		// デフォルト値
		allowedOrigins = []string{"http://localhost:3000"}
	} else {
		allowedOrigins = strings.Split(origins, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}
}

// setCORSHeaders はCORSヘッダーを設定する
//
//環境変数ALLOWED_ORIGINSで指定されたオリジンのみを許可する。
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// オリジンが許可リストに含まれているかチェック
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
	}
}

// healthCheckHandler はヘルスチェックエンドポイント
//
//	@Summary	ヘルスチェック
//	@Tags		system
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}

// helloHandler はHello Worldを返すエンドポイント
//
//	@Summary	Hello World
//	@Tags		system
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/ [get]
func helloHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Hello World from SikouLab Backend!"}`)
}

func main() {
	mux := http.NewServeMux()

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", healthCheckHandler)

	// Hello Worldエンドポイント
	mux.HandleFunc("/", helloHandler)

	// サーバー起動
	addr := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
