package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// mockSubscriptionRepository はSubscriptionRepositoryのモック。
type mockSubscriptionRepository struct {
	findActivePlansFunc              func(ctx context.Context) ([]domain.Plan, error)
	findPlanByIDFunc                 func(ctx context.Context, planID int64) (*domain.Plan, error)
	updateSubscriptionStatusFunc     func(ctx context.Context, univapaySubscriptionID, status string) error
	findByUnivaPaySubscriptionIDFunc func(ctx context.Context, univapaySubscriptionID string) (*domain.User, error)
	updateUnivaPaySubscriptionIDFunc func(ctx context.Context, userID int64, subscriptionID string) error
}

// FindActivePlans はモックのアクティブプラン取得を実行する。
func (m *mockSubscriptionRepository) FindActivePlans(ctx context.Context) ([]domain.Plan, error) {
	return m.findActivePlansFunc(ctx)
}

// FindPlanByID はモックのプランID検索を実行する。
func (m *mockSubscriptionRepository) FindPlanByID(ctx context.Context, planID int64) (*domain.Plan, error) {
	return m.findPlanByIDFunc(ctx, planID)
}

// UpdateSubscriptionStatus はモックのサブスクリプションステータス更新を実行する。
func (m *mockSubscriptionRepository) UpdateSubscriptionStatus(ctx context.Context, univapaySubscriptionID, status string) error {
	return m.updateSubscriptionStatusFunc(ctx, univapaySubscriptionID, status)
}

// FindByUnivaPaySubscriptionID はモックのサブスクリプションIDによるユーザー検索を実行する。
func (m *mockSubscriptionRepository) FindByUnivaPaySubscriptionID(ctx context.Context, univapaySubscriptionID string) (*domain.User, error) {
	return m.findByUnivaPaySubscriptionIDFunc(ctx, univapaySubscriptionID)
}

// UpdateUnivaPaySubscriptionID はモックのUnivaPayサブスクリプションID更新を実行する。
func (m *mockSubscriptionRepository) UpdateUnivaPaySubscriptionID(ctx context.Context, userID int64, subscriptionID string) error {
	return m.updateUnivaPaySubscriptionIDFunc(ctx, userID, subscriptionID)
}

// mockUnivaPayClient はUnivaPayClientのモック。
type mockUnivaPayClient struct {
	createSubscriptionFunc func(ctx context.Context, tokenID string, amount int, currency string) (string, error)
}

// CreateSubscription はモックのサブスクリプション作成を実行する。
func (m *mockUnivaPayClient) CreateSubscription(ctx context.Context, tokenID string, amount int, currency string) (string, error) {
	return m.createSubscriptionFunc(ctx, tokenID, amount, currency)
}

