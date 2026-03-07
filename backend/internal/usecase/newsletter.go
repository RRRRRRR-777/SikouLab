package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/rs/zerolog"
)

// ErrNotFound は対象リソースが存在しない場合のエラー。
var ErrNotFound = errors.New("リソースが見つかりません")

// ErrInvalidEmail はメールアドレスの形式が不正な場合のエラー。
var ErrInvalidEmail = errors.New("メールアドレスの形式が不正です")

// NewsletterRepository はニュースレター購読のDB操作インターフェース。
// テスト時にモック可能にするために定義する。
type NewsletterRepository interface {
	// FindByUserID はユーザーIDで購読レコードを検索する。見つからない場合は nil, nil を返す。
	FindByUserID(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	// Upsert は購読レコードを作成または更新する。
	Upsert(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
	// UpdateIsActive は購読状態を更新する。見つからない場合は nil, nil を返す。
	UpdateIsActive(ctx context.Context, userID int64, isActive bool) (*domain.NewsletterSubscription, error)
	// UpdateEmail はメールアドレスを更新する。見つからない場合は nil, nil を返す。
	UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
}

// NewsletterUsecase はニュースレター購読機能のユースケースを提供する。
type NewsletterUsecase struct {
	repo   NewsletterRepository
	logger zerolog.Logger
}

// NewNewsletterUsecase はNewsletterUsecaseを作成する。
func NewNewsletterUsecase(repo NewsletterRepository, logger zerolog.Logger) *NewsletterUsecase {
	return &NewsletterUsecase{
		repo:   repo,
		logger: logger,
	}
}

// validateEmail はメールアドレスの形式を検証する。
// Go標準のnet/mail.ParseAddressでRFC準拠のバリデーションを行う。
func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// GetSubscription は設定画面でユーザーの現在の購読状態を表示するために使用する。
// 未登録の場合はErrNotFoundを返す。
func (u *NewsletterUsecase) GetSubscription(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	sub, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("購読情報取得失敗: %w", err)
	}
	if sub == nil {
		return nil, ErrNotFound
	}
	return sub, nil
}

// Subscribe はニュースレターの購読を登録する。
// 既存レコードがある場合はupsertでemail更新+is_active=trueにする。
func (u *NewsletterUsecase) Subscribe(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	sub, err := u.repo.Upsert(ctx, userID, email)
	if err != nil {
		return nil, fmt.Errorf("購読登録失敗: %w", err)
	}

	u.logger.Info().
		Int64("user_id", userID).
		Msg("[Newsletter] 購読登録完了")

	return sub, nil
}

// Unsubscribe はニュースレターの購読を解除する（論理削除: is_active=false）。
// 未登録の場合はErrNotFoundを返す。
func (u *NewsletterUsecase) Unsubscribe(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error) {
	sub, err := u.repo.UpdateIsActive(ctx, userID, false)
	if err != nil {
		return nil, fmt.Errorf("購読解除失敗: %w", err)
	}
	if sub == nil {
		return nil, ErrNotFound
	}

	u.logger.Info().
		Int64("user_id", userID).
		Msg("[Newsletter] 購読解除完了")

	return sub, nil
}

// UpdateEmail はニュースレターのメールアドレスを変更する。
// 未登録の場合はErrNotFoundを返す。
func (u *NewsletterUsecase) UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	sub, err := u.repo.UpdateEmail(ctx, userID, email)
	if err != nil {
		return nil, fmt.Errorf("メールアドレス更新失敗: %w", err)
	}
	if sub == nil {
		return nil, ErrNotFound
	}

	u.logger.Info().
		Int64("user_id", userID).
		Msg("[Newsletter] メールアドレス変更完了")

	return sub, nil
}
