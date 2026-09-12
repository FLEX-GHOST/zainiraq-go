package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(zain.WithTimeout(25 * time.Second))

	sessionFile := "session.json"
	if err := client.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("Session file '%s' not found or invalid: %v\n", sessionFile, err)
		fmt.Println("Please run 'examples/01_otp_login' first to login and create a valid session.")
		return
	}

	masterWallet := "07801234567"
	if envWallet := os.Getenv("ZAIN_MASTER_WALLET"); envWallet != "" {
		masterWallet = envWallet
	}
	client.SetMasterWallet(masterWallet)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	fmt.Println("==================================================================")
	fmt.Println(" Zain Iraq Wallet & Automated Incoming Transfer Verification (Go SDK)")
	fmt.Println(" تثبيت محفظة زين والتحقق الآلي من الحوالات الواردة ومعرفة رقم المرسل")
	fmt.Println("==================================================================")
	fmt.Printf("Master Wallet Phone (رقم المحفظة المعتمد): %s\n\n", client.MasterWallet())

	fmt.Println("1. Fetching live wallet balance and status overview...")
	overview, err := client.GetWalletOverview(ctx)
	if err != nil {
		fmt.Printf("Warning: Failed to fetch wallet overview: %v\n", err)
	} else {
		fmt.Printf("   Wallet Phone : %s\n", overview.Phone)
		fmt.Printf("   Balance      : %.0f IQD\n", overview.Balance)
		fmt.Printf("   Validity     : %s\n", overview.Validity)
		fmt.Printf("   Active SIM   : %v\n", overview.IsActive)
	}

	fmt.Println("\n2. Customer Payment & Transfer Instructions:")
	sampleAmount := int64(5000)
	ussdCode := client.FormatUSSDTransfer(client.MasterWallet(), sampleAmount)
	fmt.Printf("   Send transfer of %d IQD using USSD dial code: %s\n", sampleAmount, ussdCode)

	fmt.Println("\n3. Querying recent incoming transfers from Zain Electronic Bill & Notifications...")
	transfers, err := client.GetIncomingTransfers(ctx, 10)
	if err != nil {
		fmt.Printf("Warning: Querying incoming transfers: %v\n", err)
	} else if len(transfers) > 0 {
		fmt.Printf("Found %d incoming transfer record(s):\n", len(transfers))
		for idx, tx := range transfers {
			fmt.Printf("   [%02d] Sender: %-14s | Amount: %-8s IQD | Time: %s | Title: %s\n",
				idx+1, tx.MSISDN, tx.Amount, tx.CreatedAt, tx.Title)
		}
	} else {
		fmt.Println("   No incoming transfers recorded recently.")
	}

	fmt.Println("\n4. Automated Transfer Verification (VerifyIncomingTransfer):")
	testSender := "07809998877"
	expectedAmount := 5000.0

	fmt.Printf("   Checking if '%s' sent at least %.0f IQD to master wallet...\n", testSender, expectedAmount)
	verified, matchedRecord, err := client.VerifyIncomingTransfer(ctx, testSender, expectedAmount)
	if err != nil {
		fmt.Printf("   Verification error: %v\n", err)
	} else if verified && matchedRecord != nil {
		fmt.Println("   >> [VERIFIED SUCCESS / تم التأكيد بنجاح]")
		fmt.Printf("      Sender Phone : %s\n", matchedRecord.MSISDN)
		fmt.Printf("      Amount Paid  : %s IQD\n", matchedRecord.Amount)
		fmt.Printf("      Transfer Time: %s\n", matchedRecord.CreatedAt)
	} else {
		fmt.Println("   >> [NOT FOUND / لم يتم العثور على تحويل مطابق]")
		fmt.Printf("      No incoming transfer from %s with minimum amount %.0f IQD.\n", testSender, expectedAmount)
	}

	fmt.Println("\n==================================================================")
	fmt.Println("Production Telegram Bot Architecture Flow (نمط بوتات شحن وتفعيل الرصيد):")
	fmt.Println("  1. Admin logs into Zain SIM once; bot calls client.SetMasterWallet(walletPhone).")
	fmt.Println("  2. Bot shows customer the master wallet and USSD code: *123*amount*wallet#")
	fmt.Println("  3. Customer sends balance transfer and submits their phone number.")
	fmt.Println("  4. Bot calls client.VerifyIncomingTransfer(ctx, customerPhone, amount).")
	fmt.Println("  5. If verified, bot records transaction and activates user account instantly!")
	fmt.Println("==================================================================")
}
