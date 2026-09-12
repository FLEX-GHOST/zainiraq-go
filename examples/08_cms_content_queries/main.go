package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(15*time.Second),
	)

	ctx := context.Background()

	fmt.Println("==================================================")
	fmt.Println(" Zain Iraq - CMS Direct Query Engine Demonstration")
	fmt.Println("==================================================")

	// 1. Query Dmart CMS Gateway Status
	fmt.Println("\n[1] Checking Dmart CMS Gateway Service Status...")
	status, err := client.GetDmartStatus(ctx)
	if err != nil {
		fmt.Printf("Dmart status error: %v\n", err)
	} else {
		fmt.Printf("Dmart CMS Gateway Status: %v\n", status)
	}

	// 2. Query App Configurations from CMS
	fmt.Println("\n[2] Fetching App Configurations from CMS...")
	cfg, err := client.GetCMSConfiguration(ctx, "")
	if err != nil {
		fmt.Printf("CMS app configurations error: %v\n", err)
	} else {
		fmt.Printf("Configurations Retrieved: %+v\n", cfg)
	}

	// 3. Query Dashboard Banners from CMS
	fmt.Println("\n[3] Fetching Dashboard Banners from CMS...")
	banners, err := client.GetCMSDashboardBanners(ctx, "")
	if err != nil {
		fmt.Printf("CMS dashboard banners error: %v\n", err)
	} else if banners != nil && banners.Message != nil {
		fmt.Printf("Dashboard Banner Message: %s\n", banners.Message.Title.String())
	}

	// 4. Query Subaccount Definitions from CMS
	fmt.Println("\n[4] Fetching Subaccount Schema from CMS...")
	sub, err := client.GetCMSSubaccount(ctx, "")
	if err != nil {
		fmt.Printf("CMS subaccount error: %v\n", err)
	} else {
		fmt.Printf("Subaccount Metadata: %+v\n", sub)
	}

	// 5. Query Slideshow Content from CMS
	fmt.Println("\n[5] Fetching Home Slideshow from CMS...")
	slideshow, err := client.GetSlideshow(ctx, "home_slideshow")
	if err != nil {
		fmt.Printf("Slideshow error: %v\n", err)
	} else if slideshow != nil {
		fmt.Printf("Slideshow Action: %s | Slides Count: %d\n", slideshow.CtaAction, len(slideshow.Slides))
	}
}
