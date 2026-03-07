package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// newsletterUsecase はNewsletterHandlerが依存するユースケースのインターフェース。
// テスト時にモック可能にするために定義する。
type newsletterUsecase interface {
	GetSubscription(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	Subscribe(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
	Unsubscribe(ctx context.Context, userID int64) (*domain.NewsletterSubscription, error)
	UpdateEmail(ctx context.Context, userID int64, email string) (*domain.NewsletterSubscription, error)
}

// NewsletterHandler はニュースレター購読APIハンドラーを提供する。
type NewsletterHandler struct {
	usecase newsletterUsecase
	logger  zerolog.Logger
}

// NewNewsletterHandler はNewsletterHandlerを作成する。
func NewNewsletterHandler(uc newsletterUsecase, logger zerolog.Logger) *NewsletterHandler {
	return &NewsletterHandler{
		usecase: uc,
		logger:  logger,
	}
}

// newsletterEmailRequest はメールアドレスを含むリクエストのJSON構造。
type newsletterEmailRequest struct {
	Email string `json:"email"`
}

// newsletterSubscriptionResponse はニュースレター購読情報のレスポンス構造。
type newsletterSubscriptionResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toNewsletterResponse はドメインモデルをレスポンス構造に変換する。
func toNewsletterResponse(sub *domain.NewsletterSubscription) newsletterSubscriptionResponse {
	return newsletterSubscriptionResponse{
		ID:        sub.ID,
		Email:     sub.Email,
		IsActive:  sub.IsActive,
		CreatedAt: sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sub.UpdatedAt.Format(time.RFC3339),
	}
}

// ServeGetSubscription はGET /api/v1/newsletter/subscription を処理し、購読状況を返す。
func (h *NewsletterHandler) ServeGetSubscription(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	sub, err := h.usecase.GetSubscription(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Code: codeNotFound, Message: "ニュースレター購読が登録されていません",
			})
			return
		}
		h.logger.Error().Err(err).Msg("購読情報取得失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, toNewsletterResponse(sub))
}

// ServeSubscribe はPOST /api/v1/newsletter/subscribe を処理し、購読を登録する。
func (h *NewsletterHandler) ServeSubscribe(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	var req newsletterEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "リクエストが不正です",
		})
		return
	}

	sub, err := h.usecase.Subscribe(r.Context(), user.ID, req.Email)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidEmail) {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Code: codeBadRequest, Message: "メールアドレスの形式が不正です",
			})
			return
		}
		h.logger.Error().Err(err).Msg("購読登録失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusCreated, toNewsletterResponse(sub))
}

// ServeUnsubscribe はPOST /api/v1/newsletter/unsubscribe を処理し、購読を解除する。
func (h *NewsletterHandler) ServeUnsubscribe(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	sub, err := h.usecase.Unsubscribe(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Code: codeNotFound, Message: "ニュースレター購読が登録されていません",
			})
			return
		}
		h.logger.Error().Err(err).Msg("購読解除失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, toNewsletterResponse(sub))
}

// ServeUpdateEmail はPUT /api/v1/newsletter/subscription を処理し、メールアドレスを変更する。
func (h *NewsletterHandler) ServeUpdateEmail(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	var req newsletterEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "リクエストが不正です",
		})
		return
	}

	sub, err := h.usecase.UpdateEmail(r.Context(), user.ID, req.Email)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidEmail) {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Code: codeBadRequest, Message: "メールアドレスの形式が不正です",
			})
			return
		}
		if errors.Is(err, usecase.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{
				Code: codeNotFound, Message: "ニュースレター購読が登録されていません",
			})
			return
		}
		h.logger.Error().Err(err).Msg("メールアドレス更新失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, toNewsletterResponse(sub))
}
