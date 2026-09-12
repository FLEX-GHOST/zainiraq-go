package main

import (
	"context"
	"fmt"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
	)

	// Load existing authenticated session if available
	_ = client.LoadSessionFromFile("zain_session.json")

	ctx := context.Background()

	fmt.Println("==================================================")
	fmt.Println(" Zain Iraq - Electronic Payments & ZainCash Demo")
	fmt.Println("==================================================")

	// 1. Generate Checkout ID for Card Payment (MasterCard / Visa)
	fmt.Println("\n[1] Generating Checkout ID for Card Payment...")
	checkoutReq := &zain.CheckoutReq{
		Amount:           "10000.00",
		Type:             "DB",
		MSISDN:           "7845900162",
		IsRegisteredCard: false,
		RegisterCard:     true,
	}

	checkoutResp, err := client.CreateCheckoutID(ctx, checkoutReq)
	if err != nil {
		fmt.Printf("Create checkout ID error: %v\n", err)
	} else {
		fmt.Printf("Checkout ID Created: %s\n", checkoutResp.CheckoutID)

		// 2. Refresh Gateway Payment Status
		fmt.Println("\n[2] Checking Payment Gateway Status...")
		statusReq := &zain.RefreshPaymentStatusReq{
			ResourcePath: "/v1/checkouts/" + checkoutResp.CheckoutID + "/payment",
			IsPayment:    true,
		}
		status, err := client.RefreshPaymentStatus(ctx, statusReq)
		if err != nil {
			fmt.Printf("Refresh payment status error: %v\n", err)
		} else if status != nil && status.PaymentStatus != nil {
			fmt.Printf("Payment ID: %s | Brand: %s\n",
				status.PaymentStatus.ID, status.PaymentStatus.PaymentBrand)
		}
	}

	// 3. ZainCash Direct Purchase Order
	fmt.Println("\n[3] Generating ZainCash Direct Purchase Order...")
	orderReq := &zain.InitiateOrderReq{
		OrderID:     "demo-order-7701",
		Amount:      5000,
		Currency:    "IQD",
		ServiceType: "BTL_DATA_5GB_WEEKLY",
		Recipient:   "7845900162",
		Lang:        "ar",
	}

	orderResp, err := client.InitiateZainCashPayment(ctx, orderReq)
	if err != nil {
		fmt.Printf("ZainCash purchase order error: %v\n", err)
	} else {
		fmt.Printf("Payment Gateway Redirect URL: %s\n", orderResp.PaymentURL)
	}

	// 4. Check Order Payment Status
	fmt.Println("\n[4] Querying Order Status...")
	orderStatus, err := client.GetPaymentStatus(ctx, "demo-order-token", "demo-order-7701")
	if err != nil {
		fmt.Printf("Get payment status error: %v\n", err)
	} else {
		fmt.Printf("Order Status: %s | TxID: %s\n", orderStatus.TransactionStatus, orderStatus.TransactionID)
	}
}
