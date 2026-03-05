package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// mockSubscriptionUsecase はsubscriptionUsecaseインターフェースのモック。
type mockSubscriptionUsecase struct {
	getPlansFunc          func(ctx context.Context) ([]domain.Plan, error)
	checkoutFunc          func(ctx context.Context, user *domain.User, tokenID string) error
	handleWebhookFunc     func(ctx context.Context, payload usecase.WebhookPayload) error
	getMySubscriptionFunc func(ctx context.Context, user *domain.User) (*usecase.SubscriptionInfo, error)
	generatePortalURLFunc func(ctx context.Context, user *domain.User) (string, error)
}

// GetPlans はモックのGetPlans処理を実行する。
func (m *mockSubscriptionUsecase) GetPlans(ctx context.Context) ([]domain.Plan, error) {
	return m.getPlansFunc(ctx)
}

// Checkout はモックのCheckout処理を実行する。
func (m *mockSubscriptionUsecase) Checkout(ctx context.Context, user *domain.User, tokenID string) error {
	return m.checkoutFunc(ctx, user, tokenID)
}

// HandleWebhook はモックのHandleWebhook処理を実行する。
func (m *mockSubscriptionUsecase) HandleWebhook(ctx context.Context, payload usecase.WebhookPayload) error {
	return m.handleWebhookFunc(ctx, payload)
}

// GetMySubscription はモックのGetMySubscription処理を実行する。
func (m *mockSubscriptionUsecase) GetMySubscription(ctx context.Context, user *domain.User) (*usecase.SubscriptionInfo, error) {
	if m.getMySubscriptionFunc != nil {
		return m.getMySubscriptionFunc(ctx, user)
	}
	return &usecase.SubscriptionInfo{}, nil
}

// GeneratePortalURL はモックのGeneratePortalURL処理を実行する。
func (m *mockSubscriptionUsecase) GeneratePortalURL(ctx context.Context, user *domain.User) (string, error) {
	if m.generatePortalURLFunc != nil {
		return m.generatePortalURLFunc(ctx, user)
	}
	return "", nil
}

