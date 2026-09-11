package zain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) RequestOTP(ctx context.Context, msisdn string) (*OTPRequestIdResp, error) {
	if msisdn == "" {
		return nil, fmt.Errorf("zain: msisdn cannot be empty")
	}
	c.SetMSISDN(msisdn)

	req := OTPRequestReq{MSISDN: msisdn}
	respBytes, err := c.doRequest(ctx, http.MethodPost, "api/otp/request", req, true, false)
	if err != nil {
		return nil, err
	}

	var resp OTPRequestIdResp
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("zain: unmarshaling otp request response: %w", err)
	}

	return &resp, nil
}

func (c *Client) ConfirmOTP(ctx context.Context, msisdn, code string) (*OTPConfirmationIdResp, error) {
	if msisdn == "" || code == "" {
		return nil, fmt.Errorf("zain: msisdn and code cannot be empty")
	}

	req := OTPConfirmReq{
		MSISDN:  msisdn,
		OTPCode: code,
	}

	return doAndDecode[OTPConfirmationIdResp](ctx, c, http.MethodPost, "api/otp/confirm", req, true, false)
}

func (c *Client) SignSession(ctx context.Context, msisdn, confirmation string) (*LoginSession, error) {
	if msisdn == "" || confirmation == "" {
		return nil, fmt.Errorf("zain: msisdn and confirmation cannot be empty")
	}

	req := OTPSessionSignReq{
		MSISDN:       msisdn,
		Confirmation: confirmation,
		UserSpace:    "galleon",
	}

	session, err := doAndDecode[LoginSession](ctx, c, http.MethodPost, "api/user/sign", req, true, false)
	if err != nil {
		return nil, err
	}

	if session != nil && session.AccessToken != "" {
		c.SetTokens(session.AccessToken, session.RefreshToken)
		c.SetMSISDN(msisdn)
	}

	return session, nil
}

func (c *Client) VerifyOTPAndLogin(ctx context.Context, msisdn, code string) (*LoginSession, error) {
	confirmResp, err := c.ConfirmOTP(ctx, msisdn, code)
	if err != nil {
		return nil, fmt.Errorf("zain: confirm otp failed: %w", err)
	}
	if confirmResp == nil || confirmResp.ConfirmationID == "" {
		return nil, fmt.Errorf("zain: missing otp confirmation id in response")
	}

	return c.SignSession(ctx, msisdn, confirmResp.ConfirmationID)
}

func (c *Client) LoginWithPassword(ctx context.Context, msisdn, password string) (*LoginSession, error) {
	if msisdn == "" || password == "" {
		return nil, fmt.Errorf("zain: msisdn and password cannot be empty")
	}
	c.SetMSISDN(msisdn)

	req := SignInReq{
		MSISDN:   msisdn,
		Password: password,
	}

	session, err := doAndDecode[LoginSession](ctx, c, http.MethodPost, "api/user/login", req, true, false)
	if err != nil {
		return nil, err
	}

	if session != nil && session.AccessToken != "" {
		c.SetTokens(session.AccessToken, session.RefreshToken)
	}

	return session, nil
}

func (c *Client) RefreshToken(ctx context.Context) (*LoginSession, error) {
	_, refreshToken := c.GetTokens()
	if refreshToken == "" {
		return nil, ErrUnauthorized
	}

	fullURL := c.baseURL + "/api/user/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("zain: creating refresh token request: %w", err)
	}

	c.mu.RLock()
	req.Header.Set("Skel-Accept-Language", c.language)
	req.Header.Set("Skel-Platform", c.platform)
	req.Header.Set("Skel-OS-Version", c.osVersion)
	req.Header.Set("Skel-Fix-Version", c.appVersion)
	req.Header.Set("Skel-Installation-Id", c.installationID)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("refresh-token", refreshToken)
	c.mu.RUnlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("zain: refresh token request failed: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zain: reading refresh token response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBytes),
		}
	}

	var raw APIResponse[LoginSession]
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		return nil, fmt.Errorf("zain: decoding refresh token response: %w", err)
	}

	if raw.Data.AccessToken != "" {
		c.SetTokens(raw.Data.AccessToken, raw.Data.RefreshToken)
	}

	return &raw.Data, nil
}

func (c *Client) SignUp(ctx context.Context, req *UserInfoReq) (*UserResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: request body cannot be nil")
	}
	return doAndDecode[UserResp](ctx, c, http.MethodPost, "api/user/create", req, true, false)
}

func (c *Client) ResetPassword(ctx context.Context, oldPwd, newPwd string) error {
	req := ResetPwdRequest{
		OldPassword: oldPwd,
		NewPassword: newPwd,
	}
	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/user/reset_password", req, false, false)
	return err
}

func (c *Client) ValidatePassword(ctx context.Context, pwd string) error {
	req := ValidatePwdReq{Password: pwd}
	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/user/validate", req, false, false)
	return err
}

func (c *Client) Logout(ctx context.Context) error {
	_, err := doAndDecode[any](ctx, c, http.MethodDelete, "api/user/logout", nil, false, false)
	c.SetTokens("", "")
	return err
}

func (c *Client) DeleteAccount(ctx context.Context) error {
	_, err := doAndDecode[any](ctx, c, http.MethodDelete, "api/user/delete", nil, false, false)
	c.SetTokens("", "")
	return err
}

func (c *Client) SubmitFeedback(ctx context.Context, rating int, comment string) error {
	req := FeedbackReq{Rating: rating, Comment: comment}
	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/user/feedback", req, false, false)
	return err
}
