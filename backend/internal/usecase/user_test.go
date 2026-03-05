package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

func ptrStr(s string) *string { return &s }

// mockUserProfileRepository はUserProfileRepositoryのモック。
type mockUserProfileRepository struct {
	updateDisplayNameFunc func(ctx context.Context, userID int64, displayName string) (*domain.User, error)
	updateAvatarURLFunc   func(ctx context.Context, userID int64, avatarURL string) error
	clearAvatarURLFunc    func(ctx context.Context, userID int64) error
	findByIDFunc          func(ctx context.Context, id int64) (*domain.User, error)
}

// UpdateDisplayName はモックの表示名更新を実行する。
func (m *mockUserProfileRepository) UpdateDisplayName(ctx context.Context, userID int64, displayName string) (*domain.User, error) {
	return m.updateDisplayNameFunc(ctx, userID, displayName)
}

// UpdateAvatarURL はモックのアバターURL更新を実行する。
func (m *mockUserProfileRepository) UpdateAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	return m.updateAvatarURLFunc(ctx, userID, avatarURL)
}

// ClearAvatarURL はモックのアバターURL削除を実行する。
func (m *mockUserProfileRepository) ClearAvatarURL(ctx context.Context, userID int64) error {
	return m.clearAvatarURLFunc(ctx, userID)
}

// FindByID はモックのID検索を実行する。
func (m *mockUserProfileRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.findByIDFunc(ctx, id)
}

// mockObjectStorage はObjectStorageのモック。
type mockObjectStorage struct {
	saveFunc   func(ctx context.Context, prefix string, id int64, data []byte, ext string) (string, error)
	deleteFunc func(ctx context.Context, key string) error
}

// Save はモックのオブジェクト保存を実行する。
func (m *mockObjectStorage) Save(ctx context.Context, prefix string, id int64, data []byte, ext string) (string, error) {
	return m.saveFunc(ctx, prefix, id, data, ext)
}

// Delete はモックのオブジェクト削除を実行する。
func (m *mockObjectStorage) Delete(ctx context.Context, key string) error {
	return m.deleteFunc(ctx, key)
}

// TestUserUsecase_UpdateDisplayName は表示名更新ユースケースの各パターンを検証する。
func TestUserUsecase_UpdateDisplayName(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		display string
		repo    *mockUserProfileRepository
		wantErr bool
	}{
		{
			name:    "正常な表示名で更新成功",
			userID:  1,
			display: "山田 太郎",
			repo: &mockUserProfileRepository{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return &domain.User{ID: 1, DisplayName: "山田 太郎"}, nil
				},
			},
			wantErr: false,
		},
		{
			name:    "空文字でバリデーション失敗",
			userID:  1,
			display: "",
			repo:    &mockUserProfileRepository{},
			wantErr: true,
		},
		{
			name:    "51文字以上でバリデーション失敗",
			userID:  1,
			display: "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんがぎぐげござ",
			repo:    &mockUserProfileRepository{},
			wantErr: true,
		},
		{
			name:    "空白のみでバリデーション失敗",
			userID:  1,
			display: "   ",
			repo:    &mockUserProfileRepository{},
			wantErr: true,
		},
		{
			name:    "DBエラー時",
			userID:  1,
			display: "山田 太郎",
			repo: &mockUserProfileRepository{
				updateDisplayNameFunc: func(_ context.Context, _ int64, _ string) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, nil)

			user, err := uc.UpdateDisplayName(context.Background(), tt.userID, tt.display)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if user == nil {
				t.Fatal("ユーザーがnilで返された")
				return
			}
			if user.DisplayName != tt.display {
				t.Errorf("表示名 = %q, want %q", user.DisplayName, tt.display)
			}
		})
	}
}

