package domain

import "time"

// Plan はサブスクリプションプランエンティティを表す。
//
// プランの価格・通貨・有効状態を保持し、ユーザーのサブスクリプション契約と紐付く。
type Plan struct {
	// ID はプランの一意識別子。
	ID int64 `db:"id"`
	// Name はプラン名。
	Name string `db:"name"`
	// Description はプランの説明。
	Description *string `db:"description"`
	// Amount は月額料金（最小通貨単位）。
	Amount int `db:"amount"`
	// Currency は通貨コード ISO-4217（例: JPY）。
	Currency string `db:"currency"`
	// IsActive はプランが有効かどうか。
	IsActive bool `db:"is_active"`
	// CreatedAt はレコード作成日時。
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt はレコード更新日時。
	UpdatedAt time.Time `db:"updated_at"`
}
