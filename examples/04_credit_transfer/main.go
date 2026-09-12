package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func loadClient() (*zain.Client, string, error) {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(20*time.Second),
	)

	// Try loading saved session
	if err := client.LoadSessionFromFile("session.json"); err == nil {
		wallet := "07801234567"
		client.SetMasterWallet(wallet)
		return client, wallet, nil
	}
	if err := client.LoadSessionFromFile("zain_session.json"); err == nil {
		wallet := "07801234567"
		client.SetMasterWallet(wallet)
		return client, wallet, nil
	}

	return nil, "", fmt.Errorf("no saved session found (run 01_otp_login first)")
}

func main() {
	client, wallet, err := loadClient()
	hasSession := (err == nil)
	if !hasSession {
		client = zain.NewClient()
		wallet = "07801234567"
		client.SetMasterWallet(wallet)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== Zain Iraq Credit Transfer & Verification ===")
	fmt.Printf("Master Wallet: %s\n\n", wallet)

	fmt.Println("1. Recording mock incoming transfer in client ledger...")
	sampleSender := "07809876543"
	sampleAmount := 5000.0
	client.RecordIncomingTransfer(zain.IncomingTransferRecord{
		MSISDN:    sampleSender,
		Amount:    "5000",
		CreatedAt: time.Now().Format("2006-01-02 15:04"),
	})
	fmt.Println("Recorded successfully.")

	fmt.Println("\n2. Verifying incoming transfer via VerifyIncomingTransfer...")
	found, rec, err := client.VerifyIncomingTransfer(ctx, sampleSender, sampleAmount)
	if err != nil {
		fmt.Printf("Verify error: %v\n", err)
	} else if found && rec != nil {
		fmt.Printf("Transfer Verified: True\n    Sender: %s | Amount: %s IQD | Date: %s\n",
			rec.MSISDN, rec.Amount, rec.CreatedAt)
	} else {
		fmt.Println("Transfer not found.")
	}

	if hasSession {
		fmt.Println("\n3. Testing live verification against network bill & notifications...")
		found, _, _ := client.VerifyIncomingTransfer(ctx, "07800000000", 5000)
		fmt.Printf("Non-existent sender verified: %v\n", found)
	}

	fmt.Println("\n3. Generating USSD transfer code for customer dial:")
	code := client.FormatUSSDTransfer(wallet, int64(sampleAmount))
	fmt.Printf("   Dial: %s\n", code)

	fmt.Println("\n4. Direct Outbound Credit Transfer (with SMS OTP):")
	fmt.Println("   // Step 1: Request Transfer OTP:")
	fmt.Println("   // otp, err := client.RequestCreditTransferOTP(ctx, senderPhone)")
	fmt.Println("   // Step 2: Confirm OTP:")
	fmt.Println("   // conf, err := client.ConfirmCreditTransferOTP(ctx, \"123456\", senderPhone)")
	fmt.Println("   // Step 3: Execute Transfer:")
	fmt.Println("   // err := client.CreditTransfer(ctx, recipientPhone, 5000, conf.Data.ConfirmationID)")
}
