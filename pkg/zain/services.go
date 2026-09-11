package zain

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetOffers(ctx context.Context, shortName, offerIDs, subscribedOffers string) (*OfferDetails, error) {
	params := url.Values{}
	if shortName != "" {
		params.Set("query_shortname", shortName)
	}
	if offerIDs != "" {
		params.Set("offer_ids", offerIDs)
	}
	if subscribedOffers != "" {
		params.Set("subscribed_offers", subscribedOffers)
	}

	path := "api/offers"
	if query := params.Encode(); query != "" {
		path += "?" + query
	}

	return doAndDecode[OfferDetails](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetSubscriptions(ctx context.Context, msisdn ...string) ([]SubscriptionContent, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/subscriptions"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	resp, err := doAndDecode[[]SubscriptionContent](ctx, c, http.MethodGet, path, nil, false, false)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []SubscriptionContent{}, nil
	}
	return *resp, nil
}

func (c *Client) SubscribeOffer(ctx context.Context, offerID string, msisdn ...string) error {
	if offerID == "" {
		return fmt.Errorf("zain: offerID cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := OfferSubUnSubReq{
		MSISDN:  target,
		OfferID: offerID,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/subscribe", req, false, false)
	return err
}

func (c *Client) UnsubscribeOffer(ctx context.Context, offerID string, msisdn ...string) error {
	if offerID == "" {
		return fmt.Errorf("zain: offerID cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := OfferSubUnSubReq{
		MSISDN:  target,
		OfferID: offerID,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodDelete, "api/number/unsubscribe", req, false, false)
	return err
}

func (c *Client) GetPersonalizedOffersATL(ctx context.Context, msisdn ...string) (*PersonalizedOffersATL, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/personalized-offers?provider_type=ATL"
	if target != "" {
		path += "&msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[PersonalizedOffersATL](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetPersonalizedOffersBTL(ctx context.Context, msisdn ...string) (*PersonalizedOffersBTL, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/personalized-offers?provider_type=BTL"
	if target != "" {
		path += "&msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[PersonalizedOffersBTL](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) ClaimDailyGift(ctx context.Context, msisdn ...string) error {
	target := c.resolveMSISDN(msisdn)
	req := NumberReq{MSISDN: target}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/claim-daily-gift", req, false, false)
	return err
}

func (c *Client) SendBundleGift(ctx context.Context, recipient, offerID string, sender ...string) error {
	if recipient == "" || offerID == "" {
		return fmt.Errorf("zain: recipient and offerID cannot be empty")
	}

	targetSender := c.resolveMSISDN(sender)
	req := SendBundleGiftReq{
		Sender:    targetSender,
		Recipient: recipient,
		OfferID:   offerID,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/send-gift", req, false, false)
	return err
}

func (c *Client) RedeemRegistrationGift(ctx context.Context, offerID string, msisdn ...string) (*OfferResp, error) {
	if offerID == "" {
		return nil, fmt.Errorf("zain: offerID cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := OfferRegistrationReq{
		MSISDN:  target,
		OfferID: offerID,
	}

	return doAndDecode[OfferResp](ctx, c, http.MethodPost, "api/number/redeem-registration-gift", req, false, false)
}

func (c *Client) InviteToKafoo(ctx context.Context, invitedMSISDN string, msisdn ...string) error {
	if invitedMSISDN == "" {
		return fmt.Errorf("zain: invitedMSISDN cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := SendInviteReq{
		MSISDN:  target,
		Invited: invitedMSISDN,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/kafoo_invite", req, false, false)
	return err
}

func (c *Client) GetElectronicBillItems(ctx context.Context, msisdn ...string) (*ElectronicBillItems, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/electronic-bill-items"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[ElectronicBillItems](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetFlexLimits(ctx context.Context, msisdn ...string) (*FlexLimitCms, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/flex-limits"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[FlexLimitCms](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) GetFlexStatus(ctx context.Context, msisdn ...string) (*FlexStatusCms, error) {
	target := c.resolveMSISDN(msisdn)
	path := "api/number/flex-status"
	if target != "" {
		path += "?msisdn=" + url.QueryEscape(target)
	}

	return doAndDecode[FlexStatusCms](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) MigrateToFlex(ctx context.Context, targetOffer string, msisdn ...string) error {
	if targetOffer == "" {
		return fmt.Errorf("zain: targetOffer cannot be empty")
	}

	target := c.resolveMSISDN(msisdn)
	req := MigrateFlexReq{
		MSISDN:      target,
		TargetOffer: targetOffer,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/number/migrate-to-flex", req, false, false)
	return err
}

func (c *Client) GetFlexBundles(ctx context.Context, sourceOffer string, bundleType ...string) ([]FlexUpsellCms, error) {
	if sourceOffer == "" {
		return nil, fmt.Errorf("zain: sourceOffer cannot be empty")
	}

	t := "flex_upsell"
	if len(bundleType) > 0 && bundleType[0] != "" {
		t = bundleType[0]
	}

	path := fmt.Sprintf("api/offers/%s?type=%s", url.PathEscape(sourceOffer), url.QueryEscape(t))
	resp, err := doAndDecode[[]FlexUpsellCms](ctx, c, http.MethodGet, path, nil, false, false)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []FlexUpsellCms{}, nil
	}
	return *resp, nil
}

func (c *Client) SharingAddMember(ctx context.Context, memberMSISDN, offerID string, quota int64, ownerMSISDN ...string) error {
	if memberMSISDN == "" || offerID == "" {
		return fmt.Errorf("zain: memberMSISDN and offerID cannot be empty")
	}

	owner := c.resolveMSISDN(ownerMSISDN)
	req := SharingAddMemberReq{
		OwnerMSISDN:  owner,
		MemberMSISDN: memberMSISDN,
		OfferID:      offerID,
		Quota:        quota,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/sharing/addrsc/", req, false, false)
	return err
}

func (c *Client) SharingQueryMembers(ctx context.Context, offerID string, msisdn ...string) (*QueryMemberResp, error) {
	target := c.resolveMSISDN(msisdn)
	params := url.Values{}
	if target != "" {
		params.Set("msisdn", target)
	}
	if offerID != "" {
		params.Set("offer_id", offerID)
	}

	path := "api/sharing/query"
	if q := params.Encode(); q != "" {
		path += "?" + q
	}

	return doAndDecode[QueryMemberResp](ctx, c, http.MethodGet, path, nil, false, false)
}

func (c *Client) SharingRemoveMember(ctx context.Context, memberMSISDN, offerID string, ownerMSISDN ...string) error {
	if memberMSISDN == "" || offerID == "" {
		return fmt.Errorf("zain: memberMSISDN and offerID cannot be empty")
	}

	owner := c.resolveMSISDN(ownerMSISDN)
	req := SharingRemoveMemberReq{
		OwnerMSISDN:  owner,
		MemberMSISDN: memberMSISDN,
		OfferID:      offerID,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/sharing/remove/", req, false, false)
	return err
}

func (c *Client) SharingTransferUnits(ctx context.Context, memberMSISDN string, units int64, ownerMSISDN ...string) error {
	if memberMSISDN == "" {
		return fmt.Errorf("zain: memberMSISDN cannot be empty")
	}
	if units <= 0 {
		return fmt.Errorf("zain: units must be greater than zero")
	}

	owner := c.resolveMSISDN(ownerMSISDN)
	req := SharingTransferMemberReq{
		OwnerMSISDN:  owner,
		MemberMSISDN: memberMSISDN,
		Units:        units,
	}

	_, err := doAndDecode[any](ctx, c, http.MethodPost, "api/sharing/transfer_unit", req, false, false)
	return err
}
