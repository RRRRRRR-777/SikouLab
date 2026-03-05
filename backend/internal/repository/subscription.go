package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// SubscriptionRepository はプランとサブスクリプション状態のDB操作を提供する。
type SubscriptionRepository struct {
	db *sqlx.DB
}

// NewSubscriptionRepository はSubscriptionRepositoryを作成する。
func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// FindActivePlans は is_active=true のプランを全件返す。
func (r *SubscriptionRepository) FindActivePlans(ctx context.Context) ([]domain.Plan, error) {
	var plans []domain.Plan
	query := `SELECT id, name, description, amount, currency, is_active, created_at, updated_at FROM plans WHERE is_active = true ORDER BY id`
	if err := r.db.SelectContext(ctx, &plans, query); err != nil {
		return nil, fmt.Errorf("アクティブプラン取得失敗: %w", err)
	}
	return plans, nil
}

// FindPlanByID はCheckout時にプランの存在と金額を検証するために使用する。
// 見つからない場合は nil, nil を返す（既存の FindByUnivaPaySubscriptionID と同じパターン）。
func (r *SubscriptionRepository) FindPlanByID(ctx context.Context, planID int64) (*domain.Plan, error) {
	var p domain.Plan
	query := `SELECT id, name, description, amount, currency, is_active, created_at, updated_at FROM plans WHERE id = $1`
	err := r.db.GetContext(ctx, &p, query, planID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("プラン取得失敗(ID=%d): %w", planID, err)
	}
	return &p, nil
}

// UpdateSubscriptionStatus はWebhookイベントに応じてユーザーの課金状態を反映するために使用する。
func (r *SubscriptionRepository) UpdateSubscriptionStatus(ctx context.Context, univapaySubscriptionID, status string) error {
	query := `UPDATE users SET subscription_status = $1, updated_at = NOW() WHERE univapay_customer_id = $2`
	result, err := r.db.ExecContext(ctx, query, status, univapaySubscriptionID)
	if err != nil {
		return fmt.Errorf("サブスクリプションステータス更新失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("影響行数の取得失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("該当するユーザーが見つかりません(univapay_customer_id=%s)", univapaySubscriptionID)
	}
	return nil
}

// FindByUnivaPaySubscriptionID は univapay_customer_id でユーザーを検索する。
// 見つからない場合は nil, nil を返す。
func (r *SubscriptionRepository) FindByUnivaPaySubscriptionID(ctx context.Context, univapaySubscriptionID string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, oauth_provider, oauth_user_id, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at FROM users WHERE univapay_customer_id = $1`
	err := r.db.GetContext(ctx, &u, query, univapaySubscriptionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索失敗(UnivaPaySubscriptionID): %w", err)
	}
	return &u, nil
}

// UpdateUnivaPaySubscriptionID は userID のユーザーに univapay_customer_id を保存する。
func (r *SubscriptionRepository) UpdateUnivaPaySubscriptionID(ctx context.Context, userID int64, subscriptionID string) error {
	query := `UPDATE users SET univapay_customer_id = $1, updated_at = NOW() WHERE id = $2`
	if _, err := r.db.ExecContext(ctx, query, subscriptionID, userID); err != nil {
		return fmt.Errorf("UnivaPayサブスクリプションID更新失敗: %w", err)
	}
	return nil
}
