package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// mockNewsletterRepository はNewsletterRepositoryのモック。
type mockNewsletterRepository struct {
	findByUserIDFunc   func(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	upsertFunc         func(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
	updateIsActiveFunc func(ctx context.Context, userID int64, isActive bool) (*domain.NewsletterSubscription, error)
	updateEmailFunc    func(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
}

// FindByUserID はモックの購読レコード検索を実行する。
func (m *mockNewsletterRepository) FindByUserID(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	return m.findByUserIDFunc(ctx, userID)
}

// Upsert はモックの購読レコード作成・更新を実行する。
func (m *mockNewsletterRepository) Upsert(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	return m.upsertFunc(ctx, userID, email)
}

// UpdateIsActive はモックの購読状態更新を実行する。
func (m *mockNewsletterRepository) UpdateIsActive(ctx context.Context, userID int64, isActive bool) (*domain.NewsletterSubscription, error) {
	return m.updateIsActiveFunc(ctx, userID, isActive)
}

// UpdateEmail はモックのメールアドレス更新を実行する。
func (m *mockNewsletterRepository) UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	return m.updateEmailFunc(ctx, userID, email)
}

// newTestSubscription はテスト用の購読レコードを返すヘルパー。
func newTestSubscription() *domain.NewsletterSubscription {
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

// TestNewsletterUsecase_GetSubscription はGetSubscriptionユースケースの各パターンを検証する。
func TestNewsletterUsecase_GetSubscription(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockNewsletterRepository
		wantErr error
	}{
		{
			name: "正常: 購読情報が返る",
			repo: &mockNewsletterRepository{
				findByUserIDFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return newTestSubscription(), nil
				},
			},
			wantErr: nil,
		},
		{
			name: "未登録でErrNotFound",
			repo: &mockNewsletterRepository{
				findByUserIDFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, nil
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name: "DBエラー",
			repo: &mockNewsletterRepository{
				findByUserIDFunc: func(_ context.Context, _ int64) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNewsletterUsecase(tt.repo, nopLogger)
			sub, err := uc.GetSubscription(context.Background(), 1)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if errors.Is(tt.wantErr, ErrNotFound) && !errors.Is(err, ErrNotFound) {
					t.Errorf("ErrNotFoundが期待されたが、異なるエラー: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if sub == nil {
				t.Fatal("購読情報がnilだが、期待される値がある")
				return
			}
			if sub.Email != "test@example.com" {
				t.Errorf("Email = %q, want %q", sub.Email, "test@example.com")
			}
		})
	}
}

// TestNewsletterUsecase_Subscribe はSubscribeユースケースの各パターンを検証する。
func TestNewsletterUsecase_Subscribe(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		repo    *mockNewsletterRepository
		wantErr error
	}{
		{
			name:  "正常: 新規登録が成功",
			email: "user@example.com",
			repo: &mockNewsletterRepository{
				upsertFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := newTestSubscription()
					sub.Email = email
					return sub, nil
				},
			},
			wantErr: nil,
		},
		{
			name:  "正常: upsert（既存レコード更新）が成功",
			email: "updated@example.com",
			repo: &mockNewsletterRepository{
				upsertFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := newTestSubscription()
					sub.Email = email
					sub.IsActive = true
					return sub, nil
				},
			},
			wantErr: nil,
		},
		{
			name:    "メールバリデーション失敗: 空文字",
			email:   "",
			repo:    &mockNewsletterRepository{},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "メールバリデーション失敗: 不正な形式",
			email:   "not-an-email",
			repo:    &mockNewsletterRepository{},
			wantErr: ErrInvalidEmail,
		},
		{
			name:  "DBエラー",
			email: "user@example.com",
			repo: &mockNewsletterRepository{
				upsertFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNewsletterUsecase(tt.repo, nopLogger)
			sub, err := uc.Subscribe(context.Background(), 1, tt.email)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if errors.Is(tt.wantErr, ErrInvalidEmail) && !errors.Is(err, ErrInvalidEmail) {
					t.Errorf("ErrInvalidEmailが期待されたが、異なるエラー: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if sub == nil {
				t.Fatal("購読情報がnilだが、期待される値がある")
				return
			}
			if sub.Email != tt.email {
				t.Errorf("Email = %q, want %q", sub.Email, tt.email)
			}
		})
	}
}

// TestNewsletterUsecase_Unsubscribe はUnsubscribeユースケースの各パターンを検証する。
func TestNewsletterUsecase_Unsubscribe(t *testing.T) {
	tests := []struct {
		name    string
		repo    *mockNewsletterRepository
		wantErr error
	}{
		{
			name: "正常: is_active=falseに更新成功",
			repo: &mockNewsletterRepository{
				updateIsActiveFunc: func(_ context.Context, _ int64, isActive bool) (*domain.NewsletterSubscription, error) {
					if isActive {
						t.Error("is_activeがtrueで呼ばれてはいけない")
					}
					sub := newTestSubscription()
					sub.IsActive = false
					return sub, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "未登録でErrNotFound",
			repo: &mockNewsletterRepository{
				updateIsActiveFunc: func(_ context.Context, _ int64, _ bool) (*domain.NewsletterSubscription, error) {
					return nil, nil
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name: "DBエラー",
			repo: &mockNewsletterRepository{
				updateIsActiveFunc: func(_ context.Context, _ int64, _ bool) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNewsletterUsecase(tt.repo, nopLogger)
			sub, err := uc.Unsubscribe(context.Background(), 1)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if errors.Is(tt.wantErr, ErrNotFound) && !errors.Is(err, ErrNotFound) {
					t.Errorf("ErrNotFoundが期待されたが、異なるエラー: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if sub == nil {
				t.Fatal("購読情報がnilだが、期待される値がある")
				return
			}
			if sub.IsActive {
				t.Error("IsActive = true, want false")
			}
		})
	}
}

// TestNewsletterUsecase_UpdateEmail はUpdateEmailユースケースの各パターンを検証する。
func TestNewsletterUsecase_UpdateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		repo    *mockNewsletterRepository
		wantErr error
	}{
		{
			name:  "正常: メールアドレス更新成功",
			email: "new@example.com",
			repo: &mockNewsletterRepository{
				updateEmailFunc: func(_ context.Context, _ int64, email string) (*domain.NewsletterSubscription, error) {
					sub := newTestSubscription()
					sub.Email = email
					return sub, nil
				},
			},
			wantErr: nil,
		},
		{
			name:  "未登録でErrNotFound",
			email: "new@example.com",
			repo: &mockNewsletterRepository{
				updateEmailFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, nil
				},
			},
			wantErr: ErrNotFound,
		},
		{
			name:    "メールバリデーション失敗: 不正な形式",
			email:   "invalid",
			repo:    &mockNewsletterRepository{},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "メールバリデーション失敗: 空文字",
			email:   "",
			repo:    &mockNewsletterRepository{},
			wantErr: ErrInvalidEmail,
		},
		{
			name:  "DBエラー",
			email: "new@example.com",
			repo: &mockNewsletterRepository{
				updateEmailFunc: func(_ context.Context, _ int64, _ string) (*domain.NewsletterSubscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewNewsletterUsecase(tt.repo, nopLogger)
			sub, err := uc.UpdateEmail(context.Background(), 1, tt.email)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if errors.Is(tt.wantErr, ErrNotFound) && !errors.Is(err, ErrNotFound) {
					t.Errorf("ErrNotFoundが期待されたが、異なるエラー: %v", err)
				}
				if errors.Is(tt.wantErr, ErrInvalidEmail) && !errors.Is(err, ErrInvalidEmail) {
					t.Errorf("ErrInvalidEmailが期待されたが、異なるエラー: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if sub == nil {
				t.Fatal("購読情報がnilだが、期待される値がある")
				return
			}
			if sub.Email != tt.email {
				t.Errorf("Email = %q, want %q", sub.Email, tt.email)
			}
		})
	}
}
