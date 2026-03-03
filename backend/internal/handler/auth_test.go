package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// mockAuthUsecase はauthUsecaseインターフェースのモック。
type mockAuthUsecase struct {
	loginFunc               func(ctx context.Context, idToken string) (*domain.User, bool, error)
	getCurrentUserFunc      func(ctx context.Context, sessionToken string) (*domain.User, error)
	createSessionCookieFunc func(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
}

// Login はモックのLogin処理を実行する。
func (m *mockAuthUsecase) Login(ctx context.Context, idToken string) (*domain.User, bool, error) {
	return m.loginFunc(ctx, idToken)
}

// GetCurrentUser はモックのGetCurrentUser処理を実行する。
func (m *mockAuthUsecase) GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, error) {
	return m.getCurrentUserFunc(ctx, sessionToken)
}

// CreateSessionCookie はモックのセッションCookie生成を実行する。
func (m *mockAuthUsecase) CreateSessionCookie(_ context.Context, _ string, _ time.Duration) (string, error) {
	if m.createSessionCookieFunc != nil {
		return m.createSessionCookieFunc(context.Background(), "", 0)
	}
	return "mock-session-cookie", nil
}

// TestAuthHandler_ServeLogin はServeLoginハンドラーの各パターンを検証する。
func TestAuthHandler_ServeLogin(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	validUser := &domain.User{
		ID: 1, OAuthProvider: "google.com", Name: "Test User",
		Role: "user", SubscriptionStatus: "active", CreatedAt: now, UpdatedAt: now,
	}

	tests := []struct {
		name         string
		body         string
		uc           *mockAuthUsecase
		wantStatus   int
		wantCookie   bool
		secureCookie bool
		wantSecure   bool
	}{
		{
			name: "本番環境ではSecure=trueのCookieが返される",
			body: `{"id_token": "valid-token"}`,
			uc: &mockAuthUsecase{
				loginFunc: func(_ context.Context, _ string) (*domain.User, bool, error) {
					return validUser, false, nil
				},
			},
			wantStatus:   http.StatusOK,
			wantCookie:   true,
			secureCookie: true,
			wantSecure:   true,
		},
		{
			name: "開発環境ではSecure=falseのCookieが返される",
			body: `{"id_token": "valid-token"}`,
			uc: &mockAuthUsecase{
				loginFunc: func(_ context.Context, _ string) (*domain.User, bool, error) {
					return validUser, false, nil
				},
			},
			wantStatus:   http.StatusOK,
			wantCookie:   true,
			secureCookie: false,
			wantSecure:   false,
		},
		{
			name:       "id_tokenが空の場合は400が返される",
			body:       `{"id_token": ""}`,
			uc:         &mockAuthUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "不正なJSONボディの場合は400が返される",
			body:       `invalid-json`,
			uc:         &mockAuthUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "ErrInvalidTokenの場合は401が返される",
			body: `{"id_token": "invalid-token"}`,
			uc: &mockAuthUsecase{
				loginFunc: func(_ context.Context, _ string) (*domain.User, bool, error) {
					return nil, false, usecase.ErrInvalidToken
				},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "その他のエラーの場合は500が返される",
			body: `{"id_token": "valid-token"}`,
			uc: &mockAuthUsecase{
				loginFunc: func(_ context.Context, _ string) (*domain.User, bool, error) {
					return nil, false, errors.New("unexpected error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{usecase: tt.uc, secureCookie: tt.secureCookie}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.ServeLogin(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantCookie {
				cookies := rec.Result().Cookies()
				var found bool
				for _, c := range cookies {
					if c.Name == sessionCookieName {
						found = true
						if !c.HttpOnly {
							t.Error("Cookie HttpOnly = false, want true")
						}
						if c.Secure != tt.wantSecure {
							t.Errorf("Cookie Secure = %v, want %v", c.Secure, tt.wantSecure)
						}
					}
				}
				if !found {
					t.Errorf("Cookie '%s' が見つからない", sessionCookieName)
				}
			}
		})
	}
}

// TestAuthHandler_ServeMe はServeMeハンドラーの各パターンを検証する。
func TestAuthHandler_ServeMe(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	validUser := &domain.User{
		ID: 1, OAuthProvider: "google.com", Name: "Test User",
		Role: "user", SubscriptionStatus: "active", CreatedAt: now, UpdatedAt: now,
	}

	tests := []struct {
		name       string
		cookie     *http.Cookie
		uc         *mockAuthUsecase
		wantStatus int
	}{
		{
			name:   "有効なCookieで200とユーザー情報が返される",
			cookie: &http.Cookie{Name: sessionCookieName, Value: "valid-session"},
			uc: &mockAuthUsecase{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return validUser, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Cookieがない場合は401が返される",
			cookie:     nil,
			uc:         &mockAuthUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "ErrInvalidTokenの場合は401が返される",
			cookie: &http.Cookie{Name: sessionCookieName, Value: "invalid-session"},
			uc: &mockAuthUsecase{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return nil, usecase.ErrInvalidToken
				},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "その他のエラーの場合は500が返される",
			cookie: &http.Cookie{Name: sessionCookieName, Value: "valid-session"},
			uc: &mockAuthUsecase{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return nil, errors.New("unexpected error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{usecase: tt.uc}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			rec := httptest.NewRecorder()

			h.ServeMe(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// TestAuthHandler_ServeLogout はログアウト処理のCookie削除を検証する。
func TestAuthHandler_ServeLogout(t *testing.T) {
	tests := []struct {
		name           string
		wantStatus     int
		wantCookieName string
		wantMaxAge     int
		secureCookie   bool
		wantSecure     bool
	}{
		{
			name:           "本番環境のログアウトCookieはSecure=true",
			wantStatus:     http.StatusNoContent,
			wantCookieName: "session",
			wantMaxAge:     -1,
			secureCookie:   true,
			wantSecure:     true,
		},
		{
			name:           "開発環境のログアウトCookieはSecure=false",
			wantStatus:     http.StatusNoContent,
			wantCookieName: "session",
			wantMaxAge:     -1,
			secureCookie:   false,
			wantSecure:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{secureCookie: tt.secureCookie}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
			rec := httptest.NewRecorder()

			h.ServeLogout(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			cookies := rec.Result().Cookies()
			var found bool
			for _, c := range cookies {
				if c.Name == tt.wantCookieName {
					found = true
					if c.MaxAge != tt.wantMaxAge {
						t.Errorf("Cookie MaxAge = %d, want %d", c.MaxAge, tt.wantMaxAge)
					}
					if !c.HttpOnly {
						t.Error("Cookie HttpOnly = false, want true")
					}
					if c.Secure != tt.wantSecure {
						t.Errorf("Cookie Secure = %v, want %v", c.Secure, tt.wantSecure)
					}
				}
			}
			if !found {
				t.Errorf("Cookie '%s' が見つからない", tt.wantCookieName)
			}
		})
	}
}
