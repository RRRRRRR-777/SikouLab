package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/firebase"
)

// ErrInvalidToken は無効なIDトークンを表すエラー。
var ErrInvalidToken = errors.New("無効なIDトークンです")

// UserRepository はユーザーのDB操作インターフェース。
// テスト時にモック可能にするために定義する。
type UserRepository interface {
	FindByOAuth(ctx context.Context, provider, oauthUserID string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}

// AuthUsecase は認証ユースケースを提供する。
type AuthUsecase struct {
	firebaseClient firebase.TokenVerifier
	userRepo       UserRepository
}

// NewAuthUsecase はAuthUsecaseを作成する。
func NewAuthUsecase(firebaseClient firebase.TokenVerifier, userRepo UserRepository) *AuthUsecase {
	return &AuthUsecase{
		firebaseClient: firebaseClient,
		userRepo:       userRepo,
	}
}

// Login はID Tokenを検証し、ユーザーを取得または作成する。
// 初回ログイン時はisFirstLogin=trueを返す。
func (u *AuthUsecase) Login(ctx context.Context, idToken string) (user *domain.User, isFirstLogin bool, err error) {
	// IDトークンを検証
	ft, err := u.firebaseClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %s", ErrInvalidToken, err.Error())
	}

	// 既存ユーザーを検索
	existing, err := u.userRepo.FindByOAuth(ctx, ft.Provider, ft.UID)
	if err != nil {
		return nil, false, fmt.Errorf("ユーザー検索失敗: %w", err)
	}

	// 既存ユーザーが見つかった場合
	if existing != nil {
		return existing, false, nil
	}

	// 新規ユーザーを作成
	newUser := &domain.User{
		OAuthProvider:      ft.Provider,
		OAuthUserID:        ft.UID,
		Name:               ft.Name,
		DisplayName:        ft.Name,
		AvatarURL:          ft.Picture,
		Role:               "user",
		SubscriptionStatus: "inactive",
	}

	created, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, false, fmt.Errorf("ユーザー作成失敗: %w", err)
	}

	return created, true, nil
}

// CreateSessionCookie はID Tokenからサーバー側セッション Cookie を生成する。
func (u *AuthUsecase) CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	return u.firebaseClient.CreateSessionCookie(ctx, idToken, expiresIn)
}

// GetCurrentUser はセッション Cookie を検証してユーザーを返す。
func (u *AuthUsecase) GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, error) {
	ft, err := u.firebaseClient.VerifySessionCookie(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidToken, err.Error())
	}

	user, err := u.userRepo.FindByOAuth(ctx, ft.Provider, ft.UID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索失敗: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	return user, nil
}
