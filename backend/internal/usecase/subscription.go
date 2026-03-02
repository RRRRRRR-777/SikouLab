package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/rs/zerolog"
)

// UnivaPay Webhook イベント名定数。
const (
	EventSubscriptionPayment  = "subscription_payment"
	EventSubscriptionFailure  = "subscription_failure"
	EventSubscriptionCanceled = "subscription_canceled"
)

// UnivaPay サブスクリプションステータス定数。
const (
	UnivaPayStatusCurrent   = "current"
	UnivaPayStatusUnpaid    = "unpaid"
	UnivaPayStatusUnconfirmed = "unconfirmed"
)

// DB保存用 サブスクリプションステータス定数。
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusPastDue   = "past_due"
	SubscriptionStatusCanceled  = "canceled"
	SubscriptionStatusTrialing  = "trialing"
)

// ErrAlreadySubscribed は既にアクティブなサブスクリプションを持つユーザーへのエラー。
var ErrAlreadySubscribed = errors.New("既にサブスクリプションが有効です")

// SubscriptionRepository はプランとsubscription_statusのDB操作インターフェース。
type SubscriptionRepository interface {
	// FindActivePlans はアクティブなプラン一覧を取得する。
	FindActivePlans(ctx context.Context) ([]domain.Plan, error)
	// FindPlanByID はプランIDでプランを取得する。
	FindPlanByID(ctx context.Context, planID int64) (*domain.Plan, error)
	// UpdateSubscriptionStatus はサブスクリプションステータスを更新する。
	UpdateSubscriptionStatus(ctx context.Context, univapaySubscriptionID, status string) error
	// FindByUnivaPaySubscriptionID はUnivaPayサブスクリプションIDでユーザーを検索する。
	FindByUnivaPaySubscriptionID(ctx context.Context, univapaySubscriptionID string) (*domain.User, error)
	// UpdateUnivaPaySubscriptionID はユーザーのUnivaPayサブスクリプションIDを更新する。
	UpdateUnivaPaySubscriptionID(ctx context.Context, userID int64, subscriptionID string) error
}

// UnivaPayClient はUnivaPay APIのインターフェース。
type UnivaPayClient interface {
	// CreateSubscription はUnivaPayでサブスクリプションを作成し、サブスクリプションIDとストアUUIDを返す。
	CreateSubscription(ctx context.Context, tokenID string, amount int, currency string) (string, string, error)
	// GetSubscription はUnivaPayサブスクリプションのステータスを取得する。storeIDはストアUUID。
	GetSubscription(ctx context.Context, storeID, subscriptionID string) (string, error)
}

// WebhookPayload はUnivaPay WebhookのペイロードJSON構造。
//
// UnivaPay の実際のペイロード形式に基づく:
//
//	{ "event": "subscription_payment", "data": { "id": "...", "status": "..." } }
type WebhookPayload struct {
	// Event はWebhookイベント名（例: "subscription_payment"）。
	Event string `json:"event"`
	// Data はWebhookのデータ本体（サブスクリプション情報が直接格納される）。
	Data WebhookData `json:"data"`
}

// WebhookData はWebhookペイロードのdataフィールド。
//
// UnivaPay はサブスクリプション情報を data 直下に格納する。
type WebhookData struct {
	// ID はUnivaPayサブスクリプションID。
	ID string `json:"id"`
	// Status はサブスクリプションのステータス（例: "current", "unconfirmed"）。
	Status string `json:"status"`
}

// SubscriptionUsecase はサブスクリプション機能のユースケースを提供する。
type SubscriptionUsecase struct {
	repo   SubscriptionRepository
	client UnivaPayClient
	logger zerolog.Logger
}

// NewSubscriptionUsecase はSubscriptionUsecaseを作成する。
func NewSubscriptionUsecase(repo SubscriptionRepository, client UnivaPayClient, logger zerolog.Logger) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

