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
	fmt.Println(" Zain Iraq - Bundle Sharing Group Management Demo")
	fmt.Println("==================================================")

	offerID := "OFFER_FAMILY_SHARE_50GB"
	familyMember := "7809988776"

	// 1. Query Current Sharing Group Status
	fmt.Println("\n[1] Querying Active Bundle Sharing Members...")
	group, err := client.SharingQueryMembers(ctx, offerID)
	if err != nil {
		fmt.Printf("Sharing query error: %v\n", err)
	} else {
		fmt.Printf("Total Group Members: %d\n", len(group.Members))
		for _, m := range group.Members {
			fmt.Printf(" - Member: %s (%s) | Allocated: %d MB | Consumed: %d MB\n",
				m.MSISDN, m.Role, m.Allocated, m.Consumed)
		}
	}

	// 2. Add a Member with Quota
	fmt.Printf("\n[2] Adding %s with 5,120 MB Quota to Sharing Group...\n", familyMember)
	err = client.SharingAddMember(ctx, familyMember, offerID, 5120)
	if err != nil {
		fmt.Printf("Add member error: %v\n", err)
	} else {
		fmt.Println("Member added to bundle sharing successfully!")
	}

	// 3. Transfer Extra Units to Member
	fmt.Printf("\n[3] Transferring 1,024 MB Extra Units to %s...\n", familyMember)
	err = client.SharingTransferUnits(ctx, familyMember, 1024)
	if err != nil {
		fmt.Printf("Transfer units error: %v\n", err)
	} else {
		fmt.Println("Units transferred successfully!")
	}

	// 4. Remove a Member from Group
	fmt.Printf("\n[4] Removing Member %s from Sharing Group...\n", familyMember)
	err = client.SharingRemoveMember(ctx, familyMember, offerID)
	if err != nil {
		fmt.Printf("Remove member error: %v\n", err)
	} else {
		fmt.Println("Member removed from sharing group successfully!")
	}
}
