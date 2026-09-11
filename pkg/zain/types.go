package zain

import (
	"encoding/json"
	"time"
)

type APIResponse[T any] struct {
	Status string         `json:"status"`
	Data   T              `json:"data,omitempty"`
	Error  *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Type             string        `json:"type,omitempty"`
	Code             int           `json:"code,omitempty"`
	Message          string        `json:"message,omitempty"`
	LocalizedMessage *LokaliseText `json:"localized_message,omitempty"`
	CorrelationID    string        `json:"correlation_id,omitempty"`
	ServiceName      string        `json:"service_name,omitempty"`
	Topic            string        `json:"topic,omitempty"`
	APIPath          string        `json:"api_path,omitempty"`
}

type LokaliseText struct {
	AR string `json:"ar,omitempty"`
	EN string `json:"en,omitempty"`
	KD string `json:"kd,omitempty"`
}

func (l *LokaliseText) String() string {
	if l == nil {
		return ""
	}
	if l.AR != "" {
		return l.AR
	}
	if l.EN != "" {
		return l.EN
	}
	return l.KD
}

type OTPRequestReq struct {
	MSISDN string `json:"msisdn"`
}

type OTPRequestIdResp struct {
	Status string `json:"status,omitempty"`
	Data   struct {
		RequestID string `json:"request_id,omitempty"`
	} `json:"data,omitempty"`
}

type OTPConfirmReq struct {
	MSISDN  string `json:"msisdn"`
	OTPCode string `json:"code"`
}

type OTPConfirmationIdResp struct {
	ConfirmationID string `json:"otp_confirmation"`
}

type OTPSessionSignReq struct {
	MSISDN       string `json:"msisdn"`
	Confirmation string `json:"otp_confirmation"`
	UserSpace    string `json:"user_space"`
}

type SignInReq struct {
	MSISDN   string `json:"msisdn"`
	Password string `json:"password"`
}

type LoginSession struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

type UserInfoReq struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdateProfileReq struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Language string `json:"language,omitempty"`
}

type ResetPwdRequest struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password"`
}

type ValidatePwdReq struct {
	Password string `json:"password"`
}

type FeedbackReq struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment,omitempty"`
}

type UserResp struct {
	ID                  string `json:"id,omitempty"`
	MSISDN              string `json:"msisdn"`
	Name                string `json:"name,omitempty"`
	Email               string `json:"email,omitempty"`
	PicURL              string `json:"profile_pic_url,omitempty"`
	PrimaryOfferingID   string `json:"primary_offering_id,omitempty"`
	CustomerBillingType string `json:"customer_billing_type"`
	CustomerType        string `json:"customer_type,omitempty"`
	UnifiedSIMStatus    string `json:"unified_sim_status,omitempty"`
	Is4GCompatible      bool   `json:"is_4g_compatible"`
	IsDebugEnabled      bool   `json:"is_debug_enabled"`
	FlexStatus          string `json:"flex_status,omitempty"`
	LoyaltyTier         string `json:"loyalty_tier,omitempty"`
	LoyaltyStatus       string `json:"loyalty_status,omitempty"`
	ImtiyazEligible     bool   `json:"imtiyaz_eligible"`
	IsGiftRedeemed      bool   `json:"is_gift_redeemed"`
	CreatedAt           string `json:"created_at,omitempty"`
	NextSubscribeTime   string `json:"next_subscribe_enabled,omitempty"`
}

type BalanceResp struct {
	Balance    *BalanceDetails `json:"balance,omitempty"`
	Loan       *BalanceDetails `json:"loan,omitempty"`
	ServerTime string          `json:"serverTime,omitempty"`
}

type BalanceDetails struct {
	Value  int64  `json:"value"`
	Expiry string `json:"expiry,omitempty"`
}

type BillDetailsResp struct {
	AdvancePayment float64 `json:"advancePayment"`
	PastDue        float64 `json:"pastDue"`
	TotalBill      float64 `json:"totalBill"`
	UnbilledAmount float64 `json:"unbilledAmount"`
	IsDunning      bool    `json:"isDunning"`
}

