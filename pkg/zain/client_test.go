package zain

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientOptionsAndHeaders(t *testing.T) {
	var capturedHeader http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(APIResponse[UserResp]{
			Status: "success",
			Data: UserResp{
				MSISDN: "9647801234567",
				Name:   "Ali",
			},
		})
	}))
	defer ts.Close()

	client := NewClient(
		WithBaseURL(ts.URL),
		WithLanguage("en"),
		WithInstallationID("custom-uuid-1234"),
		WithTokens("access-token-xyz", "refresh-token-abc"),
		WithMSISDN("9647801234567"),
		WithTimeout(5*time.Second),
	)

	if client.GetMSISDN() != "9647801234567" {
		t.Fatalf("expected msisdn 9647801234567, got %s", client.GetMSISDN())
	}

	acc, ref := client.GetTokens()
	if acc != "access-token-xyz" || ref != "refresh-token-abc" {
		t.Fatalf("unexpected tokens: %s, %s", acc, ref)
	}

	ctx := context.Background()
	profile, err := client.GetProfile(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Name != "Ali" {
		t.Fatalf("expected Ali, got %s", profile.Name)
	}

	if capturedHeader.Get("Skel-Accept-Language") != "en" {
		t.Errorf("expected Skel-Accept-Language en, got %s", capturedHeader.Get("Skel-Accept-Language"))
	}
	if capturedHeader.Get("Skel-Installation-Id") != "custom-uuid-1234" {
		t.Errorf("expected installation id custom-uuid-1234, got %s", capturedHeader.Get("Skel-Installation-Id"))
	}
	if capturedHeader.Get("gal-msisdn") != "9647801234567" {
		t.Errorf("expected gal-msisdn 9647801234567, got %s", capturedHeader.Get("gal-msisdn"))
	}
	if capturedHeader.Get("Authorization") != "Bearer access-token-xyz" {
		t.Errorf("expected Bearer access-token-xyz, got %s", capturedHeader.Get("Authorization"))
	}
}

func TestAuthFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/otp/request":
			var req OTPRequestReq
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.MSISDN == "9647800000000" {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(APIResponse[any]{
					Status: "error",
					Error: &ErrorResponse{
						Code:    4001,
						Message: "Invalid phone number",
					},
				})
				return
			}
			_ = json.NewEncoder(w).Encode(OTPRequestIdResp{
				Status: "success",
				Data: struct {
					RequestID string `json:"request_id,omitempty"`
				}{RequestID: "req-123"},
			})
		case "/api/otp/confirm":
			_ = json.NewEncoder(w).Encode(APIResponse[OTPConfirmationIdResp]{
				Status: "success",
				Data: OTPConfirmationIdResp{
					ConfirmationID: "confirm-token-456",
				},
			})
		case "/api/user/sign":
			var req OTPSessionSignReq
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Confirmation != "confirm-token-456" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(APIResponse[LoginSession]{
				Status: "success",
				Data: LoginSession{
					AccessToken:  "jwt-token-111",
					RefreshToken: "refresh-token-222",
				},
			})
		case "/api/user/token":
			refHeader := r.Header.Get("refresh-token")
			if refHeader != "refresh-token-222" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(APIResponse[LoginSession]{
				Status: "success",
				Data: LoginSession{
					AccessToken:  "jwt-token-333",
					RefreshToken: "refresh-token-444",
				},
			})
		case "/api/user/login":
			_ = json.NewEncoder(w).Encode(APIResponse[LoginSession]{
				Status: "success",
				Data: LoginSession{
					AccessToken:  "jwt-token-pwd",
					RefreshToken: "refresh-token-pwd",
				},
			})
		case "/api/user/logout", "/api/user/delete":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{
				Status: "success",
			})
		case "/api/user/reset_password", "/api/user/validate", "/api/user/feedback":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{
				Status: "success",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	ctx := context.Background()

	// 1. Request OTP Error case
	_, err := client.RequestOTP(ctx, "9647800000000")
	if err == nil {
		t.Fatal("expected error on invalid phone, got nil")
	}

	// 2. Request OTP Success
	otpReq, err := client.RequestOTP(ctx, "9647801234567")
	if err != nil {
		t.Fatalf("request otp failed: %v", err)
	}
	if otpReq.Data.RequestID != "req-123" {
		t.Errorf("expected req-123, got %s", otpReq.Data.RequestID)
	}

	// 3. Confirm OTP and Login combined flow
	session, err := client.VerifyOTPAndLogin(ctx, "9647801234567", "1234")
	if err != nil {
		t.Fatalf("verify otp failed: %v", err)
	}
	if session.AccessToken != "jwt-token-111" {
		t.Errorf("expected jwt-token-111, got %s", session.AccessToken)
	}

	acc, ref := client.GetTokens()
	if acc != "jwt-token-111" || ref != "refresh-token-222" {
		t.Fatalf("tokens not saved in client: acc=%s ref=%s", acc, ref)
	}

	// 4. Refresh Token
	newSession, err := client.RefreshToken(ctx)
	if err != nil {
		t.Fatalf("refresh token failed: %v", err)
	}
	if newSession.AccessToken != "jwt-token-333" {
		t.Errorf("expected jwt-token-333, got %s", newSession.AccessToken)
	}

	// 5. Password Login
	pwdSession, err := client.LoginWithPassword(ctx, "9647801234567", "secret123")
	if err != nil {
		t.Fatalf("password login failed: %v", err)
	}
	if pwdSession.AccessToken != "jwt-token-pwd" {
		t.Errorf("expected jwt-token-pwd, got %s", pwdSession.AccessToken)
	}

	// 6. Reset & Validate password
	if err := client.ResetPassword(ctx, "old", "new"); err != nil {
		t.Errorf("reset password failed: %v", err)
	}
	if err := client.ValidatePassword(ctx, "new"); err != nil {
		t.Errorf("validate password failed: %v", err)
	}

	// 7. Feedback & Logout
	if err := client.SubmitFeedback(ctx, 5, "Excellent app"); err != nil {
		t.Errorf("feedback failed: %v", err)
	}
	if err := client.Logout(ctx); err != nil {
		t.Errorf("logout failed: %v", err)
	}
	accAfter, _ := client.GetTokens()
	if accAfter != "" {
		t.Errorf("expected empty token after logout, got %s", accAfter)
	}
}

