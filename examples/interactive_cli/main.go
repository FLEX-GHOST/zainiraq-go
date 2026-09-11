package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(30 * time.Second),
	)

	sessionFile := "zain_session.json"
	_ = client.LoadSessionFromFile(sessionFile)

	reader := bufio.NewReader(os.Stdin)

	for {
		msisdn := client.GetMSISDN()
		if msisdn == "" {
			msisdn = "Not Logged In"
		}

		fmt.Println("\n========================================")
		fmt.Println("    Zain Iraq Official CLI Console     ")
		fmt.Printf("    Current Session: %s\n", msisdn)
		fmt.Println("========================================")
		fmt.Println("1. Login with OTP (Request & Confirm)")
		fmt.Println("2. View Profile & SIM Information")
		fmt.Println("3. Check Balance & Bill Details")
		fmt.Println("4. Recharge Voucher Card")
		fmt.Println("5. Transfer Balance to Another Number")
		fmt.Println("6. Browse Active Subscriptions & Offers")
		fmt.Println("7. Claim Daily Free Gift")
		fmt.Println("8. Check Loyalty Points & Tier")
		fmt.Println("9. Check In-App Notifications")
		fmt.Println("10. Find Stores Near Me")
		fmt.Println("0. Exit")
		fmt.Print("\nSelect an action: ")

		input, _ := reader.ReadString('\n')
		action := strings.TrimSpace(input)

		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)

		switch action {
		case "1":
			fmt.Print("Enter Zain Phone Number: ")
			ph, _ := reader.ReadString('\n')
			ph = strings.TrimSpace(ph)
			if strings.HasPrefix(ph, "07") {
				ph = "964" + ph[1:]
			}
			otpReq, err := client.RequestOTP(ctx, ph)
			if err != nil {
				fmt.Printf("OTP Request failed: %v\n", err)
				cancel()
				continue
			}
			fmt.Printf("SMS sent! Request ID: %s\n", otpReq.Data.RequestID)
			fmt.Print("Enter OTP code: ")
			code, _ := reader.ReadString('\n')
			code = strings.TrimSpace(code)
			sess, err := client.VerifyOTPAndLogin(ctx, ph, code)
			if err != nil {
				fmt.Printf("Login failed: %v\n", err)
				cancel()
				continue
			}
			_ = client.SaveSessionToFile(sessionFile)
			fmt.Printf("Logged in successfully! Token: %s...\n", sess.AccessToken[:20])

		case "2":
			profile, err := client.GetProfile(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nName: %s | Number: %s | Billing: %s | Flex: %s\n", profile.Name, profile.MSISDN, profile.CustomerBillingType, profile.FlexStatus)
			}

		case "3":
			bal, err := client.GetBalance(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if bal.Balance != nil {
				fmt.Printf("\nBalance: %d IQD | Expiry: %s\n", bal.Balance.Value, bal.Balance.Expiry)
			}

		case "4":
			fmt.Print("Enter 16-digit voucher number: ")
			v, _ := reader.ReadString('\n')
			v = strings.TrimSpace(v)
			err := client.RechargeVoucher(ctx, v)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Voucher recharge completed successfully!")
			}

		case "5":
			fmt.Print("Receiver phone number: ")
			rPh, _ := reader.ReadString('\n')
			fmt.Print("Amount in IQD: ")
			amtStr, _ := reader.ReadString('\n')
			amt, _ := strconv.ParseInt(strings.TrimSpace(amtStr), 10, 64)
			fmt.Print("Enter OTP Confirmation Token: ")
			otpConf, _ := reader.ReadString('\n')
			err := client.CreditTransfer(ctx, strings.TrimSpace(rPh), amt, strings.TrimSpace(otpConf))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Credit transfer completed successfully!")
			}

		case "6":
			subs, err := client.GetSubscriptions(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nActive Subscriptions (%d):\n", len(subs))
				for _, s := range subs {
					title := s.ID
					if s.CMS != nil && s.CMS.Title != nil {
						title = s.CMS.Title.String()
					}
					fmt.Printf(" - %s (Status: %s | Expires: %s)\n", title, s.Status, s.ExpireTime)
				}
			}

		case "7":
			err := client.ClaimDailyGift(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Daily gift claimed successfully!")
			}

		case "8":
			loyalty, err := client.GetLoyaltyInfo(ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nTier: %s | Total Points: %d | Spendable: %d\n", loyalty.Tier, loyalty.TotalPoints, loyalty.TotalSpendablePoints)
			}

		case "9":
			notifs, err := client.GetNotifications(ctx, 0, 5, nil)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nRecent Notifications (%d):\n", len(notifs.Notifications))
				for _, n := range notifs.Notifications {
					fmt.Printf(" - [%s] %s: %s\n", n.CreatedAt, n.Title.String(), n.Body.String())
				}
			}

		case "10":
			fmt.Print("Enter city/search keyword (e.g. Baghdad): ")
			term, _ := reader.ReadString('\n')
			stores, err := client.GetStoreNearMe(ctx, strings.TrimSpace(term))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("\nFound %d stores nearby:\n", len(stores.Stores))
				for _, s := range stores.Stores {
					fmt.Printf(" - %s (%s): %s\n", s.Name, s.City, s.Address)
				}
			}

		case "0":
			fmt.Println("Goodbye!")
			cancel()
			return

		default:
			fmt.Println("Unknown choice, please try again.")
		}

		cancel()
	}
}
