package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"

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

// ポータルURLのデフォルト値。
const defaultPortalURL = "https://widget.univapay.com/portal"

// SubscriptionInfo は認証済みユーザーのサブスクリプション状態をハンドラーに返すための構造体。
type SubscriptionInfo struct {
	// PlanName はプラン名（サブスク未登録の場合は空文字）。
	PlanName string
	// Amount は月額料金（最小通貨単位、サブスク未登録の場合は0）。
	Amount int
	// Currency は通貨コード ISO-4217（サブスク未登録の場合は空文字）。
	Currency string
	// Status はサブスクリプション状態。
	Status string
}

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

// WebhookPayload はUnivaPay Webhookのリクエストボディをデシリアライズするための構造体。
//
// UnivaPay の実際のペイロード形式に基づく。
// 例: { "event": "subscription_payment", "data": { "id": "...", "status": "..." } }。
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
	repo          SubscriptionRepository
	client        UnivaPayClient
	logger        zerolog.Logger
	portalBaseURL string
}

// NewSubscriptionUsecase はSubscriptionUsecaseを作成する。
// portalBaseURLが空の場合はデフォルトのUnivaPayポータルURLを使用する。
func NewSubscriptionUsecase(repo SubscriptionRepository, client UnivaPayClient, logger zerolog.Logger, portalBaseURL ...string) *SubscriptionUsecase {
	u := &SubscriptionUsecase{
		repo:   repo,
		client: client,
		logger: logger,
	}
	if len(portalBaseURL) > 0 {
		u.portalBaseURL = portalBaseURL[0]
	}
	return u
}

// GetPlans はフロントエンドのプラン選択画面に表示するプラン一覧を提供する。
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
	if user.SubscriptionStatus == SubscriptionStatusActive {
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

// HandleWebhook はUnivaPayからの決済通知に基づいてユーザーの課金状態を同期する。
//
// 未知イベントや対応ユーザー不在の場合はエラーなしで無視する。
// 冪等性を保証し、同一イベントを複数回受信しても安全に処理する。
//
// イベントマッピング:
//   - subscription_payment + current → active
//   - subscription_payment + その他 → past_due
//   - subscription_failure → past_due
//   - subscription_canceled → canceled
//   - その他 → 無視。
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

// GetMySubscription は認証済みユーザーのサブスクリプション状態を返す。
//
// plan_idがNULL（サブスク未登録）の場合はゼロ値のSubscriptionInfoを返し、エラーにはしない。
func (u *SubscriptionUsecase) GetMySubscription(ctx context.Context, user *domain.User) (*SubscriptionInfo, error) {
	info := &SubscriptionInfo{
		Status: user.SubscriptionStatus,
	}

	if user.PlanID == nil {
		return info, nil
	}

	plan, err := u.repo.FindPlanByID(ctx, *user.PlanID)
	if err != nil {
		return nil, fmt.Errorf("プラン情報取得失敗: %w", err)
	}

	if plan != nil {
		info.PlanName = plan.Name
		info.Amount = plan.Amount
		info.Currency = plan.Currency
	}

	return info, nil
}

// GeneratePortalURL はUnivaPayカスタマーポータルURLを生成する。
//
// univapay_customer_idがNULL（サブスク未登録）の場合はErrNotFoundを返す。
func (u *SubscriptionUsecase) GeneratePortalURL(_ context.Context, user *domain.User) (string, error) {
	if user.UnivaPayCustomerID == nil {
		return "", ErrNotFound
	}

	baseURL := u.portalBaseURL
	if baseURL == "" {
		baseURL = defaultPortalURL
	}

	portalURL := baseURL + "?customer=" + url.QueryEscape(*user.UnivaPayCustomerID)

	u.logger.Info().
		Int64("user_id", user.ID).
		Msg("[Portal] ポータルURL生成")

	return portalURL, nil
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
