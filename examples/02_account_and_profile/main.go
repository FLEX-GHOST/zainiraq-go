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
	fmt.Println("   Zain Iraq - Account & Profile Query  ")
	fmt.Println("========================================")

	// 1. Fetch User Profile
	fmt.Println("\n[1] Fetching User Profile...")
	profile, err := client.GetProfile(ctx)
	if err != nil {
		fmt.Printf("Failed to get profile: %v\n", err)
	} else {
		fmt.Printf("Name            : %s\n", profile.Name)
		fmt.Printf("MSISDN          : %s\n", profile.MSISDN)
		fmt.Printf("Billing Type    : %s\n", profile.CustomerBillingType)
		fmt.Printf("Customer Type   : %s\n", profile.CustomerType)
		fmt.Printf("Flex Status     : %s\n", profile.FlexStatus)
		fmt.Printf("SIM Status      : %s\n", profile.UnifiedSIMStatus)
		fmt.Printf("Email           : %s\n", profile.Email)
		fmt.Printf("4G Compatible   : %t\n", profile.Is4GCompatible)
	}

	// 2. Fetch Wallet & Main Balance
	fmt.Println("\n[2] Querying Wallet & Balances...")
	balance, err := client.GetBalance(ctx)
	if err != nil {
		fmt.Printf("Failed to get balance: %v\n", err)
	} else if balance.Balance != nil {
		fmt.Printf("Main Balance    : %d IQD\n", balance.Balance.Value)
		fmt.Printf("Expiry Date     : %s\n", balance.Balance.Expiry)
	}

	// 3. Query Bill Details (Postpaid)
	fmt.Println("\n[3] Querying Bill Details...")
	bill, err := client.GetBillDetails(ctx)
	if err != nil {
		fmt.Printf("Failed to get bill details: %v\n", err)
	} else {
		fmt.Printf("Total Bill      : %.2f IQD\n", bill.TotalBill)
		fmt.Printf("Unbilled Amount : %.2f IQD\n", bill.UnbilledAmount)
		fmt.Printf("Advance Payment : %.2f IQD\n", bill.AdvancePayment)
		fmt.Printf("Past Due        : %.2f IQD\n", bill.PastDue)
	}

	// 4. Query Electronic Bill Breakdown
	fmt.Println("\n[4] Querying Electronic Bill Items...")
	ebill, err := client.GetElectronicBillItems(ctx)
	if err != nil {
		fmt.Printf("Failed to get electronic bill: %v\n", err)
	} else {
		fmt.Printf("Eligible: %t | Total Items: %d\n", ebill.Eligible, len(ebill.Items))
		for i, item := range ebill.Items {
			if i >= 5 {
				fmt.Printf("... and %d more items\n", len(ebill.Items)-5)
				break
			}
			fmt.Printf(" - [%s] %s | Vol: %s | Charge: %s IQD\n", item.StartTime, item.ServiceTypeName, item.RatingVolume, item.ChargeAmount)
		}
	}

	// 5. Query Linked Subaccounts
	fmt.Println("\n[5] Querying Subaccounts...")
	subaccounts, err := client.GetSubAccounts(ctx)
	if err != nil {
		fmt.Printf("Failed to get subaccounts: %v\n", err)
	} else {
		fmt.Printf("Found %d subaccount items\n", len(subaccounts))
		for _, sub := range subaccounts {
			fmt.Printf(" - Type: %s | Amount: %d | Initial: %d | Expiry: %s\n", sub.AccountType, sub.Amount, sub.InitialAmount, sub.ExpiryDate)
		}
	}

	// 6. Check SIM Pre-Registration Status
	fmt.Println("\n[6] Checking SIM Pre-Registration Status...")
	preReg, err := client.CheckPreRegistration(ctx)
	if err != nil {
		fmt.Printf("Failed to check pre-registration: %v\n", err)
	} else {
		fmt.Printf("Associated: %t | Unified SIM Status: %s\n", preReg.AssociatedWithUser, preReg.UnifiedSIMStatus)
	}
}
