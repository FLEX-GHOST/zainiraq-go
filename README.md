<div align="center">

# zain-go

**Production-grade, zero-dependency Go SDK for Zain Iraq APIs**

[![Go Version](https://img.shields.io/badge/Go-1.26+-18181b?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Dependencies](https://img.shields.io/badge/Dependencies-Zero%20(Stdlib)-18181b?style=flat-square)](https://pkg.go.dev/)
[![License](https://img.shields.io/badge/License-MIT-18181b?style=flat-square)](LICENSE)

<br />

مكتبة برمجية متكاملة، احترافية، وعالية الأداء مكتوبة بلغة **Go (Golang)** للتعامل مع واجهات برمجة تطبيقات شركة **زين العراق (Zain Iraq)**.  
المكتبة مبنية بنسبة **100% بالاعتماد على المكتبة القياسية للغة Go** وبدون أي مكاتب أو اعتمادات خارجية نهائياً.

</div>

---

## بنية المشروع (Project Structure)

```
zain-go/
├── pkg/
│   └── zain/                   # حزمة الـ SDK الأساسية (بدون اعتمادات خارجية)
│       ├── client.go           # إعداد العميل والترويسات وإدارة الجلسات
│       ├── auth.go             # المصادقة، OTP، والكابتشا
│       ├── profile.go          # الحساب، الرصيد، والباقات الفعالة
│       ├── services.go         # الباقات، الخدمات، وعمليات التحويل والشحن
│       └── types.go            # هياكل البيانات ونماذج الاستجابة JSON
├── examples/                   # أمثلة عملية مستقلة لكل خدمة
├── ENDPOINTS.md                # التوثيق التقني لنقاط النهاية
├── README.md
└── go.mod
```
