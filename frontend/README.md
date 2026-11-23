# Luna Collections — Frontend Application

Modern, standalone **Angular 19** Single Page Application (SPA) designed for high performance, accessibility, and luxury e-commerce showcases.

---

## 1. Key Technologies & Architecture

- **Angular 19 Standalone Architecture**: Built entirely with standalone components (`standalone: true`), eliminating legacy `NgModule` boilerplate.
- **Angular Signals**: Modern reactive state management using `signal()`, `computed()`, and `effect()` for predictable, fine-grained UI updates.
- **Dynamic Theming**: The application loads design tokens (colors, gradients, typography, fonts) from the `/api/site` endpoint and injects them as CSS custom properties on `:root`.
- **Reusable Confirmation Dialog**: Global `ConfirmDialogService` provides modal confirmations without using blocking native browser `alert()` or `confirm()` dialogs.
- **Interactive 360° Product Viewer**: Custom multi-angle product viewer allowing customers to inspect vehicles, furniture, shoes, and luxury items with smooth drag and touch rotation.
- **Multi-Strategy Checkout**: Dynamic checkout tabs that adapt in real time to whichever payment methods the administrator has enabled (M-Pesa STK Push, Till/Paybill, Card, Cash, Bank Transfer / KCB Buni).

---

## 2. Directory Structure

```
frontend/src/app/
├── admin/                     # CMS Administration Module (Lazy-loaded)
│   ├── layout/                # Admin navigation bar & sidebar shell
│   ├── pages/
│   │   ├── dashboard/         # Analytics overview & quick stats
│   │   ├── login/             # Secure JWT login with password eye toggle
│   │   ├── media-admin/       # Photo & video asset library manager
│   │   ├── orders-admin/      # Customer order fulfillment & payment tracking
│   │   ├── payments-admin/    # Gateway toggles, Till/Paybill numbers & simulator
│   │   ├── presets-admin/     # 1-click industry switcher with confirmation modal
│   │   ├── sections-admin/    # Drag-and-drop homepage section layout manager
│   │   ├── services-admin/    # Catalog CRUD with specs & 360 image sets
│   │   └── site-editor/       # Branding, copy, and color theme editor
│   └── services/              # Admin HTTP API service
├── core/                      # Global Core Singletons & Data Models
│   ├── api.service.ts         # Public REST API client
│   ├── auth.guard.ts          # Route guard for admin paths
│   ├── auth.service.ts        # JWT session storage & auth state
│   └── models.ts              # TypeScript interfaces & payment models
├── layout/                    # Public Layout Shell
│   ├── header/                # Responsive navigation & mobile slide-out drawer
│   ├── footer/                # Branded footer with industry-specific copyright
│   └── shell/                 # Layout wrapper
├── pages/                     # Public Storefront Pages
│   ├── about/                 # Studio story & highlights
│   ├── checkout/              # Multi-strategy payment & STK polling checkout
│   ├── collections/           # Filterable product catalog grid
│   ├── contact/               # Interactive inquiry & booking form
│   ├── gallery/               # Visual lookbook & media gallery
│   ├── home/                  # Dynamic section-driven landing page
│   ├── house-calls/           # Delivery & white-glove service details
│   ├── pricing/               # Transparent pricing cards
│   └── reviews/               # Verified customer testimonials
└── shared/                    # Reusable UI Components
    ├── booking-form/          # Customer order & inquiry form
    ├── confirm-dialog/        # Reusable modal confirmation dialog & service
    ├── faq-list/              # Accordion FAQ component
    ├── product-image-viewer/  # Interactive 360° image inspection viewer
    ├── section-heading/       # Stylized section headers with badges
    └── service-grid/          # Responsive item grid with specs & pricing
```

---

## 3. Development Workflow

### Prerequisites
- Node.js 18+ or 20+
- npm 9+

### Installing Dependencies
```bash
npm install
```

### Running the Development Server
```bash
npm start
```
Runs the application on `http://localhost:4204` proxied to the Go API running on port `8084` via `proxy.conf.runtime.json`.

---

## 4. Production Build

To compile optimized production artifacts:

```bash
npm run build
```

- Output directory: `dist/frontend/browser/`
- Budgets: Configured in `angular.json` (1MB initial warning, 2MB max error; 15kB component styles warning, 25kB max error).
- Optimization: Minified JS/CSS with tree shaking and dead-code elimination.
- Deployment: The compiled directory is automatically served by the Go backend (`backend/backend`) in production.

---

## 5. Key Frontend Features

### Strategy-Driven Checkout & Real-Time Polling
The checkout component (`pages/checkout`):
1. Calls `/api/payments/methods` to display only the payment options currently enabled by the admin.
2. When the customer submits an order via **M-Pesa Express**:
   - Sends `POST /api/payments/initiate` to trigger an STK Push.
   - Displays a pulsing waiting animation requesting the user enter their M-Pesa PIN on their phone.
   - Automatically polls `GET /api/payments/order-status?orderId=...` every 2 seconds.
   - Automatically transitions to the *Payment Confirmed* receipt card upon webhook reconciliation without requiring a page reload.
   - Includes a *Simulate M-Pesa PIN* testing button for development and staging environments.

### Password Visibility Eye Toggles
All password and secret key fields (Admin Login, Consumer Secret, Passkey, Card Secret) include an accessible eye toggle button enabling one-click show/hide functionality.

### Global Confirm Dialog
Instead of browser-native `window.confirm()` or `window.alert()`, components inject `ConfirmDialogService`:
```typescript
const confirmed = await this.confirmDialog.confirm({
  title: 'Apply Preset',
  message: 'This will reconfigure branding, hero media, and catalog items.',
  confirmText: 'Apply Preset',
  cancelText: 'Cancel',
  variant: 'warning'
});
if (confirmed) {
  // execute action
}
```
