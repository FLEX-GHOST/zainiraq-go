package zain

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	digitsOnlyRegex = regexp.MustCompile(`[^\d]`)
	amountRegex     = regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	phoneRegex      = regexp.MustCompile(`(07[789]\d{8}|7[789]\d{8}|9647[789]\d{8})`)
)

func convertEasternNumerals(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= '٠' && r <= '٩':
			b.WriteRune('0' + (r - '٠'))
		case r >= '۰' && r <= '۹':
			b.WriteRune('0' + (r - '۰'))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func cleanDigits(s string) string {
	return digitsOnlyRegex.ReplaceAllString(convertEasternNumerals(s), "")
}

func parseAmountString(s string) float64 {
	clean := strings.ReplaceAll(s, ",", "")
	match := amountRegex.FindString(clean)
	if match == "" {
		return 0
	}
	val, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}
	return val
}

func phonesMatch(a, b string, toleranceDigits int) bool {
	cleanA := cleanDigits(a)
	cleanB := cleanDigits(b)
	if len(cleanA) < toleranceDigits || len(cleanB) < toleranceDigits {
		return cleanA == cleanB && cleanA != ""
	}
	tailA := cleanA[len(cleanA)-toleranceDigits:]
	tailB := cleanB[len(cleanB)-toleranceDigits:]
	return tailA == tailB
}

func (c *Client) GetWalletOverview(ctx context.Context) (*WalletOverview, error) {
	phone := c.MasterWallet()
	if phone == "" {
		phone = c.GetMSISDN()
	}

	balResp, err := c.GetBalance(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("zain: getting balance overview: %w", err)
	}

	overview := &WalletOverview{
		Phone:    phone,
		IsActive: true,
	}

	if balResp != nil && balResp.Balance != nil {
		overview.Balance = float64(balResp.Balance.Value)
		overview.Validity = balResp.Balance.Expiry
	}

	return overview, nil
}

func (c *Client) GetIncomingTransfers(ctx context.Context, limit ...int) ([]IncomingTransferRecord, error) {
	maxItems := 30
	if len(limit) > 0 && limit[0] > 0 {
		maxItems = limit[0]
	}

	var results []IncomingTransferRecord

	inMem := c.GetRecordedIncomingTransfers()
	results = append(results, inMem...)

	phone := c.MasterWallet()
	billItems, err := c.GetElectronicBillItems(ctx, phone)
	if err == nil && billItems != nil {
		for _, it := range billItems.Items {
			amtVal := parseAmountString(it.ChargeAmount)
			if amtVal <= 0 {
				continue
			}

			senderNumber := cleanDigits(it.BNumber)
			if senderNumber == "" {
				continue
			}

			results = append(results, IncomingTransferRecord{
				MSISDN:      senderNumber,
				Amount:      fmt.Sprintf("%.0f", amtVal),
				CreatedAt:   it.StartTime,
				Title:       it.ServiceTypeName,
				ServiceType: it.ServiceTypeID,
				Raw:         it,
			})
		}
	}

	notifPage, err := c.GetNotifications(ctx, 0, maxItems, nil)
	if err == nil && notifPage != nil && len(notifPage.Notifications) > 0 {
		for _, n := range notifPage.Notifications {
			title := ""
			body := ""
			if n.Title != nil {
				title = n.Title.AR
				if title == "" {
					title = n.Title.EN
				}
			}
			if n.Body != nil {
				body = n.Body.AR
				if body == "" {
					body = n.Body.EN
				}
			}

			fullText := title + " " + body
			if strings.Contains(fullText, "تحويل") || strings.Contains(fullText, "رصيد") ||
				strings.Contains(strings.ToLower(fullText), "transfer") || strings.Contains(strings.ToLower(fullText), "credit") {

				extractedPhone := phoneRegex.FindString(fullText)
				extractedAmt := parseAmountString(fullText)

				if extractedPhone != "" && extractedAmt > 0 {
					results = append(results, IncomingTransferRecord{
						MSISDN:      cleanDigits(extractedPhone),
						Amount:      fmt.Sprintf("%.0f", extractedAmt),
						CreatedAt:   n.CreatedAt,
						Title:       title,
						ServiceType: "notification",
						Raw:         n,
					})
				}
			}
		}
	}

	if len(results) > maxItems {
		results = results[:maxItems]
	}

	return results, nil
}

func (c *Client) GetCDRTransferHistory(ctx context.Context, limit ...int) ([]IncomingTransferRecord, error) {
	return c.GetIncomingTransfers(ctx, limit...)
}

