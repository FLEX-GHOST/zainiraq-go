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

## جدول نقاط النهاية المعتمدة (Endpoints Matrix - 92 APIs)

المكتبة تغطي **92 واجهة برمجية (Endpoints)** رسمية تم فحصها وتحليلها هندسياً من تطبيق زين العراق الرسمي:

| # | الطريقة | المسار (Endpoint Path) | دالة Go SDK المقابلة | الوصف التفصيلي باللغة العربية |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `POST` | `/api/otp/request` | `client.RequestOTP(ctx, msisdn)` | طلب إرسال رمز التحقق OTP برسالة نصية SMS وتوليد `request_id` |
| **02** | `POST` | `/api/otp/confirm` | `client.ConfirmOTP(ctx, msisdn, code)` | تأكيد رمز التحقق المدخل من المشترك وتوليد معرف التوثيق `confirmation_id` |
| **03** | `POST` | `/api/user/sign` | `client.SignSession(ctx, msisdn, conf)` | استبدال معرف التوثيق بتوكنات الجلسة الكاملة (`access_token` و `refresh_token`) |
| **04** | `POST` | `/api/user/login` | `client.LoginWithPassword(ctx, msisdn, pwd)` | تسجيل الدخول المباشر بالرقم وكلمة المرور المشفرة |
| **05** | `POST` | `/api/user/token` | `client.RefreshToken(ctx)` | تجديد توكن الوصول المنتهي تلقائياً في الخلفية باستخدام توكن التحديث |
| **06** | `GET` | `/api/v2/user/profile` | `client.GetProfile(ctx)` | جلب الملف الشخصي الكامل: الاسم، الرقم، نوع الخط، ومعرف Galleon |
| **07** | `PATCH`| `/api/v2/user/profile` | `client.UpdateProfile(ctx, req)` | تعديل وتحديث بيانات المشترك (الاسم، البريد، الإعدادات) |
| **08** | `POST` | `/api/user/create` | `client.SignUp(ctx, req)` | تسجيل حساب جديد لمشترك زين وتعيين بيانات الدخول |
| **09** | `POST` | `/api/user/validate` | `client.ValidatePassword(ctx, pwd)` | التحقق من مطابقة كلمة المرور الحالية للمشترك قبل العمليات الحساسة |
| **10** | `POST` | `/api/user/reset_password` | `client.ResetPassword(ctx, old, new)` | تغيير وإعادة تعيين كلمة مرور الحساب |
| **11** | `DELETE`| `/api/user/logout` | `client.Logout(ctx)` | تسجيل الخروج وإبطال صلاحية التوكن على الخادم السحابي |
| **12** | `DELETE`| `/api/user/delete` | `client.DeleteAccount(ctx)` | حذف وإغلاق حساب المشترك نهائياً من نظام التطبيق |
| **13** | `POST` | `/api/user/feedback` | `client.SubmitFeedback(ctx, rate, comment)` | إرسال تقييم المشترك وملاحظاته حول تجربة التطبيق |
| **14** | `GET` | `/api/number/summary` | `client.GetSummary(ctx)` | ملخص شامل للحساب: الرصيد، الحصص الفعالة، والخدمات المشتركة |
| **15** | `GET` | `/api/number/wallet` | `client.GetBalance(ctx)` / `client.GetWalletBalance(ctx)` | استعلام رصيد المحفظة الأساسي وتاريخ انتهاء صلاحية الخط بدقة |
| **16** | `GET` | `/api/number/loan` | `client.GetLoan(ctx)` | استعلام تفاصيل سلفة الرصيد للطوارئ والمبلغ المستحق |
| **17** | `GET` | `/api/number/postpaid-history`| `client.GetPostpaidHistory(ctx)` | سجل الفواتير والدفعات السابقة للخطوط الآجلة الدفع |
| **18** | `GET` | `/api/number/subaccounts` | `client.GetSubaccounts(ctx)` | جلب الحسابات الفرعية ورصيد البيانات (إنترنت، مكالمات، رسائل) |
| **19** | `GET` | `/api/number/dashboard_message` | `client.GetDashboardMessages(ctx)` | استعلام رسائل التنبيه والاشعارات الموجهة للخط في الواجهة |
| **20** | `GET` | `/api/number/query-bill` | `client.GetBill(ctx)` | استعلام الفاتورة الحالية، المبالغ غير المفوترة، والمستحقات السابقة |
| **21** | `GET` | `/api/number/query-bill-items` | `client.GetBillItems(ctx)` | تفاصيل بنود وبنود استهلاك الفاتورة بالتفصيل |
| **22** | `GET` | `/api/number/query-unbilled` | `client.GetUnbilled(ctx)` | استعلام الاستهلاك المفتوح خارج الفاتورة قبل صدورها |
| **23** | `GET` | `/api/number/advance-payment` | `client.GetAdvancePayment(ctx)` | تفاصيل المبالغ المدفوعة مقدماً ورصيد التسديد المستقبلي |
| **24** | `POST` | `/api/number/change-language` | `client.ChangeLanguage(ctx, lang)` | تغيير لغة الإشعارات والرسائل النصية للنظام (عربي، كردي، إنكليزي) |
| **25** | `GET` | `/api/number/electronic-bill-items` | `client.GetCDRTransferHistory(ctx)` / `client.GetElectronicBillItems(ctx)` | كشف الحساب وسجل تحويلات الرصيد الواردة (CDR) لمعرفة رقم المرسل والمبلغ وتأكيد الدفع التلقائي بدون تسجيل دخول الزبون |
| **26** | `POST` | `/api/number/charge-voucher` | `client.RechargeVoucher(ctx, pin)` | شحن وتعبئة الرصيد بكارت الشحن الورقي (16 رقماً) |
| **27** | `POST` | `/api/number/credit-transfer` | `client.CreditTransfer(ctx, to, amt, otp)` | تحويل رصيد نقدي من رقم إلى رقم آخر في شبكة زين العراق |
| **28** | `POST` | `/api/number/extend-validity` | `client.ExtendValidity(ctx, amount)` | تمديد صلاحية استقبال وإرسال الخط بخصم من الرصيد |
| **29** | `POST` | `/api/payment/create-checkout-id` | `client.CreateCheckoutID(ctx, req)` | توليد معرّف الدفع Checkout ID لبوابة الدفع الإلكتروني والبطاقات |
| **30** | `POST` | `/api/payment/refresh-payment-status` | `client.RefreshPaymentStatus(ctx, chkId)` | التحقق من نجاح أو فشل معاملة الدفع الإلكتروني بعد اكتمالها |
| **31** | `POST` | `/api/payment/save-card` | `client.SaveCard(ctx, req)` | حفظ بطاقة الدفع (ماستركارد/فيزا) المشفرة لاستخدامها مستقبلاً |
| **32** | `POST` | `/api/payment/set-default-card` | `client.SetDefaultCard(ctx, cardId)` | تعيين بطاقة دفع معينة كخيار افتراضي في الحساب |
| **33** | `GET` | `/api/payment/user-cards` | `client.GetUserCards(ctx)` | جلب قائمة البطاقات المصرفية المحفوظة للمشترك |
| **34** | `DELETE`| `/api/payment/delete-card` | `client.DeleteCard(ctx, cardId)` | فك ربط وحذف بطاقة مصرفية محفوظة من الحساب |
| **35** | `POST` | `/api/payment/purchase-order` | `client.CreatePurchaseOrder(ctx, req)` | إنشاء أمر شراء ودفع مباشر عبر محفظة زين كاش (ZainCash) |
| **36** | `GET` | `/api/offers` | `client.GetOffersCMS(ctx, queryShortname)` | جلب كتالوج العروض والأسعار من محرك المحتوى CMS |
| **37** | `GET` | `/api/offers/{source_offer}` | `client.GetFlexBundles(ctx, offer, type)` | استعلام خيارات ترقية الباقات المتاحة للعرض الحالي (Flex Upsell) |
| **38** | `GET` | `/api/number/personalized-offers?provider_type=ATL` | `client.GetPersonalizedOffersATL(ctx)` | جلب العروض الترويجية العامة فوق الخط (Above-The-Line) |
| **39** | `GET` | `/api/number/personalized-offers?provider_type=BTL` | `client.GetPersonalizedOffersBTL(ctx)` | جلب العروض المخصصة الموجهة خصيصاً لرقم المشترك (BTL) |
| **40** | `GET` | `/api/number/subscriptions` | `client.GetSubscriptions(ctx)` | استعلام جميع الاشتراكات والباقات الفعالة على الخط ومواعيد تجديدها |
| **41** | `POST` | `/api/number/subscribe` | `client.SubscribeOffer(ctx, offerId)` | تفعيل والاشتراك الفوري في باقة معينة وخصم قيمتها من الرصيد |
| **42** | `DELETE`| `/api/number/unsubscribe` | `client.UnsubscribeOffer(ctx, offerId)` | إلغاء الاشتراك وإيقاف التجديد التلقائي للباقة الفعالة |
| **43** | `POST` | `/api/number/claim-daily-gift` | `client.ClaimDailyGift(ctx)` | استلام والمطالبة بالهدية المجانية اليومية المتاحة للرقم |
| **44** | `POST` | `/api/number/send-gift` | `client.SendBundleGift(ctx, to, offerId)` | إهداء باقة إنترنت أو دقائق لرقم مشترك آخر على شبكة زين |
| **45** | `POST` | `/api/number/redeem-registration-gift` | `client.RedeemRegistrationGift(ctx, offerId)` | استلام هدية التسجيل والترحيب للمشتركين الجدد |
| **46** | `POST` | `/api/number/kafoo_invite` | `client.InviteToKafoo(ctx, invitedMSISDN)` | إرسال دعوة برنامج كفو (Kafoo Referral) لرقم صديق |
| **47** | `GET` | `/api/number/flex-status` | `client.GetFlexStatus(ctx)` | استعلام حالة خط فليكس التراكمي، نقاط الاستهلاك، وصلاحية الباقة |
| **48** | `GET` | `/api/number/flex-limits` | `client.GetFlexLimits(ctx)` | استعلام حدود وسقوف استهلاك الوحدات المسموح بها لباقة فليكس |
| **49** | `POST` | `/api/number/migrate-to-flex` | `client.MigrateToFlex(ctx, targetOffer)` | تحويل خط المشترك إلى باقة فليكس المتكاملة |
| **50** | `POST` | `/api/sharing/addrsc/` | `client.SharingAddMember(ctx, to, off, quota)` | إضافة رقم مشترك جديد لمجموعة مشاركة باقة الإنترنت العائلية |
| **51** | `GET` | `/api/sharing/query` | `client.SharingQueryMembers(ctx, offerId)` | استعلام قائمة الأرقام المنضمة للمجموعة والحصص المحددة والمستهلكة |
| **52** | `POST` | `/api/sharing/remove/` | `client.SharingRemoveMember(ctx, to, offerId)` | إزالة وحذف رقم من مجموعة مشاركة الباقة |
| **53** | `POST` | `/api/sharing/transfer_unit` | `client.SharingTransferUnits(ctx, to, units)` | تحويل ونقل وحدات ميغابايت إضافية لأحد أعضاء المجموعة |
| **54** | `GET` | `/api/loyalty/info` | `client.GetLoyaltyInfo(ctx)` | استعلام رصيد نقاط الولاء، فئة المشترك (Gold/Platinum)، والنقاط القابلة للصرف |
| **55** | `GET` | `/api/loyalty/history` | `client.GetLoyaltyHistory(ctx)` | سجل وتاريخ عمليات اكتساب واستبدال نقاط برنامج الولاء |
| **56** | `GET` | `/api/loyalty/hot-bundles` | `client.GetLoyaltyHotBundles(ctx)` | الباقات والعروض المميزة المتاحة للاستبدال بالنقاط مباشرة |
| **57** | `POST` | `/api/loyalty/rewards/credit` | `client.RedeemLoyaltyCredit(ctx, req)` | استبدال نقاط الولاء برصيد نقدي يضاف لمجمل رصيد الشريحة |
| **58** | `POST` | `/api/loyalty/rewards/offers` | `client.RedeemLoyaltyOffer(ctx, req)` | استبدال نقاط الولاء بباقات إنترنت ومكالمات مجانية |
| **59** | `POST` | `/api/loyalty/rewards/promo-code` | `client.RedeemLoyaltyPromoCode(ctx, req)` | استبدال نقاط الولاء بكوبونات خصم وقسائم شراء رقمية |
| **60** | `GET` | `/api/loyalty/faqs` | `client.GetLoyaltyFAQs(ctx)` | الأسئلة الشائعة وقواعد استخدام برنامج ولاء زين باللغات الثلاث |
| **61** | `GET` | `/api/imtiyaz/categories` | `client.GetImtiyazCategories(ctx)` | تصنيفات وأقسام شركاء برنامج امتياز (مطاعم، تسوق، فنادق، صحة) |
| **62** | `GET` | `/api/imtiyaz/merchants` | `client.GetImtiyazMerchants(ctx, category)` | قائمة المتاجر والشركاء المعتمدين ونسب الخصومات المقدمة |
| **63** | `GET` | `/api/imtiyaz/merchant/{id}` | `client.GetMerchant(ctx, cat, id)` | تفاصيل المتجر المعتمد، الفروع، العناوين، وشروط العرض |
| **64** | `POST` | `/api/imtiyaz/redeem` | `client.RedeemImtiyazDiscount(ctx, merchantId)` | توليد رمز الخصم الحصري (QR / Barcode) للمشترك لإبرازه في المتجر |
| **65** | `GET` | `/api/complaints/categories` | `client.GetComplaintCategories(ctx)` | جلب تصنيفات وأقسام الشكاوى الفنية المعتمدة في زين العراق |
| **66** | `GET` | `/api/complaints/sub-categories` | `client.GetComplaintSubCategories(ctx, catId)` | جلب التصنيفات الفرعية وأسباب المشاكل التابعة لكل قسم |
| **67** | `GET` | `/api/complaints/tickets` | `client.GetTickets(ctx)` | استعلام سجل تذاكر الشكاوى المفتوحة والسابقة وحالتها الحالية |
| **68** | `GET` | `/api/complaints/ticket/{id}` | `client.GetTicketDetails(ctx, ticketId)` | تتبع مسار معالجة تذكرة محددة، ردود فريق الدعم، وتحديثات الحل |
| **69** | `POST` | `/api/complaints/create-ticket` | `client.CreateTicket(ctx, req)` | تقديم تذكرة شكوى رسمية جديدة مع إمكانية رفع مرفقات وصور فنية |
| **70** | `POST` | `/api/complaints/reopen-ticket` | `client.ReopenTicket(ctx, ticketId, reason)` | إعادة فتح تذكرة شكوى مغلقة في حال عدم حل المشكلة بصورة مرضية |
| **71** | `GET` | `/api/dashboard/notifications` | `client.GetNotifications(ctx, off, lim, read)` | جلب صندوق الإشعارات والتنبيهات الموجهة للمستخدم مع التصفية |
| **72** | `POST` | `/api/dashboard/notifications/mark-read` | `client.MarkNotificationsRead(ctx, ids)` | تحديد وتحديث حالة الإشعارات كمقروءة على الخادم |
| **73** | `GET` | `/api/dashboard/stories` | `client.GetStories(ctx)` | جلب القصص الإعلانية والتفاعلية القصيرة (In-App Stories) |
| **74** | `GET` | `/api/dashboard/banners` | `client.GetBanners(ctx)` | جلب اللوحات الإعلانية والبانرات الترويجية لشاشة التطبيق الرئيسية |
| **75** | `GET` | `/api/nearme/shops` | `client.GetNearMeShops(ctx, lat, lng, rad)` | البحث عن مراكز وفروع زين وموزعيها المعتمدين الأقرب جغرافياً |
| **76** | `GET` | `/api/nearme/cities` | `client.GetCities(ctx)` | قائمة بجميع المحافظات والمدن العراقية المدعومة ومراكزها |
| **77** | `GET` | `/api/digital-services` | `client.GetDigitalServices(ctx)` | قائمة الخدمات الرقمية والترفيهية (شاهد، أنغامي، ألعاب) |
| **78** | `POST` | `/api/digital-services/subscribe` | `client.SubscribeDigitalService(ctx, svcId)` | الاشتراك في خدمة رقمية مع خصم الرسوم من الرصيد أو الفاتورة |
| **79** | `POST` | `/api/digital-services/unsubscribe` | `client.UnsubscribeDigitalService(ctx, svcId)` | إلغاء الاشتراك في الخدمة الترفيهية الفعالة وإيقاف التجديد |
| **80** | `POST` | `/dmart/public/excute/query/galleon` | `client.QueryCMS(ctx, "galleon", body)` | تنفيذ استعلام مباشر على مساحة Galleon في محرك المحتوى CMS |
| **81** | `POST` | `/dmart/public/excute/query/products` | `client.QueryCMS(ctx, "products", body)` | استعلام كتالوج المنتجات والأسعار والعروض المحدثة لحظياً |
| **82** | `POST` | `/dmart/public/excute/query/app_configurations` | `client.GetAppConfigurations(ctx)` | استخراج متغيرات وبيانات ضبط التطبيق وقيم الميزات (Feature Flags) |
| **83** | `POST` | `/dmart/public/excute/query/subaccounts` | `client.GetSubaccountsCMS(ctx)` | استعلام مخطط وتفاصيل الحسابات الفرعية وتعريفات الحصص من CMS |
| **84** | `POST` | `/dmart/public/excute/query/offers` | `client.GetOffersCatalogCMS(ctx)` | استعلام شامل لكامل شجرة العروض والخصومات المتاحة على الشبكة |
| **85** | `POST` | `/dmart/public/excute/query/faqs` | `client.GetFAQsCMS(ctx)` | جلب قاعدة المعرفة والأسئلة الأكثر شيوعاً باللغات الثلاث |
| **86** | `POST` | `/dmart/public/excute/query/roaming` | `client.GetRoamingCMS(ctx)` | استعلام قائمة الدول، الشبكات الشريكة، وتعرفة وباقات التجوال الدولي |
| **87** | `GET` | `/api/number/friends-and-family` | `client.GetFriendsAndFamily(ctx)` | استعلام قائمة أرقام الأصدقاء والعائلة المضافة للاستفادة من التخفيض |
| **88** | `POST` | `/api/number/friends-and-family/add` | `client.AddFriendsAndFamily(ctx, msisdn)` | إضافة رقم مفضل جديد لقائمة الأصدقاء والعائلة |
| **89** | `DELETE`| `/api/number/friends-and-family/remove`| `client.RemoveFriendsAndFamily(ctx, msisdn)` | حذف رقم من قائمة الأصدقاء والعائلة |
| **90** | `GET` | `/api/number/esim-details` | `client.GetESIMDetails(ctx)` | جلب بيانات الشريحة الإلكترونية eSIM ورمز التفعيل QR |
| **91** | `POST` | `/api/number/swap-sim` | `client.RequestSIMSwap(ctx, iccid)` | طلب استبدال وتفعيل الشريحة الجديدة برقم البطاقة ICCID |
| **92** | `GET` | `/api/system/time` | `client.GetServerTime(ctx)` | استعلام التوقيت الرسمي الدقيق لخوادم زين وحساب فارق التوقيت وضبط صلاحية الـ OTP والتوكنات |

