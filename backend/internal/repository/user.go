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
	query := `SELECT id, oauth_provider, oauth_user_id, email, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at FROM users WHERE oauth_provider = $1 AND oauth_user_id = $2`
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
		INSERT INTO users (oauth_provider, oauth_user_id, email, name, display_name, avatar_url, role, subscription_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *`
	var created domain.User
	err = tx.GetContext(ctx, &created, userQuery,
		user.OAuthProvider, user.OAuthUserID, user.Email, user.Name, user.DisplayName, user.AvatarURL, user.Role, user.SubscriptionStatus)
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
	query := `SELECT id, oauth_provider, oauth_user_id, email, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索失敗(ID): %w", err)
	}
	return &u, nil
}

// UpdateDisplayName は表示名を更新し、更新後のユーザーを返す。
func (r *UserRepository) UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error) {
	var u domain.User
	query := `UPDATE users SET display_name = $1, updated_at = NOW() WHERE id = $2 RETURNING id, oauth_provider, oauth_user_id, email, name, display_name, avatar_url, role, plan_id, univapay_customer_id, subscription_status, created_at, updated_at`
	err := r.db.GetContext(ctx, &u, query, displayName, userID)
	if err != nil {
		return nil, fmt.Errorf("表示名更新失敗: %w", err)
	}
	return &u, nil
}

// UpdateEmail はユーザーのメールアドレスを新しい値に置き換える。
func (r *UserRepository) UpdateEmail(ctx context.Context, userID int64, email string) error {
	query := `UPDATE users SET email = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, email, userID)
	if err != nil {
		return fmt.Errorf("メールアドレス更新失敗: %w", err)
	}
	return nil
}

// UpdateAvatarURL はアバター画像のURLを新しい値に置き換える。
func (r *UserRepository) UpdateAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	if err != nil {
		return fmt.Errorf("アバターURL更新失敗: %w", err)
	}
	return nil
}

// ClearAvatarURL はアバターURLをNULLに更新する。
func (r *UserRepository) ClearAvatarURL(ctx context.Context, userID int64) error {
	query := `UPDATE users SET avatar_url = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("アバターURL削除失敗: %w", err)
	}
	return nil
}