// VerifyIncomingTransferFromMemory inspects only the in-memory recorded transfers (e.g. from SMS)
// with zero network latency and zero external HTTP calls.
func (c *Client) VerifyIncomingTransferFromMemory(senderPhone string, minAmount float64) (bool, *IncomingTransferRecord) {
	cleanSender := cleanDigits(senderPhone)
	if len(cleanSender) < 9 {
		return false, nil
	}

	for _, rec := range c.GetRecordedIncomingTransfers() {
		if phonesMatch(rec.MSISDN, cleanSender, 9) {
			amtVal := parseAmountString(rec.Amount)
			if minAmount <= 0 || amtVal >= minAmount {
				matched := rec
				return true, &matched
			}
		}
	}
	return false, nil
}

func (c *Client) queryCloudTransfers(ctx context.Context, cleanSender string, minAmount float64) (bool, *IncomingTransferRecord, error) {
	liveTransfers, err := c.GetIncomingTransfers(ctx, 50)
	if err != nil {
		return false, nil, fmt.Errorf("zain: querying incoming transfers: %w", err)
	}

	for _, rec := range liveTransfers {
		if phonesMatch(rec.MSISDN, cleanSender, 9) {
			amtVal := parseAmountString(rec.Amount)
			if minAmount <= 0 || amtVal >= minAmount {
				matched := rec
				return true, &matched, nil
			}
		}
	}

	return false, nil, nil
}

func (c *Client) VerifyIncomingTransfer(ctx context.Context, senderPhone string, minAmount float64) (bool, *IncomingTransferRecord, error) {
	cleanSender := cleanDigits(senderPhone)
	if len(cleanSender) < 9 {
		return false, nil, fmt.Errorf("zain: invalid sender phone number (minimum 9 digits required)")
	}

	if ok, rec := c.VerifyIncomingTransferFromMemory(cleanSender, minAmount); ok && rec != nil {
		return true, rec, nil
	}

	return c.queryCloudTransfers(ctx, cleanSender, minAmount)
}

// NormalizeMSISDN standardizes any Iraqi phone format (local 078..., 78..., international +964..., or Arabic numerals)
// into the standard 13-digit international format required by Zain Iraq APIs (9647XXXXXXXX).
func NormalizeMSISDN(input string) (string, error) {
	digits := cleanDigits(input)

	if strings.HasPrefix(digits, "00964") {
		digits = strings.TrimPrefix(digits, "00")
	}

	switch {
	case len(digits) == 13 && strings.HasPrefix(digits, "9647"):
		return digits, nil
	case len(digits) == 11 && strings.HasPrefix(digits, "07"):
		return "964" + digits[1:], nil
	case len(digits) == 10 && strings.HasPrefix(digits, "7"):
		return "964" + digits, nil
	default:
		return "", fmt.Errorf("zain: invalid Iraqi phone number format: %q", input)
	}
}

// FormatLocalMSISDN converts a phone number into local Iraqi format (07XXXXXXXX).
func FormatLocalMSISDN(phone string) string {
	clean := cleanDigits(phone)
	if strings.HasPrefix(clean, "964") && len(clean) == 13 {
		return "0" + clean[3:]
	}
	if strings.HasPrefix(clean, "7") && len(clean) == 10 {
		return "0" + clean
	}
	return clean
}

type cloudPollResult struct {
	ok       bool
	rec      *IncomingTransferRecord
	err      error
	interval time.Duration
}

