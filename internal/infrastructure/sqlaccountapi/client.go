package sqlaccountapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"shwetaik-sqlacc-stock-api/internal/config"
)

type Client struct {
	baseURL      string
	accessKey    string
	secretKey    string
	region       string
	service      string
	sessionToken string
	httpClient   *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:      cfg.SQLAccountAPIHost,
		accessKey:    cfg.SQLAccountAPIAccessKey,
		secretKey:    cfg.SQLAccountAPISecretKey,
		region:       cfg.SQLAccountAPIRegion,
		service:      cfg.SQLAccountAPIService,
		sessionToken: cfg.SQLAccountAPISessionToken,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) CreatePaymentVoucher(ctx context.Context, payload map[string]any) (map[string]any, error) {
	return c.doSignedJSON(ctx, http.MethodPost, "/paymentvoucher", payload)
}

func (c *Client) PutStockItemPrice(ctx context.Context, dockey int, payload map[string]any) (map[string]any, error) {
	return c.doSignedJSON(ctx, http.MethodPut, fmt.Sprintf("/stockitem/%d", dockey), payload)
}

func (c *Client) doSignedJSON(ctx context.Context, method, path string, payload map[string]any) (map[string]any, error) {
	if c.baseURL == "" || c.accessKey == "" || c.secretKey == "" || c.region == "" || c.service == "" {
		return nil, fmt.Errorf("sqlaccountapi: vendor API credentials are not configured")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("sqlaccountapi: marshal request body: %w", err)
	}

	url := strings.TrimSuffix(c.baseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("sqlaccountapi: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	signRequest(req, body, c.accessKey, c.secretKey, c.region, c.service, c.sessionToken, time.Now())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sqlaccountapi: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sqlaccountapi: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sqlaccountapi: vendor API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		return map[string]any{}, nil
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("sqlaccountapi: decode response: %w", err)
	}
	return result, nil
}
