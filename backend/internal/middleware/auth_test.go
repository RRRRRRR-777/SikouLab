package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
)

// mockSessionVerifier は authSessionVerifier インターフェースのモック。
type mockSessionVerifier struct {
	getCurrentUserFunc func(ctx context.Context, sessionToken string) (*domain.User, error)
}

// GetCurrentUser はモックの GetCurrentUser 処理を実行する。
func (m *mockSessionVerifier) GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, error) {
	return m.getCurrentUserFunc(ctx, sessionToken)
}

// newTestUser はテスト用のユーザーを生成するヘルパー。
func newTestUser(role, subscriptionStatus string) *domain.User {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return &domain.User{
		ID:                 1,
		OAuthProvider:      "google.com",
		Name:               "テストユーザー",
		Role:               role,
		SubscriptionStatus: subscriptionStatus,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// okHandler は常に 200 OK を返すテスト用ハンドラー。
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// TestAuthMiddleware_RequireAuth は RequireAuth ミドルウェアの各パターンを検証する。
func TestAuthMiddleware_RequireAuth(t *testing.T) {
	validUser := newTestUser("user", "active")

	tests := []struct {
		name       string
		cookie     *http.Cookie
		uc         *mockSessionVerifier
		wantStatus int
		wantUser   bool
	}{
		{
			name:   "Cookieあり・有効なトークンの場合はnextが呼ばれる",
			cookie: &http.Cookie{Name: "session", Value: "valid-token"},
			uc: &mockSessionVerifier{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return validUser, nil
				},
			},
			wantStatus: http.StatusOK,
			wantUser:   true,
		},
		{
			name:       "Cookieなしの場合は401が返される",
			cookie:     nil,
			uc:         &mockSessionVerifier{},
			wantStatus: http.StatusUnauthorized,
			wantUser:   false,
		},
		{
			name:   "Cookieあり・無効なトークン（検証失敗）の場合は401が返される",
			cookie: &http.Cookie{Name: "session", Value: "invalid-token"},
			uc: &mockSessionVerifier{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return nil, context.DeadlineExceeded // トークン検証失敗を想定
				},
			},
			wantStatus: http.StatusUnauthorized,
			wantUser:   false,
		},
		{
			name:   "Cookieあり・トークンは有効だがDBにユーザーが存在しない場合は401",
			cookie: &http.Cookie{Name: "session", Value: "valid-token-no-user"},
			uc: &mockSessionVerifier{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return nil, nil
				},
			},
			wantStatus: http.StatusUnauthorized,
			wantUser:   false,
		},
		{
			name:   "contextにUserが設定されること",
			cookie: &http.Cookie{Name: "session", Value: "valid-token"},
			uc: &mockSessionVerifier{
				getCurrentUserFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return validUser, nil
				},
			},
			wantStatus: http.StatusOK,
			wantUser:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// contextからユーザーを取り出して検証するハンドラー
			var capturedUser *domain.User
			checkHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 未実装の UserFromContext でユーザーを取り出す（Redフェーズ: コンパイルエラー想定）
				capturedUser = middleware.UserFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			// 未実装の RequireAuth ミドルウェアを使用（Redフェーズ: コンパイルエラー想定）
			h := middleware.RequireAuth(tt.uc)(checkHandler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantUser && capturedUser == nil {
				t.Error("contextにUserが設定されていない")
			}
			if !tt.wantUser && capturedUser != nil {
				t.Error("contextにUserが設定されているが、設定されていないことを期待")
			}
		})
	}
}

// TestAuthMiddleware_RequireRole は RequireRole ミドルウェアの各パターンを検証する。
func TestAuthMiddleware_RequireRole(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		allowedRoles []string
		wantStatus   int
	}{
		{
			name:         "admin roleでadminエンドポイントにアクセスできる",
			userRole:     "admin",
			allowedRoles: []string{"admin"},
			wantStatus:   http.StatusOK,
		},
		{
			name:         "user roleでadminエンドポイントにアクセスすると403",
			userRole:     "user",
			allowedRoles: []string{"admin"},
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "writer roleでadminエンドポイントにアクセスすると403",
			userRole:     "writer",
			allowedRoles: []string{"admin"},
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "writer roleでwriterエンドポイントにアクセスできる",
			userRole:     "writer",
			allowedRoles: []string{"writer"},
			wantStatus:   http.StatusOK,
		},
		{
			name:         "user roleでwriterエンドポイントにアクセスすると403",
			userRole:     "user",
			allowedRoles: []string{"writer"},
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "admin roleはwriter権限エンドポイントにもアクセスできる（admin > writer > user）",
			userRole:     "admin",
			allowedRoles: []string{"writer"},
			wantStatus:   http.StatusOK,
		},
		{
			name:         "contextにUserがない場合は401",
			userRole:     "", // Userなし（contextに設定しない）
			allowedRoles: []string{"admin"},
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "不明なroleの場合は403",
			userRole:     "unknown",
			allowedRoles: []string{"admin"},
			wantStatus:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 未実装の RequireRole ミドルウェアを使用（Redフェーズ: コンパイルエラー想定）
			h := middleware.RequireRole(tt.allowedRoles...)(okHandler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin", nil)

			// contextにUserがない場合はそのまま、ある場合は設定する
			if tt.userRole != "" {
				user := newTestUser(tt.userRole, "active")
				// 未実装の contextWithUser でcontextにuserを設定する想定
				// 実際のテストではRequireAuthが設定したcontextを利用するが、
				// ここでは直接contextにuserを設定するためのヘルパーが必要
				ctx := middleware.ContextWithUser(req.Context(), user)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// TestAuthMiddleware_RequireSubscription は RequireSubscription ミドルウェアの各パターンを検証する。
func TestAuthMiddleware_RequireSubscription(t *testing.T) {
	tests := []struct {
		name               string
		subscriptionStatus string
		hasUser            bool
		wantStatus         int
	}{
		{
			name:               "subscription_status=activeのユーザーはアクセスできる",
			subscriptionStatus: "active",
			hasUser:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "subscription_status=canceledのユーザーは403",
			subscriptionStatus: "canceled",
			hasUser:            true,
			wantStatus:         http.StatusForbidden,
		},
		{
			name:               "subscription_status=past_dueのユーザーは403",
			subscriptionStatus: "past_due",
			hasUser:            true,
			wantStatus:         http.StatusForbidden,
		},
		{
			name:               "contextにUserがない場合は401",
			subscriptionStatus: "",
			hasUser:            false,
			wantStatus:         http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 未実装の RequireSubscription ミドルウェアを使用（Redフェーズ: コンパイルエラー想定）
			h := middleware.RequireSubscription()(okHandler)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription-required", nil)

			if tt.hasUser {
				user := newTestUser("user", tt.subscriptionStatus)
				ctx := middleware.ContextWithUser(req.Context(), user)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
