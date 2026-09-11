package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetProfile(ctx context.Context) (*UserResp, error) {
	return doAndDecode[UserResp](ctx, c, http.MethodGet, "api/v2/user/profile", nil, false, false)
}

func (c *Client) UpdateProfile(ctx context.Context, req *UpdateProfileReq) (*UserResp, error) {
	if req == nil {
		return nil, fmt.Errorf("zain: request body cannot be nil")
	}
	return doAndDecode[UserResp](ctx, c, http.MethodPatch, "api/v2/user/profile", req, false, false)
}

func (c *Client) resolveMSISDN(msisdn []string) string {
	if len(msisdn) > 0 && msisdn[0] != "" {
		return msisdn[0]
	}
	return c.GetMSISDN()
}

func (c *Client) GetBalance(ctx context.Context, msisdn ...string) (*BalanceResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/wallet"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[BalanceResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetBillDetails(ctx context.Context, msisdn ...string) (*BillDetailsResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/query-bill"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[BillDetailsResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) CheckPreRegistration(ctx context.Context, msisdn ...string) (*PreRegisterStatusResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/v2/number/pre-registration-check"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[PreRegisterStatusResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetSubAccounts(ctx context.Context, msisdn ...string) ([]SubAccountContent, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/subaccounts"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	resp, err := doAndDecode[[]SubAccountContent](ctx, c, http.MethodGet, path, nil, false, false)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []SubAccountContent{}, nil
	}
	return *resp, nil
}
