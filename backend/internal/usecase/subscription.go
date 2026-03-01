package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// ErrAlreadySubscribed は既にアクティブなサブスクリプションを持つユーザーへのエラー。
var ErrAlreadySubscribed = errors.New("既にサブスクリプションが有効です")

// SubscriptionRepository はプランとsubscription_statusのDB操作インターフェース。
type SubscriptionRepository interface {
	// FindActivePlans はアクティブなプラン一覧を取得する。
	FindActivePlans(ctx context.Context) ([]domain.Plan, error)
	// UpdateSubscriptionStatus はサブスクリプションステータスを更新する。
	UpdateSubscriptionStatus(ctx context.Context, univapaySubscriptionID, status string) error
	// FindByUnivaPaySubscriptionID はUnivaPayサブスクリプションIDでユーザーを検索する。
	FindByUnivaPaySubscriptionID(ctx context.Context, univapaySubscriptionID string) (*domain.User, error)
	// UpdateUnivaPaySubscriptionID はユーザーのUnivaPayサブスクリプションIDを更新する。
	UpdateUnivaPaySubscriptionID(ctx context.Context, userID int64, subscriptionID string) error
}

// UnivaPayClient はUnivaPay APIのインターフェース。
type UnivaPayClient interface {
	// CreateSubscription はUnivaPayでサブスクリプションを作成し、サブスクリプションIDを返す。
	CreateSubscription(ctx context.Context, tokenID string, planID int64) (string, error)
}

// WebhookPayload はUnivaPay WebhookのペイロードJSON構造。
type WebhookPayload struct {
	// Event はWebhookイベント名（例: "SUBSCRIPTION_PAYMENT"）。
	Event string `json:"event"`
	// Data はWebhookのデータ本体。
	Data WebhookData `json:"data"`
}

// WebhookData はWebhookペイロードのdataフィールド。
type WebhookData struct {
	// Subscriptions はサブスクリプション情報。
	Subscriptions WebhookSubscription `json:"subscriptions"`
}

// WebhookSubscription はWebhookペイロード内のサブスクリプション情報。
type WebhookSubscription struct {
	// ID はUnivaPayサブスクリプションID。
	ID string `json:"id"`
	// Status はサブスクリプションのステータス（例: "successful", "failed"）。
	Status string `json:"status"`
}

// SubscriptionUsecase はサブスクリプション機能のユースケースを提供する。
type SubscriptionUsecase struct {
	repo   SubscriptionRepository
	client UnivaPayClient
}

// NewSubscriptionUsecase はSubscriptionUsecaseを作成する。
func NewSubscriptionUsecase(repo SubscriptionRepository, client UnivaPayClient) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		repo:   repo,
		client: client,
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

// Checkout はUnivaPayサブスクリプションを作成し、サブスクリプションIDをDBに保存する。
//
// userのSubscriptionStatusが"active"の場合はErrAlreadySubscribedを返す。
// UnivaPay APIエラー時はDB更新を行わない。
func (u *SubscriptionUsecase) Checkout(ctx context.Context, user *domain.User, tokenID string) error {
	// 既にアクティブなサブスクリプションを持つユーザーは新規作成不可
	if user.SubscriptionStatus == "active" {
		return ErrAlreadySubscribed
	}

	// UnivaPayでサブスクリプションを作成（planIDは現フェーズでは1固定）
	subscriptionID, err := u.client.CreateSubscription(ctx, tokenID, 1)
	if err != nil {
		return fmt.Errorf("UnivaPayサブスクリプション作成失敗: %w", err)
	}

	// サブスクリプションIDをDBに保存
	if err := u.repo.UpdateUnivaPaySubscriptionID(ctx, user.ID, subscriptionID); err != nil {
		return fmt.Errorf("サブスクリプションID保存失敗: %w", err)
	}

	return nil
}

// HandleWebhook はUnivaPay Webhookイベントを処理してsubscription_statusを更新する。
//
// 未知イベントや対応ユーザー不在の場合はエラーなしで無視する。
// 冪等性: 同一イベントを複数回受信しても安全に処理する（DB側でUPDATE冪等）。
//
// イベントマッピング:
//   - SUBSCRIPTION_PAYMENT + successful → active
//   - SUBSCRIPTION_PAYMENT + failed     → past_due
//   - SUBSCRIPTION_FAILED               → past_due
//   - SUBSCRIPTION_CANCELED             → canceled
//   - その他                              → 無視（エラーなし）
func (u *SubscriptionUsecase) HandleWebhook(ctx context.Context, payload WebhookPayload) error {
	// イベントと決済ステータスからsubscription_statusへのマッピングを決定
	newStatus, ok := resolveSubscriptionStatus(payload.Event, payload.Data.Subscriptions.Status)
	if !ok {
		// 未知イベントは無視する
		return nil
	}

	subscriptionID := payload.Data.Subscriptions.ID

	// サブスクリプションIDに対応するユーザーを検索
	user, err := u.repo.FindByUnivaPaySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("ユーザー検索失敗: %w", err)
	}

	// 対応するユーザーが存在しない場合は無視する
	if user == nil {
		return nil
	}

	// subscription_statusを更新
	if err := u.repo.UpdateSubscriptionStatus(ctx, subscriptionID, newStatus); err != nil {
		return fmt.Errorf("サブスクリプションステータス更新失敗: %w", err)
	}

	return nil
}

// resolveSubscriptionStatus はWebhookイベントと決済ステータスからDB保存用のステータスを解決する。
//
// 戻り値のboolがfalseの場合は未知イベントを示す。
func resolveSubscriptionStatus(event, paymentStatus string) (string, bool) {
	switch event {
	case "SUBSCRIPTION_PAYMENT":
		switch paymentStatus {
		case "successful":
			return "active", true
		case "failed":
			return "past_due", true
		default:
			return "past_due", true
		}
	case "SUBSCRIPTION_FAILED":
		return "past_due", true
	case "SUBSCRIPTION_CANCELED":
		return "canceled", true
	default:
		return "", false
	}
}
