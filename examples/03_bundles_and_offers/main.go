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
	fmt.Println("   Zain Iraq - Bundles, Offers & Flex   ")
	fmt.Println("========================================")

	// 1. Query Active Subscriptions
	fmt.Println("\n[1] Querying Active Subscriptions...")
	subs, err := client.GetSubscriptions(ctx)
	if err != nil {
		fmt.Printf("Failed to get subscriptions: %v\n", err)
	} else {
		fmt.Printf("Active Subscriptions Count: %d\n", len(subs))
		for _, sub := range subs {
			title := sub.ID
			if sub.CMS != nil && sub.CMS.Title != nil {
				title = sub.CMS.Title.String()
			}
			fmt.Printf(" - ID: %s | Title: %s | Status: %s | Expires: %s\n", sub.ID, title, sub.Status, sub.ExpireTime)
		}
	}

	// 2. Query Offers Catalog
	fmt.Println("\n[2] Fetching Offers Catalog...")
	offers, err := client.GetOffers(ctx, "data_bundles", "", "")
	if err != nil {
		fmt.Printf("Failed to get offers: %v\n", err)
	} else {
		fmt.Printf("Total Offers: %d\n", len(offers.Offers))
		for i, o := range offers.Offers {
			if i >= 5 {
				fmt.Printf("... and %d more offers\n", len(offers.Offers)-5)
				break
			}
			price := int64(0)
			if o.Specs != nil {
				price = o.Specs.Price
			}
			fmt.Printf(" - [%s] %s | Price: %d IQD\n", o.CBSID, o.Title.String(), price)
		}
	}

	// 3. Query Personalized BTL (Below The Line) Offers
	fmt.Println("\n[3] Checking Personalized (BTL) Exclusive Offers...")
	btlOffers, err := client.GetPersonalizedOffersBTL(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch BTL offers: %v\n", err)
	} else {
		fmt.Printf("Exclusive Offers Available: %d\n", len(btlOffers.Offers))
		for _, o := range btlOffers.Offers {
			fmt.Printf(" - [Special Offer] %s (CBS ID: %s)\n", o.Title.String(), o.CBSID)
		}
	}

	// 4. Claim Daily Free Gift
	fmt.Println("\n[4] Claiming Daily Free Gift...")
	err = client.ClaimDailyGift(ctx)
	if err != nil {
		fmt.Printf("Claim daily gift: %v (might have already been claimed today)\n", err)
	} else {
		fmt.Println("Daily gift claimed successfully!")
	}

	// 5. Query Flex Limits & Status
	fmt.Println("\n[5] Checking Flex Status & Limits...")
	flexStatus, err := client.GetFlexStatus(ctx)
	if err != nil {
		fmt.Printf("Failed to get flex status: %v\n", err)
	} else {
		fmt.Printf("Flex Status: %s | Offer: %s | Active Until: %s\n", flexStatus.Status, flexStatus.OfferName, flexStatus.ActiveUntil)
	}

	flexLimits, err := client.GetFlexLimits(ctx)
	if err != nil {
		fmt.Printf("Failed to get flex limits: %v\n", err)
	} else {
		fmt.Printf("Total Points: %d | Remaining Points: %d\n", flexLimits.TotalPoints, flexLimits.RemainingPoints)
	}

	// 6. Query Bundle Sharing Members
	fmt.Println("\n[6] Querying Bundle Sharing Group...")
	sharing, err := client.SharingQueryMembers(ctx, "")
	if err != nil {
		fmt.Printf("Failed to query bundle sharing: %v\n", err)
	} else {
		fmt.Printf("Group Members: %d\n", len(sharing.Members))
		for _, m := range sharing.Members {
			fmt.Printf(" - %s (%s): Allocated %d | Consumed: %d\n", m.MSISDN, m.Role, m.Allocated, m.Consumed)
		}
	}
}