type PreRegisterStatusResp struct {
	AssociatedWithUser bool   `json:"associated_with_user"`
	UnifiedSIMStatus   string `json:"unified_sim_status"`
}

type SubAccountContent struct {
	AccountType   string         `json:"account_type"`
	Amount        int64          `json:"amount"`
	InitialAmount int64          `json:"initial_amount"`
	ExpiryDate    string         `json:"expiry_date,omitempty"`
	CMS           *SubAccountCms `json:"cms,omitempty"`
}

type SubAccountCms struct {
	AccountType string        `json:"account_type,omitempty"`
	Title       *LokaliseText `json:"title,omitempty"`
	Unit        string        `json:"unit,omitempty"`
}

type ValidityOptionBody struct {
	ValidityOptions []ValidityOption `json:"validity_options"`
}

type ValidityOption struct {
	Duration string `json:"duration"`
	Price    string `json:"price"`
}

type ExtendValidityReq struct {
	MSISDN string `json:"msisdn"`
	Amount int64  `json:"amount"`
}

type ExtendValidityResp struct {
	NewActiveStop  string `json:"newActiveStop,omitempty"`
	NewSuspendStop string `json:"newSuspendStop,omitempty"`
	NewDisableStop string `json:"newDisableStop,omitempty"`
	Balance        int64  `json:"balance"`
}

type CreditTransferReq struct {
	Sender          string `json:"sender"`
	Recipient       string `json:"recipient"`
	Amount          int64  `json:"amount"`
	OTPConfirmation string `json:"otp_confirmation"`
}

type VoucherRechargeReq struct {
	MSISDN  string `json:"msisdn,omitempty"`
	PinCode string `json:"pincode"`
}

type OfferDetails struct {
	Offers []OfferDetailCms `json:"offers"`
}

type OfferDetailCms struct {
	CBSID                 string         `json:"cbs_id"`
	CRMID                 string         `json:"crm_id,omitempty"`
	TOMSID                string         `json:"toms_id,omitempty"`
	Channel               string         `json:"channel,omitempty"`
	Title                 *LokaliseText  `json:"title,omitempty"`
	Desc                  *LokaliseText  `json:"desc,omitempty"`
	ShortDescription      *LokaliseText  `json:"short_description,omitempty"`
	OfferTypeLocalization *LokaliseText  `json:"offer_type_localization,omitempty"`
	OfferTagLocalization  *LokaliseText  `json:"offer_tag_localization,omitempty"`
	Specs                 *OfferSpecs    `json:"specs,omitempty"`
	Services              []OfferService `json:"services,omitempty"`
	Subscribed            bool           `json:"subscribed"`
	IsHotBundle           bool           `json:"is_hot_bundle"`
	SharedUsers           int            `json:"shared_users,omitempty"`
	SharedCredit          string         `json:"shared_credit_account_id,omitempty"`
	CUGAccount            string         `json:"cug_account_id,omitempty"`
}

type OfferSpecs struct {
	Price          int64       `json:"price"`
	Validity       int         `json:"validity"`
	CanUnsubscribe string      `json:"can_unsubscribe,omitempty"`
	SpecialOffer   interface{} `json:"special_offer,omitempty"`
	Giftable       bool        `json:"giftable"`
}

type OfferService struct {
	Type               string        `json:"type"`
	Value              int64         `json:"value,omitempty"`
	Data               int64         `json:"data,omitempty"`
	Voice              int64         `json:"voice,omitempty"`
	SMS                int64         `json:"sms,omitempty"`
	OnNet              int64         `json:"on_net,omitempty"`
	OffNet             int64         `json:"off_net,omitempty"`
	CrossNet           int64         `json:"cross_net,omitempty"`
	UL                 string        `json:"ul,omitempty"`
	Restriction        *LokaliseText `json:"restriction,omitempty"`
	FlexPoints         string        `json:"flex_points,omitempty"`
	LoyaltyPointAmount string        `json:"loyalty_point_amount,omitempty"`
}

