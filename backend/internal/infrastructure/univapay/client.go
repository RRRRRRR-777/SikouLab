// Package univapay はUnivaPay API連携を提供する。
package univapay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client はUnivaPay APIクライアント。
//
// UnivaPay REST API を使用してサブスクリプションの作成を行う。
// 認証は Bearer トークン（{storeSecret}.{storeID}）を使用する。
type Client struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

// NewClient はUnivaPay APIクライアントを作成する。
//
// storeID と storeSecret から Bearer トークンを生成する。
// baseURL が空の場合はデフォルト（https://api.univapay.com）を使用する。
func NewClient(storeID, storeSecret string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    "https://api.univapay.com",
		authToken:  storeSecret + "." + storeID,
	}
}

// subscriptionRequest はUnivaPay サブスクリプション作成リクエストのJSON構造。
type subscriptionRequest struct {
	TransactionTokenID string `json:"transaction_token_id"`
	Amount             int    `json:"amount"`
	Currency           string `json:"currency"`
	Period             string `json:"period"`
}

// subscriptionResponse はUnivaPay サブスクリプション作成レスポンスのJSON構造。
type subscriptionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// univaPayErrorResponse はUnivaPay APIエラーレスポンスのJSON構造。
type univaPayErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateSubscription はUnivaPayでサブスクリプションを作成し、サブスクリプションIDを返す。
//
// POST /subscriptions に対してリクエストを送信する。
// period は "monthly" 固定（現フェーズの仕様）。
func (c *Client) CreateSubscription(ctx context.Context, tokenID string, amount int, currency string) (string, error) {
	reqBody := subscriptionRequest{
		TransactionTokenID: tokenID,
		Amount:             amount,
		Currency:           currency,
		Period:             "monthly",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("リクエストボディのJSON変換失敗: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/subscriptions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("HTTPリクエスト作成失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("UnivaPay APIリクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp univaPayErrorResponse
		if decErr := json.NewDecoder(resp.Body).Decode(&errResp); decErr == nil && errResp.Message != "" {
			return "", fmt.Errorf("UnivaPay APIエラー(status=%d): %s", resp.StatusCode, errResp.Message)
		}
		return "", fmt.Errorf("UnivaPay APIエラー(status=%d)", resp.StatusCode)
	}

	var result subscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("レスポンスのJSON変換失敗: %w", err)
	}

	if result.ID == "" {
		return "", fmt.Errorf("UnivaPay APIレスポンスにサブスクリプションIDが含まれていません")
	}

	return result.ID, nil
}
