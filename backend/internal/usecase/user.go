package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// displayNameMaxLen は表示名の最大文字数。
const displayNameMaxLen = 50

// ErrDisplayNameEmpty は表示名が空の場合のエラー。
var ErrDisplayNameEmpty = errors.New("表示名は必須です")

// ErrDisplayNameTooLong は表示名が上限を超えた場合のエラー。
var ErrDisplayNameTooLong = errors.New("表示名は50文字以内で入力してください")

// ErrDisplayNameBlankOnly は表示名が空白のみの場合のエラー。
var ErrDisplayNameBlankOnly = errors.New("表示名に空白のみは指定できません")

// ErrStorageNotConfigured はオブジェクトストレージが未設定の場合のエラー。
var ErrStorageNotConfigured = errors.New("ストレージが設定されていません")

// UserProfileRepository はプロフィール更新に必要なDB操作インターフェース。
type UserProfileRepository interface {
	// UpdateDisplayName は表示名を更新し、更新後のユーザーを返す。
	UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error)
	// UpdateAvatarURL はアバターURLを更新する。
	UpdateAvatarURL(ctx context.Context, userID int64, avatarURL string) error
	// ClearAvatarURL はアバターURLをNULLに更新する。
	ClearAvatarURL(ctx context.Context, userID int64) error
	// FindByID はIDでユーザーを検索する。
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}

// ObjectStorage はオブジェクトの保存・削除インターフェース。
type ObjectStorage interface {
	// Save はオブジェクトを保存し、オブジェクトキーを返す。
	Save(ctx context.Context, prefix string, id int64, data []byte, ext string) (string, error)
	// Delete はオブジェクトを削除する。存在しない場合はエラーなしで完了する。
	Delete(ctx context.Context, key string) error
}

// UserUsecase はユーザープロフィール機能のユースケースを提供する。
type UserUsecase struct {
	repo    UserProfileRepository
	storage ObjectStorage
}

// NewUserUsecase はUserUsecaseを作成する。
func NewUserUsecase(repo UserProfileRepository, storage ObjectStorage) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		storage: storage,
	}
}

// UpdateDisplayName はユーザーが他のユーザーから識別されるための表示名を安全に更新する。
//
// 不正な入力からデータ品質を守るため、以下のバリデーションを行う:
//   - 空文字は不可（ErrDisplayNameEmpty）
//   - 51文字以上は不可（ErrDisplayNameTooLong）
//   - 空白のみは不可（ErrDisplayNameBlankOnly）。
func (u *UserUsecase) UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error) {
	// バリデーション
	if displayName == "" {
		return nil, ErrDisplayNameEmpty
	}
	if utf8.RuneCountInString(displayName) > displayNameMaxLen {
		return nil, ErrDisplayNameTooLong
	}
	if strings.TrimSpace(displayName) == "" {
		return nil, ErrDisplayNameBlankOnly
	}

	user, err := u.repo.UpdateDisplayName(ctx, userID, displayName)
	if err != nil {
		return nil, fmt.Errorf("表示名更新失敗: %w", err)
	}
	return user, nil
}

// contentTypeToExt はContent-Typeから拡張子を解決する。
func contentTypeToExt(contentType string) (string, bool) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/gif":
		return ".gif", true
	default:
		return "", false
	}
}

// UploadAvatar はユーザーのプロフィールを個性化するためにアバター画像を保存しURLを永続化する。
//
// ストレージにゴミファイルが残らないよう、既存アバターがあれば先に削除してから新画像を保存する。
// Content-Typeに基づいて拡張子を決定する。
// storageがnilの場合はErrStorageNotConfiguredを返す。
func (u *UserUsecase) UploadAvatar(ctx context.Context, userID int64, fileData []byte, contentType string) (string, error) {
	if u.storage == nil {
		return "", ErrStorageNotConfigured
	}

	ext, ok := contentTypeToExt(contentType)
	if !ok {
		return "", fmt.Errorf("未対応の画像形式です: %s", contentType)
	}

	// 既存アバターを取得して削除
	existing, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("ユーザー取得失敗: %w", err)
	}
	if existing != nil && existing.AvatarURL != nil && *existing.AvatarURL != "" {
		// 古いファイルの削除失敗は無視する（ゴミファイルとして残る）
		_ = u.storage.Delete(ctx, *existing.AvatarURL)
	}

	// ストレージに保存
	key, err := u.storage.Save(ctx, "users", userID, fileData, ext)
	if err != nil {
		return "", fmt.Errorf("アバター保存失敗: %w", err)
	}

	// DBのavatar_urlをオブジェクトキーとして更新
	if err := u.repo.UpdateAvatarURL(ctx, userID, key); err != nil {
		return "", fmt.Errorf("アバターURL更新失敗: %w", err)
	}

	return key, nil
}

// DeleteAvatar はアバター画像を削除し、DBのavatar_urlをNULLにする。
//
// ファイルが存在しない場合もDB更新は行う。
// storageがnilの場合はストレージ削除をスキップし、DB更新のみ行う。
func (u *UserUsecase) DeleteAvatar(ctx context.Context, userID int64) error {
	// 現在のアバターURLを取得
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("ユーザー取得失敗: %w", err)
	}

	// ファイルが存在する場合は削除
	if u.storage != nil && user != nil && user.AvatarURL != nil && *user.AvatarURL != "" {
		_ = u.storage.Delete(ctx, *user.AvatarURL)
	}

	// DBのavatar_urlをNULLに更新
	if err := u.repo.ClearAvatarURL(ctx, userID); err != nil {
		return fmt.Errorf("アバターURL削除失敗: %w", err)
	}

	return nil
}