---

## تفاصيل وميزات المكتبة (Features Breakdown)

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

### 3. كشف الحساب والتحقق التلقائي من تحويلات الرصيد بدون تسجيل دخول (CDR & Transfer Verification)

توفر المكتبة منظومة متكاملة تتيح للمتاجر والأنظمة السحابية والتطبيقات التحقق البرمجي التلقائي والفوري من استلام حوالات الرصيد من الزبائن **دون الحاجة لتسجيل دخول الزبون** وبدون أي تدخل يدوي للأدمن، بالاعتماد على كشف حساب الشريحة (CDR) وسجل الفواتير والإشعارات اللحظية:

| # | الطريقة | المسار (Endpoint Path) | دالة Go SDK المقابلة | الوصف التفصيلي باللغة العربية |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `GET` | `/api/number/electronic-bill-items` | `client.GetCDRTransferHistory(ctx, limit)`<br>`client.GetElectronicBillItems(ctx)` | كشف حساب سجل تحويلات الرصيد الواردة للشريحة (CDR) مع رقم المرسل والمبلغ والتاريخ الدقيق بالثانية لتأكيد الدفع التلقائي دون تسجيل دخول الزبون |
| **02** | `GET` | `/api/dashboard/notifications` | `client.GetNotifications(ctx, off, lim, read)`<br>`client.GetIncomingTransfers(ctx, limit)` | جلب إشعارات وتنبيهات وصول الرصيد من النظام فورياً وفحص الحوالات الجديدة المكتملة |
| **03** | `GET` | `/api/number/wallet` | `client.GetBalance(ctx)`<br>`client.GetWalletBalance(ctx)` | استعلام رصيد المحفظة الأساسي للشريحة وتاريخ انتهاء الصلاحية بدقة متناهية |
| **04** | `GET` | `/api/number/summary` | `client.GetSummary(ctx)`<br>`client.GetWalletOverview(ctx)` | نظرة عامة شاملة لرصيد المحفظة الأساسي، الصلاحية، ونوع الخط والخدمات |
| **05** | `POST` | `/api/number/charge-voucher` | `client.RechargeVoucher(ctx, voucherPIN)` | شحن وتعبئة الرصيد الفوري بكارت الشحن الورقي (16 رقماً) للخط الحالي أو لرقم آخر |
| **06** | `POST` | `/api/number/credit-transfer` | `client.CreditTransfer(ctx, recipient, amount, otp)` | تحويل رصيد نقدي مباشر من الخط إلى رقم آخر في شبكة زين العراق |
| **07** | `POST` | `/api/otp/request` | `client.RequestCreditTransferOTP(ctx, senderMSISDN)` | طلب رمز التحقق OTP عبر رسالة SMS لعمليات تحويل الرصيد المباشر |
| **08** | `POST` | `/api/otp/confirm` | `client.ConfirmCreditTransferOTP(ctx, otpCode, senderMSISDN)` | تأكيد رمز التحقق واستخراج توكن المصادقة والتفويض لإتمام التحويل |
| **09** | `POST` | `/api/number/extend-validity` | `client.ExtendValidity(ctx, amount)` | تمديد صلاحية استقبال وإرسال الخط بخصم من رصيد الحساب |
| **10** | `GET` | `/api/number/validity-options` | `client.GetValidityOptions(ctx)` | استعلام خيارات وأسعار تمديد صلاحية الخط المتاحة رسمياً |
| **11** | `POST` | `/api/payment/create-checkout-id` | `client.CreateCheckoutID(ctx, req)` | توليد معرّف الدفع Checkout ID لبوابة الدفع الإلكتروني والبطاقات المصرفية |
| **12** | `POST` | `/api/payment/refresh-payment-status` | `client.RefreshPaymentStatus(ctx, checkoutID)` | التحقق من حالة إتمام معاملة الدفع الإلكتروني وتأكيد نجاحها |
| **13** | `POST` | `/api/payment/save-card` | `client.SaveCard(ctx, req)` | حفظ وتشفير بيانات البطاقة المصرفية للعمليات القادمة |
| **14** | `POST` | `/api/payment/set-default-card` | `client.SetDefaultCard(ctx, cardID)` | تعيين بطاقة دفع معينة كخيار افتراضي في الحساب |
| **15** | `GET` | `/api/payment/user-cards` | `client.GetUserCards(ctx)` | استعلام قائمة البطاقات المصرفية المحفوظة للمشترك |
| **16** | `DELETE` | `/api/payment/delete-card` | `client.DeleteCard(ctx, cardID)` | حذف وفك ربط بطاقة مصرفية محفوظة من الحساب |
| **17** | `POST` | `/api/payment/purchase-order` | `client.CreatePurchaseOrder(ctx, req)`<br>`client.InitiateZainCashPayment(ctx, req)` | إنشاء أمر دفع مباشر وشراء عبر محفظة زين كاش (ZainCash) |