func TestProfileAndBalance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/user/profile":
			if r.Method == http.MethodPatch {
				_ = json.NewEncoder(w).Encode(APIResponse[UserResp]{
					Status: "success",
					Data: UserResp{
						Name:  "Ali Updated",
						Email: "ali@example.com",
					},
				})
				return
			}
			_ = json.NewEncoder(w).Encode(APIResponse[UserResp]{
				Status: "success",
				Data: UserResp{
					MSISDN:              "9647801234567",
					Name:                "Ali",
					CustomerBillingType: "PREPAID",
				},
			})
		case "/api/number/wallet":
			_ = json.NewEncoder(w).Encode(APIResponse[BalanceResp]{
				Status: "success",
				Data: BalanceResp{
					Balance: &BalanceDetails{
						Value:  5000,
						Expiry: "2026-12-31",
					},
				},
			})
		case "/api/number/query-bill":
			_ = json.NewEncoder(w).Encode(APIResponse[BillDetailsResp]{
				Status: "success",
				Data: BillDetailsResp{
					TotalBill: 25000,
				},
			})
		case "/api/v2/number/pre-registration-check":
			_ = json.NewEncoder(w).Encode(APIResponse[PreRegisterStatusResp]{
				Status: "success",
				Data: PreRegisterStatusResp{
					AssociatedWithUser: true,
					UnifiedSIMStatus:   "ACTIVE",
				},
			})
		case "/api/number/subaccounts":
			_ = json.NewEncoder(w).Encode(APIResponse[[]SubAccountContent]{
				Status: "success",
				Data: []SubAccountContent{
					{AccountType: "DATA", Amount: 2048},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithMSISDN("9647801234567"))
	ctx := context.Background()

	profile, err := client.GetProfile(ctx)
	if err != nil || profile.CustomerBillingType != "PREPAID" {
		t.Fatalf("get profile failed: %v", err)
	}

	updated, err := client.UpdateProfile(ctx, &UpdateProfileReq{Name: "Ali Updated"})
	if err != nil || updated.Name != "Ali Updated" {
		t.Fatalf("update profile failed: %v", err)
	}

	bal, err := client.GetBalance(ctx)
	if err != nil || bal.Balance.Value != 5000 {
		t.Fatalf("get balance failed: %v", err)
	}

	bill, err := client.GetBillDetails(ctx)
	if err != nil || bill.TotalBill != 25000 {
		t.Fatalf("get bill failed: %v", err)
	}

	preReg, err := client.CheckPreRegistration(ctx)
	if err != nil || !preReg.AssociatedWithUser {
		t.Fatalf("check prereg failed: %v", err)
	}

	subAccs, err := client.GetSubAccounts(ctx)
	if err != nil || len(subAccs) != 1 || subAccs[0].Amount != 2048 {
		t.Fatalf("get subaccounts failed: %v", err)
	}
}

func TestRechargeAndTransfer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/number/charge-voucher":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/number/credit-transfer":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/number/extend-validity":
			_ = json.NewEncoder(w).Encode(APIResponse[ExtendValidityResp]{
				Status: "success",
				Data: ExtendValidityResp{
					Balance: 4000,
				},
			})
		case "/api/number/validity-options":
			_ = json.NewEncoder(w).Encode(APIResponse[ValidityOptionBody]{
				Status: "success",
				Data: ValidityOptionBody{
					ValidityOptions: []ValidityOption{
						{Duration: "30_DAYS", Price: "1000"},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithMSISDN("9647801234567"))
	ctx := context.Background()

	if err := client.RechargeVoucher(ctx, "1234567890123456"); err != nil {
		t.Fatalf("recharge voucher failed: %v", err)
	}

	if err := client.CreditTransfer(ctx, "9647809999999", 5000, "otp-conf-ok"); err != nil {
		t.Fatalf("credit transfer failed: %v", err)
	}

	ext, err := client.ExtendValidity(ctx, 1000)
	if err != nil || ext.Balance != 4000 {
		t.Fatalf("extend validity failed: %v", err)
	}

	opts, err := client.GetValidityOptions(ctx)
	if err != nil || len(opts.ValidityOptions) != 1 {
		t.Fatalf("get validity options failed: %v", err)
	}
}

func TestServicesOffersAndSharing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/offers":
			_ = json.NewEncoder(w).Encode(APIResponse[OfferDetails]{
				Status: "success",
				Data: OfferDetails{
					Offers: []OfferDetailCms{
						{CBSID: "offer-1", Title: &LokaliseText{AR: "باقة 5 جيجا", EN: "5GB Bundle"}},
					},
				},
			})
		case "/api/number/subscriptions":
			_ = json.NewEncoder(w).Encode(APIResponse[[]SubscriptionContent]{
				Status: "success",
				Data: []SubscriptionContent{
					{ID: "sub-1", Status: "ACTIVE"},
				},
			})
		case "/api/number/subscribe", "/api/number/unsubscribe":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/number/claim-daily-gift", "/api/number/send-gift", "/api/number/kafoo_invite":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/number/personalized-offers":
			pt := r.URL.Query().Get("provider_type")
			if pt == "ATL" {
				_ = json.NewEncoder(w).Encode(APIResponse[PersonalizedOffersATL]{
					Status: "success",
					Data: PersonalizedOffersATL{
						Offers: []OfferDetailCms{{CBSID: "atl-1"}},
					},
				})
			} else {
				_ = json.NewEncoder(w).Encode(APIResponse[PersonalizedOffersBTL]{
					Status: "success",
					Data: PersonalizedOffersBTL{
						Offers: []OfferDetailCms{{CBSID: "btl-1"}},
					},
				})
			}
		case "/api/number/flex-limits":
			_ = json.NewEncoder(w).Encode(APIResponse[FlexLimitCms]{
				Status: "success",
				Data:   FlexLimitCms{RemainingPoints: 100, TotalPoints: 500},
			})
		case "/api/number/flex-status":
			_ = json.NewEncoder(w).Encode(APIResponse[FlexStatusCms]{
				Status: "success",
				Data:   FlexStatusCms{Status: "ACTIVE", OfferName: "Flex 500"},
			})
		case "/api/number/migrate-to-flex":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/sharing/query":
			_ = json.NewEncoder(w).Encode(APIResponse[QueryMemberResp]{
				Status: "success",
				Data: QueryMemberResp{
					Members: []SharingMember{{MSISDN: "9647801111111", Role: "MEMBER", Allocated: 1024}},
				},
			})
		case "/api/sharing/addrsc/", "/api/sharing/remove/", "/api/sharing/transfer_unit":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithMSISDN("9647801234567"))
	ctx := context.Background()

	offers, err := client.GetOffers(ctx, "", "", "")
	if err != nil || len(offers.Offers) != 1 {
		t.Fatalf("get offers failed: %v", err)
	}

	subs, err := client.GetSubscriptions(ctx)
	if err != nil || len(subs) != 1 {
		t.Fatalf("get subscriptions failed: %v", err)
	}

	if err := client.SubscribeOffer(ctx, "offer-1"); err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	if err := client.UnsubscribeOffer(ctx, "offer-1"); err != nil {
		t.Fatalf("unsubscribe failed: %v", err)
	}

	atl, err := client.GetPersonalizedOffersATL(ctx)
	if err != nil || len(atl.Offers) != 1 {
		t.Fatalf("get atl failed: %v", err)
	}

	btl, err := client.GetPersonalizedOffersBTL(ctx)
	if err != nil || len(btl.Offers) != 1 {
		t.Fatalf("get btl failed: %v", err)
	}

	if err := client.ClaimDailyGift(ctx); err != nil {
		t.Fatalf("claim daily gift failed: %v", err)
	}
	if err := client.SendBundleGift(ctx, "9647809999999", "offer-1"); err != nil {
		t.Fatalf("send bundle gift failed: %v", err)
	}
	if err := client.InviteToKafoo(ctx, "9647809999999"); err != nil {
		t.Fatalf("invite kafoo failed: %v", err)
	}

	// Flex
	limits, err := client.GetFlexLimits(ctx)
	if err != nil || limits.RemainingPoints != 100 {
		t.Fatalf("flex limits failed: %v", err)
	}
	flexStatus, err := client.GetFlexStatus(ctx)
	if err != nil || flexStatus.Status != "ACTIVE" {
		t.Fatalf("flex status failed: %v", err)
	}
	if err := client.MigrateToFlex(ctx, "flex-plan-1"); err != nil {
		t.Fatalf("migrate flex failed: %v", err)
	}

	// Sharing
	members, err := client.SharingQueryMembers(ctx, "offer-1")
	if err != nil || len(members.Members) != 1 {
		t.Fatalf("query members failed: %v", err)
	}
	if err := client.SharingAddMember(ctx, "9647802222222", "offer-1", 512); err != nil {
		t.Fatalf("add member failed: %v", err)
	}
	if err := client.SharingTransferUnits(ctx, "9647802222222", 256); err != nil {
		t.Fatalf("transfer units failed: %v", err)
	}
	if err := client.SharingRemoveMember(ctx, "9647802222222", "offer-1"); err != nil {
		t.Fatalf("remove member failed: %v", err)
	}
}

func TestLoyaltyAndImtiyaz(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/loyalty/info":
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyResp]{
				Status: "success",
				Data: LoyaltyResp{
					LoyaltyStatus: "ACTIVE",
					TotalPoints:   1500,
					Tier:          "GOLD",
				},
			})
		case "/api/loyalty/hot-bundles":
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyHotBundlesResp]{
				Status: "success",
				Data: LoyaltyHotBundlesResp{
					Bundles: []LoyaltyBundleOfferItem{{ID: "hot-1", Points: 200}},
				},
			})
		case "/api/loyalty/rewards/offers":
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
				return
			}
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyBundleOffers]{
				Status: "success",
				Data: LoyaltyBundleOffers{
					Offers: []LoyaltyBundleOfferItem{{ID: "offer-1", Points: 300}},
				},
			})
		case "/api/loyalty/rewards/promo-codes":
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyPromoCodesResp]{
				Status: "success",
				Data: LoyaltyPromoCodesResp{
					PromoCodes: []PromoCodeItem{{ID: "promo-1", Points: 100}},
				},
			})
		case "/api/v2/loyalty/statement":
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyStatementResp]{
				Status: "success",
				Data: LoyaltyStatementResp{
					Transactions: []LoyaltyTransaction{{ID: "tx-1", Points: 50}},
				},
			})
		case "/api/loyalty/rewards/history":
			_ = json.NewEncoder(w).Encode(APIResponse[RewardsHistoryResp]{
				Status: "success",
				Data: RewardsHistoryResp{
					Rewards: []RewardsHistoryItem{{ID: "rew-1", Points: 50}},
				},
			})
		case "/api/loyalty/faq":
			_ = json.NewEncoder(w).Encode(APIResponse[LoyaltyFAQResp]{
				Status: "success",
				Data: LoyaltyFAQResp{
					FAQs: []LoyaltyFAQItem{{Question: &LokaliseText{AR: "سؤال"}, Answer: &LokaliseText{AR: "جواب"}}},
				},
			})
		case "/api/loyalty/rewards/credit":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/loyalty/rewards/promo-code":
			_ = json.NewEncoder(w).Encode(APIResponse[RedeemPromoCodeResp]{
				Status: "success",
				Data: RedeemPromoCodeResp{
					Code:       "ZAIN-DISCOUNT-50",
					ExpiryDate: "2026-12-31",
				},
			})
		case "/api/v2/imtiyaz/merchant":
			_ = json.NewEncoder(w).Encode(APIResponse[MerchantDetail]{
				Status: "success",
				Data: MerchantDetail{
					ID:   "merch-1",
					Name: &LokaliseText{EN: "Burger King"},
				},
			})
		case "/api/imtiyaz/voucher-history":
			_ = json.NewEncoder(w).Encode(APIResponse[VoucherHistoryBody]{
				Status: "success",
				Data: VoucherHistoryBody{
					History: []VoucherHistoryItem{{VoucherCode: "VOUCHER-123"}},
				},
			})
		case "/api/v2/imtiyaz/reveal":
			_ = json.NewEncoder(w).Encode(APIResponse[RevealVoucherResp]{
				Status: "success",
				Data: RevealVoucherResp{
					VoucherCode: "REVEALED-CODE-999",
				},
			})
		case "/api/imtiyaz/report":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithMSISDN("9647801234567"))
	ctx := context.Background()

	info, err := client.GetLoyaltyInfo(ctx)
	if err != nil || info.TotalPoints != 1500 {
		t.Fatalf("get loyalty info failed: %v", err)
	}

	hot, err := client.GetLoyaltyHotBundles(ctx)
	if err != nil || len(hot.Bundles) != 1 {
		t.Fatalf("get hot bundles failed: %v", err)
	}

	bundles, err := client.GetLoyaltyBundles(ctx)
	if err != nil || len(bundles.Offers) != 1 {
		t.Fatalf("get loyalty bundles failed: %v", err)
	}

	promos, err := client.GetLoyaltyPromoCodes(ctx, "", "", "", "", 10, "")
	if err != nil || len(promos.PromoCodes) != 1 {
		t.Fatalf("get promos failed: %v", err)
	}

	stmt, err := client.GetLoyaltyStatement(ctx, "", "")
	if err != nil || len(stmt.Transactions) != 1 {
		t.Fatalf("get statement failed: %v", err)
	}

	hist, err := client.GetRewardsHistory(ctx)
	if err != nil || len(hist.Rewards) != 1 {
		t.Fatalf("get rewards history failed: %v", err)
	}

	faq, err := client.GetLoyaltyFAQ(ctx)
	if err != nil || len(faq.FAQs) != 1 {
		t.Fatalf("get faq failed: %v", err)
	}

	if err := client.RedeemCredit(ctx, 100); err != nil {
		t.Fatalf("redeem credit failed: %v", err)
	}
	if err := client.RedeemLoyaltyBundle(ctx, "bundle-1"); err != nil {
		t.Fatalf("redeem loyalty bundle failed: %v", err)
	}
	promoResp, err := client.RedeemPromoCode(ctx, "promo-1")
	if err != nil || promoResp.Code != "ZAIN-DISCOUNT-50" {
		t.Fatalf("redeem promo failed: %v", err)
	}

	// Imtiyaz
	merch, err := client.GetMerchant(ctx, "food", "burger_king")
	if err != nil || merch.ID != "merch-1" {
		t.Fatalf("get merchant failed: %v", err)
	}
	vHist, err := client.GetVoucherHistory(ctx)
	if err != nil || len(vHist.History) != 1 {
		t.Fatalf("get voucher history failed: %v", err)
	}
	revealed, err := client.RevealVoucher(ctx, "merch-1", "food")
	if err != nil || revealed.VoucherCode != "REVEALED-CODE-999" {
		t.Fatalf("reveal voucher failed: %v", err)
	}
	if err := client.ReportMerchant(ctx, "merch-1", "closed", "store was closed"); err != nil {
		t.Fatalf("report merchant failed: %v", err)
	}
}

