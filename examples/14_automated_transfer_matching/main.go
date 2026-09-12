package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

// AccountUser represents an account or customer in the application ledger.
type AccountUser struct {
	AccountID   string
	AccountName string
	SenderPhone string
	Balance     float64
	IsActive    bool
	ActivatedAt string
}

// LedgerTransaction represents a processed transfer to prevent double-spending / replay attacks.
type LedgerTransaction struct {
	AccountID string
	Phone     string
	Amount    float64
	TxTime    string
}

// MemoryLedger simulates a transactional storage system (e.g. SQLite / PostgreSQL) for accounting.
type MemoryLedger struct {
	mu           sync.Mutex
	users        map[string]*AccountUser
	transactions map[string]LedgerTransaction // key: senderPhone:amount:txTime
}

func NewMemoryLedger() *MemoryLedger {
	return &MemoryLedger{
		users:        make(map[string]*AccountUser),
		transactions: make(map[string]LedgerTransaction),
	}
}

func (l *MemoryLedger) GetOrCreateUser(accountID, accountName string) *AccountUser {
	l.mu.Lock()
	defer l.mu.Unlock()

	user, exists := l.users[accountID]
	if !exists {
		user = &AccountUser{
			AccountID:   accountID,
			AccountName: accountName,
			Balance:     0.0,
			IsActive:    false,
		}
		l.users[accountID] = user
	}
	return user
}