#### محرك المطابقة والتحقق الذاتي ومعالجة الرسائل (Verification & SMS Engine):

| # | النوع | الوظيفة / العملية | دالة Go SDK المقابلة | الوصف التفصيلي باللغة العربية |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `محرك محلي` | فحص ومطابقة الحوالة | `client.VerifyIncomingTransfer(ctx, senderPhone, minAmount)` | التحقق البرمجي التلقائي والفوري من استلام حوالة رصيد من زبون **دون الحاجة لتسجيل دخول الزبون** وبدون أي تدخل يدوي للأدمن |
| **02** | `محرك محلي` | الانتظار الذكي (Smart Polling) | `client.WaitForIncomingTransfer(ctx, phone, amt, interval)` | فحص دوري متكرر كل X ثوانٍ حتى وصول الحوالة فعلياً في كشف الحساب وتفعيل الطلب آلياً |
| **03** | `محلل SMS` | قراءة رسائل الـ SMS | `zain.ParseTransferSMS(smsText)` | استخراج رقم المرسل والمبلغ المالي من نص رسائل زين العراق الرسمية (يدعم الأرقام الشرقية `٠١٢٣٤` والغربية `01234`) |
| **04** | `سجل محلي` | تسجيل فوري للحوالة | `client.RecordIncomingTransferFromSMS(smsText)` | تسجيل الحوالة المقروءة من رسالة الـ SMS تلقائياً في دفتر المطابقة بالذاكرة مع حماية منع الازدواجية |
| **05** | `إعدادات` | تثبيت محفظة النظام | `client.SetMasterWallet(phone)`<br>`client.MasterWallet()` | تثبيت واسترجاع رقم الشريحة المعتمدة لاستقبال الأموال والرصيد في النظام |
| **06** | `توليد USSD` | كود التحويل السريع | `client.FormatUSSDTransfer(recipient, amount)` | توليد كود التحويل المباشر لزين العراق (`*123*amount*recipient#`) لإرساله للزبون للتحويل فورياً |
| **07** | `معالجة أرقام` | توحيد صيغ الأرقام العراقية | `zain.NormalizeMSISDN(phone)`<br>`zain.FormatLocalMSISDN(phone)` | تحويل الأرقام للصيغة المعيارية الدولية والمحلية ومطابقة آخر 9 أرقام لتجاوز اختلافات الصيغ |

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

