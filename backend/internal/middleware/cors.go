package middleware

import (
	"net/http"
)

// CORS はAllowedOriginsに基づくCORSミドルウェアを返す。
// オリジン検証とOPTIONSプリフライトリクエストに対応する。
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	// 高速なルックアップのためmapに変換
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if _, ok := originSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				// withCredentials=trueのリクエスト（Cookie送信）に必要
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// プリフライトリクエストは即座に返す
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
