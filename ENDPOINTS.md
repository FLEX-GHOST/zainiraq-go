<div align="center">

<img src="assets/zain_icon.png" alt="أيقونة تطبيق زين العراق الرسمية" width="105" />
<br />
<img src="assets/zain.svg" alt="شعار زين العراق" width="280" />

# مرجع نقاط نهاية واجهة برمجة تطبيقات زين العراق (Zain Iraq API Endpoints)

**المواصفات الهندسية الشاملة لجميع نقاط نهاية HTTP الرسمية لتطبيق زين العراق (Zain Iraq v6.6.0) باللغتين العربية والإنكليزية مع نماذج الطلب والاستجابة وهياكل البيانات.**

<br />

[![Specification](https://img.shields.io/badge/Specification-100%25%20Verified-00A3A6?style=flat-square)](ENDPOINTS.md)
[![Endpoints](https://img.shields.io/badge/Endpoints-92%20Reverse--Engineered-18181b?style=flat-square)](ENDPOINTS.md)
[![Protocol](https://img.shields.io/badge/Protocol-HTTPS%2FREST%20%2B%20CMS-18181b?style=flat-square)](ENDPOINTS.md)
[![Client](https://img.shields.io/badge/Go%20Client-zainiraq--go-18181b?style=flat-square)](https://github.com/FLEX-GHOST/zainiraq-go)
[![Zero-Dependencies](https://img.shields.io/badge/Dependencies-Zero%20(Stdlib%20Only)-brightgreen?style=flat-square)](go.mod)

<br />

جميع نقاط النهاية الموضحة في هذا المستند تم تفكيكها وتحليلها هندسياً من تطبيق أندرويد الرسمي لزين العراق (`mobi.foo.zain` الإصدار 6.6.0) من خلال فحص ملفات الديكس (`classes.dex` إلى `classes7.dex`) وحزم الكوتلن المترجمة تحت الحزمة `org.skelets.app`، ومطابقتها مقابل خوادم الوسيط الرسمي (`mw-mobile.iq.zain.com`) وخادم إدارة المحتوى (`cms-mobile.iq.zain.com`).

</div>

---

## فهرس المحتويات

1. [جدول نقاط النهاية المعتمدة (Endpoints Matrix)](#1-جدول-نقاط-النهاية-المعتمدة-endpoints-matrix)
2. [البيئة، النطاقات والأمان المتقدم (Security & Environment)](#2-البيئة-النطاقات-والأمان-المتقدم-security--environment)
3. [المصادقة وإدارة الحسابات (Authentication & Account)](#3-المصادقة-وإدارة-الحسابات-authentication--account)
4. [الرصيد، المحفظة والفوترة الآجلة (Balance, Wallet & Billing)](#4-الرصيد-المحفظة-والفوترة-الآجلة-balance-wallet--billing)
5. [شحن الرصيد وتحويل الأموال وبوابة الدفع (Recharge, Transfer & Payments)](#5-شحن-الرصيد-وتحويل-الأموال-وبوابة-الدفع-recharge-transfer--payments)
6. [العروض والباقات ونظام فليكس (Offers, Bundles & Flex)](#6-العروض-والباقات-ونظام-فليكس-offers-bundles--flex)
7. [مشاركة الباقات والأرقام العائلية (Bundle Sharing & FNF)](#7-مشاركة-الباقات-والأرقام-العائلية-bundle-sharing--fnf)
8. [برنامج المكافآت ونقاط امتياز (Loyalty & Imtiyaz)](#8-برنامج-المكافآت-ونقاط-امتياز-loyalty--imtiyaz)
9. [مركز الدعم وتتبع الشكاوى الفنية (Support & Complaints)](#9-مركز-الدعم-وتتبع-الشكاوى-الفنية-support--complaints)
10. [لوحة التحكم والإشعارات والمتاجر القريبة (Dashboard, Notifications & Near Me)](#10-لوحة-التحكم-والإشعارات-والمتاجر-القريبة-dashboard-notifications--near-me)
11. [محرك استعلامات المحتوى CMS (Direct CMS Query Engine)](#11-محرك-استعلامات-المحتوى-cms-direct-cms-query-engine)

---

## 1. جدول نقاط النهاية المعتمدة (Endpoints Matrix)

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
| **15** | `GET` | `/api/number/balance` | `client.GetBalance(ctx)` | استعلام رصيد المحفظة الأساسي وتاريخ انتهاء صلاحية الخط |
| **16** | `GET` | `/api/number/loan` | `client.GetLoan(ctx)` | استعلام تفاصيل سلفة الرصيد للطوارئ والمبلغ المستحق |
| **17** | `GET` | `/api/number/postpaid-history`| `client.GetPostpaidHistory(ctx)` | سجل الفواتير والدفعات السابقة للخطوط الآجلة الدفع |
| **18** | `GET` | `/api/number/subaccounts` | `client.GetSubaccounts(ctx)` | جلب الحسابات الفرعية ورصيد البيانات (إنترنت، مكالمات، رسائل) |
| **19** | `GET` | `/api/number/dashboard-message` | `client.GetDashboardMessage(ctx)` | استعلام رسائل التنبيه والاشعارات الموجهة للخط في الواجهة |
| **20** | `GET` | `/api/number/query-bill` | `client.GetBill(ctx)` | استعلام الفاتورة الحالية، المبالغ غير المفوترة، والمستحقات السابقة |
| **21** | `GET` | `/api/number/query-bill-items` | `client.GetBillItems(ctx)` | تفاصيل بنود وبنود استهلاك الفاتورة بالتفصيل |
| **22** | `GET` | `/api/number/query-unbilled` | `client.GetUnbilled(ctx)` | استعلام الاستهلاك المفتوح خارج الفاتورة قبل صدورها |
| **23** | `GET` | `/api/number/advance-payment` | `client.GetAdvancePayment(ctx)` | تفاصيل المبالغ المدفوعة مقدماً ورصيد التسديد المستقبلي |
| **24** | `POST` | `/api/number/change-language` | `client.ChangeLanguage(ctx, lang)` | تغيير لغة الإشعارات والرسائل النصية للنظام (عربي، كردي، إنكليزي) |
| **25** | `GET` | `/api/number/electronic-bill-items` | `client.GetElectronicBillItems(ctx)` | استعلام بنود الفاتورة الإلكترونية المعتمدة للطباعة والتوثيق |
| **26** | `POST` | `/api/payment/voucher` | `client.RechargeVoucher(ctx, pin)` | شحن وتعبئة الرصيد بكارت الشحن الورقي (16 رقماً) |
| **27** | `POST` | `/api/payment/credit-transfer` | `client.CreditTransfer(ctx, to, amt, otp)` | تحويل رصيد نقدي من رقم إلى رقم آخر في شبكة زين العراق |
| **28** | `POST` | `/api/payment/validity-extension` | `client.ExtendValidity(ctx, amount)` | تمديد صلاحية استقبال وإرسال الخط بخصم من الرصيد |
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
| **92** | `GET` | `/api/system/time` | `client.GetServerTime(ctx)` | مزامنة الوقت الدقيق مع خوادم زين لضبط المعاملات الحساسة للوقت |

---

## 2. البيئة، النطاقات والأمان المتقدم (Security & Environment)

### 2.1 النطاقات الرسمية (Official Base URLs)

| الغرض والبيئة | عنوان الخادم (Base URL) |
| :--- | :--- |
| **خادم الوسيط الأساسي (Middleware API)** | `https://mw-mobile.iq.zain.com` |
| **خادم إدارة المحتوى والكتالوج (CMS Engine)** | `https://cms-mobile.iq.zain.com` |
| **بوابة زين كاش (ZainCash Gateway)** | `https://api.zaincash.iq` |
| **بوابة الدفع المصرفي (HyperPay / OPPWA)** | `https://oppwa.iq.zain.com` (أو خادم ACI المعين) |

### 2.2 بصمات تشفير الربط الأمني (SSL Pinning Public Keys)

تم استخراج بصمات التشفير الرسمية (SHA-256 Hashes) لشهادات SSL/TLS من تكوين شبكة أمان التطبيق (`res/xml/network_security_config.xml`):

```text
sha256/PRF8+1xgYILmmEynfgeHZX4uYGXTh7AgwSJ75b9U5Mo=
sha256/Wec45nQiFwKvHtuHxSAMGkt19k+uPSw9JlEkxhvYPHk=
sha256/i7WTqTvh0OioIruIfFR4kMPnBqrS2rdiVPl/s2uC/CY=
```

### 2.3 ترويسات بروتوكول Skelets الإلزامية (Mandatory Headers)

| اسم الترويسة (Header) | القيمة الافتراضية / الوصف | مثال |
| :--- | :--- | :--- |
| `Skel-Accept-Language` | لغة المحتوى المطلوبة (`ar`, `en`, `kd`) | `ar` |
| `Skel-Platform` | نظام التشغيل | `Android` |
| `Skel-OS-Version` | إصدار واجهة برمجة تطبيقات النظام API Level | `14` (أو `34`) |
| `Skel-Fix-Version` | إصدار تطبيق زين العراق المثبت | `6.6.0` |
| `Skel-Installation-Id` | معرف تثبيت فريد بصيغة UUIDv4 لكل جلسة | `9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d` |
| `gal-msisdn` | رقم خط المشترك بدون الصفر أو المقدمة الدولية | `7845900162` |
| `Authorization` | توكن الوصول بصيغة Bearer | `Bearer eyJhbGciOi...` |
| `refresh-token` | توكن التحديث يُرسل فقط عند تجديد الجلسة | `dGhpcy1pcy1yZWZyZXNo...` |
| `Content-Type` | نوع حمولة البيانات المرسلة | `application/json` |
| `User-Agent` | وكيل المستخدم القياسي المتوافق مع التطبيق | `okhttp/4.12.0 (Zain-Iraq/6.6.0)` |

---

## 3. المصادقة وإدارة الحسابات (Authentication & Account)

### 3.1 طلب رمز التحقق (Request OTP)
* **المسار**: `POST /api/otp/request`
* **الترويسات**: ترويسات `Skel-*` الإلزامية، `gal-msisdn: 78XXXXXXXX`
* **دالة Go SDK**: `client.RequestOTP(ctx, "78XXXXXXXX")`

**طلب JSON**:
```json
{
  "msisdn": "7845900162"
}
```

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "request_id": "req-a1b2c3d4-e5f6-7890"
  }
}
```

---

### 3.2 تأكيد رمز التحقق (Confirm OTP)
* **المسار**: `POST /api/otp/confirm`
* **دالة Go SDK**: `client.ConfirmOTP(ctx, "78XXXXXXXX", "123456")`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "code": "123456"
}
```

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "confirmation_id": "conf-9876-5432-10fe-dcba"
  }
}
```

---

### 3.3 استبدال التوثيق بتوكنات الجلسة (Sign Session)
* **المسار**: `POST /api/user/sign`
* **دالة Go SDK**: `client.SignSession(ctx, "78XXXXXXXX", confirmationID)`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "confirmation": "conf-9876-5432-10fe-dcba",
  "user_space": "galleon"
}
```

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcy1pcy1hLXZhbGlkLXJlZnJlc2gtdG9rZW4...",
    "user_id": "14265463-f2e6-44e5-9006-f1cd3c33ab50",
    "expires_in": 86400
  }
}
```

---

### 3.4 استعلام الملف الشخصي (Get User Profile)
* **المسار**: `GET /api/v2/user/profile`
* **الترويسات**: `Authorization: Bearer <access_token>`
* **دالة Go SDK**: `client.GetProfile(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "id": "14265463-f2e6-44e5-9006-f1cd3c33ab50",
    "user_space": "galleon",
    "name": "Ahmed Ali",
    "msisdn": "7845900162",
    "is_debug_enabled": false,
    "is_4g_compatible": true,
    "is_gift_redeemed": true,
    "unified_sim_status": "NORMAL",
    "primary_offering_id": 40008,
    "customer_billing_type": "prepaid",
    "credit_cards": []
  }
}
```

---

## 4. الرصيد، المحفظة والفوترة الآجلة (Balance, Wallet & Billing)

### 4.1 استعلام رصيد المحفظة والصلاحية (Get Wallet Balance)
* **المسار**: `GET /api/number/balance`
* **دالة Go SDK**: `client.GetBalance(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "balance": {
      "value": 5221,
      "expiry": "2026-07-21T00:00:00"
    },
    "loan": {
      "amount": 0,
      "due_date": ""
    }
  }
}
```

---

### 4.2 استعلام الحسابات الفرعية واستهلاك الباقات (Subaccounts)
* **المسار**: `GET /api/number/subaccounts`
* **دالة Go SDK**: `client.GetSubaccounts(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": [
    {
      "account_type": 5443,
      "amount": 3600,
      "expiry_date": "2026-09-21T00:00:00"
    },
    {
      "account_type": 56259,
      "amount": 1048576,
      "expiry_date": "2026-10-01T00:00:00"
    }
  ]
}
```

---

### 4.3 استعلام الفاتورة للخطوط الآجلة الدفع (Query Bill)
* **المسار**: `GET /api/number/query-bill`
* **دالة Go SDK**: `client.GetBill(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "total_bill": 19992.0,
    "advance_payment": 8.0,
    "past_due": 0.0,
    "unbilled_amount": 20000.0,
    "is_dunning": false
  }
}
```

---

## 5. شحن الرصيد وتحويل الأموال وبوابة الدفع (Recharge, Transfer & Payments)

### 5.1 شحن الرصيد بكارت الشحن (Recharge Voucher)
* **المسار**: `POST /api/payment/voucher`
* **دالة Go SDK**: `client.RechargeVoucher(ctx, "1234567890123456")`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "voucher_number": "1234567890123456"
}
```

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "transaction_id": "tx-recharge-778899",
    "amount": 10000,
    "new_balance": 15221
  }
}
```

---

### 5.2 تحويل رصيد نقدي بين الخطوط (Credit Transfer)
* **المسار**: `POST /api/payment/credit-transfer`
* **دالة Go SDK**: `client.CreditTransfer(ctx, "7801234567", 5000, "")`

**طلب JSON**:
```json
{
  "sender": "7845900162",
  "recipient": "7801234567",
  "amount": 5000,
  "otp": ""
}
```

---

### 5.3 تمديد صلاحية الخط (Extend Validity)
* **المسار**: `POST /api/payment/validity-extension`
* **دالة Go SDK**: `client.ExtendValidity(ctx, 3000)`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "amount": 3000
}
```

---

### 5.4 إنشاء معرّف الدفع بالبطاقات (Create Checkout ID)
* **المسار**: `POST /api/payment/create-checkout-id`
* **دالة Go SDK**: `client.CreateCheckoutID(ctx, req)`

**طلب JSON**:
```json
{
  "amount": "10000.00",
  "currency": "IQD",
  "payment_type": "DB",
  "merchant_transaction_id": "order-101",
  "customer_email": "user@example.com"
}
```

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "checkout_id": "B12C34D56E78F90A1B2C3D4E5F67890A.prod01-vm-tx02"
  }
}
```

---

## 6. العروض والباقات ونظام فليكس (Offers, Bundles & Flex)

### 6.1 تفعيل واشتراك في باقة (Subscribe Offer)
* **المسار**: `POST /api/number/subscribe`
* **دالة Go SDK**: `client.SubscribeOffer(ctx, "OFFER_4G_MONTHLY_10GB")`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "offer_id": "OFFER_4G_MONTHLY_10GB"
}
```

---

### 6.2 استلام الهدية اليومية المجانية (Claim Daily Gift)
* **المسار**: `POST /api/number/claim-daily-gift`
* **دالة Go SDK**: `client.ClaimDailyGift(ctx)`

**طلب JSON**:
```json
{
  "msisdn": "7845900162"
}
```

---

### 6.3 استعلام الاشتراكات الفعالة (Get Subscriptions)
* **المسار**: `GET /api/number/subscriptions`
* **دالة Go SDK**: `client.GetSubscriptions(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": [
    {
      "offer_id": "OFFER_4G_MONTHLY_10GB",
      "name": "باقة 10 جيجابايت الشهرية",
      "expiry_date": "2026-10-11T23:59:59",
      "auto_renew": true,
      "status": "ACTIVE"
    }
  ]
}
```

---

## 7. مشاركة الباقات والأرقام العائلية (Bundle Sharing & FNF)

### 7.1 إضافة رقم لمجموعة المشاركة (Sharing Add Member)
* **المسار**: `POST /api/sharing/addrsc/`
* **دالة Go SDK**: `client.SharingAddMember(ctx, "7809998877", "OFFER_SHARE_50GB", 5120)`

**طلب JSON**:
```json
{
  "owner_msisdn": "7845900162",
  "member_msisdn": "7809998877",
  "offer_id": "OFFER_SHARE_50GB",
  "quota": 5120
}
```

---

### 7.2 استعلام أعضاء المجموعة والحصص (Query Sharing Members)
* **المسار**: `GET /api/sharing/query?msisdn=7845900162&offer_id=OFFER_SHARE_50GB`
* **دالة Go SDK**: `client.SharingQueryMembers(ctx, "OFFER_SHARE_50GB")`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "offer_id": "OFFER_SHARE_50GB",
    "total_pool_mb": 51200,
    "members": [
      {
        "msisdn": "7809998877",
        "role": "MEMBER",
        "allocated_mb": 5120,
        "consumed_mb": 1200
      }
    ]
  }
}
```

---

## 8. برنامج المكافآت ونقاط امتياز (Loyalty & Imtiyaz)

### 8.1 استعلام رصيد النقاط وفئة المشترك (Get Loyalty Info)
* **المسار**: `GET /api/loyalty/info`
* **دالة Go SDK**: `client.GetLoyaltyInfo(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "loyalty_status": "active",
    "imtiyaz_eligible": true,
    "tier": "gold",
    "localized_tier": {
      "en": "Gold",
      "ar": "الذهبية",
      "kd": "زێڕین"
    },
    "next_tier": "platinum",
    "total_points": 10000,
    "total_spendable_points": 9200,
    "total_points_to_next_tier": 20000
  }
}
```

---

### 8.2 استبدال النقاط برصيد نقدي (Redeem Credit)
* **المسار**: `POST /api/loyalty/rewards/credit`
* **دالة Go SDK**: `client.RedeemLoyaltyCredit(ctx, req)`

**طلب JSON**:
```json
{
  "msisdn": "7845900162",
  "points_amount": 5000,
  "credit_value": 5000
}
```

---

## 9. مركز الدعم وتتبع الشكاوى الفنية (Support & Complaints)

### 9.1 فتح تذكرة شكوى جديدة مع مرفقات (Create Ticket Multipart)
* **المسار**: `POST /api/complaints/create-ticket`
* **نوع المحتوى**: `multipart/form-data`
* **دالة Go SDK**: `client.CreateTicket(ctx, req)`

**أجزاء الطلب المتعدد (Multipart Form)**:
* `category_id`: `CAT_INTERNET_4G`
* `sub_category_id`: `SUBCAT_SLOW_SPEED`
* `msisdn`: `7845900162`
* `description`: `انقطاع متكرر لخدمة الإنترنت في منطقة المنصور`
* `latitude`: `33.3152`
* `longitude`: `44.3661`
* `attachment_1`: (بيانات صورة الفحص الفني الثنائية - ثنائية)

---

### 9.2 استعلام قائمة التذاكر وتتبع الحل (Get Tickets)
* **المسار**: `GET /api/complaints/tickets`
* **دالة Go SDK**: `client.GetTickets(ctx)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": [
    {
      "ticket_id": "TICK-2026-99881",
      "category": "خدمة الإنترنت",
      "status": "IN_PROGRESS",
      "created_at": "2026-09-10T14:30:00Z",
      "last_update": "تم إحالة الشكوى للفريق الفني الميداني للمنطقة"
    }
  ]
}
```

---

## 10. لوحة التحكم والإشعارات والمتاجر القريبة (Dashboard, Notifications & Near Me)

### 10.1 استعلام الإشعارات (Get Notifications)
* **المسار**: `GET /api/dashboard/notifications?offset=0&limit=10`
* **دالة Go SDK**: `client.GetNotifications(ctx, 0, 10, nil)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": {
    "unread_count": 2,
    "notifications": [
      {
        "id": "notif-101",
        "title": {
          "ar": "عرض الجمعة المميز",
          "en": "Special Friday Offer",
          "kd": "پێشنیاری تایبەتی هەینی"
        },
        "body": {
          "ar": "احصل على ضعف الرصيد عند الشحن بكارت 10,000 دينار",
          "en": "Get double credit on 10,000 IQD voucher recharge",
          "kd": "دوو ئەوەندە باڵانس وەربگرە"
        },
        "is_read": false,
        "created_at": "2026-09-11T12:00:00Z"
      }
    ]
  }
}
```

---

### 10.2 البحث عن الفروع ومراكز الخدمة القريبة (Near Me Shops)
* **المسار**: `GET /api/nearme/shops?latitude=33.3152&longitude=44.3661&radius=10`
* **دالة Go SDK**: `client.GetNearMeShops(ctx, 33.3152, 44.3661, 10)`

**استجابة خادم زين**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "SHOP_BAGHDAD_MANSOUR",
      "name": {
        "ar": "فرع المنصور الرئيسي",
        "en": "Al-Mansour Main Branch",
        "kd": "لقی سەرەکی مەنسوور"
      },
      "address": "شارع 14 رمضان، قرب ساحة الرواد",
      "latitude": 33.3145,
      "longitude": 44.3650,
      "is_open": true,
      "working_hours": "08:30 - 20:30"
    }
  ]
}
```

