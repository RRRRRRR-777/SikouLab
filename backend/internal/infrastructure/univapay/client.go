// Package univapay はUnivaPay API連携を提供する。
package univapay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client はUnivaPay APIクライアント。
//
// UnivaPay REST API を使用してサブスクリプションの作成を行う。
// 認証は Bearer トークン（{storeSecret}.{appToken}）を使用する。
// URLパスに必要なストアUUIDはAPIレスポンスから取得する。
type Client struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

// NewClient はUnivaPay APIクライアントを作成する。
//
// appToken と storeSecret から Bearer トークンを生成する。
// appToken はUnivaPay管理画面で発行されるアプリケーショントークン（JWT）。
// URLパスに使用するストアUUIDとは異なるため、ストアUUIDはAPIレスポンスから取得する。
func NewClient(appToken, storeSecret string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL:   "https://api.univapay.com",
		authToken: storeSecret + "." + appToken,
	}
}

// subscriptionRequest はUnivaPay サブスクリプション作成リクエストのJSON構造。
type subscriptionRequest struct {
	TransactionTokenID string `json:"transaction_token_id"`
	Amount             int    `json:"amount"`
	Currency           string `json:"currency"`
	Period             string `json:"period"`
}

// subscriptionResponse はUnivaPay サブスクリプション作成/取得レスポンスのJSON構造。
type subscriptionResponse struct {
	ID      string `json:"id"`
	StoreID string `json:"store_id"`
	Status  string `json:"status"`
}

// univaPayErrorResponse はUnivaPay APIエラーレスポンスのJSON構造。
type univaPayErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetSubscription はWebhook受信時にサブスクリプションの最新状態を確認するために使用する。
//
// GET /stores/{storeID}/subscriptions/{id} に対してリクエストを送信し、ステータス文字列を返す。
// storeID はCreateSubscriptionのレスポンスから取得したストアUUIDを使用する。
func (c *Client) GetSubscription(ctx context.Context, storeID, subscriptionID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/stores/"+storeID+"/subscriptions/"+subscriptionID, nil)
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト作成失敗: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("UnivaPay APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rawBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("UnivaPay APIエラー(status=%d): %s", resp.StatusCode, string(rawBody))
	}

	var result subscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("レスポンスのJSON変換失敗: %w", err)
	}

	return result.Status, nil
}

// CreateSubscription はユーザーの決済トークンを使って月額課金を開始する。
//
// POST /subscriptions に対してリクエストを送信し、サブスクリプションIDとストアUUIDを返す。
// period は "monthly" 固定（現フェーズの仕様）。
// 戻り値の storeID は後続の GetSubscription 呼び出しに必要となる。
func (c *Client) CreateSubscription(ctx context.Context, tokenID string, amount int, currency string) (string, string, error) {
	reqBody := subscriptionRequest{
		TransactionTokenID: tokenID,
		Amount:             amount,
		Currency:           currency,
		Period:             "monthly",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("リクエストボディのJSON変換失敗: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/subscriptions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", "", fmt.Errorf("HTTPリクエスト作成失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("UnivaPay APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		rawBody, _ := io.ReadAll(resp.Body)
		var errResp univaPayErrorResponse
		if decErr := json.Unmarshal(rawBody, &errResp); decErr == nil && errResp.Message != "" {
			return "", "", fmt.Errorf("UnivaPay APIエラー(status=%d): %s", resp.StatusCode, errResp.Message)
		}
		return "", "", fmt.Errorf("UnivaPay APIエラー(status=%d): %s", resp.StatusCode, string(rawBody))
	}

	var result subscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("レスポンスのJSON変換失敗: %w", err)
	}

	if result.ID == "" {
		return "", "", fmt.Errorf("UnivaPay APIレスポンスにサブスクリプションIDが含まれていません")
	}

	return result.ID, result.StoreID, nil
}
