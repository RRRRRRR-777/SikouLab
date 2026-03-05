package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/firebase"
)

// mockTokenVerifier はFirebase TokenVerifierのモック。
type mockTokenVerifier struct {
	verifyFunc              func(ctx context.Context, idToken string) (*firebase.FirebaseToken, error)
	createSessionCookieFunc func(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
	verifySessionCookieFunc func(ctx context.Context, sessionCookie string) (*firebase.FirebaseToken, error)
}

// VerifyIDToken はモックのIDトークン検証を実行する。
func (m *mockTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*firebase.FirebaseToken, error) {
	return m.verifyFunc(ctx, idToken)
}

// CreateSessionCookie はモックのセッションCookie生成を実行する。
func (m *mockTokenVerifier) CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	if m.createSessionCookieFunc != nil {
		return m.createSessionCookieFunc(ctx, idToken, expiresIn)
	}
	return "mock-session-cookie", nil
}

// VerifySessionCookie はモックのセッションCookie検証を実行する。
func (m *mockTokenVerifier) VerifySessionCookie(ctx context.Context, sessionCookie string) (*firebase.FirebaseToken, error) {
	if m.verifySessionCookieFunc != nil {
		return m.verifySessionCookieFunc(ctx, sessionCookie)
	}
	// デフォルトはVerifyIDTokenと同じ動作
	return m.verifyFunc(ctx, sessionCookie)
}

// mockUserRepo はUserRepositoryのモック。
type mockUserRepo struct {
	findByOAuthFunc func(ctx context.Context, provider, oauthUserID string) (*domain.User, error)
	createFunc      func(ctx context.Context, user *domain.User) (*domain.User, error)
	findByIDFunc    func(ctx context.Context, id int64) (*domain.User, error)
}

// FindByOAuth はモックのOAuth検索を実行する。
func (m *mockUserRepo) FindByOAuth(ctx context.Context, provider, oauthUserID string) (*domain.User, error) {
	return m.findByOAuthFunc(ctx, provider, oauthUserID)
}

// Create はモックのユーザー作成を実行する。
func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	return m.createFunc(ctx, user)
}

// FindByID はモックのID検索を実行する。
func (m *mockUserRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.findByIDFunc(ctx, id)
}

// TestAuthUsecase_Login はログインユースケースの各パターンを検証する。
func TestAuthUsecase_Login(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	validToken := &firebase.FirebaseToken{
		UID:      "firebase-uid-123",
		Email:    "test@example.com",
		Name:     "Test User",
		Picture:  "https://example.com/avatar.png",
		Provider: "google.com",
	}

	existingUser := &domain.User{
		ID:                 1,
		OAuthProvider:      "google.com",
		OAuthUserID:        "firebase-uid-123",
		Name:               "Test User",
		DisplayName:        "Test User",
		AvatarURL:          "https://example.com/avatar.png",
		Role:               "user",
		SubscriptionStatus: "active",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	tests := []struct {
		name           string
		idToken        string
		verifier       *mockTokenVerifier
		repo           *mockUserRepo
		wantFirstLogin bool
		wantErr        bool
		wantErrIs      error
	}{
		{
			name:    "有効なIDトークンで既存ユーザーがログインする",
			idToken: "valid-token",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return existingUser, nil
				},
			},
			wantFirstLogin: false,
			wantErr:        false,
		},
		{
			name:    "無効なIDトークンでエラーが返される",
			idToken: "invalid-token",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return nil, errors.New("token verification failed")
				},
			},
			repo:      &mockUserRepo{},
			wantErr:   true,
			wantErrIs: ErrInvalidToken,
		},
		{
			name:    "FindByOAuthがエラーを返す場合はエラーが返される",
			idToken: "valid-token",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name:    "Createがエラーを返す場合はエラーが返される",
			idToken: "valid-token",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return nil, nil
				},
				createFunc: func(_ context.Context, _ *domain.User) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name:    "初回ログインでユーザーが新規作成される",
			idToken: "valid-token",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return nil, nil
				},
				createFunc: func(_ context.Context, u *domain.User) (*domain.User, error) {
					return &domain.User{
						ID:                 2,
						OAuthProvider:      u.OAuthProvider,
						OAuthUserID:        u.OAuthUserID,
						Name:               u.Name,
						DisplayName:        u.DisplayName,
						AvatarURL:          u.AvatarURL,
						Role:               u.Role,
						SubscriptionStatus: u.SubscriptionStatus,
						CreatedAt:          now,
						UpdatedAt:          now,
					}, nil
				},
			},
			wantFirstLogin: true,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &AuthUsecase{
				firebaseClient: tt.verifier,
				userRepo:       tt.repo,
			}

			user, isFirstLogin, err := uc.Login(context.Background(), tt.idToken)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("エラー種別が一致しない: got %v, want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if user == nil {
				t.Fatal("ユーザーがnilで返された")
			}
			if isFirstLogin != tt.wantFirstLogin {
				t.Errorf("isFirstLogin = %v, want %v", isFirstLogin, tt.wantFirstLogin)
			}
		})
	}
}

// TestAuthUsecase_GetCurrentUser はセッション確認ユースケースの各パターンを検証する。
func TestAuthUsecase_GetCurrentUser(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	validToken := &firebase.FirebaseToken{
		UID:      "firebase-uid-123",
		Provider: "google.com",
	}

	existingUser := &domain.User{
		ID:                 1,
		OAuthProvider:      "google.com",
		OAuthUserID:        "firebase-uid-123",
		Name:               "Test User",
		Role:               "user",
		SubscriptionStatus: "active",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	tests := []struct {
		name         string
		sessionToken string
		verifier     *mockTokenVerifier
		repo         *mockUserRepo
		wantErr      bool
		wantErrIs    error
	}{
		{
			name:         "有効なセッショントークンでユーザー情報が返される",
			sessionToken: "valid-session",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return existingUser, nil
				},
			},
			wantErr: false,
		},
		{
			name:         "DBにユーザーが存在しない場合はErrInvalidTokenが返される",
			sessionToken: "valid-session",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return nil, nil
				},
			},
			wantErr:   true,
			wantErrIs: ErrInvalidToken,
		},
		{
			name:         "FindByOAuthがエラーを返す場合はエラーが返される",
			sessionToken: "valid-session",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return validToken, nil
				},
			},
			repo: &mockUserRepo{
				findByOAuthFunc: func(_ context.Context, _, _ string) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name:         "無効なセッショントークンでエラーが返される",
			sessionToken: "invalid-session",
			verifier: &mockTokenVerifier{
				verifyFunc: func(_ context.Context, _ string) (*firebase.FirebaseToken, error) {
					return nil, errors.New("token expired")
				},
			},
			repo:      &mockUserRepo{},
			wantErr:   true,
			wantErrIs: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &AuthUsecase{
				firebaseClient: tt.verifier,
				userRepo:       tt.repo,
			}

			user, err := uc.GetCurrentUser(context.Background(), tt.sessionToken)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("エラー種別が一致しない: got %v, want %v", err, tt.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if user == nil {
				t.Fatal("ユーザーがnilで返された")
			}
		})
	}
}
