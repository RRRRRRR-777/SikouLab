// Package repository はデータベースアクセスを提供する。
//
// PostgreSQLへの接続プール管理を担当する。
package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewDB はPostgreSQLへの接続プールを作成する。
// 接続確認のためPingを実行し、失敗した場合はエラーを返す。
func NewDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("DB接続オープン失敗: %w", err)
	}

	// 接続プール設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB接続確認失敗: %w", err)
	}

	return db, nil
}
