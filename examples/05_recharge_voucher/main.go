package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func loadClient() (*zain.Client, error) {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(20*time.Second),
	)

	if err := client.LoadSessionFromFile("session.json"); err == nil {
		return client, nil
	}
	if err := client.LoadSessionFromFile("zain_session.json"); err == nil {
		return client, nil
	}

	return nil, fmt.Errorf("no saved session found (run 01_otp_login first)")
}

func main() {
	client, err := loadClient()
	hasSession := (err == nil)
	if !hasSession {
		client = zain.NewClient()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println("=== Zain Iraq Voucher Recharge & Validity ===")

	if hasSession {
		fmt.Println("1. Fetching Current Account Balance...")
		balance, err := client.GetBalance(ctx)
		if err != nil {
			fmt.Printf("Balance error: %v\n", err)
		} else if balance != nil && balance.Balance != nil {
			fmt.Printf("Current Balance: %d IQD | Expiry: %s\n", balance.Balance.Value, balance.Balance.Expiry)
		}

		fmt.Println("\n2. Fetching Line Validity Extension Options...")
		opts, err := client.GetValidityOptions(ctx)
		if err != nil {
			fmt.Printf("Validity options error: %v\n", err)
		} else if opts != nil {
			for i, opt := range opts.ValidityOptions {
				fmt.Printf("[%d] Duration: %s | Price: %s IQD\n", i+1, opt.Duration, opt.Price)
			}
		}
	} else {
		fmt.Println("1. To recharge a 16-digit voucher card on the current authenticated line:")
		fmt.Println("   err := client.RechargeVoucher(ctx, \"1234567890123456\")")
		fmt.Println("   // err == nil -> Recharge successful")

		fmt.Println("\n2. To recharge a voucher for another line:")
		fmt.Println("   err := client.RechargeVoucher(ctx, \"1234567890123456\", \"9647801234567\")")

		fmt.Println("\n3. To extend SIM validity with balance:")
		fmt.Println("   resp, err := client.ExtendValidity(ctx, 1000) // 1,000 IQD")

		fmt.Println("\nNote: Run 'go run examples/01_otp_login/main.go' first to test with a live session.")
	}
}
