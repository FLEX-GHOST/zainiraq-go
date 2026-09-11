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
	client := zain.NewClient(zain.WithLanguage("ar"))

	sessionFile := "zain_session.json"
	if err := client.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("Session file '%s' not found. Please run 01_otp_login first.\n", sessionFile)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("========================================")
	fmt.Println("  Zain Iraq - Recharge & Credit Operations")
	fmt.Println("========================================")
	fmt.Println("1. Recharge Line with Voucher Card")
	fmt.Println("2. Transfer Balance / Credit to Another Number")
	fmt.Println("3. View Validity Extension Options")
	fmt.Println("4. Extend SIM Validity")
	fmt.Print("\nSelect an operation [1-4]: ")

	choiceInput, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(choiceInput)

	switch choice {
	case "1":
		fmt.Print("Enter 16-digit voucher number: ")
		voucherInput, _ := reader.ReadString('\n')
		voucher := strings.TrimSpace(voucherInput)
		if voucher == "" {
			fmt.Println("Voucher number cannot be empty.")
			return
		}

		fmt.Println("Processing voucher recharge...")
		err := client.RechargeVoucher(ctx, voucher)
		if err != nil {
			fmt.Printf("Recharge failed: %v\n", err)
			return
		}
		fmt.Println("Recharge successful!")

	case "2":
		fmt.Print("Enter receiver MSISDN (e.g. 9647801234567): ")
		receiverInput, _ := reader.ReadString('\n')
		receiver := strings.TrimSpace(receiverInput)

		fmt.Print("Enter amount in IQD (e.g. 1000): ")
		amountInput, _ := reader.ReadString('\n')
		amount, _ := strconv.ParseInt(strings.TrimSpace(amountInput), 10, 64)

		fmt.Print("Enter OTP Confirmation token (from SMS): ")
		otpInput, _ := reader.ReadString('\n')
		otpConf := strings.TrimSpace(otpInput)

		fmt.Printf("Transferring %d IQD to %s...\n", amount, receiver)
		err := client.CreditTransfer(ctx, receiver, amount, otpConf)
		if err != nil {
			fmt.Printf("Credit transfer failed: %v\n", err)
			return
		}
		fmt.Println("Credit transfer successful!")

	case "3":
		fmt.Println("Fetching validity extension options...")
		optionsBody, err := client.GetValidityOptions(ctx)
		if err != nil {
			fmt.Printf("Failed to get validity options: %v\n", err)
			return
		}
		fmt.Println("\nAvailable Validity Options:")
		for i, opt := range optionsBody.ValidityOptions {
			fmt.Printf(" [%d] Duration: %s | Cost: %s IQD\n", i+1, opt.Duration, opt.Price)
		}

	case "4":
		fmt.Print("Enter amount in IQD to purchase validity: ")
		amtInput, _ := reader.ReadString('\n')
		amt, _ := strconv.ParseInt(strings.TrimSpace(amtInput), 10, 64)

		fmt.Println("Extending line validity...")
		resp, err := client.ExtendValidity(ctx, amt)
		if err != nil {
			fmt.Printf("Extend validity failed: %v\n", err)
			return
		}
		fmt.Printf("Validity extended successfully! New active until: %s | Balance: %d IQD\n", resp.NewActiveStop, resp.Balance)

	default:
		fmt.Println("Invalid option selected.")
	}
}
