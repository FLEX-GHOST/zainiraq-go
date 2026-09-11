package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) InitiateZainCashPayment(ctx context.Context, req *InitiateOrderReq) (*InitiateOrderResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: initiate payment request cannot be nil")
	}
	return doAndDecode[InitiateOrderResp](ctx, c, http.MethodPost, "api/payment/zaincash/initiate", req, false, false)
}

func (c *Client) InitiateZainCashPaymentV2(ctx context.Context, req *InitiateOrderReq) (*InitiateOrderResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: initiate payment request cannot be nil")
	}
	return doAndDecode[InitiateOrderResp](ctx, c, http.MethodPost, "api/v2/payment/zaincash/initiate", req, false, false)
}

func (c *Client) InitiateCreditCardOrder(ctx context.Context, req *InitiateOrderReq) (*InitiateOrderResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: initiate credit card order request cannot be nil")
	}
	return doAndDecode[InitiateOrderResp](ctx, c, http.MethodPost, "api/payment/creditcard/initiate", req, false, false)
}

func (c *Client) InitiateOrder(ctx context.Context, req *InitiateOrderReq) (*InitiateOrderResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: initiate order request cannot be nil")
	}
	return doAndDecode[InitiateOrderResp](ctx, c, http.MethodPost, "api/purchase/order", req, false, false)
}

func (c *Client) GetPaymentStatus(ctx context.Context, token, transactionID string) (*OrderPaymentStatusResp, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if transactionID != "" {
		params.Set("transaction_id", transactionID)
	}

	path := "api/payment/status"
	if q := params.Encode(); q != "" {
		path += "?" + q
	}

	return doAndDecode[OrderPaymentStatusResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetPaymentStatusV2(ctx context.Context, token, transactionID string) (*OrderPaymentStatusResp, error) {
	params := url.Values{}
	if token != "" {
		params.Set("token", token)
	}
	if transactionID != "" {
		params.Set("transaction_id", transactionID)
	}

	path := "api/v2/payment/status"
	if q := params.Encode(); q != "" {
		path += "?" + q
	}

	return doAndDecode[OrderPaymentStatusResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) CreateCheckoutID(ctx context.Context, req *CheckoutReq) (*CheckoutResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: checkout request cannot be nil")
	}
	return doAndDecode[CheckoutResp](ctx, c, http.MethodPost, "api/payment/create-checkout-id", req, false, false)
}

func (c *Client) SaveCreditCard(ctx context.Context, req *SaveCreditCardReq) error {
	if req == nil {
		return fmt.Errorf("zain: save credit card request cannot be nil")
	}
	if req.SpaceName == "" {
		req.SpaceName = "galleon"
	}
	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/payment/save-card", req, false, false)
	return err
}

func (c *Client) SetDefaultCard(ctx context.Context, cardToken string) error {
	if cardToken == "" {
		return fmt.Errorf("zain: cardToken cannot be empty")
	}

	req := CardTokenReq{
		CardRegistrationToken: cardToken,
		SpaceName:             "galleon",
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/payment/set-default-card", req, false, false)
	return err
}

func (c *Client) DeleteCard(ctx context.Context, cardToken string) error {
	if cardToken == "" {
		return fmt.Errorf("zain: cardToken cannot be empty")
	}

	req := CardTokenReq{
		CardRegistrationToken: cardToken,
		SpaceName:             "galleon",
	}

	_, err := doAndDecode[any](ctx, c, http.MethodDelete, "api/payment/delete-card", req, false, false)
	return err
}

func (c *Client) RefreshPaymentStatus(ctx context.Context, req *RefreshPaymentStatusReq) (*GatewayPaymentStatusResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: refresh payment status request cannot be nil")
	}
	return doAndDecode[GatewayPaymentStatusResp](ctx, c, http.MethodPost, "api/payment/refresh-payment-status", req, false, false)
}
