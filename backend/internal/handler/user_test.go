package handler

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/middleware"
	"github.com/RRRRRRR-777/SicouLab/backend/internal/usecase"
)

// mockUserUsecase はuserUsecaseインターフェースのモック。
type mockUserUsecase struct {
	updateDisplayNameFunc func(ctx context.Context, userID int64, displayName string) (*domain.User, error)
	uploadAvatarFunc      func(ctx context.Context, userID int64, fileData []byte, contentType string) (string, error)
	deleteAvatarFunc      func(ctx context.Context, userID int64) error
}

// UpdateDisplayName はモックの表示名更新処理を実行する。
func (m *mockUserUsecase) UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error) {
	return m.updateDisplayNameFunc(ctx, userID, displayName)
}

// UploadAvatar はモックのアバターアップロード処理を実行する。
func (m *mockUserUsecase) UploadAvatar(ctx context.Context, userID int64, fileData []byte, contentType string) (string, error) {
	return m.uploadAvatarFunc(ctx, userID, fileData, contentType)
}

// DeleteAvatar はモックのアバター削除処理を実行する。
func (m *mockUserUsecase) DeleteAvatar(ctx context.Context, userID int64) error {
	return m.deleteAvatarFunc(ctx, userID)
}

// TestUserHandler_ServeUpdateProfile は PATCH /api/v1/users/me の各パターンを検証する。
func TestUserHandler_ServeUpdateProfile(t *testing.T) {
	validUser := &domain.User{
		ID:   1,
		Name: "Test User",
		Role: "user",
	}

	updatedUser := &domain.User{
		ID:          1,
		Name:        "Test User",
		DisplayName: "山田 太郎",
		Role:        "user",
	}

	tests := []struct {
		name       string
		body       string
		injectUser *domain.User
		uc         *mockUserUsecase
		wantStatus int
		wantBody   string
	}{
		{
			name:       "正常な表示名で200が返る",
			body:       `{"display_name":"山田 太郎"}`,
			injectUser: validUser,
			uc: &mockUserUsecase{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return updatedUser, nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"display_name":"山田 太郎"`,
		},
		{
			name:       "空文字で400が返る",
			body:       `{"display_name":""}`,
			injectUser: validUser,
			uc: &mockUserUsecase{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return nil, usecase.ErrDisplayNameEmpty
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "51文字以上で400が返る",
			body:       `{"display_name":"` + strings.Repeat("あ", 51) + `"}`,
			injectUser: validUser,
			uc: &mockUserUsecase{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return nil, usecase.ErrDisplayNameTooLong
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "空白のみで400が返る",
			body:       `{"display_name":"   "}`,
			injectUser: validUser,
			uc: &mockUserUsecase{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return nil, usecase.ErrDisplayNameBlankOnly
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "未認証で401が返る",
			body:       `{"display_name":"山田 太郎"}`,
			injectUser: nil,
			uc:         &mockUserUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "usecaseエラーで500が返る",
			body:       `{"display_name":"山田 太郎"}`,
			injectUser: validUser,
			uc: &mockUserUsecase{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(tt.uc, zerolog.Nop(), "http://test-storage")
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeUpdateProfile(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if !strings.Contains(rec.Body.String(), tt.wantBody) {
					t.Errorf("レスポンスボディに %q が含まれない: got %q", tt.wantBody, rec.Body.String())
				}
			}
		})
	}
}

// TestUserHandler_ServeUploadAvatar は POST /api/v1/users/avatar の各パターンを検証する。
func TestUserHandler_ServeUploadAvatar(t *testing.T) {
	validUser := &domain.User{
		ID:   1,
		Name: "Test User",
		Role: "user",
	}

	// テスト用の有効なJPEGヘッダー（最小限のJFIFシグネチャ）
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	// テスト用のPNGヘッダー（8バイトのPNGシグネチャ）
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	tests := []struct {
		name        string
		injectUser  *domain.User
		fileContent []byte
		fileName    string
		contentType string
		uc          *mockUserUsecase
		wantStatus  int
		wantBody    string
	}{
		{
			name:        "有効なJPEGで200が返る",
			injectUser:  validUser,
			fileContent: append(jpegHeader, make([]byte, 100)...),
			fileName:    "avatar.jpg",
			contentType: "image/jpeg",
			uc: &mockUserUsecase{
				uploadAvatarFunc: func(_ context.Context, _ int64, _ []byte, _ string) (string, error) {
					return "users/1_12345.jpg", nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   "avatar_url",
		},
		{
			name:        "有効なPNGで200が返る",
			injectUser:  validUser,
			fileContent: append(pngHeader, make([]byte, 100)...),
			fileName:    "avatar.png",
			contentType: "image/png",
			uc: &mockUserUsecase{
				uploadAvatarFunc: func(_ context.Context, _ int64, _ []byte, _ string) (string, error) {
					return "users/1_12345.png", nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   "avatar_url",
		},
		{
			name:        "5MB超過で400が返る",
			injectUser:  validUser,
			fileContent: make([]byte, 5<<20+1),
			fileName:    "big.jpg",
			contentType: "image/jpeg",
			uc:          &mockUserUsecase{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "不正な形式（text/plain）で400が返る",
			injectUser:  validUser,
			fileContent: []byte("not an image"),
			fileName:    "test.txt",
			contentType: "text/plain",
			uc:          &mockUserUsecase{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "未認証で401が返る",
			injectUser:  nil,
			fileContent: append(jpegHeader, make([]byte, 100)...),
			fileName:    "avatar.jpg",
			contentType: "image/jpeg",
			uc:          &mockUserUsecase{},
			wantStatus:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(tt.uc, zerolog.Nop(), "http://test-storage")

			// multipart/form-data リクエストを構築
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			part, err := writer.CreateFormFile("image", tt.fileName)
			if err != nil {
				t.Fatalf("multipart作成失敗: %v", err)
			}
			if _, err := part.Write(tt.fileContent); err != nil {
				t.Fatalf("ファイル書き込み失敗: %v", err)
			}
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/avatar", &buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeUploadAvatar(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d, body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

			if tt.wantBody != "" {
				if !strings.Contains(rec.Body.String(), tt.wantBody) {
					t.Errorf("レスポンスボディに %q が含まれない: got %q", tt.wantBody, rec.Body.String())
				}
			}
		})
	}
}

// TestUserHandler_ServeDeleteAvatar は DELETE /api/v1/users/avatar の各パターンを検証する。
func TestUserHandler_ServeDeleteAvatar(t *testing.T) {
	validUser := &domain.User{
		ID:   1,
		Name: "Test User",
		Role: "user",
	}

	tests := []struct {
		name       string
		injectUser *domain.User
		uc         *mockUserUsecase
		wantStatus int
	}{
		{
			name:       "正常に204が返る",
			injectUser: validUser,
			uc: &mockUserUsecase{
				deleteAvatarFunc: func(_ context.Context, _ int64) error {
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "未認証で401が返る",
			injectUser: nil,
			uc:         &mockUserUsecase{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "usecaseエラーで500が返る",
			injectUser: validUser,
			uc: &mockUserUsecase{
				deleteAvatarFunc: func(_ context.Context, _ int64) error {
					return errors.New("delete error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(tt.uc, zerolog.Nop(), "http://test-storage")
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/avatar", nil)

			if tt.injectUser != nil {
				ctx := middleware.ContextWithUser(req.Context(), tt.injectUser)
				req = req.WithContext(ctx)
			}

			rec := httptest.NewRecorder()
			h.ServeDeleteAvatar(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("ステータスコード = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
