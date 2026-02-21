package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

// Recovery はpanicを捕捉してERRORログを出力し、500レスポンスを返すミドルウェア。
// アプリケーションの安定稼働のため、ハンドラ内で発生したpanicをプロセス停止させずに処理する。
func Recovery(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().
						Interface("error", err).
						Str("stack", string(debug.Stack())).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Msg("パニック発生")

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"code":"INTERNAL_SERVER_ERROR","message":"internal server error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
