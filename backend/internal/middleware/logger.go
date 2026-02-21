// Package middleware はHTTPミドルウェアを提供する。
//
// リクエストログ、パニックリカバリ、CORSなど
// 横断的関心事を処理する。
package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// responseWriter はステータスコードを記録するためのラッパー。
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader はステータスコードを記録してから委譲する。
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger はリクエストログを出力するミドルウェアを返す。
// メソッド、パス、ステータスコード、レイテンシをJSON形式でINFOログ出力する。
func Logger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.statusCode).
				Dur("latency", time.Since(start)).
				Msg("リクエスト処理完了")
		})
	}
}