type SubscriptionContent struct {
	ID         string          `json:"id"`
	Status     string          `json:"status,omitempty"`
	ExpireTime string          `json:"expire_time,omitempty"`
	CMS        *OfferDetailCms `json:"cms,omitempty"`
}

type OfferSubUnSubReq struct {
	MSISDN  string `json:"msisdn"`
	OfferID string `json:"offer_id"`
	Action  string `json:"action,omitempty"`
}

type PersonalizedOffersATL struct {
	Offers []OfferDetailCms `json:"offers"`
}

type PersonalizedOffersBTL struct {
	Offers []OfferDetailCms `json:"offers"`
}

type NumberReq struct {
	MSISDN string `json:"msisdn"`
}

type SendInviteReq struct {
	MSISDN  string `json:"msisdn"`
	Invited string `json:"invited_msisdn"`
}

type OfferRegistrationReq struct {
	MSISDN  string `json:"msisdn"`
	OfferID string `json:"offer_id"`
}

type OfferResp struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type SendBundleGiftReq struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	OfferID   string `json:"offer_id"`
}

type MigrateFlexReq struct {
	MSISDN      string `json:"msisdn"`
	TargetOffer string `json:"target_offer"`
}

type FlexLimitCms struct {
	RemainingPoints int64  `json:"remaining_points"`
	TotalPoints     int64  `json:"total_points"`
	ResetDate       string `json:"reset_date,omitempty"`
}

type FlexStatusCms struct {
	Status      string `json:"status"`
	OfferName   string `json:"offer_name,omitempty"`
	ActiveUntil string `json:"active_until,omitempty"`
}

type FlexUpsellCms struct {
	OfferID string        `json:"offer_id"`
	Title   *LokaliseText `json:"title,omitempty"`
	Price   int64         `json:"price"`
	Points  int64         `json:"points"`
}

type QueryMemberResp struct {
	Members []SharingMember `json:"members"`
}

type SharingMember struct {
	MSISDN    string `json:"msisdn"`
	Role      string `json:"role"`
	Allocated int64  `json:"allocated,omitempty"`
	Consumed  int64  `json:"consumed,omitempty"`
}

type SharingAddMemberReq struct {
	OwnerMSISDN  string `json:"owner_msisdn"`
	MemberMSISDN string `json:"member_msisdn"`
	OfferID      string `json:"offer_id"`
	Quota        int64  `json:"quota,omitempty"`
}

type SharingRemoveMemberReq struct {
	OwnerMSISDN  string `json:"owner_msisdn"`
	MemberMSISDN string `json:"member_msisdn"`
	OfferID      string `json:"offer_id"`
}

type SharingTransferMemberReq struct {
	OwnerMSISDN  string `json:"owner_msisdn"`
	MemberMSISDN string `json:"member_msisdn"`
	Units        int64  `json:"units"`
}

type LoyaltyResp struct {
	LoyaltyStatus                  string           `json:"loyalty_status"`
	ImtiyazEligible                bool             `json:"imtiyaz_eligible"`
	Tier                           string           `json:"tier,omitempty"`
	LocalizedTier                  *LokaliseText    `json:"localized_tier,omitempty"`
	NextTier                       string           `json:"next_tier,omitempty"`
	LocalizedNextTier              *LokaliseText    `json:"localized_next_tier,omitempty"`
	TotalPoints                    int              `json:"total_points"`
	TotalSpendablePoints           int              `json:"total_spendable_points"`
	PointsToNextTier               int              `json:"points_to_next_tier"`
	TotalPointsToNextTier          int              `json:"total_points_to_next_tier"`
	EarliestSpendablePointValidity string           `json:"earliest_spendable_point_validity,omitempty"`
	SpendablePoints                []SpendablePoint `json:"spendable_points,omitempty"`
}

type SpendablePoint struct {
	Points int    `json:"points"`
	Expiry string `json:"expiry,omitempty"`
}

type LoyaltyHotBundlesResp struct {
	Bundles []LoyaltyBundleOfferItem `json:"bundles"`
}

