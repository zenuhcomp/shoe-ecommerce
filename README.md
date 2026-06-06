# Luna Collections

A full-stack e-commerce showcase for **Luna Collections** — a state-of-the-art ladies official shoe vendor. Built with **Go** (REST API) and **Angular 19** (SPA). Cloned from the eye-lashes studio template and customized for premium office footwear: pumps, oxfords, loafers, block heels, and more.

## Project structure

```
luna-collections/
├── backend/     # Go REST API (chi router)
└── frontend/    # Angular 19 standalone app
```

## Run locally

Copy and edit environment variables:

```bash
cp backend/.env.example backend/.env
```

Start **both** backend and frontend:

```bash
./run.sh
```

Default ports (override in `.env`):

| Variable | Default | Purpose |
|----------|---------|---------|
| `BACKEND_PORT` | `8084` | Go API |
| `FRONTEND_PORT` | `4204` | Angular dev server |

Open [http://localhost:4204](http://localhost:4204).

### Database mode

| Value | Engine | Use case |
|-------|--------|----------|
| `test` | SQLite (`./data/luna_collections.db`) | Local dev |
| `production` | PostgreSQL | Production |

## CMS (content management)

Admin panel: **http://localhost:4204/admin**

Default credentials (change in `backend/.env`):

| Variable | Default |
|----------|---------|
| `CMS_ADMIN_EMAIL` | `admin@lunacollections.local` |
| `CMS_ADMIN_PASSWORD` | `luna-admin-change-me` |

### What you can edit

- **Site settings** — brand, hero, contact, theme, about, delivery, shoe care, pricing banner
- **Shoe styles** — full CRUD with pricing and featured flag (12 sample official styles seeded)
- **Features, press, reviews, FAQ, gallery, navigation**
- **Order inquiries** — view and delete submissions from the contact form

## Branding

| Element | Luna Collections |
|---------|------------------|
| Palette | Midnight navy, champagne gold, warm cream |
| Fonts | Playfair Display + Cormorant Garamond + DM Sans |
| Focus | Ladies official & office-ready footwear |

Theme tokens are served from `/api/site` and applied as CSS variables in the browser.

## Responsive design

- Fluid typography and spacing via `clamp()` and CSS custom properties
- Mobile-first navigation with touch-friendly 44px tap targets
- Safe-area insets for notched devices
- Adaptive grids from 320px to ultrawide
- Reduced-motion support for accessibility

## Build for production

```bash
cd backend && go build -o bin/api .
cd frontend && npm run build
```

Serve `frontend/dist/frontend/browser` behind a reverse proxy and route `/api` to the Go service.

## systemd deployment

The included unit at `deploy/luna.service` runs the app as the `app` user from
`/srv/test/shoe-ecommerce`. Before starting the service, make sure that account
owns the deployed project path because `run.sh` writes runtime files to `.run/`,
`backend/data/`, and `frontend/proxy.conf.runtime.json`.

```bash
sudo useradd --system --create-home --shell /usr/sbin/nologin app
sudo chown -R app:app /srv/test/shoe-ecommerce
sudo chmod +x /srv/test/shoe-ecommerce/run.sh
sudo cp /srv/test/shoe-ecommerce/deploy/luna.service /etc/systemd/system/luna.service
sudo systemctl daemon-reload
sudo systemctl enable --now luna.service
```

If the `app` user already exists, `useradd` may report that and can be skipped.

## API endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/health` | Health check |
| `GET /api/home` | Homepage payload |
| `GET /api/site` | Branding, theme, contact |
| `GET /api/services` | Shoe collections & prices |
| `GET /api/gallery` | Lookbook images |
| `GET /api/reviews` | Customer reviews |
| `GET /api/faq` | FAQ items |
| `POST /api/bookings` | Submit order inquiry |
