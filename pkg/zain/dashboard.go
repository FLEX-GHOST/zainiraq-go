package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) GetDashboardMessages(ctx context.Context, msisdn ...string) (*DashboardMsgContent, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/dashboard_message"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[DashboardMsgContent](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetNotifications(ctx context.Context, offset, limit int, read *bool) (*InAppMSGPageContent, error) {
	params := url.Values{}
	params.Set("offset", strconv.Itoa(offset))
	if limit <= 0 {
		limit = 20
	}
	params.Set("limit", strconv.Itoa(limit))
	if read != nil {
		params.Set("read", strconv.FormatBool(*read))
	}

	path := "api/notifications?" + params.Encode()
	return doAndDecode[InAppMSGPageContent](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) ReadNotification(ctx context.Context, shortname string) error {
	if shortname == "" {
		return fmt.Errorf("zain: shortname cannot be empty")
	}

	path := fmt.Sprintf("api/notifications/read-notification/%s", url.PathEscape(shortname))
	_, err := doAndDecode[any](ctx, c, http.MethodPatch, path, nil, false, false)
	return err
}

func (c *Client) ReadAllNotifications(ctx context.Context) error {
	_, err := doAndDecode[any](ctx, c, http.MethodPatch, "api/notifications/read-all-notifications", nil, false, false)
	return err
}

func (c *Client) DeleteNotification(ctx context.Context, shortname string) error {
	if shortname == "" {
		return fmt.Errorf("zain: shortname cannot be empty")
	}

	path := fmt.Sprintf("api/notifications/delete-notification/%s", url.PathEscape(shortname))
	_, err := doAndDecode[any](ctx, c, http.MethodDelete, path, nil, false, false)
	return err
}

func (c *Client) GetStoreNearMe(ctx context.Context, searchTerm ...string) (*ZainStoreNearMeContent, error) {
	path := "api/near_me_store/near-me-store"
	if len(searchTerm) > 0 && searchTerm[0] != "" {
		path += "?searchTerm=" + url.QueryEscape(searchTerm[0])
	}

	return doAndDecode[ZainStoreNearMeContent](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetDigitalServices(ctx context.Context, category, search, billingType string) (*DigitalServiceContent, error) {
	params := url.Values{}
	params.Set("application", "galleon")
	params.Set("platform", "android")
	if billingType != "" {
		params.Set("customer_billing_type", billingType)
	} else {
		params.Set("customer_billing_type", "PREPAID_NORMAL")
	}
	if category != "" {
		params.Set("category", category)
	}
	if search != "" {
		params.Set("search", search)
	}

	path := "api/digital_services_catalogue_management?" + params.Encode()
	return doAndDecode[DigitalServiceContent](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetReferralCount(ctx context.Context) (*ReferralResp, error) {
	return doAndDecode[ReferralResp](ctx, c, http.MethodGet, "api/referral/count", nil, false, false)
}

func (c *Client) ReportNewInstall(ctx context.Context, installID, referralCode string) error {
	if installID == "" {
		c.mu.RLock()
		installID = c.installationID
		c.mu.RUnlock()
	}

	req := NewInstallReq{
		InstallationID: installID,
		ReferralCode:   referralCode,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/referral/new-app-install", req, true, false)
	return err
}

func (c *Client) GetRandomError(ctx context.Context) error {
	_, err := doAndDecode[any](ctx, c, http.MethodGet, "api/errors/random", nil, false, false)
	return err
}

