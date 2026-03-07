package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// NewsletterRepository はニュースレター購読のDB操作を提供する。
type NewsletterRepository struct {
	db *sqlx.DB
}

// NewNewsletterRepository はNewsletterRepositoryを作成する。
func NewNewsletterRepository(db *sqlx.DB) *NewsletterRepository {
	return &NewsletterRepository{db: db}
}

// FindByUserID はユーザーIDで購読レコードを検索する。
// 見つからない場合は nil, nil を返す（既存の FindByUnivaPaySubscriptionID と同じパターン）。
func (r *NewsletterRepository) FindByUserID(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	var sub domain.NewsletterSubscription
	query := `SELECT id, user_id, email, is_active, created_at, updated_at FROM newsletter_subscriptions WHERE user_id = $1`
	err := r.db.GetContext(ctx, &sub, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("購読レコード検索失敗(user_id=%d): %w", userID, err)
	}
	return &sub, nil
}

// Upsert は購読レコードを作成または更新する。
// user_idのUNIQUE制約を利用し、既存レコードがあればemail更新+is_active=trueにする。
func (r *NewsletterRepository) Upsert(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	var sub domain.NewsletterSubscription
	query := `
		INSERT INTO newsletter_subscriptions (user_id, email, is_active, created_at, updated_at)
		VALUES ($1, $2, true, NOW(), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET email = EXCLUDED.email, is_active = true, updated_at = NOW()
		RETURNING id, user_id, email, is_active, created_at, updated_at`
	err := r.db.GetContext(ctx, &sub, query, userID, email)
	if err != nil {
		return nil, fmt.Errorf("購読レコードupsert失敗(user_id=%d): %w", userID, err)
	}
	return &sub, nil
}

// UpdateIsActive は購読解除・再開時に論理削除フラグを切り替えるために使用する。
// 見つからない場合は nil, nil を返す。
func (r *NewsletterRepository) UpdateIsActive(ctx context.Context, userID int64, isActive bool) (*domain.NewsletterSubscription, error) {
	var sub domain.NewsletterSubscription
	query := `
		UPDATE newsletter_subscriptions
		SET is_active = $1, updated_at = NOW()
		WHERE user_id = $2
		RETURNING id, user_id, email, is_active, created_at, updated_at`
	err := r.db.GetContext(ctx, &sub, query, isActive, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("購読状態更新失敗(user_id=%d): %w", userID, err)
	}
	return &sub, nil
}

// UpdateEmail はユーザーがニュースレターの配信先を変更する際に使用する。
// 見つからない場合は nil, nil を返す。
func (r *NewsletterRepository) UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	var sub domain.NewsletterSubscription
	query := `
		UPDATE newsletter_subscriptions
		SET email = $1, updated_at = NOW()
		WHERE user_id = $2
		RETURNING id, user_id, email, is_active, created_at, updated_at`
	err := r.db.GetContext(ctx, &sub, query, email, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("メールアドレス更新失敗(user_id=%d): %w", userID, err)
	}
	return &sub, nil
}
