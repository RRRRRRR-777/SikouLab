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
	getPlansFunc      func(ctx context.Context) ([]domain.Plan, error)
	checkoutFunc      func(ctx context.Context, user *domain.User, tokenID string) error
	handleWebhookFunc func(ctx context.Context, payload usecase.WebhookPayload) error
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
			body:       `{"transaction_token_id":"tok_xxx"}`,
			injectUser: nil,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "transaction_token_id が空文字",
			body:       `{"transaction_token_id":""}`,
			injectUser: validUser,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "transaction_token_id フィールドなし",
			body:       `{}`,
			injectUser: validUser,
			uc:         &mockSubscriptionUsecase{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "正常",
			body:       `{"transaction_token_id":"tok_xxx"}`,
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
			body:       `{"transaction_token_id":"tok_xxx"}`,
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
			body:       `{"transaction_token_id":"tok_xxx"}`,
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
		Event: "SUBSCRIPTION_PAYMENT",
		Data: usecase.WebhookData{
			Subscriptions: usecase.WebhookSubscription{
				ID:     "sub_abc123",
				Status: "successful",
			},
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
