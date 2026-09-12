package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(zain.WithLanguage("ar"))

	sessionFile := "zain_session.json"
	if err := client.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("Session file '%s' not found. Please run 01_otp_login first.\n", sessionFile)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("========================================")
	fmt.Println("   Zain Iraq - Loyalty & Imtiyaz Program ")
	fmt.Println("========================================")

	// 1. Query Loyalty Points & Tier
	fmt.Println("\n[1] Fetching Zain Loyalty Status...")
	loyalty, err := client.GetLoyaltyInfo(ctx)
	if err != nil {
		fmt.Printf("Failed to get loyalty info: %v\n", err)
	} else {
		fmt.Printf("Current Tier          : %s\n", loyalty.Tier)
		fmt.Printf("Total Points          : %d\n", loyalty.TotalPoints)
		fmt.Printf("Spendable Points      : %d\n", loyalty.TotalSpendablePoints)
		fmt.Printf("Earliest Points Expiry: %s\n", loyalty.EarliestSpendablePointValidity)
	}

	// 2. Query Loyalty Statement & Points History
	fmt.Println("\n[2] Fetching Loyalty Statement...")
	statement, err := client.GetLoyaltyStatement(ctx, "", "")
	if err != nil {
		fmt.Printf("Failed to get loyalty statement: %v\n", err)
	} else {
		fmt.Printf("Recent Points Transactions: %d\n", len(statement.Transactions))
		for i, item := range statement.Transactions {
			if i >= 5 {
				break
			}
			fmt.Printf(" - [%s] %s : %d pts (%s)\n", item.Date, item.Description, item.Points, item.Type)
		}
	}

	// 3. Query Hot Rewards Bundles
	fmt.Println("\n[3] Checking Loyalty Rewards Bundles...")
	bundles, err := client.GetLoyaltyHotBundles(ctx)
	if err != nil {
		fmt.Printf("Failed to get hot bundles: %v\n", err)
	} else {
		fmt.Printf("Available Rewards Bundles: %d\n", len(bundles.Bundles))
		for _, b := range bundles.Bundles {
			fmt.Printf(" - [%s] %s | Points Required: %d\n", b.ID, b.Title.String(), b.Points)
		}
	}

	// 4. Query Imtiyaz Partner Discounts
	fmt.Println("\n[4] Querying Imtiyaz Partner Merchants...")
	merchant, err := client.GetMerchant(ctx, "restaurants", "burger_king")
	if err != nil {
		fmt.Printf("Failed to fetch merchant details: %v\n", err)
	} else {
		fmt.Printf("Partner: %s\n", merchant.Name.String())
		fmt.Printf("Discount Offer: %s\n", merchant.Discount)
		fmt.Printf("Description: %s\n", merchant.Description.String())
	}

	// 5. Query Imtiyaz Claimed Voucher History
	fmt.Println("\n[5] Checking Imtiyaz Claimed Vouchers...")
	vouchers, err := client.GetVoucherHistory(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch voucher history: %v\n", err)
	} else {
		fmt.Printf("Claimed Vouchers: %d\n", len(vouchers.History))
		for _, v := range vouchers.History {
			fmt.Printf(" - [%s] Code: %s | Date: %s | Status: %s\n", v.Merchant, v.VoucherCode, v.Date, v.Status)
		}
	}
}
