<div align="center">

<img src="assets/zain.svg" alt="Zain Iraq Logo" width="240" />

# zainiraq-go

**Production-grade, zero-dependency Go SDK for Zain Iraq APIs**

[![Go Version](https://img.shields.io/badge/Go-1.26+-18181b?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Verified Endpoints](https://img.shields.io/badge/Endpoints-92%20Verified-833AB4?style=flat-square)](ENDPOINTS.md)
[![Dependencies](https://img.shields.io/badge/Dependencies-Zero%20(Stdlib)-18181b?style=flat-square)](https://pkg.go.dev/)
[![License](https://img.shields.io/badge/License-MIT-18181b?style=flat-square)](LICENSE)

<br />

مكتبة برمجية متكاملة، احترافية، وعالية الأداء مكتوبة بلغة **Go (Golang)** للتعامل مع واجهات برمجة تطبيقات شركة **زين العراق (Zain Iraq)**.  
المكتبة مبنية بنسبة **100% بالاعتماد على المكتبة القياسية للغة Go** وبدون أي مكاتب أو اعتمادات خارجية نهائياً.

</div>

---

## بنية المشروع (Project Structure)

```
zainiraq-go/
├── pkg/
│   └── zain/                    # حزمة الـ SDK الأساسية (بدون اعتمادات خارجية)
│       ├── client.go            # إعداد العميل، الترويسات، وإدارة الجلسات
│       ├── auth.go              # المصادقة، OTP، وتسجيل الدخول
│       ├── profile.go           # الملف الشخصي، الرصيد، والحسابات
│       ├── dashboard.go         # لوحة التحكم، الإشعارات، والفوترة
│       ├── recharge.go          # شحن الكروت، تحويل الرصيد، وتمديد الصلاحية
│       ├── verification.go      # تثبيت المحفظة والتحقق الآلي من التحويلات
│       ├── services.go          # باقات وعروض زين ونظام فليكس
│       ├── loyalty.go           # برنامج المكافآت وعروض امتياز
│       ├── payments.go          # بوابة الدفع الإلكتروني ومحفظة زين كاش
│       ├── support.go           # مركز الدعم وتذاكر الشكاوى الفنية
│       ├── content.go           # محرك استعلامات المحتوى CMS
│       └── types.go             # هياكل البيانات ونماذج الاستجابة JSON
├── examples/                    # أمثلة عملية مستقلة وتطبيق تفاعلي CLI
├── assets/                      # شعار زين العراق الرسمي
├── ENDPOINTS.md                 # التوثيق التقني لجميع نقاط النهاية الـ 92
├── README.md
└── go.mod
```

---

## ميزات وقدرات المكتبة (Features)

المكتبة تغطي **92 واجهة برمجية (Endpoints)** رسمية تم فحصها وتحليلها هندسياً من تطبيق زين العراق الرسمي:

### 1. المصادقة وإدارة الجلسات (Authentication & Sessions)
* **طلب رمز التحقق (SMS OTP)**: `RequestOTP(ctx, msisdn)` لإرسال رمز الدخول مباشرة للهاتف مع توليد `request_id`.
* **تأكيد رمز الدخول**: `ConfirmOTP(ctx, msisdn, code)` واستخراج معرف التوثيق `confirmation_id`.
* **استبدال التوثيق بتوكنات الدخول**: `SignSession(ctx, msisdn, confID)` لتوليد توكنات الجلسة الكاملة (`AccessToken` و `RefreshToken`).
* **تسجيل الدخول المجمع بخطوة واحدة**: `VerifyOTPAndLogin(ctx, msisdn, code)`.
* **تسجيل الدخول بكلمة المرور**: `LoginWithPassword(ctx, msisdn, password)`.
* **التجديد التلقائي للتوكن**: `RefreshToken(ctx)` لتجديد الجلسة دون انقطاع عبر ترويسة `refresh-token`.
* **إدارة الحساب وكلمة المرور**: `ResetPassword(ctx, oldPwd, newPwd)` و `ValidatePassword(ctx, pwd)`.
* **تخزين الجلسات واستعادتها**: `SaveSessionToFile(path)` و `LoadSessionFromFile(path)` لحفظ واستعادة الجلسات كملفات JSON دون الحاجة لطلب رمز SMS مجدداً.
* **تصدير واستيراد الكائنات**: `ExportSession()` و `ImportSession(session)` للتعامل البرمجي والبيئات الموزعة.
* **تسجيل الخروج وحذف الحساب**: `Logout(ctx)` لإبطال التوكن سحابياً و `DeleteAccount(ctx)`.

### 2. الرصيد، المحفظة والفوترة الآجلة (Balance, Wallet & Billing)
* **الملف الشخصي الشامل**: `GetProfile(ctx)` لجلب الاسم، الرقم، نوع الخط، ومعرف Galleon وحالة 4G.
* **ملخص الحساب المتكامل**: `GetSummary(ctx)` لرصيد الحساب، صلاحية الخط، الحصص الفعالة والخدمات.
* **رصيد المحفظة والصلاحية**: `GetBalance(ctx)` لمعرفة الرصيد النقدي وتاريخ انتهاء الخط بدقة.
* **سلفة الرصيد للطوارئ**: `GetLoan(ctx)` للاستعلام عن قيمة السلفة المستحقة وتاريخ استحقاقها.
* **الحسابات الفرعية ورصيد الباقات**: `GetSubaccounts(ctx)` لتفصيل رصيد الإنترنت (4G Capped/Unlimited)، الدقائق والرسائل.
* **فواتير الخطوط الآجلة الدفع (Postpaid)**:
  * `GetBill(ctx)`: إجمالي الفاتورة الحالية والمبالغ غير المفوترة والمستحقات السابقة.
  * `GetBillItems(ctx)`: تفاصيل استهلاك بنود الفاتورة.
  * `GetUnbilled(ctx)`: الاستهلاك المفتوح خارج الفاتورة.
  * `GetAdvancePayment(ctx)`: الدفعات النقدية المقدمة.
  * `GetElectronicBillItems(ctx)`: بنود الفاتورة الإلكترونية المعتمدة.
  * `GetPostpaidHistory(ctx)`: سجل الفواتير والدفعات السابقة.
* **إدارة لغة الخط**: `ChangeLanguage(ctx, lang)` لتعيين لغة الرسائل النصية والإشعارات (`ar`, `en`, `kd`).

### 3. شحن الرصيد وتحويل الأموال وتثبيت المحفظة (Recharge, Transfer & Master Wallet)
* **تثبيت رقم المحفظة المعتمد (Master Wallet)**:
  * `client.SetMasterWallet(phone)`: تثبيت رقم محفظة البوت لاستقبال التحويلات.
  * `client.MasterWallet()`: جلب رقم المحفظة الحالي.
  * `client.GetWalletOverview(ctx)`: نظرة شاملة لرصيد المحفظة، الصلاحية، ونوع الخط.
  * `client.GetWalletBalance(ctx)`: استعلام رصيد المحفظة المباشر عبر `api/number/wallet`.
* **التحقق الآلي من الحوالات الواردة ومعرفة رقم المرسل (Automated Incoming Verification)**:
  * `client.VerifyIncomingTransfer(ctx, senderPhone, minAmount)`: فحص فوري ومطابقة تلقائية لسجلات التحويل الوارد من رقم المشترك وتأكيد دفع الطلب آلياً بدون أي تدخل بشري.
  * `client.GetIncomingTransfers(ctx, limit)`: استخراج كافة الحوالات الواردة من بنود الفاتورة الإلكترونية (`api/number/electronic-bill-items`) وإشعارات الرسائل (`api/notifications`).
  * `client.FormatUSSDTransfer(recipient, amount)`: توليد كود التحويل السريع لزين العراق (`*123*amount*recipient#`).
* **تحويل الرصيد النقدي (P2P Credit Transfer)**:
  * `client.RequestCreditTransferOTP(ctx, senderMSISDN)`: طلب كود تحقق SMS لعملية التحويل.
  * `client.ConfirmCreditTransferOTP(ctx, otpCode, senderMSISDN)`: تأكيد الرمز واستخراج توكن التحويل.
  * `client.CreditTransfer(ctx, recipient, amount, otpConfirmation)`: إرسال الرصيد الفعلي للمشترك الآخر.
* **شحن كروت الرصيد الورقية**: `RechargeVoucher(ctx, voucherPIN)` لشحن الكروت ذات الـ 16 رقماً مع استلام الرصيد الجديد فورياً.
* **تمديد صلاحية استقبال وإرسال الخط**: `ExtendValidity(ctx, amount)` و `GetValidityOptions(ctx)`.
* **بوابة الدفع الإلكتروني (Card & Checkout)**:
  * `CreateCheckoutID(ctx, req)`: إنشاء معرف Checkout ID لبطاقات ماستركارد وفيزا.
  * `RefreshPaymentStatus(ctx, checkoutID)`: فحص حالة عملية الدفع والتأكد من إتمامها.
  * `SaveCard(ctx, req)`: حفظ وتشفير بيانات البطاقة لاستخدامها مستقبلاً.
  * `SetDefaultCard(ctx, cardID)` و `GetUserCards(ctx)` و `DeleteCard(ctx, cardID)`.
* **محفظة زين كاش (ZainCash Direct Purchase)**:
  * `InitiateZainCashPayment(ctx, req)` و `InitiateZainCashPaymentV2(ctx, req)`: بدء الدفع المباشر عبر ZainCash.

### 4. العروض والباقات ونظام فليكس (Offers, Bundles & Flex)
* **كتالوج العروض المعتمد**: `GetOffersCMS(ctx, queryShortname)` لجلب باقات الإنترنت والمكالمات.
* **عروض فليكس والترقية**: `GetFlexBundles(ctx, sourceOffer, bundleType)` لاستعلام خيارات الترقية (Flex Upsell).
* **عروض مخصصة فوق وتحت الخط (ATL / BTL)**:
  * `GetPersonalizedOffersATL(ctx)`: العروض الترويجية العامة فوق الخط.
  * `GetPersonalizedOffersBTL(ctx)`: العروض المصممة خصيصاً لرقم المشترك وسجل استهلاكه.
* **الاشتراكات الفعالة وإلغاؤها**:
  * `GetSubscriptions(ctx)`: قائمة بجميع الباقات المفعلة ومواعيد تجديدها التلقائي.
  * `SubscribeOffer(ctx, offerID)`: تفعيل الباقة الفوري بخصم من الرصيد.
  * `UnsubscribeOffer(ctx, offerID)`: إلغاء الاشتراك وإيقاف التجديد التلقائي للباقة.
* **الهدايا والمكافآت وبرنامج كفو**:
  * `ClaimDailyGift(ctx)`: استلام الهدية اليومية المجانية المتاحة للمشترك.
  * `SendBundleGift(ctx, recipient, offerID)`: إهداء باقة لرقم آخر مع الخصم من رصيد الحساب.
  * `RedeemRegistrationGift(ctx, offerID)`: استلام هدية الترحيب للمشتركين الجدد.
  * `InviteToKafoo(ctx, invitedMSISDN)`: إرسال دعوة برنامج كفو (Kafoo Referral) للأصدقاء.
* **إدارة خطوط فليكس التراكمية**:
  * `GetFlexStatus(ctx)`: استعلام رصيد نقاط وحصص فليكس.
  * `GetFlexLimits(ctx)`: استعلام السقوف القصوى المسموح بها للاستهلاك.
  * `MigrateToFlex(ctx, targetOffer)`: تحويل الخط إلى نظام باقات فليكس.

### 5. مشاركة الباقات والأرقام المفضلة (Bundle Sharing & FNF)
* **مشاركة سعة الإنترنت (Bundle Sharing)**:
  * `SharingAddMember(ctx, memberMSISDN, offerID, quotaMB)`: إضافة خط للمجموعة وتحديد حصته.
  * `SharingQueryMembers(ctx, offerID)`: استعلام الأعضاء والحصص المحددة والمستهلكة.
  * `SharingRemoveMember(ctx, memberMSISDN, offerID)`: إزالة رقم من مجموعة المشاركة.
  * `SharingTransferUnits(ctx, memberMSISDN, units)`: تحويل ميغابايت إضافية لعضو المجموعة.
* **الأرقام المفضلة (Friends & Family)**:
  * `GetFriendsAndFamily(ctx)`: استعلام قائمة أرقام الأصدقاء والعائلة.
  * `AddFriendsAndFamily(ctx, msisdn)`: إضافة رقم جديد للاستفادة من التخفيض.
  * `RemoveFriendsAndFamily(ctx, msisdn)`: حذف رقم من القائمة.

### 6. برنامج المكافآت ونقاط امتياز (Loyalty & Imtiyaz)
* **رصيد برنامج الولاء**: `GetLoyaltyInfo(ctx)` لاستعلام رصيد النقاط، الفئة الحالية (Gold/Platinum)، والنقاط القابلة للاستبدال.
* **سجل المعاملات**: `GetLoyaltyHistory(ctx)` لتتبع عمليات اكتساب وصرف النقاط.
* **عروض الاستبدال المميزة**: `GetLoyaltyHotBundles(ctx)`.
* **استبدال النقاط**:
  * `RedeemLoyaltyCredit(ctx, req)`: استبدال النقاط برصيد نقدي على الشريحة.
  * `RedeemLoyaltyOffer(ctx, req)`: استبدال النقاط بباقات إنترنت ومكالمات مجانية.
  * `RedeemLoyaltyPromoCode(ctx, req)`: استبدال النقاط بقسائم خصم تجارية.
* **برنامج امتياز (Imtiyaz Partners)**:
  * `GetImtiyazCategories(ctx)`: أقسام الشركاء (مطاعم، فنادق، مراكز تجارية، ترفيه).
  * `GetImtiyazMerchants(ctx, category)`: قائمة المتاجر ونسب الخصومات المتاحة.
  * `GetMerchant(ctx, category, id)`: تفاصيل فروع المتجر، ساعات العمل والشروط.
  * `RedeemImtiyazDiscount(ctx, merchantID)`: توليد كود أو باركود الخصم الفوري.

### 7. مركز الدعم وتتبع الشكاوى الفنية (Support & Complaints)
* **أقسام الشكاوى**: `GetComplaintCategories(ctx)` و `GetComplaintSubCategories(ctx, catID)`.
* **سجل التذاكر المفتوحة**: `GetTickets(ctx)` لاستعلام جميع الشكاوى السابقة وحالتها.
* **تتبع تفاصيل التذكرة**: `GetTicketDetails(ctx, ticketID)` لمتابعة ردود الدعم الفني ومسار الحل.
* **فتح تذكرة شكوى جديدة (Multipart)**: `CreateTicket(ctx, req)` مع دعم رفع مرفقات وصور فنية للأعطال وإحداثيات الموقع (GPS).
* **إعادة فتح التذكرة**: `ReopenTicket(ctx, ticketID, reason)`.

### 8. لوحة التحكم والإشعارات والمتاجر القريبة (Dashboard, Notifications & Near Me)
* **صندوق الإشعارات**: `GetNotifications(ctx, offset, limit, isRead)` مع دعم تصفية المقروء وغير المقروء.
* **تحديث حالة الإشعار**: `MarkNotificationsRead(ctx, ids)`.
* **القصص والبانرات التفاعلية**: `GetStories(ctx)` و `GetBanners(ctx)`.
* **فروع وموزعي زين القريبين**: `GetNearMeShops(ctx, lat, lng, radius)` للبحث الجغرافي بالأقرب مسافة.
* **المدن والمحافظات**: `GetCities(ctx)`.
* **الخدمات الرقمية والترفيهية**: `GetDigitalServices(ctx)` و `SubscribeDigitalService(ctx, svcID)` و `UnsubscribeDigitalService(ctx, svcID)`.

### 9. محرك استعلامات المحتوى CMS (Direct CMS Query Engine)
* استعلامات سحابية مباشرة فائقة السرعة على خادم `https://cms-mobile.iq.zain.com`:
  * `QueryCMS(ctx, space, body)`: تنفيذ استعلام حر على مساحات Galleon أو Products.
  * `GetAppConfigurations(ctx)`: قراءة متغيرات ضبط التطبيق والـ Feature Flags الحية.
  * `GetSubaccountsCMS(ctx)`: استعلام شجرة تصنيفات الحسابات الفرعية وتعريفات الحصص.
  * `GetOffersCatalogCMS(ctx)`: جلب شجرة كتالوج العروض الكاملة.
  * `GetFAQsCMS(ctx)`: جلب الأسئلة الشائعة باللغات الثلاث (عربي، إنكليزي، كردي).
  * `GetRoamingCMS(ctx)`: استعلام تعرفة وشبكات التجوال الدولي للدول حول العالم.

---

## البدء السريع (Quick Start)

### 1. تثبيت الحزمة

```bash
go get github.com/FLEX-GHOST/zainiraq-go/pkg/zain
```

### 2. مثال تسجيل الدخول وحفظ الجلسة

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client, err := zain.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// طلب إرسال رمز التحقق
	otpResp, err := client.RequestOTP(ctx, "07801234567")
	if err != nil {
		log.Fatalf("request otp failed: %v", err)
	}
	fmt.Printf("OTP sent successfully! Request ID: %s\n", otpResp.Data.RequestID)

	// تأكيد الرمز المكون من 6 أرقام
	session, err := client.VerifyOTPAndLogin(ctx, "07801234567", "123456")
	if err != nil {
		log.Fatalf("verification failed: %v", err)
	}
	fmt.Printf("Logged in successfully! User Space: %s\n", session.UserSpace)

	// حفظ الجلسة للاستخدام المستقبلي
	if err := client.SaveSessionToFile("session.json"); err != nil {
		log.Fatalf("failed to save session: %v", err)
	}
	fmt.Println("Session saved to session.json!")
}
```

### 3. فحص الرصيد والباقات بجلسة سابقة

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FLEX-GHOST/zainiraq-go/pkg/zain"
)

func main() {
	client, err := zain.NewClient()
	if err != nil {
		log.Fatalf("client error: %v", err)
	}

	if err := client.LoadSessionFromFile("session.json"); err != nil {
		log.Fatalf("session load error: %v", err)
	}

	ctx := context.Background()
	balance, err := client.GetBalance(ctx)
	if err != nil {
		log.Fatalf("balance error: %v", err)
	}

	fmt.Printf("الرصيد الحالي: %d د.ع\n", balance.Balance.Value)
	fmt.Printf("تاريخ الصلاحية: %s\n", balance.Balance.Expiry)
}
```

---

## الأمان والمواصفات الهندسية (Architecture & Security)

* **صفر اعتمادات خارجية (100% Standard Library)**: الاعتماد الحصري على مكتبات Go القياسية (`net/http`, `crypto/tls`, `context`, `encoding/json`, `sync`) مما يضمن سرعة البرمجة والترجمة وأعلى درجات الأمان من ثغرات سلاسل التوريد.
* **إعادة استخدام الاتصالات وحماية مقابس TCP (Socket Leak Prevention)**: كل استجابة HTTP يتم تصريفها بالكامل عبر `io.Copy(io.Discard, resp.Body)` قبل إغلاقها `defer resp.Body.Close()` لضمان إعادة استخدام اتصالات TCP/TLS دون تسريب المقابس.
* **مجمّع اتصالات متقدم (HTTP/2 Connection Pooling)**: استخدام `http.Transport` مخصص مع ضبط `MaxIdleConns: 100`, `IdleConnTimeout: 90s`, و `TLSHandshakeTimeout: 10s`.
* **انضباط المهلات الزمنية (Context-Bound I/O)**: جميع عمليات الإدخال والإخراج تدعم تمرير `context.Context` مع إمكانية إلغاء الطلبات وتعيين مهلات قصوى (`Timeout`).
* **أمان التزامن الكامل (Concurrency-Safe)**: العميل محمي بواسطة `sync.RWMutex` لإتاحة مشاركة نفس الكائن بأمان تام عبر آلاف الـ Goroutines المتزامنة في بوتات التلغرام وتطبيقات الخوادم.
* **معالجة الأخطاء النموذجية (Strongly-Typed Errors)**: تمثيل أخطاء الخادم عبر نوع بيانات موحد `*APIError` يحتوي على كود الحالة HTTP، رسالة الخطأ الأصلية، والترجمة المحلية.

---

## دليل الأمثلة الجاهزة (Examples Suite)

| المجلد | الوصف | أمر التشغيل المباشر |
| :--- | :--- | :--- |
| **`01_otp_login`** | تسجيل الدخول عبر رمز التحقق SMS، حفظ الجلسة واستعادتها. | `go run examples/01_otp_login/main.go` |
| **`02_account_and_profile`** | الاستعلام عن الرصيد، تفاصيل الحساب، والحصص الفعالة. | `go run examples/02_account_and_profile/main.go` |
| **`03_recharge_and_transfer`** | شحن الرصيد بكروت التعبئة، تحويل الرصيد بالـ OTP، وتمديد الصلاحية. | `go run examples/03_recharge_and_transfer/main.go` |
| **`04_bundles_and_offers`** | استعراض الباقات، العروض المخصصة، ونظام فليكس. | `go run examples/04_bundles_and_offers/main.go` |
| **`05_loyalty_and_imtiyaz`** | نقاط برنامج المكافآت، استبدال النقاط، وعروض امتياز. | `go run examples/05_loyalty_and_imtiyaz/main.go` |
| **`06_support_and_tickets`** | تذاكر الدعم الفني، إرفاق الصور، وتتبع المعالجة. | `go run examples/06_support_and_tickets/main.go` |
| **`07_cms_content_queries`** | استعلامات محرك المحتوى CMS لكتالوج العروض والضبط. | `go run examples/07_cms_content_queries/main.go` |
| **`08_payments_and_zaincash`** | بوابات الدفع الإلكتروني، البطاقات، ومحفظة زين كاش. | `go run examples/08_payments_and_zaincash/main.go` |
| **`09_bundle_sharing_and_fnf`** | مشاركة الباقات العائلية، وتحديد الحصص ونقل الوحدات. | `go run examples/09_bundle_sharing_and_fnf/main.go` |
| **`10_nearme_and_notifications`** | فروع زين القريبة، الإشعارات، والخدمات الرقمية. | `go run examples/10_nearme_and_notifications/main.go` |
| **`11_wallet_and_incoming_transfer_verification`** | تثبيت المحفظة والتحقق الآلي من تحويلات الرصيد لبوتات التليغرام. | `go run examples/11_wallet_and_incoming_transfer_verification/main.go` |
| **`12_daily_gift_and_rewards_automation`** | أتمتة سحب الهدايا اليومية، تحويل نقاط المكافآت، وفحص الصلاحية. | `go run examples/12_daily_gift_and_rewards_automation/main.go` |
| **`interactive_cli`** | تطبيق تيرمينال تفاعلي شامل يتيح تجربة جميع ميزات المكتبة عبر قائمة نصية مرئية. | `go run examples/interactive_cli/main.go` |

---

## الاختبارات وضمان الجودة (Testing & Quality)

```bash
# تشغيل جميع اختبارات الحزمة
go test -v ./...

# التحقق من فحص مفسر Go القياسي
go vet ./...
```

---

## الترخيص وإخلاء المسؤولية (License & Legal Disclaimer)

- **الترخيص**: هذا المشروع مرخص ومفتوح المصدر بموجب رخصة [MIT](LICENSE).
- **العلامة التجارية**: اسم "Zain" وشعارها علامتان تجاريتان مسجلتان لمجموعة زين للاتصالات (Zain Group).
- **إخلاء المسؤولية**: هذا المشروع (`zainiraq-go`) هو مكتبة برمجية مستقلة غير رسمية تم تطويرها لأغراض تعليمية وتطويرية، وليست تابعة لشركة زين العراق أو معتمدة منها بشكل رسمي.
