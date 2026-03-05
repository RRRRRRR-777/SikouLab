package domain

import "time"

// User はユーザーエンティティを表す。
//
// OAuthプロバイダ経由で認証されたユーザー情報を保持する。
// ロール（admin/writer/user）とプランに基づいて権限が決定される。
type User struct {
	// ID はユーザーの一意識別子。
	ID int64 `db:"id"`
	// OAuthProvider はOAuthプロバイダ識別子（"google.com"等）。
	OAuthProvider string `db:"oauth_provider"`
	// OAuthUserID はOAuthプロバイダ側のユーザーID。
	OAuthUserID string `db:"oauth_user_id"`
	// Name はユーザー名。
	Name string `db:"name"`
	// DisplayName は表示名。
	DisplayName string `db:"display_name"`
	// AvatarURL はアバター画像のURL。NULLの場合はnil。
	AvatarURL *string `db:"avatar_url"`
	// Role はユーザーのロール（"admin", "writer", "user"）。
	Role string `db:"role"`
	// PlanID は契約プランのID。
	PlanID *int64 `db:"plan_id"`
	// UnivaPayCustomerID はUnivaPay顧客ID。
	UnivaPayCustomerID *string `db:"univapay_customer_id"`
	// SubscriptionStatus はサブスクリプション状態（"active", "canceled"等）。
	SubscriptionStatus string `db:"subscription_status"`
	// CreatedAt はレコード作成日時。
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt はレコード更新日時。
	UpdatedAt time.Time `db:"updated_at"`
}

// UserSettings はユーザー設定エンティティを表す。
//
// サイドバーの展開状態など、UI設定を保持する。
type UserSettings struct {
	// ID は設定レコードの一意識別子。
	ID int64 `db:"id"`
	// UserID は対象ユーザーのID。
	UserID int64 `db:"user_id"`
	// SidebarArticleExpanded はサイドバーの記事セクションの展開状態。
	SidebarArticleExpanded bool `db:"sidebar_article_expanded"`
	// SidebarAdminExpanded はサイドバーの管理セクションの展開状態。
	SidebarAdminExpanded bool `db:"sidebar_admin_expanded"`
	// CreatedAt はレコード作成日時。
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt はレコード更新日時。
	UpdatedAt time.Time `db:"updated_at"`
}