### 4. شحن كارت الرصيد الورقي (Voucher Card Recharge)

شحن كروت الرصيد المكونة من 16 رقماً مع التحقق من نجاح العملية:

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

	// كارت شحن زين العراق (16 رقماً)
	voucherPIN := "1234567890123456"

	// شحن الكارت على رقم الخط المسجل
	if err := client.RechargeVoucher(ctx, voucherPIN); err != nil {
		log.Fatalf("فشل شحن الكارت: %v", err)
	}
	fmt.Println("تم شحن كارت الرصيد بنجاح!")

	// فحص الرصيد الجديد المحدث
	balance, err := client.GetBalance(ctx)
	if err == nil {
		fmt.Printf("الرصيد الجديد: %d د.ع | الصلاحية: %s\n",
			balance.Balance.Value, balance.Balance.Expiry)
	}
}
```

### 5. تحويل الرصيد المباشر بالـ OTP (P2P Credit Transfer)

دورة تحويل رصيد من الشريحة لرقم آخر عبر كود التحقق SMS:

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
	senderPhone := client.GetMSISDN()
	recipientPhone := "07809876543"
	amountIQD := int64(2000)

	// 1. طلب رمز التحقق SMS لعملية التحويل
	otpResp, err := client.RequestCreditTransferOTP(ctx, senderPhone)
	if err != nil {
		log.Fatalf("فشل طلب رمز التحويل: %v", err)
	}
	fmt.Printf("تم إرسال رمز التأكيد SMS! Request ID: %s\n", otpResp.Data.RequestID)

	// 2. تأكيد الرمز المستلم واستخراج توكن التأكيد
	smsCode := "123456"
	confResp, err := client.ConfirmCreditTransferOTP(ctx, smsCode, senderPhone)
	if err != nil {
		log.Fatalf("فشل تأكيد الرمز: %v", err)
	}

	// 3. إتمام عملية تحويل الرصيد
	if err := client.CreditTransfer(ctx, recipientPhone, amountIQD, confResp.Data.ConfirmationID); err != nil {
		log.Fatalf("فشل تحويل الرصيد: %v", err)
	}

	fmt.Printf("تم تحويل %d د.ع إلى الرقم %s بنجاح!\n", amountIQD, recipientPhone)
}
```

