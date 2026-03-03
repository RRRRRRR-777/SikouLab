package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// UserRepository はユーザーのDB操作を提供する。
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository はUserRepositoryを作成する。
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByOAuth はOAuthプロバイダとユーザーIDでユーザーを検索する。
// 見つからない場合はnil, nilを返す（エラーではない）。
func (r *UserRepository) FindByOAuth(ctx context.Context, provider, oauthUserID string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, oauth_provider, oauth_user_id, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at FROM users WHERE oauth_provider = $1 AND oauth_user_id = $2`
	err := r.db.GetContext(ctx, &u, query, provider, oauthUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索失敗(OAuth): %w", err)
	}
	return &u, nil
}

// Create はユーザーとuser_settingsをトランザクションで作成する。
func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("トランザクション開始失敗: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// ユーザー作成
	userQuery := `
		INSERT INTO users (oauth_provider, oauth_user_id, name, display_name, avatar_url, role, subscription_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`
	var created domain.User
	err = tx.GetContext(ctx, &created, userQuery,
		user.OAuthProvider, user.OAuthUserID, user.Name, user.DisplayName, user.AvatarURL, user.Role, user.SubscriptionStatus)
	if err != nil {
		return nil, fmt.Errorf("ユーザー作成失敗: %w", err)
	}

	// user_settings作成（デフォルト値で初期化）
	settingsQuery := `INSERT INTO user_settings (user_id) VALUES ($1)`
	_, err = tx.ExecContext(ctx, settingsQuery, created.ID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー設定作成失敗: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("トランザクションコミット失敗: %w", err)
	}

	return &created, nil
}

// FindByID はIDでユーザーを検索する。
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, oauth_provider, oauth_user_id, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索失敗(ID): %w", err)
	}
	return &u, nil
}
