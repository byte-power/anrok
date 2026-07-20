package anrok

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

const defaultBaseURL = "https://api.anrok.com"

var defaultHTTPClient = &http.Client{
	Timeout: 2 * time.Minute,
}

// Client Anrok API 客户端（仅实现 Transactions：createOrUpdate、createEphemeral）
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	logger     Logger
}

// NewClient 创建新的 Anrok 客户端
// apiKey: Anrok API 密钥，对应请求头 Authorization: Bearer {apiKey}
// logger: 日志记录器，可以为 nil（使用 NopLogger）
func NewClient(apiKey string, logger Logger) *Client {
	if logger == nil {
		logger = &NopLogger{}
	}
	return &Client{
		httpClient: defaultHTTPClient,
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
		logger:     logger,
	}
}

// NewClientWithHTTPClient 使用自定义 HTTP 客户端创建 Anrok 客户端
func NewClientWithHTTPClient(apiKey string, httpClient *http.Client, logger Logger) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	if logger == nil {
		logger = &NopLogger{}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
		logger:     logger,
	}
}

// NewClientWithBaseURL 用于测试或代理场景，可指定非默认的 API 根地址（须含 scheme，无尾部斜杠）
func NewClientWithBaseURL(apiKey, baseURL string, httpClient *http.Client, logger Logger) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	if logger == nil {
		logger = &NopLogger{}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
		logger:     logger,
	}
}

func (c *Client) doRequest(path string, requestBody any, response any) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return &RequestError{Op: "marshal request", Err: err}
		}
		bodyReader = bytes.NewReader(bodyBytes)
		c.logger.Debug("anrok_request",
			String("url", url),
			String("body", string(bodyBytes)))
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return &RequestError{Op: "create request", Err: err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("anrok_request_error",
			String("url", url),
			ErrorField(err))
		return &RequestError{Op: "send request", Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("anrok_read_response_error",
			String("url", url),
			ErrorField(err))
		return &RequestError{Op: "read response", Err: err}
	}

	c.logger.Debug("anrok_response",
		String("url", url),
		Number("status_code", resp.StatusCode),
		String("body", string(respBody)))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.buildAPIError(resp, string(respBody))
	}

	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return &RequestError{Op: "unmarshal response", Err: err}
		}
	}

	return nil
}

func (c *Client) buildAPIError(resp *http.Response, body string) error {
	base := APIError{StatusCode: resp.StatusCode, Body: body}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return &RateLimitError{APIError: base, RetryAfter: retryAfter}
	}

	var typed struct {
		Type ErrorType `json:"type"`
	}
	if json.Unmarshal([]byte(body), &typed) == nil && typed.Type != "" {
		return &TypedError{APIError: base, Type: typed.Type}
	}

	return &base
}

// CreateOrUpdateTransaction 根据发票明细计算销售税并在 Anrok 中保存交易（用于报税与阈值监控）
func (c *Client) CreateOrUpdateTransaction(req CreateOrUpdateTransactionRequest) (*CreateOrUpdateTransactionResponse, error) {
	var out CreateOrUpdateTransactionResponse
	if err := c.doRequest("/v1/seller/transactions/createOrUpdate", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateEphemeralTransaction 计算销售税但不写入 Anrok（草稿/预览场景）
func (c *Client) CreateEphemeralTransaction(req CreateEphemeralTransactionRequest) (*CreateEphemeralTransactionResponse, error) {
	var out CreateEphemeralTransactionResponse
	if err := c.doRequest("/v1/seller/transactions/createEphemeral", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
