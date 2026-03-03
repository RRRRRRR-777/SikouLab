// Package firebase はFirebase Admin SDKとの連携を提供する。
//
// IDトークンの検証を通じて、OAuthプロバイダ経由の認証を処理する。
package firebase

import (
	"context"
	"fmt"
	"time"

	firebaseSDK "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseToken はFirebase ID Tokenを検証した結果を表す。
type FirebaseToken struct {
	// UID はFirebaseユーザーの一意識別子。
	UID string
	// Email はユーザーのメールアドレス。
	Email string
	// Name はユーザーの表示名。
	Name string
	// Picture はユーザーのアバターURL。
	Picture string
	// Provider はOAuthプロバイダ識別子（"google.com", "apple.com"等）。
	Provider string
}

// TokenVerifier はIDトークン検証のインターフェース。
// テスト時にモック可能にするために定義する。
type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error)
	// CreateSessionCookie はID Tokenからセッション Cookie を生成する。
	// expiresIn はセッションの有効期限（最大14日）。
	CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
	// VerifySessionCookie はセッション Cookie を検証し、ユーザー情報を返す。
	VerifySessionCookie(ctx context.Context, sessionCookie string) (*FirebaseToken, error)
}

// Client はFirebase Admin SDKクライアントを表す。
type Client struct {
	authClient *auth.Client
}

// NewClient はFirebase Admin SDKクライアントを初期化する。
// serviceAccountJSON が空の場合はApplication Default Credentials (ADC) を使用する。
// projectID はIDトークン検証時のaudience確認に必要。
func NewClient(ctx context.Context, projectID, serviceAccountJSON string) (*Client, error) {
	var app *firebaseSDK.App
	var err error

	appConfig := &firebaseSDK.Config{ProjectID: projectID}

	if serviceAccountJSON != "" {
		opt := option.WithCredentialsJSON([]byte(serviceAccountJSON))
		app, err = firebaseSDK.NewApp(ctx, appConfig, opt)
	} else {
		app, err = firebaseSDK.NewApp(ctx, appConfig)
	}
	if err != nil {
		return nil, fmt.Errorf("Firebase App初期化失敗: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("Firebase Auth Client初期化失敗: %w", err)
	}

	return &Client{authClient: authClient}, nil
}

// CreateSessionCookie はFirebase ID Tokenからサーバー側セッション Cookie を生成する。
// Firebase Admin SDKが生成するセッション Cookie は指定した有効期限まで有効（最大14日）。
func (c *Client) CreateSessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	cookie, err := c.authClient.SessionCookie(ctx, idToken, expiresIn)
	if err != nil {
		return "", fmt.Errorf("セッションCookie生成失敗: %w", err)
	}
	return cookie, nil
}

// VerifySessionCookie はFirebaseセッション Cookie を検証し、ユーザー情報を返す。
func (c *Client) VerifySessionCookie(ctx context.Context, sessionCookie string) (*FirebaseToken, error) {
	token, err := c.authClient.VerifySessionCookieAndCheckRevoked(ctx, sessionCookie)
	if err != nil {
		return nil, fmt.Errorf("セッションCookie検証失敗: %w", err)
	}
	return c.tokenToFirebaseToken(token), nil
}

// VerifyIDToken はFirebase ID Tokenを検証し、ユーザー情報を返す。
func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*FirebaseToken, error) {
	token, err := c.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("IDトークン検証失敗: %w", err)
	}
	return c.tokenToFirebaseToken(token), nil
}

// tokenToFirebaseToken は auth.Token からFirebaseTokenに変換する共通ヘルパー。
func (c *Client) tokenToFirebaseToken(token *auth.Token) *FirebaseToken {
	ft := &FirebaseToken{
		UID: token.UID,
	}

	// Firebase tokensのclaimsからユーザー情報を取得
	if email, ok := token.Claims["email"].(string); ok {
		ft.Email = email
	}
	if name, ok := token.Claims["name"].(string); ok {
		ft.Name = name
	}
	if picture, ok := token.Claims["picture"].(string); ok {
		ft.Picture = picture
	}

	// サインインプロバイダを取得
	if firebase, ok := token.Claims["firebase"].(map[string]interface{}); ok {
		if provider, ok := firebase["sign_in_provider"].(string); ok {
			ft.Provider = provider
		}
	}

	return ft
}
