package zain

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseURL   = "https://mw-mobile.iq.zain.com"
	DefaultCMSURL    = "https://cms-mobile.iq.zain.com"
	DefaultLanguage  = "ar"
	DefaultAppVer    = "6.6.0"
	DefaultPlatform  = "Android"
	DefaultOSVersion = "14"
	DefaultUserAgent = "okhttp/4.12.0"
)

type Option func(*Client)

type Client struct {
	httpClient     *http.Client
	baseURL        string
	cmsURL         string
	language       string
	platform       string
	osVersion      string
	appVersion     string
	installationID string
	userAgent      string
	accessToken               string
	refreshToken              string
	msisdn                    string
	masterWallet              string
	recordedIncomingTransfers []IncomingTransferRecord
	onTokenUpdate             func(*SessionData)
	mu                        sync.RWMutex
}

func NewClient(opts ...Option) *Client {
	installationID, _ := generateUUID()

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	c := &Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		baseURL:        DefaultBaseURL,
		cmsURL:         DefaultCMSURL,
		language:       DefaultLanguage,
		platform:       DefaultPlatform,
		osVersion:      DefaultOSVersion,
		appVersion:     DefaultAppVer,
		installationID: installationID,
		userAgent:      DefaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithLanguage(lang string) Option {
	return func(c *Client) {
		if lang != "" {
			c.language = lang
		}
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.httpClient.Timeout = d
		}
	}
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.baseURL = strings.TrimRight(url, "/")
		}
	}
}

func WithCMSURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.cmsURL = strings.TrimRight(url, "/")
		}
	}
}

func WithInstallationID(id string) Option {
	return func(c *Client) {
		if id != "" {
			c.installationID = id
		}
	}
}

func WithTokens(accessToken, refreshToken string) Option {
	return func(c *Client) {
		c.accessToken = accessToken
		c.refreshToken = refreshToken
	}
}

func WithMSISDN(msisdn string) Option {
	return func(c *Client) {
		c.msisdn = msisdn
	}
}

func WithMasterWallet(wallet string) Option {
	return func(c *Client) {
		c.masterWallet = wallet
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

func WithOnTokenUpdate(fn func(*SessionData)) Option {
	return func(c *Client) {
		c.onTokenUpdate = fn
	}
}

func (c *Client) SetTokens(accessToken, refreshToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = accessToken
	c.refreshToken = refreshToken
	if c.onTokenUpdate != nil {
		c.onTokenUpdate(&SessionData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			MSISDN:       c.msisdn,
		})
	}
}

func (c *Client) GetTokens() (accessToken, refreshToken string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accessToken, c.refreshToken
}

func (c *Client) SetMSISDN(msisdn string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msisdn = msisdn
}

func (c *Client) GetMSISDN() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.msisdn
}

func (c *Client) SetMasterWallet(wallet string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.masterWallet = wallet
}

func (c *Client) MasterWallet() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.masterWallet != "" {
		return c.masterWallet
	}
	return c.msisdn
}

func (c *Client) RecordIncomingTransfer(rec IncomingTransferRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedIncomingTransfers = append(c.recordedIncomingTransfers, rec)
}

func (c *Client) GetRecordedIncomingTransfers() []IncomingTransferRecord {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]IncomingTransferRecord, len(c.recordedIncomingTransfers))
	copy(res, c.recordedIncomingTransfers)
	return res
}

func (c *Client) ClearRecordedIncomingTransfers() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedIncomingTransfers = nil
}

func (c *Client) ExportSession() *SessionData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return &SessionData{
		AccessToken:  c.accessToken,
		RefreshToken: c.refreshToken,
		MSISDN:       c.msisdn,
	}
}

func (c *Client) ImportSession(s *SessionData) {
	if s == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = s.AccessToken
	c.refreshToken = s.RefreshToken
	c.msisdn = s.MSISDN
}