type LoyaltyBundleOffers struct {
	Offers []LoyaltyBundleOfferItem `json:"offers"`
}

type LoyaltyBundleOfferItem struct {
	ID          string        `json:"id"`
	Points      int           `json:"points"`
	Title       *LokaliseText `json:"title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
	Validity    int           `json:"validity"`
}

type LoyaltyPromoCodesResp struct {
	PromoCodes []PromoCodeItem `json:"promo_codes"`
}

type PromoCodeItem struct {
	ID          string        `json:"id"`
	Title       *LokaliseText `json:"title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
	Points      int           `json:"points"`
	Merchant    string        `json:"merchant,omitempty"`
	Category    string        `json:"category,omitempty"`
}

type LoyaltyStatementResp struct {
	Transactions []LoyaltyTransaction `json:"transactions"`
}

type LoyaltyTransaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	Type        string `json:"type"`
}

type RewardsHistoryResp struct {
	Rewards []RewardsHistoryItem `json:"rewards"`
}

type RewardsHistoryItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Points int    `json:"points"`
	Status string `json:"status"`
}

type LoyaltyFAQResp struct {
	FAQs []LoyaltyFAQItem `json:"faqs"`
}

type LoyaltyFAQItem struct {
	Question *LokaliseText `json:"question"`
	Answer   *LokaliseText `json:"answer"`
}

type RedeemCreditReq struct {
	MSISDN string `json:"msisdn"`
	Points int    `json:"points"`
}

type RedeemBundleReq struct {
	MSISDN   string `json:"msisdn"`
	BundleID string `json:"bundle_id"`
}

type RedeemPromoCodeReq struct {
	MSISDN      string `json:"msisdn"`
	PromoCodeID string `json:"promo_code_id"`
}

type RedeemPromoCodeResp struct {
	Code       string `json:"code"`
	ExpiryDate string `json:"expiry_date,omitempty"`
}