func (l *MemoryLedger) RecordTransaction(accountID, phone string, amount float64, txTime string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := fmt.Sprintf("%s:%.0f:%s", phone, amount, txTime)
	if _, exists := l.transactions[key]; exists {
		// Reject duplicate / already claimed transfer
		return false
	}

	l.transactions[key] = LedgerTransaction{
		AccountID: accountID,
		Phone:     phone,
		Amount:    amount,
		TxTime:    txTime,
	}

	// Update customer balance and activate account
	if user, ok := l.users[accountID]; ok {
		user.Balance += amount
		user.IsActive = true
		user.SenderPhone = phone
		user.ActivatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	return true
}

func main() {
	fmt.Println("================================================================================")
	fmt.Println("   Zain Iraq - Automated Incoming Transfer & SMS Verification Engine")
	fmt.Println("   محرك التحقق الآلي من تحويلات رصيد زين العراق ومطابقتها من الرسائل والكشف السحابي")
	fmt.Println("================================================================================")

	client := zain.NewClient(zain.WithTimeout(15 * time.Second))
	ledger := NewMemoryLedger()

	// 1. Configure System Master Wallet
	masterWallet := "07801234567"
	client.SetMasterWallet(masterWallet)
	fmt.Printf("[1] Master Wallet Configured : %s\n", client.MasterWallet())

	// 2. Customer payment instructions
	depositPrice := 5000.0 // 5,000 IQD required
	ussdCode := client.FormatUSSDTransfer(client.MasterWallet(), int64(depositPrice))
	fmt.Printf("[2] USSD Dial Code for Users  : %s\n", ussdCode)
	fmt.Printf("    (الزبون يقوم بالطلب المباشر من هاتفه دون الحاجة لتسجيل دخوله في النظام)\n\n")

	// 3. Simulate Incoming SMS Messages & In-App Notifications from Zain Iraq
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("[3] محاكاة وصول رسائل SMS وإشعارات استلام رصيد من شبكة زين العراق:")
	fmt.Println("--------------------------------------------------------------------------------")

	incomingSMSList := []string{
		"تم استلام رصيد بقيمة 5,000 د.ع من الرقم 07809876543 بنجاح. رصيدك الحالي هو 25,000 د.ع",
		"تم تحويل مبلغ ١٠٠٠٠ دينار من الرقم ٠٧٨٠١١٢٢٣٣٤ بنجاح",
		"You have received 5,000 IQD credit from 9647805566778",
	}

	for i, sms := range incomingSMSList {
		rec, err := client.RecordIncomingTransferFromSMS(sms)
		if err != nil {
			fmt.Printf("    ❌ فشل معالجة الرسالة %d: %v\n", i+1, err)
			continue
		}
		fmt.Printf("    📩 رسالة واردة [%d]:\n", i+1)
		fmt.Printf("       - النص الأصلي : \"%s\"\n", sms)
		fmt.Printf("       - رقم المرسل  : %s\n", rec.MSISDN)
		fmt.Printf("       - المبلغ المحول: %s د.ع\n", rec.Amount)
		fmt.Printf("       - التوقيت     : %s\n\n", rec.CreatedAt)
	}

	// 4. Simulate Customer Verification Requests (Zero-Customer-Login Matching)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("[4] محاكاة طلبات تأكيد الدفع وإجراء المطابقة الآلية الذكية (Verification Flow):")
	fmt.Println("--------------------------------------------------------------------------------")

	ctx := context.Background()

	testScenarios := []struct {
		AccountID   string
		AccountName string
		InputPhone  string
		Label       string
	}{
		{
			AccountID:   "usr_1001",
			AccountName: "Ahmed Ali",
			InputPhone:  "07809876543", // Matched with SMS 1
			Label:       "زبون قام بالتحويل وكتب رقمه بالصيغة المحلية (078...)",
		},
		{
			AccountID:   "usr_1002",
			AccountName: "Omar Hassan",
			InputPhone:  "7801122334", // Matched with SMS 2 (without leading zero)
			Label:       "زبون قام بالتحويل وكتب رقمه بدون صفر (78...)",
		},
		{
			AccountID:   "usr_1003",
			AccountName: "Mustafa Tariq",
			InputPhone:  "07800000000", // No transfer sent
			Label:       "زبون لم يقم بالتحويل ويحاول تأكيد الدفع بدون رصيد",
		},
		{
			AccountID:   "usr_1004",
			AccountName: "Fraud Attempt",
			InputPhone:  "07809876543", // Attempting replay of SMS 1 (already claimed by usr_1001)
			Label:       "محاولة احتيال / إعادة استخدام نفس حوالة الزبون الأول لتفعيل حساب آخر",
		},
	}

	for _, sc := range testScenarios {
		fmt.Printf("👤 المشترك: %s (معرف الحساب: %s)\n", sc.AccountName, sc.AccountID)
		fmt.Printf("   الحالة: %s\n", sc.Label)
		fmt.Printf("   الرقم المدخل: %s\n", sc.InputPhone)

		user := ledger.GetOrCreateUser(sc.AccountID, sc.AccountName)

		// Verification execution: calls client.VerifyIncomingTransfer
		checkCtx, checkCancel := context.WithTimeout(ctx, 1*time.Second)
		verified, matchedRecord, _ := client.VerifyIncomingTransfer(checkCtx, sc.InputPhone, depositPrice)
		checkCancel()

		if !verified || matchedRecord == nil {
			fmt.Println("   ❌ النتيجة: لم يتم العثور على تحويل وارد مطابق!")
			fmt.Printf("      - يرجى التأكد من تحويل الرصيد أولاً عبر الكود: %s\n\n", ussdCode)
			continue
		}

		// Parse actual transferred amount
		var depositVal float64
		fmt.Sscanf(matchedRecord.Amount, "%f", &depositVal)
		if depositVal <= 0 {
			depositVal = depositPrice
		}

		// Check database to prevent double-spending
		recorded := ledger.RecordTransaction(user.AccountID, matchedRecord.MSISDN, depositVal, matchedRecord.CreatedAt)
		if !recorded {
			fmt.Println("   ⛔ النتيجة: تم رفض العملية! هذه الحوالة تم استخدامها مسبقاً لحساب آخر.")
			fmt.Println("      (Anti-Replay Protection: الحماية التلقائية من التكرار والاحتيال)")
			fmt.Println()
			continue
		}

		fmt.Println("   🎉 النتيجة: تم التحقق بنجاح وتفعيل الحساب فورياً!")
		fmt.Printf("      - المبلغ المودع : %s د.ع\n", matchedRecord.Amount)
		fmt.Printf("      - رصيد الحساب   : %.0f د.ع\n", user.Balance)
		fmt.Printf("      - حالة التفعيل : %t\n", user.IsActive)
		fmt.Printf("      - وقت التفعيل  : %s\n\n", user.ActivatedAt)
	}

	fmt.Println("================================================================================")
	fmt.Println("✅ اكتملت دورة التحقق والمطابقة الآلية بنجاح 100% وبدون أي تدخل بشري!")
	fmt.Println("================================================================================")
}
