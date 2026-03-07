package domain

import "time"

// NewsletterSubscription はニュースレター購読エンティティを表す。
//
// 1ユーザーにつき1レコードで、論理削除（is_active=false）により購読解除を管理する。
type NewsletterSubscription struct {
	// ID は購読レコードの一意識別子。
	ID int64 `db:"id" json:"id"`
	// UserID は購読ユーザーのID。
	UserID int64 `db:"user_id" json:"-"`
	// Email はニュースレター配信先メールアドレス。
	Email string `db:"email" json:"email"`
	// IsActive は購読状態（true=購読中、false=停止中）。
	IsActive bool `db:"is_active" json:"is_active"`
	// CreatedAt はレコード作成日時。
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt はレコード更新日時。
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
