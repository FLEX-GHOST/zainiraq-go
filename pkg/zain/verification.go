package zain

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	digitsOnlyRegex = regexp.MustCompile(`[^\d]`)
	amountRegex     = regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	phoneRegex      = regexp.MustCompile(`(07[789]\d{8}|7[789]\d{8}|9647[789]\d{8})`)
)

func cleanDigits(s string) string {
	return digitsOnlyRegex.ReplaceAllString(s, "")
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

func (c *Client) VerifyIncomingTransfer(ctx context.Context, senderPhone string, minAmount float64) (bool, *IncomingTransferRecord, error) {
	cleanSender := cleanDigits(senderPhone)
	if len(cleanSender) < 9 {
		return false, nil, fmt.Errorf("zain: invalid sender phone number (minimum 9 digits required)")
	}

	for _, rec := range c.GetRecordedIncomingTransfers() {
		if phonesMatch(rec.MSISDN, cleanSender, 9) {
			amtVal := parseAmountString(rec.Amount)
			if minAmount <= 0 || amtVal >= minAmount {
				matched := rec
				return true, &matched, nil
			}
		}
	}

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
