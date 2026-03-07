package handler

import (
	"net/http"
	"strings"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
)

// エラーレスポンスのコード定数。
const (
	codeBadRequest    = "BAD_REQUEST"
	codeUnauthorized  = "UNAUTHORIZED"
	codeNotFound      = "NOT_FOUND"
	codeInternalError = "INTERNAL_ERROR"
)

// requireUser はリクエストコンテキストから認証済みユーザーを取得する。
// 未認証の場合は401レスポンスを返しnilを返す。
func requireUser(w http.ResponseWriter, r *http.Request) *domain.User {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{
			Code: codeUnauthorized, Message: "認証が必要です",
		})
		return nil
	}
	return user
}

// derefString はnilポインタを空文字列に変換する。
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// resolveStorageURL はオブジェクトキーをストレージのフルURLに変換する。
// 既にhttp(s)で始まるURLはそのまま返す。先頭スラッシュは除去する。
// GCS JSON API形式（/storage/v1/b/を含む）の場合は ?alt=media を付与してダウンロード可能にする。
func resolveStorageURL(baseURL string, key *string) string {
	if key == nil || *key == "" {
		return ""
	}
	k := *key
	if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
		return k
	}
	k = strings.TrimLeft(k, "/")
	u := strings.TrimRight(baseURL, "/") + "/" + k
	// GCS JSON API形式の場合、?alt=media がないとメタデータが返る
	if strings.Contains(u, "/storage/v1/b/") {
		u += "?alt=media"
	}
	return u
}
