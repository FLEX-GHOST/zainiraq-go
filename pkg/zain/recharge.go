package zain

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) RechargeVoucher(ctx context.Context, pin string, msisdn ...string) error {
	if pin == "" {
		return fmt.Errorf("zain: pin cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := VoucherRechargeReq{
		MSISDN:  target,
		PinCode: pin,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/charge-voucher", req, false, false)
	return err
}

func (c *Client) CreditTransfer(ctx context.Context, recipient string, amount int64, otpConfirmation string) error {
	if recipient == "" {
		return fmt.Errorf("zain: recipient cannot be empty")
	}
	if amount <= 0 {
		return fmt.Errorf("zain: amount must be greater than zero")
	}
	if otpConfirmation == "" {
		return fmt.Errorf("zain: otpConfirmation cannot be empty")
	}

	sender := c.GetMSISDN()
	req := CreditTransferReq{
		Sender:          sender,
		Recipient:       recipient,
		Amount:          amount,
		OTPConfirmation: otpConfirmation,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/credit-transfer", req, false, false)
	return err
}

func (c *Client) ExtendValidity(ctx context.Context, amount int64, msisdn ...string) (*ExtendValidityResp, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("zain: amount must be greater than zero")
	}

	target := c.resolveMSISDN(msisdn)
	req := ExtendValidityReq{
		MSISDN: target,
		Amount: amount,
	}

	return doAndDecode[ExtendValidityResp](ctx, c, http.MethodPost, "api/number/extend-validity", req, false, false)
}

func (c *Client) GetValidityOptions(ctx context.Context) (*ValidityOptionBody, error) {
	return doAndDecode[ValidityOptionBody](ctx, c, http.MethodGet, "api/number/validity-options", nil, false, false)
}