// GetPlans はアクティブなプラン一覧を取得する。
func (u *SubscriptionUsecase) GetPlans(ctx context.Context) ([]domain.Plan, error) {
	plans, err := u.repo.FindActivePlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("アクティブプラン取得失敗: %w", err)
	}
	return plans, nil
}

// Checkout はUnivaPayサブスクリプションIDをDBに保存する。
//
// ウィジェット（checkout: "payment"モード）がサブスクリプション作成と3Dセキュア認証を
// 一括処理するため、バックエンドはIDの保存のみを行う。
// ステータスの更新はWebhookで処理される。
//
// userのSubscriptionStatusが"active"の場合はErrAlreadySubscribedを返す。
func (u *SubscriptionUsecase) Checkout(ctx context.Context, user *domain.User, subscriptionID string) error {
	// 既にアクティブなサブスクリプションを持つユーザーは新規作成不可
	if user.SubscriptionStatus == "active" {
		return ErrAlreadySubscribed
	}

	// サブスクリプションIDをDBに保存
	if err := u.repo.UpdateUnivaPaySubscriptionID(ctx, user.ID, subscriptionID); err != nil {
		return fmt.Errorf("サブスクリプションID保存失敗: %w", err)
	}

	u.logger.Info().
		Int64("user_id", user.ID).
		Str("subscription_id", subscriptionID).
		Msg("[Checkout] サブスクリプションID保存完了")

	return nil
}

// HandleWebhook はUnivaPay Webhookイベントを処理してsubscription_statusを更新する。
//
// 未知イベントや対応ユーザー不在の場合はエラーなしで無視する。
// 冪等性: 同一イベントを複数回受信しても安全に処理する（DB側でUPDATE冪等）。
//
// イベントマッピング:
//   - subscription_payment + current → active
//   - subscription_payment + その他   → past_due
//   - subscription_failure           → past_due
//   - subscription_canceled          → canceled
//   - その他                           → 無視（エラーなし）
func (u *SubscriptionUsecase) HandleWebhook(ctx context.Context, payload WebhookPayload) error {
	u.logger.Debug().
		Str("event", payload.Event).
		Str("subscription_id", payload.Data.ID).
		Str("status", payload.Data.Status).
		Msg("[Webhook] 受信")

	// イベントと決済ステータスからsubscription_statusへのマッピングを決定
	newStatus, ok := resolveSubscriptionStatus(payload.Event, payload.Data.Status)
	if !ok {
		u.logger.Info().
			Str("event", payload.Event).
			Msg("[Webhook] 未知イベントを無視")
		return nil
	}

	subscriptionID := payload.Data.ID

	// サブスクリプションIDに対応するユーザーを検索
	user, err := u.repo.FindByUnivaPaySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}

	// 対応するユーザーが存在しない場合は無視する
	if user == nil {
		u.logger.Info().
			Str("subscription_id", subscriptionID).
			Msg("[Webhook] 対応するユーザーが見つからないため無視")
		return nil
	}

	u.logger.Info().
		Int64("user_id", user.ID).
		Str("new_status", newStatus).
		Msg("[Webhook] ステータス更新")

	// subscription_statusを更新
	if err := u.repo.UpdateSubscriptionStatus(ctx, subscriptionID, newStatus); err != nil {
		return fmt.Errorf("サブスクリプションステータス更新失敗: %w", err)
	}

	return nil
}

// resolveSubscriptionStatus はWebhookイベントと決済ステータスからDB保存用のステータスを解決する。
//
// 戻り値のboolがfalseの場合は未知イベントを示す。
func resolveSubscriptionStatus(event, subscriptionStatus string) (string, bool) {
	switch event {
	case EventSubscriptionPayment:
		switch subscriptionStatus {
		case UnivaPayStatusCurrent:
			return SubscriptionStatusActive, true
		default:
			return SubscriptionStatusPastDue, true
		}
	case EventSubscriptionFailure:
		return SubscriptionStatusPastDue, true
	case EventSubscriptionCanceled:
		return SubscriptionStatusCanceled, true
	default:
		return "", false
	}
}
