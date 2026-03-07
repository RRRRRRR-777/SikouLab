package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// mockNewsletterUsecase はnewsletterUsecaseインターフェースのモック。
type mockNewsletterUsecase struct {
	getSubscriptionFunc func(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	subscribeFunc       func(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
	unsubscribeFunc     func(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	updateEmailFunc     func(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
}

// GetSubscription はモックのGetSubscription処理を実行する。
func (m *mockNewsletterUsecase) GetSubscription(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	return m.getSubscriptionFunc(ctx, userID)
}

// Subscribe はモックのSubscribe処理を実行する。
func (m *mockNewsletterUsecase) Subscribe(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	return m.subscribeFunc(ctx, userID, email)
}

// Unsubscribe はモックのUnsubscribe処理を実行する。
func (m *mockNewsletterUsecase) Unsubscribe(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	return m.unsubscribeFunc(ctx, userID)
}

// UpdateEmail はモックのUpdateEmail処理を実行する。
func (m *mockNewsletterUsecase) UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	return m.updateEmailFunc(ctx, userID, email)
}

// testUser はテスト用の認証済みユーザーを返すヘルパー。
func testUser() *domain.User {
	return &domain.User{
		ID:   1,
		Name: "Test User",
		Role: "user",
	}
}

// testSubscription はテスト用の購読レコードを返すヘルパー。
func testSubscription() *domain.NewsletterSubscription {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return &domain.NewsletterSubscription{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TestNewsletterHandler_ServeGetSubscription はServeGetSubscriptionハンドラーの各パターンを検証する。
func TestNewsletterHandler_ServeGetSubscription(t *testing.T) {
	tests := []struct {
		name       string
		injectUser *domain.User
		uc         *mockNewsletterUsecase
		wantStatus int
		wantEmail  string
	}{
		{
			name:       "正常: 購読登録済みで200が返る",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				getSubscriptionFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return testSubscription(), nil
				},
			},
			wantStatus: http.StatusOK,
			wantEmail:  "test@example.com",
		},
		{
			name:       "未登録で404が返る",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				getSubscriptionFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "未認証で401が返る",
			injectUser: nil,
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "サーバーエラーで500が返る",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				getSubscriptionFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewNewsletterHandler(tt.uc, zerolog.Nop())
			req := httptest.NewRequest(http.MethodGet, "/api/v1/newsletter/subscription", nil)

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeGetSubscription(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantEmail != "" {
				var resp newsletterSubscriptionResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("レスポンスのデコード失敗: %v", err)
				}
				if resp.Email != tt.wantEmail {
					t.Errorf("Email = %q, want %q", resp.Email, tt.wantEmail)
				}
			}
		})
	}
}

// TestNewsletterHandler_ServeSubscribe はServeSubscribeハンドラーの各パターンを検証する。
func TestNewsletterHandler_ServeSubscribe(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		injectUser *domain.User
		uc         *mockNewsletterUsecase
		wantStatus int
	}{
		{
			name:       "正常: 正常なメールで201が返る",
			body:       `{"email":"user@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				subscribeFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := testSubscription()
					sub.Email = email
					return sub, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "不正なメール形式で400が返る",
			body:       `{"email":"not-an-email"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				subscribeFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrInvalidEmail
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "空文字で400が返る",
			body:       `{"email":""}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				subscribeFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrInvalidEmail
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "既存レコードがあればupsertで201が返る",
			body:       `{"email":"updated@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				subscribeFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := testSubscription()
					sub.Email = email
					sub.IsActive = true
					return sub, nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "未認証で401が返る",
			body:       `{"email":"user@example.com"}`,
			injectUser: nil,
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "サーバーエラーで500が返る",
			body:       `{"email":"user@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				subscribeFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "不正なJSONで400が返る",
			body:       `{invalid json`,
			injectUser: testUser(),
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewNewsletterHandler(tt.uc, zerolog.Nop())
			req := httptest.NewRequest(http.MethodPost, "/api/v1/newsletter/subscribe", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeSubscribe(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// TestNewsletterHandler_ServeUnsubscribe はServeUnsubscribeハンドラーの各パターンを検証する。
func TestNewsletterHandler_ServeUnsubscribe(t *testing.T) {
	tests := []struct {
		name       string
		injectUser *domain.User
		uc         *mockNewsletterUsecase
		wantStatus int
		wantActive bool
	}{
		{
			name:       "正常: 購読中から解除で200が返る（is_active=false）",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				unsubscribeFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					sub := testSubscription()
					sub.IsActive = false
					return sub, nil
				},
			},
			wantStatus: http.StatusOK,
			wantActive: false,
		},
		{
			name:       "未登録で404が返る",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				unsubscribeFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "未認証で401が返る",
			injectUser: nil,
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "サーバーエラーで500が返る",
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				unsubscribeFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewNewsletterHandler(tt.uc, zerolog.Nop())
			req := httptest.NewRequest(http.MethodPost, "/api/v1/newsletter/unsubscribe", nil)

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeUnsubscribe(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp newsletterSubscriptionResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("レスポンスのデコード失敗: %v", err)
				}
				if resp.IsActive != tt.wantActive {
					t.Errorf("IsActive = %v, want %v", resp.IsActive, tt.wantActive)
				}
			}
		})
	}
}

// TestNewsletterHandler_ServeUpdateEmail はServeUpdateEmailハンドラーの各パターンを検証する。
func TestNewsletterHandler_ServeUpdateEmail(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		injectUser *domain.User
		uc         *mockNewsletterUsecase
		wantStatus int
		wantEmail  string
	}{
		{
			name:       "正常: 正常なメールで200が返る",
			body:       `{"email":"new@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				updateEmailFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := testSubscription()
					sub.Email = email
					return sub, nil
				},
			},
			wantStatus: http.StatusOK,
			wantEmail:  "new@example.com",
		},
		{
			name:       "不正なメール形式で400が返る",
			body:       `{"email":"invalid"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				updateEmailFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrInvalidEmail
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "未登録で404が返る",
			body:       `{"email":"new@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				updateEmailFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, usecase.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "未認証で401が返る",
			body:       `{"email":"new@example.com"}`,
			injectUser: nil,
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "サーバーエラーで500が返る",
			body:       `{"email":"new@example.com"}`,
			injectUser: testUser(),
			uc: &mockNewsletterUsecase{
				updateEmailFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "不正なJSONで400が返る",
			body:       `{invalid json`,
			injectUser: testUser(),
			uc:         &mockNewsletterUsecase{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewNewsletterHandler(tt.uc, zerolog.Nop())
			req := httptest.NewRequest(http.MethodPut, "/api/v1/newsletter/subscription", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeUpdateEmail(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantEmail != "" {
				var resp newsletterSubscriptionResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("レスポンスのデコード失敗: %v", err)
				}
				if resp.Email != tt.wantEmail {
					t.Errorf("Email = %q, want %q", resp.Email, tt.wantEmail)
				}
			}
		})
	}
}
