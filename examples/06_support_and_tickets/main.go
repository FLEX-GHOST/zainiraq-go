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
	fmt.Println("   Zain Iraq - Help & Support Tickets   ")
	fmt.Println("========================================")

	// 1. Fetch Complaint Templates & Categories
	fmt.Println("\n[1] Fetching Complaint Templates & Issue Categories...")
	templates, err := client.GetTemplates(ctx, "", "", "Network", "")
	if err != nil {
		fmt.Printf("Failed to get templates: %v\n", err)
	} else {
		fmt.Printf("Templates Found: %d\n", len(templates.Items))
		for _, item := range templates.Items {
			fmt.Printf(" - [%s] %s | Summary Code: %s\n", item.CategorizationCode, item.ItemEngTxt, item.SummariesCode)
		}
	}

	// 2. Fetch User's Open & Closed Tickets
	fmt.Println("\n[2] Fetching Submitted Support Tickets...")
	tickets, err := client.GetTickets(ctx, "ALL")
	if err != nil {
		fmt.Printf("Failed to get tickets: %v\n", err)
	} else {
		fmt.Printf("Total Submitted Tickets: %d\n", len(tickets.Tickets))
		for _, t := range tickets.Tickets {
			fmt.Printf(" - [%s] Status: %s | Created: %s | Category: %s\n", t.TxtIssueID, t.IssueStatus.String(), t.CreateTime, t.Category)
		}
	}

	// 3. Example of Submitting a Ticket (Dry-run demonstration)
	fmt.Println("\n[3] Preparing Support Ticket Submission Payload...")
	ticketReq := &zain.SubmitTicketReq{
		City:            "Baghdad",
		Governorate:     "Baghdad",
		TicketLanguage:  "ar",
		Source:          "MyZain",
		Description:     "Weak 4G signal in Karrada area during peak hours",
		SummaryID:       "NET-4G-COV",
		QuestionAnswers: `{"coverage_issue": "slow_data", "signal_bars": "1"}`,
		Attachments: []zain.Attachment{
			{
				Name: "speedtest.png",
				Data: "base64orRawDataHere",
			},
		},
	}
	fmt.Printf("Payload ready: %s (City: %s, Summary: %s)\n", ticketReq.Description, ticketReq.City, ticketReq.SummaryID)
	fmt.Println("To submit in real environment, call: client.CreateTicket(ctx, ticketReq)")

	// 4. Fetch Reopen Reasons
	fmt.Println("\n[4] Querying Ticket Reopen Reasons...")
	reasons, err := client.GetReOpenReasons(ctx)
	if err != nil {
		fmt.Printf("Failed to get reopen reasons: %v\n", err)
	} else {
		fmt.Printf("Available Reopen Reasons: %d\n", len(reasons.ReOpenReasons))
		for _, r := range reasons.ReOpenReasons {
			fmt.Printf(" - [EN: %s] [AR: %s]\n", r.ReasonNameEN, r.ReasonNameAR)
		}
	}
}