func (c *Client) SaveSessionToFile(filePath string) error {
	session := c.ExportSession()
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("zain: marshaling session: %w", err)
	}
	return os.WriteFile(filePath, data, 0600)
}

func (c *Client) LoadSessionFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("zain: reading session file: %w", err)
	}
	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return fmt.Errorf("zain: parsing session file: %w", err)
	}
	c.ImportSession(&session)
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, noAuth bool, isCMS bool) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("zain: marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}

	baseURL := c.baseURL
	if isCMS {
		baseURL = c.cmsURL
	}

	fullURL := baseURL + "/" + strings.TrimLeft(path, "/")

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("zain: creating request: %w", err)
	}

	c.mu.RLock()
	lang := c.language
	platform := c.platform
	osVer := c.osVersion
	appVer := c.appVersion
	installID := c.installationID
	token := c.accessToken
	msisdn := c.msisdn
	c.mu.RUnlock()

	req.Header.Set("Skel-Accept-Language", lang)
	req.Header.Set("Skel-Platform", platform)
	req.Header.Set("Skel-OS-Version", osVer)
	req.Header.Set("Skel-Fix-Version", appVer)
	req.Header.Set("Skel-Installation-Id", installID)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	if msisdn != "" {
		req.Header.Set("gal-msisdn", msisdn)
	}

	if !noAuth && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("zain: request failed: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zain: reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp APIResponse[any]
		if jsonErr := json.Unmarshal(respBytes, &errResp); jsonErr == nil && errResp.Error != nil {
			return nil, &APIError{
				StatusCode:       resp.StatusCode,
				Type:             errResp.Error.Type,
				Code:             errResp.Error.Code,
				Message:          errResp.Error.Message,
				LocalizedMessage: errResp.Error.LocalizedMessage,
				CorrelationID:    errResp.Error.CorrelationID,
				ServiceName:      errResp.Error.ServiceName,
				Topic:            errResp.Error.Topic,
				APIPath:          errResp.Error.APIPath,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBytes),
		}
	}

	return respBytes, nil
}

func decodeResponse(resp *http.Response, target any) error {
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("zain: reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp APIResponse[any]
		if jsonErr := json.Unmarshal(respBytes, &errResp); jsonErr == nil && errResp.Error != nil {
			return &APIError{
				StatusCode:       resp.StatusCode,
				Type:             errResp.Error.Type,
				Code:             errResp.Error.Code,
				Message:          errResp.Error.Message,
				LocalizedMessage: errResp.Error.LocalizedMessage,
				CorrelationID:    errResp.Error.CorrelationID,
				ServiceName:      errResp.Error.ServiceName,
				Topic:            errResp.Error.Topic,
				APIPath:          errResp.Error.APIPath,
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBytes),
		}
	}

	if err := json.Unmarshal(respBytes, target); err != nil {
		return fmt.Errorf("zain: decoding response JSON: %w", err)
	}

	return nil
}

func doAndDecode[T any](ctx context.Context, c *Client, method, path string, body any, noAuth bool, isCMS bool) (*T, error) {
	respBytes, err := c.doRequest(ctx, method, path, body, noAuth, isCMS)
	if err != nil {
		return nil, err
	}

	var raw APIResponse[T]
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		var direct T
		if err2 := json.Unmarshal(respBytes, &direct); err2 == nil {
			return &direct, nil
		}
		return nil, fmt.Errorf("zain: decoding response JSON: %w", err)
	}

	if raw.Status != "" && raw.Status != "success" && raw.Error != nil {
		return nil, &APIError{
			StatusCode:       200,
			Type:             raw.Error.Type,
			Code:             raw.Error.Code,
			Message:          raw.Error.Message,
			LocalizedMessage: raw.Error.LocalizedMessage,
			CorrelationID:    raw.Error.CorrelationID,
			ServiceName:      raw.Error.ServiceName,
			Topic:            raw.Error.Topic,
			APIPath:          raw.Error.APIPath,
		}
	}

	return &raw.Data, nil
}

func generateUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generating uuid: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
