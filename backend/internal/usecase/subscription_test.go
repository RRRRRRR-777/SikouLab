package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/rs/zerolog"
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

// nopLogger はテスト用のno-opロガー。
var nopLogger = zerolog.Nop()

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
			uc := NewSubscriptionUsecase(tt.repo, nil, nopLogger)

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
//
// ウィジェット（checkout: "payment"モード）がサブスクリプション作成を行うため、
// バックエンドのCheckoutはサブスクリプションIDのDB保存のみを行う。
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
		name           string
		user           *domain.User
		subscriptionID string
		repo           *mockSubscriptionRepository
		wantErr        bool
		wantErrIs      error
	}{
		{
			name:           "正常: subscription_id保存成功",
			user:           trialingUser,
			subscriptionID: "sub_abc123",
			repo: &mockSubscriptionRepository{
				updateUnivaPaySubscriptionIDFunc: func(_ context.Context, userID int64, subID string) error {
					if userID != 1 {
						t.Errorf("期待されたuserID=1, got=%d", userID)
					}
					if subID != "sub_abc123" {
						t.Errorf("期待されたsubscriptionID=sub_abc123, got=%s", subID)
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name:           "既にactiveのユーザー",
			user:           activeUser,
			subscriptionID: "sub_abc123",
			repo:           &mockSubscriptionRepository{},
			wantErr:        true,
			wantErrIs:      ErrAlreadySubscribed,
		},
		{
			name:           "DB更新エラー時",
			user:           trialingUser,
			subscriptionID: "sub_abc123",
			repo: &mockSubscriptionRepository{
				updateUnivaPaySubscriptionIDFunc: func(_ context.Context, _ int64, _ string) error {
					return errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewSubscriptionUsecase(tt.repo, nil, nopLogger)

			err := uc.Checkout(context.Background(), tt.user, tt.subscriptionID)

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
			name: "subscription_payment + current",
			payload: WebhookPayload{
				Event: "subscription_payment",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "current",
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
			name: "subscription_payment + unpaid",
			payload: WebhookPayload{
				Event: "subscription_payment",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "unpaid",
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
			name: "subscription_failure",
			payload: WebhookPayload{
				Event: "subscription_failure",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "",
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
			name: "subscription_canceled",
			payload: WebhookPayload{
				Event: "subscription_canceled",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "",
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
				Event: "unknown_event",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "",
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
				Event: "subscription_payment",
				Data: WebhookData{
					ID:     "sub_unknown",
					Status: "current",
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
				Event: "subscription_payment",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "current",
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
				Event: "subscription_payment",
				Data: WebhookData{
					ID:     "sub_abc123",
					Status: "current",
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
			uc := NewSubscriptionUsecase(tt.repo, nil, nopLogger)

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

	uc := NewSubscriptionUsecase(repo, nil, nopLogger)
	payload := WebhookPayload{
		Event: "subscription_payment",
		Data: WebhookData{
			ID:     "sub_abc123",
			Status: "current",
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
