package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// avatarMaxSize はアバター画像の最大サイズ（5MB）。
const avatarMaxSize = 5 << 20

// allowedImageTypes はアップロード可能な画像のContent-Type一覧。
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// userUsecase はUserHandlerが依存するユースケースのインターフェース。
// テスト時にモック可能にするために定義する。
type userUsecase interface {
	UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error)
	UploadAvatar(ctx context.Context, userID int64, fileData []byte, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, userID int64) error
}

// UserHandler はユーザープロフィールAPIハンドラーを提供する。
type UserHandler struct {
	usecase        userUsecase
	logger         zerolog.Logger
	storageBaseURL string
}

// NewUserHandler はUserHandlerを作成する。
func NewUserHandler(uc userUsecase, logger zerolog.Logger, storageBaseURL string) *UserHandler {
	return &UserHandler{
		usecase:        uc,
		logger:         logger,
		storageBaseURL: storageBaseURL,
	}
}

// updateProfileRequest は表示名更新リクエストのJSON構造。
type updateProfileRequest struct {
	DisplayName string `json:"display_name"`
}

// avatarResponse はアバターアップロードレスポンスのJSON構造。
type avatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

// ServeUpdateProfile はユーザーが自分の表示名を変更できるようにするためのエンドポイントを提供する。
//
// PATCH /api/v1/users/me を処理する。
// バリデーションはusecase層に委譲し、エラー種別に応じてHTTPステータスを返す。
func (h *UserHandler) ServeUpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "リクエストが不正です",
		})
		return
	}

	updated, err := h.usecase.UpdateDisplayName(r.Context(), user.ID, req.DisplayName)
	if err != nil {
		// バリデーションエラーの判定
		if errors.Is(err, usecase.ErrDisplayNameEmpty) ||
			errors.Is(err, usecase.ErrDisplayNameTooLong) ||
			errors.Is(err, usecase.ErrDisplayNameBlankOnly) {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Code: codeBadRequest, Message: err.Error(),
			})
			return
		}
		h.logger.Error().Err(err).Msg("表示名更新失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, struct {
		User userResponse `json:"user"`
	}{User: toUserResponse(updated, h.storageBaseURL)})
}

// ServeUploadAvatar は POST /api/v1/users/avatar を処理し、アバター画像をアップロードする。
//
// multipart/form-data で "image" フィールドの画像を受信する。
// バリデーション: JPEG/PNG/GIF形式、最大5MB。
func (h *UserHandler) ServeUploadAvatar(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	// リクエストボディのサイズ制限（multipartヘッダー・boundary等のオーバーヘッド分を加算）
	r.Body = http.MaxBytesReader(w, r.Body, avatarMaxSize+4096)

	// multipart/form-dataからファイルを取得
	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Warn().Err(err).Str("content_type", r.Header.Get("Content-Type")).Msg("FormFile取得失敗")
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "画像ファイルが必要です",
		})
		return
	}
	defer file.Close()

	// ファイルサイズチェック
	if header.Size > avatarMaxSize {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "画像サイズは5MB以下にしてください",
		})
		return
	}

	// ファイル内容を読み取り
	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "画像ファイルの読み込みに失敗しました",
		})
		return
	}

	// Content-Type判定: ヘッダー優先、不正ならファイル内容から検出
	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		contentType = http.DetectContentType(data)
	}
	if !allowedImageTypes[contentType] {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Code: codeBadRequest, Message: "対応形式はJPEG/PNG/GIFです",
		})
		return
	}

	avatarKey, err := h.usecase.UploadAvatar(r.Context(), user.ID, data, contentType)
	if err != nil {
		if errors.Is(err, usecase.ErrStorageNotConfigured) {
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Code: codeInternalError, Message: "ストレージが設定されていません",
			})
			return
		}
		h.logger.Error().Err(err).Msg("アバターアップロード失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	writeJSON(w, http.StatusOK, avatarResponse{AvatarURL: resolveStorageURL(h.storageBaseURL, &avatarKey)})
}

// ServeDeleteAvatar はユーザーがアバターをデフォルトに戻せるようにするためのエンドポイントを提供する。
//
// DELETE /api/v1/users/avatar を処理する。
// ストレージからの物理削除とDBのavatar_url NULL化をusecase層に委譲する。
func (h *UserHandler) ServeDeleteAvatar(w http.ResponseWriter, r *http.Request) {
	user := requireUser(w, r)
	if user == nil {
		return
	}

	if err := h.usecase.DeleteAvatar(r.Context(), user.ID); err != nil {
		h.logger.Error().Err(err).Msg("アバター削除失敗")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code: codeInternalError, Message: "サーバーエラーが発生しました",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