// WaitForIncomingTransfer polls periodically until an incoming transfer from senderPhone is verified,
// or until the context expires or is cancelled.
// It prioritizes local in-memory SMS records (0 network calls) checked rapidly in milliseconds,
// and applies exponential backoff with jitter on remote cloud CDR polling to prevent
// HTTP 429 Too Many Requests and IP rate limits on Zain's servers.
func (c *Client) WaitForIncomingTransfer(ctx context.Context, senderPhone string, minAmount float64, interval time.Duration) (*IncomingTransferRecord, error) {
	if interval <= 0 {
		interval = 3 * time.Second
	}

	normalized, err := NormalizeMSISDN(senderPhone)
	if err != nil {
		normalized = cleanDigits(senderPhone)
	}

	// 1. Immediate fast-path: in-memory SMS ledger (zero network latency)
	if ok, rec := c.VerifyIncomingTransferFromMemory(normalized, minAmount); ok && rec != nil {
		return rec, nil
	}

	currentInterval := interval
	maxInterval := 15 * time.Second
	if maxInterval < currentInterval {
		maxInterval = currentInterval * 2
	}

	localCheckFreq := 100 * time.Millisecond
	if localCheckFreq > interval {
		localCheckFreq = interval
	}
	if localCheckFreq < 10*time.Millisecond {
		localCheckFreq = 10 * time.Millisecond
	}

	localTicker := time.NewTicker(localCheckFreq)
	defer localTicker.Stop()

	cloudResultChan := make(chan cloudPollResult, 1)
	cloudInFlight := false
	nextCloudCheck := time.Now() // trigger first cloud check immediately

	for {
		// If cloud query is idle and scheduled time has arrived, trigger single non-blocking check
		if !cloudInFlight && time.Now().After(nextCloudCheck) {
			cloudInFlight = true
			go func(pollInterval time.Duration) {
				defer func() {
					if r := recover(); r != nil {
						// Goroutine panic recovery
					}
				}()
				cloudCtx, cloudCancel := context.WithTimeout(ctx, 5*time.Second)
				defer cloudCancel()
				ok, rec, queryErr := c.queryCloudTransfers(cloudCtx, normalized, minAmount)
				select {
				case cloudResultChan <- cloudPollResult{ok: ok, rec: rec, err: queryErr, interval: pollInterval}:
				case <-ctx.Done():
				}
			}(currentInterval)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-localTicker.C:
			// Rapid check for local SMS ledger (zero network overhead, zero latency)
			if ok, rec := c.VerifyIncomingTransferFromMemory(normalized, minAmount); ok && rec != nil {
				return rec, nil
			}

		case res := <-cloudResultChan:
			cloudInFlight = false
			if res.ok && res.rec != nil {
				return res.rec, nil
			}

			// Backoff logic on cloud failure or missing transfer:
			if res.err != nil && (strings.Contains(res.err.Error(), "429") || strings.Contains(res.err.Error(), "rate")) {
				currentInterval = maxInterval
			} else {
				currentInterval = time.Duration(float64(res.interval) * 1.5)
				if currentInterval > maxInterval {
					currentInterval = maxInterval
				}
			}

			jitter := time.Duration((time.Now().UnixNano()%400)-200) * time.Millisecond
			nextCloudCheck = time.Now().Add(currentInterval + jitter)
		}
	}
}

// ParseTransferSMS parses an incoming credit transfer SMS text (received via SMS gateway, Android SMS listener, or GSM modem)
// and extracts the sender phone number and transferred amount.
// Supported formats include Arabic and English Zain SMS templates:
// - "تم استلام رصيد بقيمة 5,000 د.ع من الرقم 07801234567 بنجاح"
// - "تم تحويل مبلغ 10000 دينار من الرقم 9647801234567"
// - "You have received 5,000 IQD from 07801234567"
func ParseTransferSMS(smsText string) (*IncomingTransferRecord, error) {
	if strings.TrimSpace(smsText) == "" {
		return nil, fmt.Errorf("zain: empty sms text")
	}

	converted := convertEasternNumerals(smsText)
	lower := strings.ToLower(converted)

	isTransfer := strings.Contains(converted, "تحويل") ||
		strings.Contains(converted, "استلام") ||
		strings.Contains(converted, "استلمت") ||
		strings.Contains(converted, "رصيد") ||
		strings.Contains(lower, "transfer") ||
		strings.Contains(lower, "received") ||
		strings.Contains(lower, "credit")

	if !isTransfer {
		return nil, fmt.Errorf("zain: message is not a recognized credit transfer sms")
	}

	phoneMatch := phoneRegex.FindString(converted)
	if phoneMatch == "" {
		return nil, fmt.Errorf("zain: no sender phone number found in sms")
	}

	amt := parseAmountString(converted)
	if amt <= 0 {
		return nil, fmt.Errorf("zain: no valid transfer amount found in sms")
	}

	return &IncomingTransferRecord{
		MSISDN:      cleanDigits(phoneMatch),
		Amount:      fmt.Sprintf("%.0f", amt),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		Title:       "تحويل رصيد وارد (SMS)",
		ServiceType: "sms",
		Raw:         smsText,
	}, nil
}

// RecordIncomingTransferFromSMS parses a credit transfer SMS text and automatically records
// it in the client's local ledger for immediate matching by VerifyIncomingTransfer.
func (c *Client) RecordIncomingTransferFromSMS(smsText string) (*IncomingTransferRecord, error) {
	rec, err := ParseTransferSMS(smsText)
	if err != nil {
		return nil, err
	}
	c.RecordIncomingTransfer(*rec)
	return rec, nil
}

