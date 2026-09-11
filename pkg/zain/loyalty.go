package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) GetLoyaltyInfo(ctx context.Context) (*LoyaltyResp, error) {
	return doAndDecode[LoyaltyResp](ctx, c, http.MethodGet, "api/loyalty/info", nil, false, false)
}

func (c *Client) GetLoyaltyHotBundles(ctx context.Context) (*LoyaltyHotBundlesResp, error) {
	return doAndDecode[LoyaltyHotBundlesResp](ctx, c, http.MethodGet, "api/loyalty/hot-bundles", nil, false, false)
}

func (c *Client) GetLoyaltyBundles(ctx context.Context) (*LoyaltyBundleOffers, error) {
	return doAndDecode[LoyaltyBundleOffers](ctx, c, http.MethodGet, "api/loyalty/rewards/offers", nil, false, false)
}

func (c *Client) GetLoyaltyPromoCodes(ctx context.Context, loyaltyType, governorate, category, cursor string, limit int, lang string) (*LoyaltyPromoCodesResp, error) {
	params := url.Values{}
	if loyaltyType != "" {
		params.Set("loyalty_type", loyaltyType)
	}
	if governorate != "" {
		params.Set("governorates", governorate)
	}
	if category != "" {
		params.Set("category", category)
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if lang != "" {
		params.Set("language", lang)
	}

	path := "api/loyalty/rewards/promo-codes"
	if q := params.Encode(); q != "" {
		path += "?" + q
	}

	return doAndDecode[LoyaltyPromoCodesResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetLoyaltyStatement(ctx context.Context, fromDate, toDate string) (*LoyaltyStatementResp, error) {
	params := url.Values{}
	if fromDate != "" {
		params.Set("from_date", fromDate)
	}
	if toDate != "" {
		params.Set("to_date", toDate)
	}

	path := "api/v2/loyalty/statement"
	if q := params.Encode(); q != "" {
		path += "?" + q
	}

	return doAndDecode[LoyaltyStatementResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetRewardsHistory(ctx context.Context) (*RewardsHistoryResp, error) {
	return doAndDecode[RewardsHistoryResp](ctx, c, http.MethodGet, "api/loyalty/rewards/history", nil, false, false)
}

func (c *Client) GetLoyaltyFAQ(ctx context.Context) (*LoyaltyFAQResp, error) {
	return doAndDecode[LoyaltyFAQResp](ctx, c, http.MethodGet, "api/loyalty/faq", nil, false, false)
}

func (c *Client) RedeemCredit(ctx context.Context, points int, msisdn ...string) error {
	if points <= 0 {
		return fmt.Errorf("zain: points must be greater than zero")
	}

	target := c.resolveMSISDN(msisdn)
	req := RedeemCreditReq{
		MSISDN: target,
		Points: points,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/loyalty/rewards/credit", req, false, false)
	return err
}

func (c *Client) RedeemLoyaltyBundle(ctx context.Context, bundleID string, msisdn ...string) error {
	if bundleID == "" {
		return fmt.Errorf("zain: bundleID cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := RedeemBundleReq{
		MSISDN:   target,
		BundleID: bundleID,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/loyalty/rewards/offers", req, false, false)
	return err
}

func (c *Client) RedeemPromoCode(ctx context.Context, promoCodeID string, msisdn ...string) (*RedeemPromoCodeResp, error) {
	if promoCodeID == "" {
		return nil, fmt.Errorf("zain: promoCodeID cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := RedeemPromoCodeReq{
		MSISDN:      target,
		PromoCodeID: promoCodeID,
	}

	return doAndDecode[RedeemPromoCodeResp](ctx, c, http.MethodPost, "api/loyalty/rewards/promo-code", req, false, false)
}

func (c *Client) GetMerchant(ctx context.Context, category, merchant string) (*MerchantDetail, error) {
	if category == "" || merchant == "" {
		return nil, fmt.Errorf("zain: category and merchant cannot be empty")
	}

	path := fmt.Sprintf("api/v2/imtiyaz/merchant?imtiyaz_category=%s&merchant=%s", url.QueryEscape(category), url.QueryEscape(merchant))
	return doAndDecode[MerchantDetail](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetVoucherHistory(ctx context.Context, category ...string) (*VoucherHistoryBody, error) {
	path := "api/imtiyaz/voucher-history"
	if len(category) > 0 && category[0] != "" {
		path += "?imtiyaz_category=" + url.QueryEscape(category[0])
	}

	return doAndDecode[VoucherHistoryBody](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) RevealVoucher(ctx context.Context, merchantID, category string) (*RevealVoucherResp, error) {
	if merchantID == "" || category == "" {
		return nil, fmt.Errorf("zain: merchantID and category cannot be empty")
	}

	req := RevealVoucherReq{
		MerchantID: merchantID,
		Category:   category,
	}

	return doAndDecode[RevealVoucherResp](ctx, c, http.MethodPost, "api/v2/imtiyaz/reveal", req, false, false)
}

func (c *Client) ReportMerchant(ctx context.Context, merchantID, reason, comment string) error {
	if merchantID == "" || reason == "" {
		return fmt.Errorf("zain: merchantID and reason cannot be empty")
	}

	req := ReportMerchantReq{
		MerchantID: merchantID,
		Reason:     reason,
		Comment:    comment,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/imtiyaz/report", req, false, false)
	return err
}
