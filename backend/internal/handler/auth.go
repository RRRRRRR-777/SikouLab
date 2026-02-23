package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// セッションCookieの設定値。
const (
	sessionCookieName   = "session"
	sessionCookieMaxAge = 604800 // 7日
)

// authUsecase はAuthHandlerが依存するユースケースのインターフェース。
// テスト時にモック可能にするために定義する。
type authUsecase interface {
	Login(ctx context.Context, idToken string) (*domain.User, bool, error)
	GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, error)
}

// AuthHandler は認証APIハンドラーを提供する。
type AuthHandler struct {
	usecase authUsecase
	logger  zerolog.Logger
}

// NewAuthHandler はAuthHandlerを作成する。
func NewAuthHandler(uc authUsecase, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{usecase: uc, logger: logger}
}

// loginRequest はログインリクエストのJSON構造。
type loginRequest struct {
	IDToken string `json:"id_token"`
}

// loginResponse はログインレスポンスのJSON構造。
type loginResponse struct {
	User         userResponse `json:"user"`
	IsFirstLogin bool         `json:"is_first_login"`
}

// userResponse はユーザー情報のJSON構造。
type userResponse struct {
	ID                 int64  `json:"id"`
	OAuthProvider      string `json:"oauth_provider"`
	Name               string `json:"name"`
	DisplayName        string `json:"display_name"`
	AvatarURL          string `json:"avatar_url"`
	Role               string `json:"role"`
	PlanID             *int64 `json:"plan_id"`
	SubscriptionStatus string `json:"subscription_status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// errorResponse はエラーレスポンスのJSON構造。
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ServeLogin は POST /api/v1/auth/login を処理する。
// Firebase ID Tokenを検証し、ユーザー情報を返す。
func (h *AuthHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: "BAD_REQUEST", Message: "リクエストが不正です",
		})
		return
	}

	if req.IDToken == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: "BAD_REQUEST", Message: "リクエストが不正です",
		})
		return
	}

	user, isFirstLogin, err := h.usecase.Login(r.Context(), req.IDToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			h.logger.Warn().Err(err).Msg("ログイン失敗: 無効なIDトークン")
			writeJSON(w, http.StatusUnauthorized, errorResponse{
				Code: "INVALID_TOKEN", Message: "無効なID Tokenです",
			})
			return
		}
		h.logger.Error().Err(err).Msg("ログイン失敗: サーバーエラー")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: "INTERNAL_ERROR", Message: "サーバーエラーが発生しました",
		})
		return
	}

	// セッションCookieを設定
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    req.IDToken,
		MaxAge:   sessionCookieMaxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, loginResponse{
		User:         toUserResponse(user),
		IsFirstLogin: isFirstLogin,
	})
}

// ServeMe は GET /api/v1/auth/me を処理する。
// セッションCookieのトークンを検証してユーザー情報を返す。
func (h *AuthHandler) ServeMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{
			Code: "UNAUTHORIZED", Message: "認証が必要です",
		})
		return
	}

	user, err := h.usecase.GetCurrentUser(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			h.logger.Warn().Err(err).Msg("セッション確認失敗: 無効なトークン")
			writeJSON(w, http.StatusUnauthorized, errorResponse{
				Code: "UNAUTHORIZED", Message: "認証が必要です",
			})
			return
		}
		h.logger.Error().Err(err).Msg("セッション確認失敗: サーバーエラー")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: "INTERNAL_ERROR", Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct {
		User userResponse `json:"user"`
	}{User: toUserResponse(user)})
}

// ServeLogout はユーザーのログアウト状態を確立するため、セッションCookieを無効化する。
// POST /api/v1/auth/logout に対応し、204を返す。
func (h *AuthHandler) ServeLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

// toUserResponse はドメインモデルをレスポンス構造に変換する。
func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:                 u.ID,
		OAuthProvider:      u.OAuthProvider,
		Name:               u.Name,
		DisplayName:        u.DisplayName,
		AvatarURL:          u.AvatarURL,
		Role:               u.Role,
		PlanID:             u.PlanID,
		SubscriptionStatus: u.SubscriptionStatus,
		CreatedAt:          u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// writeJSON はJSONレスポンスを書き込む。
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
