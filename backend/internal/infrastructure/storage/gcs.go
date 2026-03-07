// Package storage はファイルストレージの実装を提供する。
package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

// GCSStorage はGoogle Cloud Storageにオブジェクトを保存する。
type GCSStorage struct {
	bucket *storage.BucketHandle
	client *storage.Client
}

// NewGCSStorage はGCSStorageを作成する。
// STORAGE_EMULATOR_HOST環境変数が設定されている場合、GCS SDKが自動でエミュレータに接続する。
func NewGCSStorage(ctx context.Context, bucketName string) (*GCSStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("GCSクライアント作成失敗: %w", err)
	}
	return &GCSStorage{
		bucket: client.Bucket(bucketName),
		client: client,
	}, nil
}

// Save はオブジェクトをGCSに保存し、オブジェクトキーを返す。
// Content-Typeは拡張子から推定してGCSオブジェクトに設定する。
func (s *GCSStorage) Save(ctx context.Context, prefix string, id int64, data []byte, ext string) (string, error) {
	key := fmt.Sprintf("%s/%d_%d%s", prefix, id, time.Now().UnixMilli(), ext)
	w := s.bucket.Object(key).NewWriter(ctx)
	w.ContentType = extToContentType(ext)
	if _, err := w.Write(data); err != nil {
		w.Close()
		return "", fmt.Errorf("GCS書き込み失敗: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("GCSライター終了失敗: %w", err)
	}
	return key, nil
}

// Delete は指定キーのオブジェクトをGCSから除去する。
// オブジェクトが存在しない場合はエラーなしで完了する。
func (s *GCSStorage) Delete(ctx context.Context, key string) error {
	err := s.bucket.Object(key).Delete(ctx)
	if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
		return fmt.Errorf("GCS削除失敗: %w", err)
	}
	return nil
}

// extToContentType は拡張子からContent-Typeを返す。未知の拡張子はapplication/octet-streamを返す。
func extToContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// Close はGCSクライアントを解放する。
func (s *GCSStorage) Close() error {
	return s.client.Close()
}
