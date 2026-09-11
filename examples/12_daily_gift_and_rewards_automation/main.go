package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client := zain.NewClient(
		zain.WithLanguage("ar"),
		zain.WithTimeout(30*time.Second),
	)

	sessionFile := "zain_session.json"
	if err := client.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("الملف '%s' غير موجود. يرجى تسجيل الدخول أولاً عبر examples/01_otp_login\n", sessionFile)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	fmt.Println("==========================================================")
	fmt.Println("   Zain Iraq - Daily Gifts & Rewards Automation Bot       ")
	fmt.Println("==========================================================")

	// 1. المطالبة بالهدية اليومية المجانية تلقائياً
	fmt.Println("\n[1] جاري المطالبة بهدية زين اليومية المجانية (Daily Gift)...")
	err := client.ClaimDailyGift(ctx)
	if err != nil {
		fmt.Printf("⚠️ تنبيه: لم يتم استلام الهدية (ربما تم استلامها اليوم مسبقاً): %v\n", err)
	} else {
		fmt.Println("✅ تم استلام الهدية اليومية المجانية بنجاح وتفعيلها على الخط!")
	}

	// 2. فحص نقاط برنامج المكافآت (Zain Loyalty Points)
	fmt.Println("\n[2] جاري فحص رصيد نقاط المكافآت...")
	loyalty, err := client.GetLoyaltyInfo(ctx)
	if err != nil {
		fmt.Printf("❌ فشل استعلام النقاط: %v\n", err)
	} else {
		fmt.Printf(" - الفئة الحالية       : %s\n", loyalty.Tier)
		fmt.Printf(" - إجمالي النقاط      : %d نقطة\n", loyalty.TotalPoints)
		fmt.Printf(" - النقاط القابلة للصرف: %d نقطة\n", loyalty.TotalSpendablePoints)

		// تحويل النقاط لرصيد نقدي إذا كان الرصيد كافياً
		if loyalty.TotalSpendablePoints >= 1000 {
			fmt.Println(" ✨ الرصيد يسمح بالاستبدال، جاري تحويل 1000 نقطة إلى رصيد نقدي بالدينار...")
			err := client.RedeemCredit(ctx, 1000)
			if err != nil {
				fmt.Printf(" ⚠️ فشل تحويل النقاط لرصيد: %v\n", err)
			} else {
				fmt.Println(" ✅ تم تحويل 1000 نقطة إلى رصيد شريحة بنجاح!")
			}
		} else {
			fmt.Println(" ℹ️ نقاطك أقل من 1000 نقطة، لا حاجة للتحويل الآن.")
		}
	}

	// 3. فحص صلاحية الشريحة
	fmt.Println("\n[3] جاري فحص رصيد وصلاحية الخط...")
	bal, err := client.GetBalance(ctx)
	if err != nil {
		fmt.Printf("❌ فشل فحص الرصيد: %v\n", err)
	} else if bal.Balance != nil {
		fmt.Printf(" - الرصيد الحالي: %d د.ع\n", bal.Balance.Value)
		fmt.Printf(" - تاريخ الصلاحية: %s\n", bal.Balance.Expiry)
	}

	// 4. خيارات تمديد الصلاحية
	fmt.Println("\n[4] جاري فحص خيارات تمديد صلاحية الشريحة المتاحة...")
	validityOpts, err := client.GetValidityOptions(ctx)
	if err != nil {
		fmt.Printf("❌ فشل استعلام خيارات الصلاحية: %v\n", err)
	} else {
		fmt.Printf(" - الخيارات المتاحة من زين: %d خيار\n", len(validityOpts.ValidityOptions))
		for _, opt := range validityOpts.ValidityOptions {
			fmt.Printf("   * تمديد %s بمبلغ %s د.ع\n", opt.Duration, opt.Price)
		}
	}

	// 5. فحص الإشعارات الحديثة
	fmt.Println("\n[5] جاري جلب آخر إشعارات التطبيق (التأكيد والحوالات)...")
	notifs, err := client.GetNotifications(ctx, 0, 5, nil)
	if err != nil {
		fmt.Printf("❌ فشل جلب الإشعارات: %v\n", err)
	} else {
		fmt.Printf(" - تم جلب %d إشعار:\n", len(notifs.Notifications))
		for _, n := range notifs.Notifications {
			fmt.Printf("   * [%s] %s: %s\n", n.CreatedAt, n.Title.String(), n.Body.String())
		}
	}

	log.Println("\n🎉 اكتمل فحص وتنفيذ مهام الأتمتة بنجاح.")
}