// TestSubscriptionUsecase_GetPlans はGetPlansユースケースの各パターンを検証する。
func TestSubscriptionUsecase_GetPlans(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	activePlan := domain.Plan{
		ID:        1,
		Name:      "プレミアムプラン",
		Amount:    1000,
		Currency:  "JPY",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		repo      *mockSubscriptionRepository
		wantCount int
		wantErr   bool
	}{
		{
			name: "アクティブプランのみ取得できる",
			repo: &mockSubscriptionRepository{
				findActivePlansFunc: func(_ context.Context) ([]domain.Plan, error) {
					return []domain.Plan{activePlan}, nil
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "アクティブプランが0件",
			repo: &mockSubscriptionRepository{
				findActivePlansFunc: func(_ context.Context) ([]domain.Plan, error) {
					return []domain.Plan{}, nil
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "DBエラー時",
			repo: &mockSubscriptionRepository{
				findActivePlansFunc: func(_ context.Context) ([]domain.Plan, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSubscriptionUsecase(tt.repo, nil)

			plans, err := uc.GetPlans(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if len(plans) != tt.wantCount {
				t.Errorf("プラン件数 = %d, want %d", len(plans), tt.wantCount)
			}
		})
	}
}

// TestSubscriptionUsecase_Checkout はCheckoutユースケースの各パターンを検証する。
func TestSubscriptionUsecase_Checkout(t *testing.T) {
	trialingUser := &domain.User{
		ID:                 1,
		SubscriptionStatus: "trialing",
	}
	activeUser := &domain.User{
		ID:                 2,
		SubscriptionStatus: "active",
	}

	tests := []struct {
		name                           string
		user                           *domain.User
		tokenID                        string
		repo                           *mockSubscriptionRepository
		client                         *mockUnivaPayClient
		wantErr                        bool
		wantErrIs                      error
		wantUpdateSubscriptionIDCalled bool
	}{
		{
			name:    "正常: サブスク作成とsubscription_id保存",
			user:    trialingUser,
			tokenID: "tok_xxx",
			repo: &mockSubscriptionRepository{
				findPlanByIDFunc: func(_ context.Context, _ int64) (*domain.Plan, error) {
					return &domain.Plan{ID: 1, Amount: 1000, Currency: "JPY"}, nil
				},
				updateUnivaPaySubscriptionIDFunc: func(_ context.Context, _ int64, _ string) error {
					return nil
				},
			},
			client: &mockUnivaPayClient{
				createSubscriptionFunc: func(_ context.Context, _ string, _ int, _ string) (string, error) {
					return "sub_abc123", nil
				},
			},
			wantErr:                        false,
			wantUpdateSubscriptionIDCalled: true,
		},
		{
			name:      "既にactiveのユーザー",
			user:      activeUser,
			tokenID:   "tok_xxx",
			repo:      &mockSubscriptionRepository{},
			client:    &mockUnivaPayClient{},
			wantErr:   true,
			wantErrIs: ErrAlreadySubscribed,
		},
		{
			name:    "UnivaPay APIエラー時はDB更新されない",
			user:    trialingUser,
			tokenID: "tok_xxx",
			repo: &mockSubscriptionRepository{
				findPlanByIDFunc: func(_ context.Context, _ int64) (*domain.Plan, error) {
					return &domain.Plan{ID: 1, Amount: 1000, Currency: "JPY"}, nil
				},
				updateUnivaPaySubscriptionIDFunc: func(_ context.Context, _ int64, _ string) error {
					t.Error("UpdateUnivaPaySubscriptionID が呼ばれてはいけない")
					return nil
				},
			},
			client: &mockUnivaPayClient{
				createSubscriptionFunc: func(_ context.Context, _ string, _ int, _ string) (string, error) {
					return "", errors.New("univapay api error")
				},
			},
			wantErr:                        true,
			wantUpdateSubscriptionIDCalled: false,
		},
		{
			name:    "DB更新エラー時",
			user:    trialingUser,
			tokenID: "tok_xxx",
			repo: &mockSubscriptionRepository{
				findPlanByIDFunc: func(_ context.Context, _ int64) (*domain.Plan, error) {
					return &domain.Plan{ID: 1, Amount: 1000, Currency: "JPY"}, nil
				},
				updateUnivaPaySubscriptionIDFunc: func(_ context.Context, _ int64, _ string) error {
					return errors.New("db error")
				},
			},
			client: &mockUnivaPayClient{
				createSubscriptionFunc: func(_ context.Context, _ string, _ int, _ string) (string, error) {
					return "sub_abc123", nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSubscriptionUsecase(tt.repo, tt.client)

			err := uc.Checkout(context.Background(), tt.user, tt.tokenID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("エラー種別が一致しない: got %v, want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
		})
	}
}

// TestSubscriptionUsecase_HandleWebhook はHandleWebhookユースケースの各パターンを検証する。
func TestSubscriptionUsecase_HandleWebhook(t *testing.T) {
	existingUser := &domain.User{
		ID:                 1,
		SubscriptionStatus: "trialing",
	}

	tests := []struct {
		name          string
		payload       WebhookPayload
		repo          *mockSubscriptionRepository
		wantErr       bool
		wantNewStatus string
	}{
		{
			name: "SUBSCRIPTION_PAYMENT + successful",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_PAYMENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "successful",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, status string) error {
					if status != "active" {
						t.Errorf("期待されたステータス active と一致しない: got %s", status)
					}
					return nil
				},
			},
			wantErr:       false,
			wantNewStatus: "active",
		},
		{
			name: "SUBSCRIPTION_PAYMENT + failed",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_PAYMENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "failed",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, status string) error {
					if status != "past_due" {
						t.Errorf("期待されたステータス past_due と一致しない: got %s", status)
					}
					return nil
				},
			},
			wantErr:       false,
			wantNewStatus: "past_due",
		},
		{
			name: "SUBSCRIPTION_FAILED",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_FAILED",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, status string) error {
					if status != "past_due" {
						t.Errorf("期待されたステータス past_due と一致しない: got %s", status)
					}
					return nil
				},
			},
			wantErr:       false,
			wantNewStatus: "past_due",
		},
		{
			name: "SUBSCRIPTION_CANCELED",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_CANCELED",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, status string) error {
					if status != "canceled" {
						t.Errorf("期待されたステータス canceled と一致しない: got %s", status)
					}
					return nil
				},
			},
			wantErr:       false,
			wantNewStatus: "canceled",
		},
		{
			name: "未知イベント",
			payload: WebhookPayload{
				Event: "UNKNOWN_EVENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					// 未知イベントでは呼ばれない想定
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, _ string) error {
					t.Error("未知イベントでUpdateSubscriptionStatusが呼ばれてはいけない")
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "subscription_idに対応するユーザー不在",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_PAYMENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_unknown",
						Status: "successful",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return nil, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, _ string) error {
					t.Error("ユーザー不在でUpdateSubscriptionStatusが呼ばれてはいけない")
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "冪等性: 同一イベントを2回呼び出す",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_PAYMENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "successful",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, _ string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "DBエラー時",
			payload: WebhookPayload{
				Event: "SUBSCRIPTION_PAYMENT",
				Data: WebhookData{
					Subscriptions: WebhookSubscription{
						ID:     "sub_abc123",
						Status: "successful",
					},
				},
			},
			repo: &mockSubscriptionRepository{
				findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
					return existingUser, nil
				},
				updateSubscriptionStatusFunc: func(_ context.Context, _ string, _ string) error {
					return errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSubscriptionUsecase(tt.repo, nil)

			err := uc.HandleWebhook(context.Background(), tt.payload)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
		})
	}
}

// TestSubscriptionUsecase_HandleWebhook_Idempotency は冪等性を個別に検証する。
func TestSubscriptionUsecase_HandleWebhook_Idempotency(t *testing.T) {
	existingUser := &domain.User{
		ID:                 1,
		SubscriptionStatus: "trialing",
	}

	callCount := 0
	repo := &mockSubscriptionRepository{
		findByUnivaPaySubscriptionIDFunc: func(_ context.Context, _ string) (*domain.User, error) {
			return existingUser, nil
		},
		updateSubscriptionStatusFunc: func(_ context.Context, _ string, _ string) error {
			callCount++
			return nil
		},
	}

	uc := NewSubscriptionUsecase(repo, nil)
	payload := WebhookPayload{
		Event: "SUBSCRIPTION_PAYMENT",
		Data: WebhookData{
			Subscriptions: WebhookSubscription{
				ID:     "sub_abc123",
				Status: "successful",
			},
		},
	}

	// 1回目
	if err := uc.HandleWebhook(context.Background(), payload); err != nil {
		t.Fatalf("1回目でエラー: %v", err)
	}
	// 2回目
	if err := uc.HandleWebhook(context.Background(), payload); err != nil {
		t.Fatalf("2回目でエラー: %v", err)
	}
	if callCount != 2 {
		t.Errorf("UpdateSubscriptionStatus 呼び出し回数 = %d, want 2", callCount)
	}
}