---

## 11. محرك استعلامات المحتوى CMS (Direct CMS Query Engine)

يستخدم تطبيق زين العراق محرك استعلامات المحتوى السحابي (`dmart/public/excute/query/...`) على خادم `https://cms-mobile.iq.zain.com` لجلب البنى المعقدة للمنتجات والإعدادات دون استهلاك موارد المعالجة في الوسيط.

### 11.1 استعلام ضبط وميزات التطبيق (App Configurations)
* **المسار**: `POST https://cms-mobile.iq.zain.com/dmart/public/excute/query/app_configurations`
* **دالة Go SDK**: `client.GetAppConfigurations(ctx)`

**طلب JSON**:
```json
{
  "query": {
    "is_active": true
  }
}
```

**استجابة محرك CMS**:
```json
{
  "status": "success",
  "error": null,
  "records": [
    {
      "resource_type": "content",
      "uuid": "30fb1125-5828-423b-a31b-85e34e82862b",
      "shortname": "enable_esim_swap",
      "branch_name": "master",
      "subpath": "app_configurations",
      "attributes": {
        "is_active": true,
        "payload": {
          "content_type": "json",
          "schema_shortname": "app_configuration",
          "body": {
            "key": "enable_esim_swap",
            "data_type": "boolean",
            "value": "true"
          }
        }
      }
    }
  ]
}
```

---

<div align="center">

**تم إنشاء هذه المواصفات الهندسية بواسطة فريق الهندسة العكسية لمنصة Antigravity.**  
*مرخصة للاستخدام البرمجي والأمني لربط الأنظمة الموزعة وبوتات التلغرام لشركة زين العراق.*

</div>
