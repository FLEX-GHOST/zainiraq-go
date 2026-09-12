package zain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func doAndDecodeCMS[T any](ctx context.Context, c *Client, path string, body any) (*T, error) {
	respBytes, err := c.doRequest(ctx, http.MethodPost, path, body, true, true)
	if err != nil {
		return nil, err
	}

	// 1. Try records wrapper format
	var recordsResp struct {
		Records []RecordResp[T] `json:"records"`
	}
	if err := json.Unmarshal(respBytes, &recordsResp); err == nil && len(recordsResp.Records) > 0 {
		return &recordsResp.Records[0].PayloadResp.Body.Data, nil
	}

	// 2. Try standard APIResponse format {"status":"success", "data": ...}
	var apiResp APIResponse[T]
	if err := json.Unmarshal(respBytes, &apiResp); err == nil && (apiResp.Status == "success" || apiResp.Error == nil) {
		return &apiResp.Data, nil
	}

	// 3. Try ContentResolverResp format
	var crResp ContentResolverResp[T]
	if err := json.Unmarshal(respBytes, &crResp); err == nil {
		return &crResp.Data, nil
	}

	// 4. Try direct unmarshal
	var direct T
	if err := json.Unmarshal(respBytes, &direct); err == nil {
		return &direct, nil
	}

	return nil, fmt.Errorf("zain: failed to decode cms response: %s", string(respBytes))
}

func (c *Client) GetCMSConfiguration(ctx context.Context, query string) (*ConfigCMS, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[ConfigCMS](ctx, c, "dmart/public/excute/query/galleon", req)
}

func (c *Client) GetCMSDashboardBanners(ctx context.Context, query string) (*DashboardBannersResp, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[DashboardBannersResp](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetCMSDashboardMessage(ctx context.Context, query string) (*DashboardMsgResp, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[DashboardMsgResp](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetCMSGameBanners(ctx context.Context, query string) (*DashboardGameResp, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[DashboardGameResp](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetCMSNotification(ctx context.Context, query string) (*InAppMSGCms, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[InAppMSGCms](ctx, c, "dmart/public/excute/query/galleon", req)
}

func (c *Client) GetCMSOffersDetails(ctx context.Context, query string) (*OfferDetailCms, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[OfferDetailCms](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetCMSSubaccount(ctx context.Context, query string) (*SubAccountCms, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[SubAccountCms](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetCMSSubscription(ctx context.Context, query string) (*OfferDetailCms, error) {
	req := ContentRequest{Query: query}
	return doAndDecodeCMS[OfferDetailCms](ctx, c, "dmart/public/excute/query/products", req)
}

func (c *Client) GetSlideshow(ctx context.Context, slideshow string) (*SlideShowContent, error) {
	if slideshow == "" {
		return nil, fmt.Errorf("zain: slideshow name cannot be empty")
	}

	path := "api/slideshows/" + url.PathEscape(slideshow)
	return doAndDecode[SlideShowContent](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetDmartStatus(ctx context.Context) (map[string]any, error) {
	resp, err := doAndDecode[map[string]any](ctx, c, http.MethodGet, "dmart/", nil, true, true)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return map[string]any{}, nil
	}
	return *resp, nil
}

func (c *Client) QueryCMS(ctx context.Context, space string, body any) (map[string]any, error) {
	path := "dmart/public/excute/query/" + url.PathEscape(space)
	resp, err := doAndDecodeCMS[map[string]any](ctx, c, path, body)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return map[string]any{}, nil
	}
	return *resp, nil
}

func (c *Client) GetAppConfigurations(ctx context.Context) (*ConfigCMS, error) {
	return c.GetCMSConfiguration(ctx, "{ is_active: true }")
}

func (c *Client) GetSubaccountsCMS(ctx context.Context) (*SubAccountCms, error) {
	return c.GetCMSSubaccount(ctx, "{ is_active: true }")
}

func (c *Client) GetOffersCatalogCMS(ctx context.Context) (*OfferDetailCms, error) {
	return c.GetCMSOffersDetails(ctx, "{ is_active: true }")
}

func (c *Client) GetFAQsCMS(ctx context.Context) (map[string]any, error) {
	return c.QueryCMS(ctx, "faqs", map[string]any{"is_active": true})
}

func (c *Client) GetRoamingCMS(ctx context.Context) (map[string]any, error) {
	return c.QueryCMS(ctx, "roaming", map[string]any{"is_active": true})
}

