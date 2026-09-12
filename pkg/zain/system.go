package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetServerTime(ctx context.Context) (*ServerTimeResp, error) {
	return doAndDecode[ServerTimeResp](ctx, c, http.MethodGet, "api/system/time", nil, true, false)
}

func (c *Client) GetFriendsAndFamily(ctx context.Context, msisdn ...string) (*FriendsAndFamilyResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/friends-and-family"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}
	return doAndDecode[FriendsAndFamilyResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) AddFriendsAndFamily(ctx context.Context, targetMSISDN string, ownerMSISDN ...string) error {
	if targetMSISDN == "" {
		return fmt.Errorf("zain: targetMSISDN cannot be empty")
	}
	owner := c.resolveMSISDN(ownerMSISDN)
	req := FriendsAndFamilyReq{
		MSISDN:       owner,
		TargetMSISDN: targetMSISDN,
	}
	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/friends-and-family/add", req, false, false)
	return err
}

func (c *Client) RemoveFriendsAndFamily(ctx context.Context, targetMSISDN string, ownerMSISDN ...string) error {
	if targetMSISDN == "" {
		return fmt.Errorf("zain: targetMSISDN cannot be empty")
	}
	owner := c.resolveMSISDN(ownerMSISDN)
	req := FriendsAndFamilyReq{
		MSISDN:       owner,
		TargetMSISDN: targetMSISDN,
	}
	_, err := doAndDecode[any](ctx, c, http.MethodDelete, "api/number/friends-and-family/remove", req, false, false)
	return err
}

func (c *Client) GetESIMDetails(ctx context.Context, msisdn ...string) (*ESIMDetailsResp, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/esim-details"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}
	return doAndDecode[ESIMDetailsResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) RequestSIMSwap(ctx context.Context, iccid string, msisdn ...string) (*SIMSwapResp, error) {
	if iccid == "" {
		return nil, fmt.Errorf("zain: iccid cannot be empty")
	}
	target := c.resolveMSISDN(msisdn)
	req := SIMSwapReq{
		MSISDN: target,
		ICCID:  iccid,
	}
	return doAndDecode[SIMSwapResp](ctx, c, http.MethodPost, "api/number/swap-sim", req, false, false)
}
