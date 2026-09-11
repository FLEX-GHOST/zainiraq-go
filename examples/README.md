# Zain Iraq Go SDK - Examples & Tutorials

This directory contains production-grade, standalone examples demonstrating every capability of the `pkg/zain` library.

---

## Examples Overview

| Directory | Description | Run Command |
| :--- | :--- | :--- |
| **`01_otp_login`** | Authentication via official SMS OTP, session generation, and saving persistent `session.json`. | `go run examples/01_otp_login/main.go` |
| **`02_account_and_profile`** | Real-time balance query, remaining 4G data/voice/SMS subaccounts, loan status, and profile info. | `go run examples/02_account_and_profile/main.go` |
| **`03_recharge_and_transfer`**| Recharging phone balance via 16-digit voucher cards, P2P credit transfer, and validity extension. | `go run examples/03_recharge_and_transfer/main.go` |
| **`04_bundles_and_offers`** | Browsing bundles, personalized ATL/BTL offers, Flex limits/status, and bundle gifting. | `go run examples/04_bundles_and_offers/main.go` |
| **`05_loyalty_and_imtiyaz`** | Checking loyalty points & tier (Gold/Platinum), redeeming rewards, and Imtiyaz partner discounts. | `go run examples/05_loyalty_and_imtiyaz/main.go` |
| **`06_support_and_tickets`** | Technical support categories, creating tickets with binary attachments & GPS, and tracking status. | `go run examples/06_support_and_tickets/main.go` |
| **`07_cms_content_queries`** | High-performance queries against the CMS engine (app configurations, offers catalog, FAQs, roaming). | `go run examples/07_cms_content_queries/main.go` |
| **`08_payments_and_zaincash`** | Payment gateway checkout ID creation, payment status polling, tokenized cards, and ZainCash orders. | `go run examples/08_payments_and_zaincash/main.go` |
| **`09_bundle_sharing_and_fnf`** | Family bundle sharing, adding members with MB quota, unit transfers, and Friends & Family numbers. | `go run examples/09_bundle_sharing_and_fnf/main.go` |
| **`10_nearme_and_notifications`**| Nearest certified branches by GPS coordinates, opening hours, Iraqi cities, and inbox notifications. | `go run examples/10_nearme_and_notifications/main.go` |
| **`interactive_cli`** | An all-in-one terminal CLI application with an interactive text menu for all core features. | `go run examples/interactive_cli/main.go` |

---

## Quick Start Guide

### 1. Interactive CLI (All-in-One)
To test and interact with all features using an interactive terminal menu:
```bash
go run examples/interactive_cli/main.go
```

### 2. Login & Session Setup
To log in with your Zain Iraq phone number and save a persistent `zain_session.json` file:
```bash
go run examples/01_otp_login/main.go
```

### 3. Check Account Balance & Subaccounts
```bash
go run examples/02_account_and_profile/main.go
```

### 4. Voucher Recharge & Credit Transfer
```bash
go run examples/03_recharge_and_transfer/main.go
```

### 5. CMS Direct Catalog Queries
```bash
go run examples/07_cms_content_queries/main.go
```
