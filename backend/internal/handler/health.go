// Package handler はHTTPリクエストハンドラを提供する。
package handler

import (
	"net/http"
)

// HealthHandler はヘルスチェックエンドポイントのハンドラ。
// アプリケーションの死活監視に使用する。
type HealthHandler struct{}

// ServeHTTP は GET /health に対してステータスokを返す。
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
