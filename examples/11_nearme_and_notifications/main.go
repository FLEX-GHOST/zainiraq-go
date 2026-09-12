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

	_ = client.LoadSessionFromFile("zain_session.json")

	ctx := context.Background()

	fmt.Println("==================================================")
	fmt.Println(" Zain Iraq - Branches (Near Me) & Notifications")
	fmt.Println("==================================================")

	// 1. Search Certified Zain Branches & Stores
	fmt.Println("\n[1] Finding nearest Zain stores in Baghdad...")
	storesContent, err := client.GetStoreNearMe(ctx, "Baghdad")
	if err != nil {
		fmt.Printf("Near me stores error: %v\n", err)
	} else {
		fmt.Printf("Certified Branches Found: %d\n", len(storesContent.Stores))
		for i, s := range storesContent.Stores {
			if i >= 5 {
				break
			}
			fmt.Printf(" - [%s] %s | City: %s | Address: %s | Open: %t\n",
				s.ID, s.Name, s.City, s.Address, s.IsOpen)
		}
	}

	// 2. Query Inbox Notifications (Filtered to unread only)
	fmt.Println("\n[2] Fetching inbox notifications...")
	unread := false
	notifs, err := client.GetNotifications(ctx, 0, 5, &unread)
	if err != nil {
		fmt.Printf("Notifications error: %v\n", err)
	} else {
		fmt.Printf("Inbox Notifications Count: %d\n", len(notifs.Notifications))
		for _, n := range notifs.Notifications {
			fmt.Printf(" - [%s] %s: %s (Read: %t)\n",
				n.CreatedAt, n.Title.String(), n.Body.String(), n.Read)
		}
	}

	// 3. Query Active Digital Services Catalogue
	fmt.Println("\n[3] Fetching active digital entertainment services...")
	digServices, err := client.GetDigitalServices(ctx, "Entertainment", "", "PREPAID_NORMAL")
	if err != nil {
		fmt.Printf("Digital services error: %v\n", err)
	} else {
		fmt.Printf("Active Digital Services: %d\n", len(digServices.Services))
		for i, svc := range digServices.Services {
			if i >= 4 {
				break
			}
			fmt.Printf(" - %s | Price: %.2f IQD\n", svc.Title.String(), svc.Price)
		}
	}

	// 4. Query Dashboard Announcements
	fmt.Println("\n[4] Querying Dashboard Alerts & Messages...")
	dashMsgs, err := client.GetDashboardMessages(ctx)
	if err != nil {
		fmt.Printf("Dashboard messages error: %v\n", err)
	} else {
		fmt.Printf("Active Dashboard Alerts: %d\n", len(dashMsgs.Messages))
	}
}
