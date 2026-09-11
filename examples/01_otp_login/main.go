package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(30 * time.Second),
	)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("   Zain Iraq - OTP Authentication Flow  ")
	fmt.Println("========================================")
	fmt.Print("Enter Zain Iraq Phone Number (e.g. 9647801234567 or 07801234567): ")
	phoneInput, _ := reader.ReadString('\n')
	phone := strings.TrimSpace(phoneInput)
	if phone == "" {
		fmt.Println("Error: Phone number is required.")
		return
	}

	// Normalize leading 0 to 964
	if strings.HasPrefix(phone, "07") {
		phone = "964" + phone[1:]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Printf("[1/3] Requesting OTP for %s...\n", phone)
	otpReq, err := client.RequestOTP(ctx, phone)
	if err != nil {
		fmt.Printf("Request OTP failed: %v\n", err)
		return
	}
	fmt.Printf("OTP successfully requested. Request ID: %s\n", otpReq.Data.RequestID)

	fmt.Print("Enter the OTP Code received via SMS: ")
	codeInput, _ := reader.ReadString('\n')
	code := strings.TrimSpace(codeInput)
	if code == "" {
		fmt.Println("Error: OTP code cannot be empty.")
		return
	}

	fmt.Println("[2/3] Confirming OTP code...")
	confirmResp, err := client.ConfirmOTP(ctx, phone, code)
	if err != nil {
		fmt.Printf("Confirm OTP failed: %v\n", err)
		return
	}
	fmt.Printf("OTP confirmed! Confirmation ID: %s\n", confirmResp.ConfirmationID)

	fmt.Println("[3/3] Signing session and retrieving JWT tokens...")
	session, err := client.SignSession(ctx, phone, confirmResp.ConfirmationID)
	if err != nil {
		fmt.Printf("Sign session failed: %v\n", err)
		return
	}

	fmt.Println("\nAuthentication Succeeded!")
	fmt.Printf("Access Token : %s...\n", session.AccessToken[:30])
	fmt.Printf("Refresh Token: %s...\n", session.RefreshToken[:30])

	sessionFile := "zain_session.json"
	if err := client.SaveSessionToFile(sessionFile); err != nil {
		fmt.Printf("Failed to persist session: %v\n", err)
	} else {
		fmt.Printf("Session saved safely to %s (chmod 0600)\n", sessionFile)
	}

	fmt.Println("\nTesting session restoration in a brand-new client instance...")
	restoredClient := zain.NewClient()
	if err := restoredClient.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("Failed to load session: %v\n", err)
		return
	}

	profile, err := restoredClient.GetProfile(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch user profile with restored token: %v\n", err)
		return
	}

	fmt.Printf("Welcome %s! MSISDN: %s | Billing: %s\n", profile.Name, profile.MSISDN, profile.CustomerBillingType)
}
