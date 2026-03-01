package handler

import (
	"context"
	"crypto/hmac"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// subscriptionUsecase はSubscriptionHandlerが依存するユースケースのインターフェース。
// テスト時にモック可能にするために定義する。
type subscriptionUsecase interface {
	GetPlans(ctx context.Context) ([]domain.Plan, error)
	Checkout(ctx context.Context, user *domain.User, tokenID string) error
	HandleWebhook(ctx context.Context, payload usecase.WebhookPayload) error
}

// SubscriptionHandler はサブスクリプションAPIハンドラーを提供する。
type SubscriptionHandler struct {
	usecase       subscriptionUsecase
	logger        zerolog.Logger
	webhookSecret string
}

// NewSubscriptionHandler はSubscriptionHandlerを作成する。
func NewSubscriptionHandler(uc subscriptionUsecase, logger zerolog.Logger, webhookSecret string) *SubscriptionHandler {
	return &SubscriptionHandler{
		usecase:       uc,
		logger:        logger,
		webhookSecret: webhookSecret,
	}
}

// checkoutRequest はチェックアウトリクエストのJSON構造。
type checkoutRequest struct {
	TransactionTokenID string `json:"transaction_token_id"`
}

// ServeGetPlans はGET /api/v1/plans を処理し、アクティブなプラン一覧を返す。
func (h *SubscriptionHandler) ServeGetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.usecase.GetPlans(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("プラン一覧取得失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: "INTERNAL_ERROR", Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, plans)
}

// ServeCheckout はPOST /api/v1/univapay/checkout を処理し、UnivaPayサブスクリプションを作成する。
//
// 認証済みユーザーのcontextが必要。transaction_token_idが必須。
// 既にactiveの場合は409、その他エラーは500を返す。
func (h *SubscriptionHandler) ServeCheckout(w http.ResponseWriter, r *http.Request) {
	// contextからユーザーを取得（RequireAuth通過後を前提）
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{
			Code: "UNAUTHORIZED", Message: "認証が必要です",
		})
		return
	}

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: "BAD_REQUEST", Message: "リクエストが不正です",
		})
		return
	}

	// transaction_token_idが空の場合はバリデーションエラー
	if req.TransactionTokenID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: "BAD_REQUEST", Message: "transaction_token_idは必須です",
		})
		return
	}

	if err := h.usecase.Checkout(r.Context(), user, req.TransactionTokenID); err != nil {
		if errors.Is(err, usecase.ErrAlreadySubscribed) {
			writeJSON(w, http.StatusConflict, errorResponse{
				Code: "ALREADY_SUBSCRIBED", Message: "既にサブスクリプションが有効です",
			})
			return
		}
		h.logger.Error().Err(err).Msg("チェックアウト失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: "INTERNAL_ERROR", Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "pending"})
}

// ServeWebhook はPOST /api/v1/univapay/webhook を処理し、UnivaPayのWebhookイベントを受け付ける。
//
// AuthorizationヘッダーとUNIVAPAY_WEBHOOK_SECRETを定数時間比較で検証してからペイロードを処理する。
// ヘッダーなし・値不一致の場合は401を返す（タイミング攻撃対策のため定数時間比較を使用）。
func (h *SubscriptionHandler) ServeWebhook(w http.ResponseWriter, r *http.Request) {
	// Authorizationヘッダーの存在確認
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{
			Code: "UNAUTHORIZED", Message: "認証が必要です",
		})
		return
	}

	// 定数時間比較（タイミング攻撃対策）
	if !hmac.Equal([]byte(authHeader), []byte(h.webhookSecret)) {
		writeJSON(w, http.StatusUnauthorized, errorResponse{
			Code: "UNAUTHORIZED", Message: "認証が必要です",
		})
		return
	}

	// ペイロードをデコード
	var payload usecase.WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: "BAD_REQUEST", Message: "ペイロードのデコードに失敗しました",
		})
		return
	}

	// Webhookイベントを処理
	if err := h.usecase.HandleWebhook(r.Context(), payload); err != nil {
		h.logger.Error().Err(err).Msg("Webhookイベント処理失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: "INTERNAL_ERROR", Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
