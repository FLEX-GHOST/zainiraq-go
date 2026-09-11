package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) GetWalletBalance(ctx context.Context, msisdn ...string) (*BalanceResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/wallet"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}
	return doAndDecode[BalanceResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) RequestCreditTransferOTP(ctx context.Context, senderMSISDN ...string) (*OTPRequestIdResp, error) {
	target := c.resolveMSISDN(senderMSISDN)
	if target == "" {
		target = c.MasterWallet()
	}
	return c.RequestOTP(ctx, target)
}

func (c *Client) ConfirmCreditTransferOTP(ctx context.Context, otpCode string, senderMSISDN ...string) (*OTPConfirmationIdResp, error) {
	target := c.resolveMSISDN(senderMSISDN)
	if target == "" {
		target = c.MasterWallet()
	}
	return c.ConfirmOTP(ctx, target, otpCode)
}

func (c *Client) FormatUSSDTransfer(recipient string, amount int64) string {
	cleanRecipient := cleanDigits(recipient)
	return fmt.Sprintf("*123*%d*%s#", amount, cleanRecipient)
}