### 6. التحقق الآلي من الحوالات الواردة للأنظمة والمتاجر الرقمية والتطبيقات

التحقق التلقائي والفوري من تحويلات الرصيد الواردة من الزبائن دون الحاجة لتسجيل دخول الزبون (عبر مطابقة رقم الهاتف وسجل الفاتورة الإلكترونية والإشعارات):

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

	// تسجيل الدخول بجلسة المحفظة المعتمدة
	if err := client.LoadSessionFromFile("session.json"); err != nil {
		log.Fatalf("session load error: %v", err)
	}

	// 1. تثبيت رقم محفظة النظام المعتمد
	masterWallet := "07801234567"
	client.SetMasterWallet(masterWallet)

	ctx := context.Background()

	// 2. توليد كود التحويل السريع للزبون (USSD Dial Code)
	customerPhone := "07809876543"
	requiredAmount := 5000.0 // 5,000 د.ع
	ussdCode := client.FormatUSSDTransfer(masterWallet, int64(requiredAmount))
	fmt.Printf("أرسل للزبون كود التحويل التالي: %s\n", ussdCode)

	// 3. التحقق الآلي الفوري بدون أي تدخل بشري
	fmt.Printf("جاري فحص وصول حوالة من %s بمبلغ لا يقل عن %.0f د.ع...\n", customerPhone, requiredAmount)
	verified, record, err := client.VerifyIncomingTransfer(ctx, customerPhone, requiredAmount)
	if err != nil {
		log.Fatalf("خطأ أثناء التحقق: %v", err)
	}

	if verified && record != nil {
		fmt.Printf("تم تأكيد دفع الطلب آلياً وبنجاح!\n")
		fmt.Printf("   المرسل: %s\n", record.MSISDN)
		fmt.Printf("   المبلغ: %s د.ع\n", record.Amount)
		fmt.Printf("   التاريخ والوقت: %s\n", record.CreatedAt)
	} else {
		fmt.Println("لم تصل الحوالة بعد، أو المبلغ المدفوع غير كافٍ.")
	}
}
```
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
| **`03_bundles_and_offers`** | استعراض الباقات، العروض المخصصة، ونظام فليكس. | `go run examples/03_bundles_and_offers/main.go` |
| **`04_credit_transfer`** | تحويل الرصيد وتأكيده بالـ OTP، توليد كود الـ USSD، والتحقق الآلي. | `go run examples/04_credit_transfer/main.go` |
| **`05_recharge_voucher`** | شحن الرصيد بكروت التعبئة المكونة من 16 رقماً، وخيارات تمديد الصلاحية. | `go run examples/05_recharge_voucher/main.go` |
| **`06_loyalty_and_imtiyaz`** | نقاط برنامج المكافآت، استبدال النقاط، وعروض امتياز. | `go run examples/06_loyalty_and_imtiyaz/main.go` |
| **`07_support_and_tickets`** | تذاكر الدعم الفني، إرفاق الصور، وتتبع المعالجة. | `go run examples/07_support_and_tickets/main.go` |
| **`08_cms_content_queries`** | استعلامات محرك المحتوى CMS لكتالوج العروض والضبط. | `go run examples/08_cms_content_queries/main.go` |
| **`09_payments_and_zaincash`** | بوابات الدفع الإلكتروني، البطاقات، ومحفظة زين كاش. | `go run examples/09_payments_and_zaincash/main.go` |
| **`10_bundle_sharing_and_fnf`** | مشاركة الباقات العائلية، وتحديد الحصص ونقل الوحدات. | `go run examples/10_bundle_sharing_and_fnf/main.go` |
| **`11_nearme_and_notifications`** | فروع زين القريبة، الإشعارات، والخدمات الرقمية. | `go run examples/11_nearme_and_notifications/main.go` |
| **`12_wallet_and_incoming_transfer_verification`** | تثبيت المحفظة والتحقق الآلي من تحويلات الرصيد للأنظمة والمتاجر المؤتمتة. | `go run examples/12_wallet_and_incoming_transfer_verification/main.go` |
| **`13_daily_gift_and_rewards_automation`** | أتمتة سحب الهدايا اليومية، تحويل نقاط المكافآت، وفحص الصلاحية. | `go run examples/13_daily_gift_and_rewards_automation/main.go` |
| **`14_automated_transfer_matching`** | التحقق الآلي من استلام الرصيد ومطابقته فورياً من الرسائل والكشف السحابي ومنع التكرار. | `go run examples/14_automated_transfer_matching/main.go` |
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