// TestSubscriptionHandler_ServeGetPlans はServeGetPlansハンドラーの各パターンを検証する。
func TestSubscriptionHandler_ServeGetPlans(t *testing.T) {
	tests := []struct {
		name       string
		uc         *mockSubscriptionUsecase
		wantStatus int
	}{
		{
			name: "正常: プラン一覧が返る",
			uc: &mockSubscriptionUsecase{
				getPlansFunc: func(_ context.Context) ([]domain.Plan, error) {
					return []domain.Plan{
						{ID: 1, Name: "プレミアムプラン", Amount: 1000, Currency: "JPY", IsActive: true},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "usecaseエラー",
			uc: &mockSubscriptionUsecase{
				getPlansFunc: func(_ context.Context) ([]domain.Plan, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSubscriptionHandler(tt.uc, zerolog.Nop(), "test-secret")
			req := httptest.NewRequest(http.MethodGet, "/api/v1/plans", nil)
			rec := httptest.NewRecorder()

			h.ServeGetPlans(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// TestSubscriptionHandler_ServeCheckout はServeCheckoutハンドラーの各パターンを検証する。
func TestSubscriptionHandler_ServeCheckout(t *testing.T) {
	validUser := &domain.User{
		ID:                 1,
		Name:               "Test User",
		Role:               "user",
		SubscriptionStatus: "trialing",
	}

	tests := []struct {
		name       string
		body       string
		injectUser *domain.User
		uc         *mockSubscriptionUsecase
		wantStatus int
		wantBody   string
	}{
		{
			name:       "未認証（contextにUserなし）",
			body:       `{"subscription_id":"sub_xxx"}`,
			injectUser: nil,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "subscription_id が空文字",
			body:       `{"subscription_id":""}`,
			injectUser: validUser,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "subscription_id フィールドなし",
			body:       `{}`,
			injectUser: validUser,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "正常",
			body:       `{"subscription_id":"sub_xxx"}`,
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				checkoutFunc: func(_ context.Context, _ *domain.User, _ string) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"pending"`,
		},
		{
			name:       "既にactive",
			body:       `{"subscription_id":"sub_xxx"}`,
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				checkoutFunc: func(_ context.Context, _ *domain.User, _ string) error {
					return usecase.ErrAlreadySubscribed
				},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "その他usecaseエラー",
			body:       `{"subscription_id":"sub_xxx"}`,
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				checkoutFunc: func(_ context.Context, _ *domain.User, _ string) error {
					return errors.New("unexpected error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSubscriptionHandler(tt.uc, zerolog.Nop(), "test-secret")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/univapay/checkout", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// contextにユーザーを注入（認証済みの場合のみ）
			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()

			h.ServeCheckout(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if !strings.Contains(rec.Body.String(), tt.wantBody) {
					t.Errorf("レスポンスボディに %q が含まれない: got %q", tt.wantBody, rec.Body.String())
				}
			}
		})
	}
}

// TestSubscriptionHandler_ServeWebhook はServeWebhookハンドラーの各パターンを検証する。
func TestSubscriptionHandler_ServeWebhook(t *testing.T) {
	const webhookSecret = "test-webhook-secret"

	validPayload := usecase.WebhookPayload{
		Event: "subscription_payment",
		Data: usecase.WebhookData{
			ID:     "sub_abc123",
			Status: "current",
		},
	}
	validBodyBytes, _ := json.Marshal(validPayload)
	validBody := string(validBodyBytes)

	tests := []struct {
		name       string
		body       string
		authHeader string
		uc         *mockSubscriptionUsecase
		wantStatus int
	}{
		{
			name:       "Authorizationヘッダーなし",
			body:       validBody,
			authHeader: "",
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "不正なAuthorization（不一致）",
			body:       validBody,
			authHeader: "wrong-secret",
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "正しいAuthorizationで正常イベント",
			body:       validBody,
			authHeader: "test-webhook-secret",
			uc: &mockSubscriptionUsecase{
				handleWebhookFunc: func(_ context.Context, _ usecase.WebhookPayload) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "usecaseエラー",
			body:       validBody,
			authHeader: "test-webhook-secret",
			uc: &mockSubscriptionUsecase{
				handleWebhookFunc: func(_ context.Context, _ usecase.WebhookPayload) error {
					return errors.New("usecase error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSubscriptionHandler(tt.uc, zerolog.Nop(), webhookSecret)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/univapay/webhook", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			h.ServeWebhook(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

// TestSubscriptionHandler_ServeGetMySubscription はServeGetMySubscriptionハンドラーの各パターンを検証する。
func TestSubscriptionHandler_ServeGetMySubscription(t *testing.T) {
	validUser := &domain.User{
		ID:                 1,
		PlanID:             int64Ptr(1),
		SubscriptionStatus: "active",
	}

	tests := []struct {
		name       string
		injectUser *domain.User
		uc         *mockSubscriptionUsecase
		wantStatus int
		wantBody   string
	}{
		{
			name:       "正常: サブスクリプション情報が返る",
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				getMySubscriptionFunc: func(_ context.Context, _ *domain.User) (*usecase.SubscriptionInfo, error) {
					return &usecase.SubscriptionInfo{
						PlanName: "ベースプラン",
						Amount:   980,
						Currency: "JPY",
						Status:   "active",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"plan_name":"ベースプラン"`,
		},
		{
			name:       "未認証で401が返る",
			injectUser: nil,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "usecaseエラーで500が返る",
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				getMySubscriptionFunc: func(_ context.Context, _ *domain.User) (*usecase.SubscriptionInfo, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSubscriptionHandler(tt.uc, zerolog.Nop(), "test-secret")
			req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions/me", nil)

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeGetMySubscription(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if !strings.Contains(rec.Body.String(), tt.wantBody) {
					t.Errorf("レスポンスボディに %q が含まれない: got %q", tt.wantBody, rec.Body.String())
				}
			}
		})
	}
}

// TestSubscriptionHandler_ServeGeneratePortalURL はServeGeneratePortalURLハンドラーの各パターンを検証する。
func TestSubscriptionHandler_ServeGeneratePortalURL(t *testing.T) {
	validUser := &domain.User{
		ID:                 1,
		SubscriptionStatus: "active",
	}

	tests := []struct {
		name       string
		injectUser *domain.User
		uc         *mockSubscriptionUsecase
		wantStatus int
		wantBody   string
	}{
		{
			name:       "正常: ポータルURLが返る",
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				generatePortalURLFunc: func(_ context.Context, _ *domain.User) (string, error) {
					return "https://widget.univapay.com/portal?customer=cust_123", nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"portal_url"`,
		},
		{
			name:       "サブスク未登録で404が返る",
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				generatePortalURLFunc: func(_ context.Context, _ *domain.User) (string, error) {
					return "", usecase.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "未認証で401が返る",
			injectUser: nil,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "usecaseエラーで500が返る",
			injectUser: validUser,
			uc: &mockSubscriptionUsecase{
				generatePortalURLFunc: func(_ context.Context, _ *domain.User) (string, error) {
					return "", errors.New("internal error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSubscriptionHandler(tt.uc, zerolog.Nop(), "test-secret")
			req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions/portal", nil)

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeGeneratePortalURL(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if !strings.Contains(rec.Body.String(), tt.wantBody) {
					t.Errorf("レスポンスボディに %q が含まれない: got %q", tt.wantBody, rec.Body.String())
				}
			}
		})
	}
}

// int64Ptr はint64ポインタを返すヘルパー。
func int64Ptr(v int64) *int64 { return &v }
