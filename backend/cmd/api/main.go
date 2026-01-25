package main

import (
	"fmt"
	"net/http"
)

// healthCheckHandler はヘルスチェックエンドポイント
//
//	@Summary	ヘルスチェック
//	@Tags		system
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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
