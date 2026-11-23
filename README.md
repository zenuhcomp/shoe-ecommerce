# Luna Collections — E-Commerce & Multi-Industry Showcase Platform

A full-stack, enterprise-grade e-commerce showcase platform built with **Go** (REST API) and **Angular 19** (Standalone Single Page Application). 

Featuring dynamic theming, multi-industry presets, interactive 360° product inspection, modern CMS administration, and a modular payment engine implementing the **Strategy Design Pattern** with complete **Safaricom Daraja M-Pesa** integration and automated webhook callback reconciliation.

---

## Table of Contents
1. [Key Features](#key-features)
2. [Architecture & Technology Stack](#architecture--technology-stack)
3. [Project Structure](#project-structure)
4. [Local Development](#local-development)
5. [Payment Processing Engine (Strategy Pattern)](#payment-processing-engine-strategy-pattern)
6. [Safaricom Daraja Integration & Callback Reconciliation](#safaricom-daraja-integration--callback-reconciliation)
7. [Industry Presets & Catalog Seeding](#industry-presets--catalog-seeding)
8. [CMS Content Management](#cms-content-management)
9. [Zero-Build Production Deployment](#zero-build-production-deployment)
10. [Systemd Service Management](#systemd-service-management)
11. [API Reference](#api-reference)

---

## Key Features

- **Multi-Method Payment Strategy Pattern**: Modular architecture supporting Lipa Na M-Pesa Online (STK Push), M-Pesa Offline (Till & Paybill), Credit/Debit Cards, Cash on Delivery, and Bank Wire / RTGS / KCB Buni.
- **Safaricom Daraja Production Integration**: Automatically syncs production credentials from `.env` with zero hardcoding of shortcodes, Till, or Paybill numbers.
- **Automated Webhook Reconciliation**: Handles Safaricom `stkCallback` webhooks at `/api/payments/mpesa/callback` with automatic order matching, status transitions (`paid` / `payment_failed`), receipt storage, and audit trail logging.
- **Real-Time Settlement Polling**: The checkout page automatically polls `/api/payments/order-status` every 2s, updating to *Payment Confirmed* the second the customer enters their PIN.
- **5 Curated Industry Presets (10 Detailed Items Each)**:
  - **Apex Motors Prestige (`vehicles`)**: Luxury & performance motor vehicles.
  - **Vanguard Atelier (`furniture`)**: Architectural, bespoke handcrafted furniture.
  - **Aura Wellness Sanctuary (`spa`)**: Holistic wellness & therapeutic treatments.
  - **Glow Beauty & Cosmetics (`beauty`)**: Clinical & botanical skincare essentials.
  - **Luna Collections (`footwear`)**: Handcrafted official & office-ready women's shoes.
- **Safe Preset Switching**: Applying a new preset seamlessly transforms site branding, typography, color tokens, and catalog items while **strictly preserving** configured payment credentials and custom Till/Paybill numbers.
- **Interactive 360° Product Viewer**: Multi-angle image viewer with smooth drag-to-rotate interaction for items across all presets.
- **Unified Dialog System**: Replaced browser `alert()` and `confirm()` prompts with an enterprise `ConfirmDialogService` and interactive modal component.
- **Zero-Build Production Server**: Pre-compiled 100% statically linked Linux x86_64 Go binary (`backend/backend`) that natively serves both the Angular SPA (`/`) and REST API (`/api/*`) on port `8084` using only **~20MB of RAM**.

---

## Architecture & Technology Stack

```
                     ┌───────────────────────────────────────────────────┐
                     │            Client Browser / Customer              │
                     └─────────────────────────┬─────────────────────────┘
                                               │
                                 HTTP / HTTPS (:8084)
                                               │
┌──────────────────────────────────────────────▼──────────────────────────────────────────────┐
│                              Go Backend Single-Binary                                       │
│                                                                                            │
│   ┌────────────────────────────────┐         ┌──────────────────────────────────────────┐  │
│   │   Embedded Static SPA Server   │         │               Chi REST Router            │  │
│   │ (frontend/dist/frontend/browser│         │             (Public & Admin APIs)        │  │
│   └────────────────────────────────┘         └────────────────────┬─────────────────────┘  │
│                                                                   │                        │
│   ┌───────────────────────────────────────────────────────────────▼─────────────────────┐  │
│   │                          Payment Strategy Manager                                   │  │
│   │                                                                                     │  │
│   │   ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌──────────────┐ ┌──────┐  │  │
│   │   │ M-Pesa Express│ │ M-Pesa Manual │ │  Credit Card  │ │ Bank Transfer│ │ Cash │  │  │
│   │   │   (STK Push)  │ │(Till/Paybill) │ │  (Stripe/3DS) │ │  (KCB Buni)  │ │(COD) │  │  │
│   │   └───────┬───────┘ └───────────────┘ └───────────────┘ └──────────────┘ └──────┘  │  │
│   └───────────┼─────────────────────────────────────────────────────────────────────────┘  │
│               │                                                                            │
│   ┌───────────▼───────────┐                            ┌────────────────────────────────┐  │
│   │  Safaricom Daraja API │                            │  Repository / Database Layer   │  │
│   │(api.safaricom.co.ke)  │                            │     (SQLite or PostgreSQL)     │  │
│   └───────────────────────┘                            └────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────────────────────────┘
```

- **Backend**: Go 1.25, Chi Router (`go-chi/chi/v5`), ModernC SQLite (`modernc.org/sqlite` - 100% pure Go, no CGO), pgx PostgreSQL driver (`jackc/pgx/v5`), JWT authentication (`golang-jwt/jwt/v5`).
- **Frontend**: Angular 19, Standalone Components, Signals, Reactive Forms, SCSS design tokens.
- **Process Management**: Native Linux `systemd` daemon.

---

## Project Structure

```
luna-collections/
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── auth/          # JWT token generation, password hashing & middleware
│   │   ├── bootstrap/     # Database seeding, admin syncing, Safaricom env sync
│   │   ├── config/        # Environment loading & validation (.env)
│   │   ├── db/            # Database connection & schema migration
│   │   ├── handlers/      # Public, Admin, and Payment HTTP handlers
│   │   ├── models/        # Go structs, request/response models
│   │   ├── payments/      # Payment Strategy Pattern implementations
│   │   │   ├── strategy.go       # Core PaymentStrategy interface & PaymentManager
│   │   │   ├── mpesa_express.go  # Safaricom Daraja STK Push strategy
│   │   │   ├── mpesa_manual.go   # Till & Paybill manual verification strategy
│   │   │   ├── card.go           # Credit/Debit Card strategy
│   │   │   ├── bank.go           # Bank Wire / RTGS / KCB Buni strategy
│   │   │   └── cash.go           # Cash on Delivery strategy
│   │   └── repository/    # SQLite / Postgres queries, presets, orders & CMS repo
│   ├── data/              # SQLite database storage (luna_collections.db)
│   ├── backend            # Pre-compiled statically linked Linux x86_64 binary
│   ├── main.go            # Backend entry point & static SPA file server
│   └── .env.example       # Template for backend environment variables
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── admin/     # Admin CMS (Dashboard, Payments, Presets, Sections, Orders)
│   │   │   ├── core/      # API services, models, auth guards & interceptors
│   │   │   ├── layout/    # Responsive header, footer, shell navigation
│   │   │   ├── pages/     # Home, Collections, Checkout, About, Reviews
│   │   │   └── shared/    # Confirm dialog, 360 viewer, booking forms, product cards
│   │   └── styles.scss    # Global design system & typography
│   └── dist/              # Pre-compiled production frontend bundle
├── deploy/
│   ├── luna.service       # Production systemd service unit
│   └── DEPLOYMENT.md      # Detailed Linux server deployment & Nginx guide
└── run.sh                 # Intelligent orchestration script (dev & static modes)
```

---

## Local Development

### 1. Environment Configuration
Copy the template configuration:
```bash
cp backend/.env.example backend/.env
```

Adjust the ports and credentials as needed:
```bash
BACKEND_PORT=8084
FRONTEND_PORT=4204
DATABASE_MODE=test
CMS_ADMIN_EMAIL=admin@lunacollections.local
CMS_ADMIN_PASSWORD=masai123
```

### 2. Start Stack
Start both backend (Go) and frontend (Angular dev server with live reload):
```bash
./run.sh
```

- **Storefront**: [http://localhost:4204](http://localhost:4204)
- **Checkout**: [http://localhost:4204/checkout](http://localhost:4204/checkout)
- **Admin Panel**: [http://localhost:4204/admin](http://localhost:4204/admin)
- **Backend API**: [http://localhost:8084/api/health](http://localhost:8084/api/health)

To run in **static production mode locally**:
```bash
./run.sh --static
```

---

## Payment Processing Engine (Strategy Pattern)

The payment module in `backend/internal/payments/` decouples payment gateways using the **Strategy Design Pattern**:

### `PaymentStrategy` Interface
```go
type PaymentStrategy interface {
    ID() string
    Name() string
    Type() string
    Description() string
    IsEnabled(cfg models.PaymentConfig) bool
    GetPublicInfo(cfg models.PaymentConfig) map[string]any
    Process(ctx context.Context, req models.PaymentInitiateRequest, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error)
    Verify(ctx context.Context, reference string, cfg models.PaymentConfig) (*models.PaymentInitiateResponse, error)
}
```

### Supported Strategies
1. **Lipa Na M-Pesa Online (`mpesa_express`)**:
   - Automated STK Push sent directly to customer handset.
   - Configured via `ConsumerKey`, `ConsumerSecret`, `PassKey`, `ShortCode` (or `.env`).
2. **M-Pesa Offline (`mpesa_manual`)**:
   - Buy Goods Till & Business Paybill.
   - Dynamically active only when an administrator enters a Till or Paybill number and flips the toggle.
3. **Credit / Debit Card (`card`)**:
   - Visa, MasterCard, and Amex online processing.
4. **Cash on Delivery (`cash`)**:
   - In-store pickup or cash handover upon delivery.
5. **Bank Wire / RTGS / KCB Buni (`bank_transfer`)**:
   - Corporate bank transfers, RTGS wires, and KCB Buni business accounts.

---

## Safaricom Daraja Integration & Callback Reconciliation

### Environment Variables
Production Safaricom credentials in `backend/.env`:
```bash
SAFARICOM_BASE_URL=https://api.safaricom.co.ke
CONSUMER_KEY=your_consumer_key
CONSUMER_SECRET=your_consumer_secret
COUNTY_SHORTC0DE=4005953
PASSKEY=your_passkey
MPESA_CALLBACK_URL=https://your-domain.com/api/payments/mpesa/callback
```

### Callback Webhook
- **Endpoint**: `POST /api/payments/mpesa/callback`
- Matches the incoming `CheckoutRequestID` with pending orders in the database.
- Parses `ResultCode`:
  - `ResultCode: 0`: Transitions order status to `paid`, stores the `MpesaReceiptNumber`, amount, and customer phone number in the audit trail.
  - Non-zero: Transitions order status to `payment_failed` and records Safaricom's `ResultDesc`.

### Simulation & Testing
You can simulate the entire Safaricom STK Push & reconciliation workflow locally or in staging:
- **Simulation Endpoint**: `POST /api/payments/mpesa/simulate-callback`
- **UI Trigger**:
  - Available inside the Customer Checkout screen under the STK waiting dialog.
  - Available inside the Admin Payment Gateways screen (`/admin/payments`).

---

## Industry Presets & Catalog Seeding

The platform supports 5 distinct industry presets. Each preset reconfigures site branding, hero video, typography, color themes, section layouts, and seeds **10 rich curated items** (complete with descriptions, specifications, KES pricing, and 360° product viewer assets):

| Preset ID | Brand Name | Industry | Curated Items |
|---|---|---|---|
| `vehicles` | Apex Motors Prestige | Luxury & Performance Automobiles | Range Rover SV, Maybach S580, Porsche 911 GT3 RS, BMW M5 CS, etc. |
| `furniture` | Vanguard Atelier | Architectural Bespoke Furniture | Oslo 3-Seater Sofa, Kyoto Live-Edge Dining Table, Mies Travertine Table, etc. |
| `spa` | Aura Wellness Sanctuary | Holistic Spa & Treatments | Himalayan Salt & Hot Stone Ritual, Deep Ocean Detox Wrap, 24K Gold Facial, etc. |
| `beauty` | Glow Beauty & Cosmetics | Luxury Skincare & Cosmetics | Botanical Squalane Elixir, Rose Quartz Hydro Serum, Silk Collagen Mask, etc. |
| `footwear` | Luna Collections | Handcrafted Official Shoes | Pointed Court Pumps, Sovereign Slingbacks, Milano Leather Oxfords, etc. |

### Safe Preset Switching
When switching presets via `/admin/presets` or `POST /api/admin/presets/apply`, existing payment settings (active payment gateways, Safaricom credentials, and custom Till/Paybill numbers) are **automatically preserved**.

---

## CMS Content Management

Access the CMS administration panel at **`/admin`** using your administrator credentials:
- **Dashboard (`/admin/dashboard`)**: Order counts, inquiries, quick metrics.
- **Payment Gateways (`/admin/payments`)**: Toggle payment methods on/off, set Till and Paybill numbers, manage API keys (with show/hide eye toggles), and run simulated webhook callbacks.
- **Preset Switcher (`/admin/presets`)**: 1-click industry switcher with instant catalog seeding.
- **Sections Manager (`/admin/sections`)**: Reorder, enable, disable, and edit section titles, subtitles, and background styles.
- **Media Library (`/admin/media`)**: Manage product photography, hero videos, and lookbook images.
- **Orders (`/admin/orders`)**: View orders, customer contact info, delivery addresses, payment settlement status, and audit logs.

---

## Zero-Build Production Deployment

Luna Collections uses a **Zero-Build Production Workflow**:
All compilation (Go binary & Angular production build) is completed locally and tracked in git. The production server does **not** need Node.js, Go, npm, or GCC.

### Deploying to Server

On your server (e.g. `/srv/apps/zenu/shoe-ecommerce`):
```bash
# 1. Pull the pre-built application
cd /srv/apps/zenu/shoe-ecommerce
git pull origin staging

# 2. Kill any old manual processes on port 8084
sudo fuser -k -9 8084/tcp 2>/dev/null || true

# 3. Ensure permissions
sudo chown -R app:app /srv/apps/zenu/shoe-ecommerce
sudo chmod +x backend/backend run.sh

# 4. Restart systemd service
sudo systemctl restart luna.service
```

---

## Systemd Service Management

The systemd service unit (`deploy/luna.service`):

```ini
[Unit]
Description=Luna Collections frontend and backend
After=network.target

[Service]
Type=simple
User=app
Group=app
WorkingDirectory=/srv/apps/zenu/shoe-ecommerce/backend
Environment=ENV_FILE=/srv/apps/zenu/shoe-ecommerce/backend/.env
Environment=APP_ENV=production
Environment=FRONTEND_DIST_DIR=/srv/apps/zenu/shoe-ecommerce/frontend/dist/frontend/browser
ExecStart=/srv/apps/zenu/shoe-ecommerce/backend/backend
Restart=always
RestartSec=5
KillSignal=SIGTERM
TimeoutStopSec=30
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

### Service Commands

```bash
# Copy unit to systemd directory
sudo cp deploy/luna.service /etc/systemd/system/luna.service

# Reload daemon
sudo systemctl daemon-reload

# Enable service to start automatically on system reboot
sudo systemctl enable luna.service

# Start service
sudo systemctl start luna.service

# Check live status
sudo systemctl status luna.service

# Stream live production logs
sudo journalctl -u luna -f
```

---

## API Reference

### Public API Endpoints
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/health` | Health check (`{"status":"ok"}`) |
| `GET` | `/api/home` | Complete homepage payload |
| `GET` | `/api/site` | Branding, themes, contact info |
| `GET` | `/api/services` | Active product / service catalog |
| `GET` | `/api/sections` | Enabled homepage sections |
| `GET` | `/api/gallery` | Lookbook & showcase media |
| `GET` | `/api/reviews` | Customer testimonials |
| `GET` | `/api/faq` | Frequently asked questions |
| `POST` | `/api/orders` | Create customer order |
| `POST` | `/api/bookings` | Create inquiry / booking request |

### Payment Endpoints
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/payments/methods` | List enabled payment gateways (public info only) |
| `POST` | `/api/payments/initiate` | Initiate payment strategy (e.g. STK push) |
| `POST` | `/api/payments/verify` | Verify reference code (e.g. manual Till receipt) |
| `GET` | `/api/payments/order-status?orderId=...` | Poll live order settlement status |
| `POST` | `/api/payments/mpesa/callback` | Safaricom Daraja STK Push webhook |
| `POST` | `/api/payments/mpesa/simulate-callback` | Simulate Daraja webhook & reconcile order |

### Admin Endpoints (Require Bearer Token)
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/admin/auth/login` | Administrator authentication |
| `GET` | `/api/admin/auth/me` | Current admin profile |
| `GET` | `/api/admin/dashboard` | Key performance indicators & counts |
| `GET` | `/api/admin/payments` | Full gateway configuration & API keys |
| `PUT` | `/api/admin/payments` | Update payment gateways, toggles & Till numbers |
| `GET` | `/api/admin/presets` | List all 5 industry presets |
| `POST` | `/api/admin/presets/apply` | Apply industry preset (with item seeding) |
| `GET` | `/api/admin/services` | List, create, update, delete catalog items |
| `GET` | `/api/admin/sections` | Manage homepage section layouts & ordering |
| `GET` | `/api/admin/media` | Media asset upload & management |
| `GET` | `/api/admin/orders` | Manage customer orders & status |
| `PUT` | `/api/admin/site` | Global site configuration & theme colors |