// TestUserUsecase_UploadAvatar はアバターアップロードユースケースの各パターンを検証する。
func TestUserUsecase_UploadAvatar(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		fileData    []byte
		contentType string
		repo        *mockUserProfileRepository
		storage     *mockObjectStorage
		wantErr     bool
		wantURL     string
	}{
		{
			name:        "正常なJPEG画像で更新成功",
			userID:      1,
			fileData:    append([]byte{0xFF, 0xD8, 0xFF, 0xE0}, make([]byte, 100)...),
			contentType: "image/jpeg",
			repo: &mockUserProfileRepository{
				findByIDFunc: func(_ context.Context, _ int64) (*domain.User, error) {
					return &domain.User{ID: 1}, nil
				},
				updateAvatarURLFunc: func(_ context.Context, _ int64, _ string) error {
					return nil
				},
			},
			storage: &mockObjectStorage{
				saveFunc: func(_ context.Context, _ string, _ int64, _ []byte, _ string) (string, error) {
					return "users/1_12345.jpg", nil
				},
				deleteFunc: func(_ context.Context, _ string) error {
					return nil
				},
			},
			wantErr: false,
			wantURL: "users/1_12345.jpg",
		},
		{
			name:        "ストレージ保存失敗",
			userID:      1,
			fileData:    append([]byte{0xFF, 0xD8, 0xFF, 0xE0}, make([]byte, 100)...),
			contentType: "image/jpeg",
			repo: &mockUserProfileRepository{
				findByIDFunc: func(_ context.Context, _ int64) (*domain.User, error) {
					return &domain.User{ID: 1}, nil
				},
			},
			storage: &mockObjectStorage{
				saveFunc: func(_ context.Context, _ string, _ int64, _ []byte, _ string) (string, error) {
					return "", errors.New("storage error")
				},
				deleteFunc: func(_ context.Context, _ string) error {
					return nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, tt.storage)

			url, err := uc.UploadAvatar(context.Background(), tt.userID, tt.fileData, tt.contentType)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if url != tt.wantURL {
				t.Errorf("URL = %q, want %q", url, tt.wantURL)
			}
		})
	}
}

// TestUserUsecase_DeleteAvatar はアバター削除ユースケースの各パターンを検証する。
func TestUserUsecase_DeleteAvatar(t *testing.T) {
	tests := []struct {
		name    string
		userID  int64
		repo    *mockUserProfileRepository
		storage *mockObjectStorage
		wantErr bool
	}{
		{
			name:   "正常に削除成功",
			userID: 1,
			repo: &mockUserProfileRepository{
				findByIDFunc: func(_ context.Context, _ int64) (*domain.User, error) {
					return &domain.User{ID: 1, AvatarURL: ptrStr("users/1_12345.jpg")}, nil
				},
				clearAvatarURLFunc: func(_ context.Context, _ int64) error {
					return nil
				},
			},
			storage: &mockObjectStorage{
				deleteFunc: func(_ context.Context, _ string) error {
					return nil
				},
				saveFunc: func(_ context.Context, _ string, _ int64, _ []byte, _ string) (string, error) {
					return "", nil
				},
			},
			wantErr: false,
		},
		{
			name:   "アバター未設定でも正常に完了",
			userID: 1,
			repo: &mockUserProfileRepository{
				findByIDFunc: func(_ context.Context, _ int64) (*domain.User, error) {
					return &domain.User{ID: 1, AvatarURL: nil}, nil
				},
				clearAvatarURLFunc: func(_ context.Context, _ int64) error {
					return nil
				},
			},
			storage: &mockObjectStorage{
				deleteFunc: func(_ context.Context, _ string) error {
					return nil
				},
				saveFunc: func(_ context.Context, _ string, _ int64, _ []byte, _ string) (string, error) {
					return "", nil
				},
			},
			wantErr: false,
		},
		{
			name:   "DBエラー時",
			userID: 1,
			repo: &mockUserProfileRepository{
				findByIDFunc: func(_ context.Context, _ int64) (*domain.User, error) {
					return nil, errors.New("db error")
				},
			},
			storage: &mockObjectStorage{
				deleteFunc: func(_ context.Context, _ string) error {
					return nil
				},
				saveFunc: func(_ context.Context, _ string, _ int64, _ []byte, _ string) (string, error) {
					return "", nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUserUsecase(tt.repo, tt.storage)

			err := uc.DeleteAvatar(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーが期待されたが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
		})
	}
}
