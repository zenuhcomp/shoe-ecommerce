# Luna Collections — Production Deployment Guide

This guide details the deployment of Luna Collections on Ubuntu / Debian Linux servers.

---

## 1. Architecture Overview (Zero-Build Server)

Luna Collections is designed with a **Zero-Build Production Architecture**:
- **Backend**: Pre-compiled, **100% statically linked Linux x86_64 binary** (`backend/backend`). It runs with **zero dependencies on Go, GLIBC, GCC, or Node.js** on the production server.
- **Frontend**: Pre-compiled Angular 19 Single Page Application (`frontend/dist/frontend/browser/`).
- **Single-Port Serving**: The Go backend natively serves both the Angular SPA (`/`) and the REST API (`/api/*`) on port `8084` using only **~20MB of RAM**.
- **Process Manager**: Managed natively by `systemd` for auto-start on boot, automatic restarts on failure, and structured logging via `journald`.

---

## 2. Server Requirements

- **OS**: Any 64-bit Linux distribution (Ubuntu 20.04/22.04/24.04, Debian 11/12, Rocky Linux, RHEL).
- **RAM**: 512MB minimum (1GB recommended).
- **CPU**: 1 vCPU is plenty.
- **Tools on Server**: Only `git`, `curl`, and `systemd` (No Node.js, npm, or Go needed).

---

## 3. Initial Server Setup

### Step 3.1: Create Dedicated `app` System User
If the `app` user does not already exist on the server:

```bash
sudo useradd --system --create-home --shell /usr/sbin/nologin app
```

### Step 3.2: Clone or Pull Repository
Navigate to the applications directory:

```bash
sudo mkdir -p /srv/apps/zenu
cd /srv/apps/zenu
sudo git clone -b staging git@github.com-zenuhcomp:zenuhcomp/shoe-ecommerce.git shoe-ecommerce
# Or if already cloned:
cd /srv/apps/zenu/shoe-ecommerce
sudo git checkout staging
sudo git pull origin staging
```

### Step 3.3: Configure Production Environment (`.env`)
Create `backend/.env` from the template:

```bash
sudo cp backend/.env.example backend/.env
sudo nano backend/.env
```

Ensure production values are set:
```bash
APP_ENV=production
DATABASE_MODE=test
SQLITE_PATH=./data/luna_collections.db
BACKEND_PORT=8084
CMS_ADMIN_EMAIL=admin@lunacollections.local
CMS_ADMIN_PASSWORD=your_strong_admin_password
CMS_ADMIN_NAME="Site Admin"
JWT_SECRET=use-a-long-random-secret-string-at-least-32-chars

# Safaricom M-Pesa API Credentials (Production)
SAFARICOM_BASE_URL=https://api.safaricom.co.ke
CONSUMER_KEY=your_production_consumer_key
CONSUMER_SECRET=your_production_consumer_secret
COUNTY_SHORTC0DE=4005953
PASSKEY=your_production_passkey
MPESA_CALLBACK_URL=https://your-domain.com/api/payments/mpesa/callback
```

### Step 3.4: Set Ownership and Execution Permissions
```bash
sudo chown -R app:app /srv/apps/zenu/shoe-ecommerce
sudo chmod +x /srv/apps/zenu/shoe-ecommerce/backend/backend
sudo chmod +x /srv/apps/zenu/shoe-ecommerce/run.sh
```

---

## 4. Systemd Service Setup

The service unit is pre-configured in `deploy/luna.service`:

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

### Step 4.1: Install & Enable Service

```bash
# 1. Copy unit file to systemd directory
sudo cp /srv/apps/zenu/shoe-ecommerce/deploy/luna.service /etc/systemd/system/luna.service

# 2. Reload systemd daemon
sudo systemctl daemon-reload

# 3. Enable service to start automatically on system reboot
sudo systemctl enable luna.service

# 4. Kill any orphaned manual processes on port 8084 before starting
sudo fuser -k -9 8084/tcp 2>/dev/null || true

# 5. Start service
sudo systemctl start luna.service

# 6. Verify service status
sudo systemctl status luna.service
```

---

## 5. Reverse Proxy Configuration (Nginx)

Place Nginx in front of port `8084` to handle SSL termination, caching, and domain routing:

```nginx
server {
    server_name your-domain.com www.your-domain.com;

    client_max_body_size 20M;

    location / {
        proxy_pass http://127.0.0.1:8084;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Static asset caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://127.0.0.1:8084;
        proxy_set_header Host $host;
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }
}
```

Enable SSL using Let's Encrypt Certbot:
```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

---

## 6. Updating the Application in Production

Since all compilation happens locally before pushing, updating the production server takes only seconds:

```bash
cd /srv/apps/zenu/shoe-ecommerce

# 1. Pull latest compiled binaries and assets
sudo -u app git pull origin staging

# 2. Restart systemd service
sudo systemctl restart luna.service

# 3. Verify health
sudo systemctl status luna.service
curl -s http://localhost:8084/api/health
```

---

## 7. Troubleshooting & Operational Commands

| Action | Command |
|---|---|
| **View live logs** | `sudo journalctl -u luna -f` |
| **View last 100 log lines** | `sudo journalctl -u luna -n 100 --no-pager` |
| **Restart service** | `sudo systemctl restart luna` |
| **Stop service** | `sudo systemctl stop luna` |
| **Check health endpoint** | `curl -s http://localhost:8084/api/health` |
| **Check payment methods** | `curl -s http://localhost:8084/api/payments/methods` |
| **Free stuck port 8084** | `sudo fuser -k -9 8084/tcp` |
| **Fix permissions** | `sudo chown -R app:app /srv/apps/zenu/shoe-ecommerce` |