func TestPaymentsAndGateway(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/payment/zaincash/initiate", "/api/v2/payment/zaincash/initiate", "/api/payment/creditcard/initiate", "/api/purchase/order":
			_ = json.NewEncoder(w).Encode(APIResponse[InitiateOrderResp]{
				Status: "success",
				Data: InitiateOrderResp{
					PaymentURL: "https://pay.example.com/checkout/123",
				},
			})
		case "/api/payment/status", "/api/v2/payment/status":
			_ = json.NewEncoder(w).Encode(APIResponse[OrderPaymentStatusResp]{
				Status: "success",
				Data: OrderPaymentStatusResp{
					TransactionID:     "tx-pay-123",
					TransactionStatus: "SUCCESS",
				},
			})
		case "/api/payment/create-checkout-id":
			_ = json.NewEncoder(w).Encode(APIResponse[CheckoutResp]{
				Status: "success",
				Data: CheckoutResp{
					CheckoutID: "chk-uuid-999",
				},
			})
		case "/api/payment/save-card", "/api/payment/set-default-card", "/api/payment/delete-card":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/payment/refresh-payment-status":
			_ = json.NewEncoder(w).Encode(APIResponse[GatewayPaymentStatusResp]{
				Status: "success",
				Data: GatewayPaymentStatusResp{
					PaymentStatus: &GatewayPaymentInfo{
						ID: "pay-info-1",
						Result: &PaymentResult{
							Code:        "000.000.000",
							Description: "Transaction succeeded",
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL))
	ctx := context.Background()

	orderReq := &InitiateOrderReq{
		OrderID:     "order-1001",
		Amount:      10000,
		Currency:    "IQD",
		ServiceType: "RECHARGE",
	}

	zcResp, err := client.InitiateZainCashPayment(ctx, orderReq)
	if err != nil || zcResp.PaymentURL == "" {
		t.Fatalf("initiate zaincash failed: %v", err)
	}

	zcV2, err := client.InitiateZainCashPaymentV2(ctx, orderReq)
	if err != nil || zcV2.PaymentURL == "" {
		t.Fatalf("initiate zaincash v2 failed: %v", err)
	}

	ccResp, err := client.InitiateCreditCardOrder(ctx, orderReq)
	if err != nil || ccResp.PaymentURL == "" {
		t.Fatalf("initiate credit card failed: %v", err)
	}

	purResp, err := client.InitiateOrder(ctx, orderReq)
	if err != nil || purResp.PaymentURL == "" {
		t.Fatalf("initiate order failed: %v", err)
	}

	st, err := client.GetPaymentStatus(ctx, "token-1", "tx-pay-123")
	if err != nil || st.TransactionStatus != "SUCCESS" {
		t.Fatalf("get payment status failed: %v", err)
	}

	st2, err := client.GetPaymentStatusV2(ctx, "token-1", "tx-pay-123")
	if err != nil || st2.TransactionStatus != "SUCCESS" {
		t.Fatalf("get payment status v2 failed: %v", err)
	}

	chk, err := client.CreateCheckoutID(ctx, &CheckoutReq{Amount: "10000"})
	if err != nil || chk.CheckoutID != "chk-uuid-999" {
		t.Fatalf("create checkout id failed: %v", err)
	}

	if err := client.SaveCreditCard(ctx, &SaveCreditCardReq{CardToken: "tok-1"}); err != nil {
		t.Fatalf("save card failed: %v", err)
	}
	if err := client.SetDefaultCard(ctx, "tok-1"); err != nil {
		t.Fatalf("set default card failed: %v", err)
	}
	if err := client.DeleteCard(ctx, "tok-1"); err != nil {
		t.Fatalf("delete card failed: %v", err)
	}

	gwStatus, err := client.RefreshPaymentStatus(ctx, &RefreshPaymentStatusReq{ResourcePath: "/v1/checkouts"})
	if err != nil || gwStatus.PaymentStatus.ID != "pay-info-1" {
		t.Fatalf("refresh payment status failed: %v", err)
	}
}

func TestDashboardAndSupport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/number/dashboard_message":
			_ = json.NewEncoder(w).Encode(APIResponse[DashboardMsgContent]{
				Status: "success",
				Data: DashboardMsgContent{
					Messages: []DashboardMessage{{ID: "msg-1"}},
				},
			})
		case "/api/notifications":
			_ = json.NewEncoder(w).Encode(APIResponse[InAppMSGPageContent]{
				Status: "success",
				Data: InAppMSGPageContent{
					Notifications: []InAppNotification{{ID: "notif-1", Shortname: "n1"}},
					Total:         1,
				},
			})
		case "/api/notifications/read-notification/n1", "/api/notifications/read-all-notifications", "/api/notifications/delete-notification/n1":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/near_me_store/near-me-store":
			_ = json.NewEncoder(w).Encode(APIResponse[ZainStoreNearMeContent]{
				Status: "success",
				Data: ZainStoreNearMeContent{
					Stores: []ZainStore{{ID: "store-1", Name: "Baghdad Al-Mansour"}},
				},
			})
		case "/api/digital_services_catalogue_management":
			_ = json.NewEncoder(w).Encode(APIResponse[DigitalServiceContent]{
				Status: "success",
				Data: DigitalServiceContent{
					Services: []DigitalServiceItem{{ID: "srv-1", Price: 1000}},
				},
			})
		case "/api/referral/count":
			_ = json.NewEncoder(w).Encode(APIResponse[ReferralResp]{
				Status: "success",
				Data:   ReferralResp{Count: 3, Code: "REF123"},
			})
		case "/api/referral/new-app-install":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{Status: "success"})
		case "/api/complaints/create-ticket":
			_ = json.NewEncoder(w).Encode(APIResponse[SubmitTicketResponse]{
				Status: "success",
				Data: SubmitTicketResponse{
					TicketID: "TCK-1001",
					Status:   "OPEN",
				},
			})
		case "/api/complaints/questions":
			_ = json.NewEncoder(w).Encode(APIResponse[QuestionsData]{
				Status: "success",
				Data: QuestionsData{
					Questions: []Questions{{QuestionID: "q1", QuestionEngTxt: "Network issue?"}},
				},
			})
		case "/api/complaints/reopen_reason":
			_ = json.NewEncoder(w).Encode(APIResponse[ReOpenReasonsData]{
				Status: "success",
				Data: ReOpenReasonsData{
					ReOpenReasons: []ReOpenReasons{{ReasonNameEN: "Issue persisted"}},
				},
			})
		case "/api/complaints/summary":
			_ = json.NewEncoder(w).Encode(APIResponse[SummariesData]{
				Status: "success",
				Data: SummariesData{
					Summaries: []Summaries{{SummaryEngTxt: "Data failure"}},
				},
			})
		case "/api/complaints/items":
			_ = json.NewEncoder(w).Encode(APIResponse[TemplatesData]{
				Status: "success",
				Data: TemplatesData{
					Items: []Items{{ItemEngTxt: "4G Issue"}},
				},
			})
		case "/api/complaints/list-ticket":
			_ = json.NewEncoder(w).Encode(APIResponse[TicketsData]{
				Status: "success",
				Data: TicketsData{
					Tickets: []TicketsList{{TxtIssueID: "TCK-1001", Category: "Network"}},
				},
			})
		case "/api/complaints/reopen-ticket":
			_ = json.NewEncoder(w).Encode(APIResponse[ReOpenTicketsResponse]{
				Status: "success",
				Data: ReOpenTicketsResponse{
					TicketID: "TCK-1001",
					Status:   "REOPENED",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithMSISDN("9647801234567"))
	ctx := context.Background()

	dash, err := client.GetDashboardMessages(ctx)
	if err != nil || len(dash.Messages) != 1 {
		t.Fatalf("get dashboard failed: %v", err)
	}

	notifs, err := client.GetNotifications(ctx, 0, 10, nil)
	if err != nil || len(notifs.Notifications) != 1 {
		t.Fatalf("get notifs failed: %v", err)
	}
	if err := client.ReadNotification(ctx, "n1"); err != nil {
		t.Fatalf("read notif failed: %v", err)
	}
	if err := client.ReadAllNotifications(ctx); err != nil {
		t.Fatalf("read all notifs failed: %v", err)
	}
	if err := client.DeleteNotification(ctx, "n1"); err != nil {
		t.Fatalf("delete notif failed: %v", err)
	}

	stores, err := client.GetStoreNearMe(ctx)
	if err != nil || len(stores.Stores) != 1 {
		t.Fatalf("get stores failed: %v", err)
	}

	srvs, err := client.GetDigitalServices(ctx, "", "", "")
	if err != nil || len(srvs.Services) != 1 {
		t.Fatalf("get digital services failed: %v", err)
	}

	refCount, err := client.GetReferralCount(ctx)
	if err != nil || refCount.Count != 3 {
		t.Fatalf("get referral count failed: %v", err)
	}

	if err := client.ReportNewInstall(ctx, "install-123", "REF123"); err != nil {
		t.Fatalf("report install failed: %v", err)
	}

	// Complaints
	tck, err := client.CreateTicket(ctx, &SubmitTicketReq{
		City:        "Baghdad",
		Governorate: "Baghdad",
		Description: "Internet slow",
	})
	if err != nil || tck.TicketID != "TCK-1001" {
		t.Fatalf("create ticket failed: %v", err)
	}

	qs, err := client.GetQuestions(ctx, "sum-1")
	if err != nil || len(qs.Questions) != 1 {
		t.Fatalf("get questions failed: %v", err)
	}

	reasons, err := client.GetReOpenReasons(ctx)
	if err != nil || len(reasons.ReOpenReasons) != 1 {
		t.Fatalf("get reopen reasons failed: %v", err)
	}

	sums, err := client.GetSummaries(ctx, "cat-1")
	if err != nil || len(sums.Summaries) != 1 {
		t.Fatalf("get summaries failed: %v", err)
	}

	temps, err := client.GetTemplates(ctx, "", "", "", "")
	if err != nil || len(temps.Items) != 1 {
		t.Fatalf("get templates failed: %v", err)
	}

	tcks, err := client.GetTickets(ctx, "")
	if err != nil || len(tcks.Tickets) != 1 {
		t.Fatalf("get tickets failed: %v", err)
	}

	reopened, err := client.ReOpenTicket(ctx, &ReOpenTicketsReq{TicketID: "TCK-1001", ReasonID: "r1"})
	if err != nil || reopened.Status != "REOPENED" {
		t.Fatalf("reopen ticket failed: %v", err)
	}
}

func TestCMSContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/dmart/public/excute/query/galleon":
			_ = json.NewEncoder(w).Encode(ContentResolverResp[ConfigCMS]{
				Status: "success",
				Data: ConfigCMS{
					AppVersion: "6.6.0",
				},
			})
		case "/dmart/public/excute/query/products":
			_ = json.NewEncoder(w).Encode(ContentResolverResp[DashboardBannersResp]{
				Status: "success",
				Data: DashboardBannersResp{
					FlexStatus: "ACTIVE",
				},
			})
		case "/api/slideshows/flex_intro":
			_ = json.NewEncoder(w).Encode(APIResponse[SlideShowContent]{
				Status: "success",
				Data: SlideShowContent{
					Slides: []SlideShowItem{{ShortName: "slide_1"}},
				},
			})
		case "/dmart/":
			_ = json.NewEncoder(w).Encode(APIResponse[map[string]any]{
				Status: "success",
				Data:   map[string]any{"healthy": true},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithCMSURL(ts.URL))
	ctx := context.Background()

	cfg, err := client.GetCMSConfiguration(ctx, "{ config }")
	if err != nil || cfg.AppVersion != "6.6.0" {
		t.Fatalf("get cms config failed: %v", err)
	}

	banner, err := client.GetCMSDashboardBanners(ctx, "{ banners }")
	if err != nil || banner.FlexStatus != "ACTIVE" {
		t.Fatalf("get cms banners failed: %v", err)
	}

	slides, err := client.GetSlideshow(ctx, "flex_intro")
	if err != nil || len(slides.Slides) != 1 {
		t.Fatalf("get slideshow failed: %v", err)
	}

	dmart, err := client.GetDmartStatus(ctx)
	if err != nil || dmart["healthy"] != true {
		t.Fatalf("get dmart status failed: %v", err)
	}
}

func TestLokaliseTextString(t *testing.T) {
	var empty *LokaliseText
	if empty.String() != "" {
		t.Errorf("expected empty string for nil LokaliseText")
	}

	arOnly := &LokaliseText{AR: "مرحبا"}
	if arOnly.String() != "مرحبا" {
		t.Errorf("expected مرحبا, got %s", arOnly.String())
	}

	enOnly := &LokaliseText{EN: "Hello"}
	if enOnly.String() != "Hello" {
		t.Errorf("expected Hello, got %s", enOnly.String())
	}

	kdOnly := &LokaliseText{KD: "Slaw"}
	if kdOnly.String() != "Slaw" {
		t.Errorf("expected Slaw, got %s", kdOnly.String())
	}
}

func TestSessionDataExpiry(t *testing.T) {
	var nilSession *SessionData
	if !nilSession.IsExpired() {
		t.Errorf("nil session should be expired")
	}

	expiredSession := &SessionData{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(-1 * time.Hour),
	}
	if !expiredSession.IsExpired() {
		t.Errorf("past session should be expired")
	}

	validSession := &SessionData{
		AccessToken: "test-token",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}
	if validSession.IsExpired() {
		t.Errorf("future session should not be expired")
	}
}

func TestWalletAndIncomingTransferVerification(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/number/wallet":
			_ = json.NewEncoder(w).Encode(APIResponse[BalanceResp]{
				Status: "success",
				Data: BalanceResp{
					Balance: &BalanceDetails{
						Value:  25000,
						Expiry: "2026-12-31",
					},
					ServerTime: "2026-09-11T16:00:00Z",
				},
			})
		case "/api/number/electronic-bill-items":
			_ = json.NewEncoder(w).Encode(APIResponse[ElectronicBillItems]{
				Status: "success",
				Data: ElectronicBillItems{
					Eligible: true,
					Items: []HistoryItem{
						{
							BNumber:         "07801112233",
							ChargeAmount:    "5000",
							StartTime:       "2026-09-11 15:30:00",
							ServiceTypeName: "Credit Transfer",
							ServiceTypeID:   "transfer",
						},
						{
							BNumber:         "07809998877",
							ChargeAmount:    "10000",
							StartTime:       "2026-09-11 16:15:00",
							ServiceTypeName: "Credit Transfer",
							ServiceTypeID:   "transfer",
						},
					},
				},
			})
		case "/api/notifications":
			_ = json.NewEncoder(w).Encode(APIResponse[InAppMSGPageContent]{
				Status: "success",
				Data: InAppMSGPageContent{
					Notifications: []InAppNotification{
						{
							ID:        "notif-1",
							Shortname: "p2p-transfer",
							Title:     &LokaliseText{AR: "تحويل رصيد", EN: "Credit Transfer"},
							Body:      &LokaliseText{AR: "تم استلام 15000 د.ع من الرقم 07805556677", EN: "Received 15000 IQD from 07805556677"},
							CreatedAt: "2026-09-11T16:20:00Z",
						},
					},
					Total: 1,
				},
			})
		case "/api/otp/request":
			var resp OTPRequestIdResp
			resp.Status = "success"
			resp.Data.RequestID = "req-wallet-123"
			_ = json.NewEncoder(w).Encode(resp)
		case "/api/otp/confirm":
			_ = json.NewEncoder(w).Encode(APIResponse[OTPConfirmationIdResp]{
				Status: "success",
				Data: OTPConfirmationIdResp{
					ConfirmationID: "conf-wallet-abc",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(
		WithBaseURL(ts.URL),
		WithMasterWallet("07801234567"),
	)
	ctx := context.Background()

	// 1. Master Wallet verification
	if client.MasterWallet() != "07801234567" {
		t.Fatalf("expected 07801234567, got %s", client.MasterWallet())
	}
	client.SetMasterWallet("07807654321")
	if client.MasterWallet() != "07807654321" {
		t.Fatalf("expected 07807654321, got %s", client.MasterWallet())
	}

	// 2. USSD formatting
	ussd := client.FormatUSSDTransfer("07809998877", 5000)
	expectedUSSD := "*123*5000*07809998877#"
	if ussd != expectedUSSD {
		t.Fatalf("expected %s, got %s", expectedUSSD, ussd)
	}

	// 3. Wallet Balance
	wb, err := client.GetWalletBalance(ctx)
	if err != nil || wb == nil || wb.Balance == nil || wb.Balance.Value != 25000 {
		t.Fatalf("get wallet balance failed: %v", err)
	}

	// 4. Wallet Overview
	overview, err := client.GetWalletOverview(ctx)
	if err != nil || overview == nil || overview.Balance != 25000 {
		t.Fatalf("get wallet overview failed: %v", err)
	}

	// 5. OTP Transfer Flow
	otpReq, err := client.RequestCreditTransferOTP(ctx)
	if err != nil || otpReq.Data.RequestID != "req-wallet-123" {
		t.Fatalf("request transfer otp failed: %v", err)
	}
	otpConf, err := client.ConfirmCreditTransferOTP(ctx, "1234")
	if err != nil || otpConf.ConfirmationID != "conf-wallet-abc" {
		t.Fatalf("confirm transfer otp failed: %v", err)
	}

	// 6. Incoming Transfers list
	transfers, err := client.GetIncomingTransfers(ctx)
	if err != nil {
		t.Fatalf("get incoming transfers failed: %v", err)
	}
	if len(transfers) < 2 {
		t.Fatalf("expected at least 2 transfers, got %d", len(transfers))
	}

	// 7. Automated Incoming Transfer Verification (from ElectronicBillItems)
	verified, rec, err := client.VerifyIncomingTransfer(ctx, "07801112233", 5000)
	if err != nil || !verified || rec == nil {
		t.Fatalf("verify incoming transfer failed: %v", err)
	}
	if rec.Amount != "5000" {
		t.Fatalf("expected amount 5000, got %s", rec.Amount)
	}

	// Suffix phone matching: "+9647801112233" should match "07801112233"
	verifiedSuffix, recSuffix, err := client.VerifyIncomingTransfer(ctx, "+9647801112233", 5000)
	if err != nil || !verifiedSuffix || recSuffix == nil {
		t.Fatalf("suffix matching failed: %v", err)
	}

	// Verification from Notification text extraction
	verifiedNotif, recNotif, err := client.VerifyIncomingTransfer(ctx, "07805556677", 10000)
	if err != nil || !verifiedNotif || recNotif == nil {
		t.Fatalf("notification transfer verify failed: %v", err)
	}
	if recNotif.Amount != "15000" {
		t.Fatalf("expected amount 15000, got %s", recNotif.Amount)
	}

	// Min amount threshold: transfer exists for 5000, but minAmount is 10000
	verifiedLow, _, _ := client.VerifyIncomingTransfer(ctx, "07801112233", 10000)
	if verifiedLow {
		t.Fatalf("expected false when amount is below threshold")
	}

	// Non-existent sender
	verifiedNone, _, _ := client.VerifyIncomingTransfer(ctx, "07800000000", 1000)
	if verifiedNone {
		t.Fatalf("expected false for non-existent sender")
	}

	// In-memory recorded transfer
	client.RecordIncomingTransfer(IncomingTransferRecord{
		MSISDN: "07807778899",
		Amount: "50000",
		Title:  "Manual Test Deposit",
	})
	verifiedMem, recMem, err := client.VerifyIncomingTransfer(ctx, "07807778899", 50000)
	if err != nil || !verifiedMem || recMem == nil {
		t.Fatalf("in-memory verify failed: %v", err)
	}

	// Invalid short phone number
	_, _, errShort := client.VerifyIncomingTransfer(ctx, "12345", 1000)
	if errShort == nil {
		t.Fatalf("expected error for short phone number")
	}
}

func TestNormalizeMSISDNAndPolling(t *testing.T) {
	testCases := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"07801234567", "9647801234567", false},
		{"7801234567", "9647801234567", false},
		{"+9647801234567", "9647801234567", false},
		{"009647801234567", "9647801234567", false},
		{"9647801234567", "9647801234567", false},
		{"٠٧٨٠١٢٣٤٥٦٧", "9647801234567", false},
		{"٩٦٤٧٨٠١٢٣٤٥٦٧", "9647801234567", false},
		{"07701234567", "9647701234567", false},
		{"12345", "", true},
		{"invalid", "", true},
	}

	for _, tc := range testCases {
		res, err := NormalizeMSISDN(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("expected error for %q, got %q", tc.input, res)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tc.input, err)
			}
			if res != tc.expected {
				t.Errorf("NormalizeMSISDN(%q) = %q, expected %q", tc.input, res, tc.expected)
			}
		}
	}

	// FormatLocalMSISDN tests
	if local := FormatLocalMSISDN("9647801234567"); local != "07801234567" {
		t.Errorf("FormatLocalMSISDN failed: expected 07801234567, got %s", local)
	}
	if local := FormatLocalMSISDN("7801234567"); local != "07801234567" {
		t.Errorf("FormatLocalMSISDN failed: expected 07801234567, got %s", local)
	}

	// Test WaitForIncomingTransfer immediate hit with in-memory record
	client := NewClient()
	client.RecordIncomingTransfer(IncomingTransferRecord{
		MSISDN: "9647809988776",
		Amount: "25000",
		Title:  "Immediate Transfer",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rec, err := client.WaitForIncomingTransfer(ctx, "07809988776", 25000, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("WaitForIncomingTransfer failed: %v", err)
	}
	if rec.Amount != "25000" {
		t.Errorf("expected 25000, got %s", rec.Amount)
	}

	// Test timeout when transfer doesn't exist
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer timeoutCancel()

	_, errTimeout := client.WaitForIncomingTransfer(timeoutCtx, "07800000000", 10000, 50*time.Millisecond)
	if errTimeout == nil {
		t.Errorf("expected timeout error, got nil")
	}
}

func TestSystemAndSIMMethods(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/system/time":
			_ = json.NewEncoder(w).Encode(APIResponse[ServerTimeResp]{
				Status: "success",
				Data: ServerTimeResp{
					ServerTime:  "2026-09-12T05:39:22+03:00",
					Timestamp:   1789180762,
					TimestampMS: 1789180762000,
					Timezone:    "Asia/Baghdad",
					UTCOffset:   "+03:00",
				},
			})
		case "/api/number/friends-and-family":
			_ = json.NewEncoder(w).Encode(APIResponse[FriendsAndFamilyResp]{
				Status: "success",
				Data: FriendsAndFamilyResp{
					MaxAllowed: 5,
					Numbers:    []string{"07801234567", "07809876543"},
				},
			})
		case "/api/number/friends-and-family/add":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{
				Status: "success",
			})
		case "/api/number/friends-and-family/remove":
			_ = json.NewEncoder(w).Encode(APIResponse[any]{
				Status: "success",
			})
		case "/api/number/esim-details":
			_ = json.NewEncoder(w).Encode(APIResponse[ESIMDetailsResp]{
				Status: "success",
				Data: ESIMDetailsResp{
					ICCID:          "8996401234567890123",
					MatchingID:     "MATCH-12345",
					QRCodeData:     "LPA:1$smdp.zain.iq$MATCH-12345",
					ActivationCode: "ACT-9988",
					Status:         "ACTIVE",
				},
			})
		case "/api/number/swap-sim":
			_ = json.NewEncoder(w).Encode(APIResponse[SIMSwapResp]{
				Status: "success",
				Data: SIMSwapResp{
					RequestID: "REQ-SWAP-123",
					Status:    "PENDING",
					Message:   "SIM swap request submitted successfully",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := NewClient(WithBaseURL(ts.URL), WithCMSURL(ts.URL))
	ctx := context.Background()

	// 1. Test GetServerTime
	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		t.Fatalf("GetServerTime failed: %v", err)
	}
	if serverTime.Timezone != "Asia/Baghdad" || serverTime.Timestamp != 1789180762 {
		t.Errorf("unexpected server time: %+v", serverTime)
	}

	// 2. Test Friends & Family
	fnf, err := client.GetFriendsAndFamily(ctx)
	if err != nil || len(fnf.Numbers) != 2 {
		t.Fatalf("GetFriendsAndFamily failed: %v", err)
	}
	if err := client.AddFriendsAndFamily(ctx, "07801112233"); err != nil {
		t.Fatalf("AddFriendsAndFamily failed: %v", err)
	}
	if err := client.RemoveFriendsAndFamily(ctx, "07801112233"); err != nil {
		t.Fatalf("RemoveFriendsAndFamily failed: %v", err)
	}

	// 3. Test eSIM Details & SIM Swap
	esim, err := client.GetESIMDetails(ctx)
	if err != nil || esim.ICCID != "8996401234567890123" {
		t.Fatalf("GetESIMDetails failed: %v", err)
	}
	swap, err := client.RequestSIMSwap(ctx, "8996409988776655443")
	if err != nil || swap.RequestID != "REQ-SWAP-123" {
		t.Fatalf("RequestSIMSwap failed: %v", err)
	}
}

func TestParseTransferSMS(t *testing.T) {
	cases := []struct {
		name        string
		sms         string
		expectPhone string
		expectAmt   string
		shouldErr   bool
	}{
		{
			name:        "Standard Zain Arabic SMS",
			sms:         "تم استلام رصيد بقيمة 5,000 د.ع من الرقم 07801234567 بنجاح",
			expectPhone: "07801234567",
			expectAmt:   "5000",
			shouldErr:   false,
		},
		{
			name:        "Eastern Arabic Numerals SMS",
			sms:         "تم تحويل مبلغ ٥٠٠٠ دينار من الرقم ٠٧٨٠٩٨٧٦٥٤٣",
			expectPhone: "07809876543",
			expectAmt:   "5000",
			shouldErr:   false,
		},
		{
			name:        "International format Zain SMS",
			sms:         "You have received 10,000 IQD credit from 9647805554433",
			expectPhone: "9647805554433",
			expectAmt:   "10000",
			shouldErr:   false,
		},
		{
			name:      "Irrelevant spam SMS",
			sms:       "عزيزي المشترك، اشترك الآن في باقة الإنترنت اليومية",
			shouldErr: true,
		},
	}

	client := NewClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec, err := ParseTransferSMS(tc.sms)
			if tc.shouldErr {
				if err == nil {
					t.Fatalf("expected error for %s, got nil", tc.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", tc.name, err)
			}
			if rec.MSISDN != tc.expectPhone || rec.Amount != tc.expectAmt {
				t.Errorf("expected phone %s amt %s, got phone %s amt %s", tc.expectPhone, tc.expectAmt, rec.MSISDN, rec.Amount)
			}

			// Test client recording and matching
			_, err = client.RecordIncomingTransferFromSMS(tc.sms)
			if err != nil {
				t.Fatalf("RecordIncomingTransferFromSMS failed: %v", err)
			}
			ok, matched, err := client.VerifyIncomingTransfer(context.Background(), tc.expectPhone, 1000)
			if err != nil || !ok || matched == nil {
				t.Fatalf("VerifyIncomingTransfer failed to match recorded SMS: %v", err)
			}
		})
	}
}



