package zain

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized       = errors.New("zain: unauthorized or missing token")
	ErrTokenExpired       = errors.New("zain: session or token has expired")
	ErrRateLimited        = errors.New("zain: rate limited by server")
	ErrNotFound           = errors.New("zain: requested resource not found")
	ErrBadRequest         = errors.New("zain: bad request or invalid parameters")
	ErrInvalidCredentials = errors.New("zain: invalid login credentials or OTP")
	ErrServerInternal     = errors.New("zain: internal server error")
)

type APIError struct {
	StatusCode       int           `json:"status_code"`
	Type             string        `json:"type,omitempty"`
	Code             int           `json:"code,omitempty"`
	Message          string        `json:"message,omitempty"`
	LocalizedMessage *LokaliseText `json:"localized_message,omitempty"`
	CorrelationID    string        `json:"correlation_id,omitempty"`
	ServiceName      string        `json:"service_name,omitempty"`
	Topic            string        `json:"topic,omitempty"`
	APIPath          string        `json:"api_path,omitempty"`
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" && e.LocalizedMessage != nil {
		msg = e.LocalizedMessage.String()
	}
	if msg == "" {
		msg = fmt.Sprintf("HTTP status %d", e.StatusCode)
	}
	if e.Code != 0 {
		return fmt.Sprintf("zain error (code %d, http %d): %s", e.Code, e.StatusCode, msg)
	}
	return fmt.Sprintf("zain error (http %d): %s", e.StatusCode, msg)
}

func (e *APIError) Is(target error) bool {
	switch {
	case errors.Is(target, ErrUnauthorized):
		return e.StatusCode == 401 || e.Code == 401
	case errors.Is(target, ErrTokenExpired):
		return e.StatusCode == 401 || e.Type == "token_expired" || e.Code == 401
	case errors.Is(target, ErrRateLimited):
		return e.StatusCode == 429 || e.Code == 429
	case errors.Is(target, ErrNotFound):
		return e.StatusCode == 404 || e.Code == 404
	case errors.Is(target, ErrBadRequest):
		return e.StatusCode == 400 || e.Code == 400
	case errors.Is(target, ErrInvalidCredentials):
		return e.StatusCode == 401 || e.StatusCode == 403 || e.Code == 401
	case errors.Is(target, ErrServerInternal):
		return e.StatusCode >= 500
	default:
		return false
	}
}