type MerchantDetail struct {
	ID          string        `json:"id"`
	Name        *LokaliseText `json:"name,omitempty"`
	Category    string        `json:"category,omitempty"`
	LogoURL     string        `json:"logo_url,omitempty"`
	Discount    string        `json:"discount,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
	Branches    []BranchInfo  `json:"branches,omitempty"`
}

type BranchInfo struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

type VoucherHistoryBody struct {
	History []VoucherHistoryItem `json:"history"`
}

type VoucherHistoryItem struct {
	VoucherCode string `json:"voucher_code"`
	Merchant    string `json:"merchant"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type RevealVoucherReq struct {
	MerchantID string `json:"merchant_id"`
	Category   string `json:"category"`
}

type RevealVoucherResp struct {
	VoucherCode string `json:"voucher_code"`
	ExpiryDate  string `json:"expiry_date"`
}

type ReportMerchantReq struct {
	MerchantID string `json:"merchant_id"`
	Reason     string `json:"reason"`
	Comment    string `json:"comment,omitempty"`
}

type InitiatePaymentReq struct {
	MSISDN      string  `json:"msisdn"`
	Amount      float64 `json:"amount"`
	ServiceType string  `json:"service_type"`
	ReferenceID string  `json:"reference_id,omitempty"`
}

type InitiatePaymentResp struct {
	TransactionID string `json:"transaction_id"`
	PaymentURL    string `json:"payment_url,omitempty"`
	Token         string `json:"token,omitempty"`
	Status        string `json:"status,omitempty"`
}

type PaymentStatusResp struct {
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Timestamp     string  `json:"timestamp,omitempty"`
}

type InitiateOrderReq struct {
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	ServiceType string  `json:"service_type"`
	Recipient   string  `json:"recipient,omitempty"`
	Lang        string  `json:"lang,omitempty"`
}

type InitiateOrderResp struct {
	PaymentURL string `json:"payment_url"`
}

type OrderPaymentStatusResp struct {
	TransactionID     string        `json:"transaction_id,omitempty"`
	TransactionStatus string        `json:"transaction_status,omitempty"`
	OrderID           string        `json:"order_id,omitempty"`
	Reason            *LokaliseText `json:"reason,omitempty"`
}

type CheckoutReq struct {
	Amount           string `json:"amount"`
	Type             string `json:"payment_type,omitempty"`
	MSISDN           string `json:"msisdn,omitempty"`
	IsRegisteredCard bool   `json:"is_registered_card"`
	RegisterCard     bool   `json:"register_card"`
}

type CheckoutResp struct {
	CheckoutID string `json:"checkout_id"`
}

type CardTokenReq struct {
	CardRegistrationToken string `json:"card_registration_token,omitempty"`
	SpaceName             string `json:"space_name"`
}

type SaveCreditCardReq struct {
	CardToken       string `json:"card_registration_token,omitempty"`
	CardLastDigits  string `json:"card_last_4_digits,omitempty"`
	CardExpiryMonth string `json:"card_expiry_month,omitempty"`
	CardExpiryYear  string `json:"card_expiry_year,omitempty"`
	CardHolder      string `json:"card_holder,omitempty"`
	CardBrand       string `json:"card_brand,omitempty"`
	SpaceName       string `json:"space_name"`
}

type RefreshPaymentStatusReq struct {
	ResourcePath string `json:"resource_path,omitempty"`
	IsPayment    bool   `json:"is_payment"`
}

type GatewayPaymentStatusResp struct {
	PaymentStatus *GatewayPaymentInfo `json:"payment_status,omitempty"`
}

type GatewayPaymentInfo struct {
	ID             string         `json:"id,omitempty"`
	RegistrationID string         `json:"registrationId,omitempty"`
	PaymentBrand   string         `json:"paymentBrand,omitempty"`
	Result         *PaymentResult `json:"result,omitempty"`
	Card           *CardInfo      `json:"card,omitempty"`
}

type PaymentResult struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type CardInfo struct {
	Bin         string `json:"bin,omitempty"`
	Last4Digits string `json:"last4Digits,omitempty"`
	Holder      string `json:"holder,omitempty"`
	ExpiryMonth string `json:"expiryMonth,omitempty"`
	ExpiryYear  string `json:"expiryYear,omitempty"`
}

type DashboardMsgContent struct {
	Messages []DashboardMessage `json:"messages"`
}

type DashboardMessage struct {
	ID          string        `json:"id"`
	Title       *LokaliseText `json:"title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
	ActionURL   string        `json:"action_url,omitempty"`
	Severity    string        `json:"severity,omitempty"`
}

type InAppMSGPageContent struct {
	Notifications []InAppNotification `json:"notifications"`
	Total         int                 `json:"total"`
}

type InAppNotification struct {
	ID        string        `json:"id"`
	Shortname string        `json:"shortname"`
	Title     *LokaliseText `json:"title,omitempty"`
	Body      *LokaliseText `json:"body,omitempty"`
	Read      bool          `json:"read"`
	CreatedAt string        `json:"created_at"`
}

type ZainStoreNearMeContent struct {
	Stores []ZainStore `json:"stores"`
}

type ZainStore struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Phone     string  `json:"phone,omitempty"`
	IsOpen    bool    `json:"is_open"`
}

type DigitalServiceContent struct {
	Services []DigitalServiceItem `json:"services"`
}

type DigitalServiceItem struct {
	ID          string        `json:"id"`
	Title       *LokaliseText `json:"title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
	Price       float64       `json:"price"`
	Category    string        `json:"category"`
	IconURL     string        `json:"icon_url,omitempty"`
}

type ReferralResp struct {
	Count int    `json:"count"`
	Code  string `json:"referral_code,omitempty"`
}

type NewInstallReq struct {
	InstallationID string `json:"installation_id"`
	ReferralCode   string `json:"referral_code,omitempty"`
}

type SubmitTicketReq struct {
	City            string       `json:"city"`
	TicketLanguage  string       `json:"ticket_language"`
	MSISDN          string       `json:"msisdn"`
	Source          string       `json:"source"`
	Governorate     string       `json:"governorate"`
	Description     string       `json:"description"`
	SummaryID       string       `json:"summary_id"`
	QuestionAnswers string       `json:"question_answers"`
	Attachments     []Attachment `json:"attachments,omitempty"`
}

type SubmitTicketResponse struct {
	TicketID  string `json:"ticket_id"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type QuestionsData struct {
	Questions []Questions `json:"questions"`
}

type Questions struct {
	QuestionID         string        `json:"question_id,omitempty"`
	QuestionEngTxt     string        `json:"question_english_txt,omitempty"`
	TxtQuestion        *LokaliseText `json:"txt_question,omitempty"`
	QuestionsMandatory string        `json:"questions_mandatory,omitempty"`
	TxtQuestionType    string        `json:"txt_question_type,omitempty"`
	MCQType            []Mcq         `json:"mcq,omitempty"`
	Answer             string        `json:"answer,omitempty"`
}

type Mcq struct {
	Choice *LokaliseText `json:"choice,omitempty"`
	Status string        `json:"status,omitempty"`
}

type SummariesData struct {
	Summaries []Summaries `json:"summary"`
}

type Summaries struct {
	SummaryEngTxt   string        `json:"summary_english_txt,omitempty"`
	TemplateID      string        `json:"template_id,omitempty"`
	TextSummary     *LokaliseText `json:"text_summary,omitempty"`
	TextSummaryCode string        `json:"text_summary_code,omitempty"`
}

type TemplatesData struct {
	Items []Items `json:"items"`
}

type Items struct {
	ItemEngTxt         string        `json:"item_english_txt,omitempty"`
	Item               *LokaliseText `json:"item,omitempty"`
	CategorizationCode string        `json:"categorization_code,omitempty"`
	Summary            string        `json:"summary,omitempty"`
	SummariesCode      string        `json:"summaries_code,omitempty"`
}

type ReOpenReasonsData struct {
	ReOpenReasons []ReOpenReasons `json:"response_reason"`
}

type ReOpenReasons struct {
	ReasonNameAR string `json:"ar"`
	ReasonNameEN string `json:"en"`
	ReasonNameKD string `json:"kd"`
}

type TicketsData struct {
	Tickets []TicketsList `json:"list_ticket"`
}

type TicketsList struct {
	TxtIssueID    string          `json:"txt_issue_id,omitempty"`
	IssueStatus   *LokaliseText   `json:"issue_status,omitempty"`
	MSISDN        string          `json:"msisdn,omitempty"`
	Category      string          `json:"category,omitempty"`
	ProblemType   string          `json:"problem_type,omitempty"`
	CreateTime    string          `json:"create_time,omitempty"`
	TicketDetails []TicketDetails `json:"ticket_details"`
}

type TicketDetails struct {
	IssueID     string        `json:"txt_issue_id,omitempty"`
	Item        *LokaliseText `json:"txt_item,omitempty"`
	Category    string        `json:"txt_category,omitempty"`
	ProblemType string        `json:"txt_problem_type,omitempty"`
	Summary     *LokaliseText `json:"txt_summary,omitempty"`
	TicketType  string        `json:"txt_ticket_type,omitempty"`
	Description string        `json:"txt_description,omitempty"`
	IssueStatus string        `json:"ddl_issue_status,omitempty"`
	Governorate string        `json:"ronly_governorate,omitempty"`
	City        string        `json:"ronly_city,omitempty"`
	Attachments []Attachment  `json:"attachments"`
	Questions   []Questions   `json:"questions"`
}

type Attachment struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
	Size string `json:"size,omitempty"`
}

type ReOpenTicketsReq struct {
	TicketID string `json:"ticket_id"`
	ReasonID string `json:"reason_id"`
	Comment  string `json:"comment,omitempty"`
}

type ReOpenTicketsResponse struct {
	TicketID string `json:"ticket_id"`
	Status   string `json:"status"`
}

type ElectronicBillItems struct {
	Eligible bool          `json:"eligible"`
	OfferID  *int          `json:"offer_id,omitempty"`
	Items    []HistoryItem `json:"items,omitempty"`
}

type HistoryItem struct {
	StartTime       string `json:"start_time,omitempty"`
	ServiceTypeID   string `json:"service_type_id,omitempty"`
	ServiceTypeName string `json:"service_type_name,omitempty"`
	RatingVolume    string `json:"rating_volume,omitempty"`
	ChargeAmount    string `json:"charge_amount,omitempty"`
	BNumber         string `json:"b_number,omitempty"`
	Unit            string `json:"unit,omitempty"`
}

type DashboardBannersResp struct {
	Message     *DashboardMsg       `json:"message,omitempty"`
	Link        *DashboardBannerLink `json:"link,omitempty"`
	Ordinal     *int                `json:"ordinal,omitempty"`
	EnabledApps []string            `json:"enabled_apps,omitempty"`
	FlexStatus  string              `json:"flex_status,omitempty"`
	Deeplink    *DeeplinkModel      `json:"deep_link,omitempty"`
}

type DashboardMsg struct {
	Title       *LokaliseText `json:"title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
}

type DashboardBannerLink struct {
	URL  string `json:"url,omitempty"`
	Type string `json:"type,omitempty"`
}

type DeeplinkModel struct {
	URL string `json:"url,omitempty"`
}

type DashboardMsgResp struct {
	Message  *DashboardMsg      `json:"message,omitempty"`
	Link     *DashboardLink     `json:"link,omitempty"`
	LinkText *DashboardLinkText `json:"link_text,omitempty"`
}

type DashboardLink struct {
	URL string `json:"url,omitempty"`
}

type DashboardLinkText struct {
	Text *LokaliseText `json:"text,omitempty"`
}

type DashboardGameResp struct {
	GameURL  string        `json:"game_url,omitempty"`
	Title    *LokaliseText `json:"title,omitempty"`
	Subtitle *LokaliseText `json:"subtitle,omitempty"`
}

type InAppMSGCms struct {
	Title    *LokaliseText `json:"title,omitempty"`
	Subtitle *LokaliseText `json:"subtitle,omitempty"`
	Content  *LokaliseText `json:"content,omitempty"`
}

type SlideShowContent struct {
	Slides    []SlideShowItem `json:"slides"`
	CtaText   *LokaliseText   `json:"cta_text,omitempty"`
	CtaAction string          `json:"cta_action,omitempty"`
}

type SlideShowItem struct {
	ShortName   string        `json:"shortname,omitempty"`
	Headline    *LokaliseText `json:"headline,omitempty"`
	MainTitle   *LokaliseText `json:"main_title,omitempty"`
	Description *LokaliseText `json:"description,omitempty"`
}

type ContentRequest struct {
	Query string                 `json:"query"`
	Args  map[string]interface{} `json:"args,omitempty"`
}

type ContentResolverResp[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type RecordResp[T any] struct {
	ShortName   string         `json:"shortname"`
	PayloadResp PayloadResp[T] `json:"attributes"`
	Attachments *Attachments   `json:"attachments,omitempty"`
}

type PayloadResp[T any] struct {
	ID   string      `json:"id,omitempty"`
	Body BodyResp[T] `json:"body"`
}

type BodyResp[T any] struct {
	Data T `json:"data"`
}

type Attachments struct {
	Media []PayloadMedia `json:"media,omitempty"`
}

type PayloadMedia struct {
	URL string `json:"url,omitempty"`
}

type ConfigCMS struct {
	AppVersion     string            `json:"app_version"`
	Maintenance    bool              `json:"maintenance_mode"`
	Features       map[string]bool   `json:"features"`
	SupportContact map[string]string `json:"support_contact"`
}

type SessionData struct {
	AccessToken  string
	RefreshToken string
	MSISDN       string
	ExpiresAt    time.Time
}

func (s *SessionData) IsExpired() bool {
	if s == nil || s.AccessToken == "" {
		return true
	}
	return time.Now().After(s.ExpiresAt)
}

func (r *APIResponse[T]) UnmarshalJSON(data []byte) error {
	type Alias APIResponse[T]
	aux := (*Alias)(r)
	return json.Unmarshal(data, aux)
}
