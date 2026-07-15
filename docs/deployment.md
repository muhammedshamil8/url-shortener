# 🚀 Production Deployment Guide

This document describes how to deploy the Snippy URL Shortener service to production environments using Docker Compose, Cloudflare, and automated GitHub Actions.

---

## 📦 Production Architecture

The production environment runs inside a multi-container stack orchestrated via Docker Compose:
1. **Nginx (Alpine)**: Serves as the public facing reverse proxy (handles HTTPS, static asset serving, and request routing).
2. **Go API App**: Runs the high-performance Go backend service.
3. **PostgreSQL 17 (Alpine)**: Holds relation datasets (users and links) with automatic migrations.
4. **Redis 7 (Alpine)**: Caches URL redirects for near-zero latency lookups.

All components run inside a secure, private bridge network. Only the Nginx container exposes public ports (`80` and `443`).

---

## ☁️ Setting Up the VM / VPS Server

### 1. Prerequisites on the VM
Before deploying, make sure your production server has:
* **Docker & Docker Compose** installed.
* **Git** installed (for syncing configuration).
* Port `80` and `443` open in the server's firewall.

### 2. Clone the Repository
Clone the codebase to the user's home directory (`~/url-shortener`) on the VM:
```bash
git clone https://github.com/muhammedshamil8/url-shortener.git ~/url-shortener
cd ~/url-shortener
```

### 3. Create Environment File
Create a `.env` file in the root `~/url-shortener` directory containing your production variables:
```env
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_db_password_here
DB_NAME=url_shortener
DB_SSLMODE=disable

# Application Configuration
APP_PORT=8080
BASE_URL=https://snippy.shamilkp.me/
APP_ENV=production
ALLOWED_ORIGINS=https://snippy.shamilkp.me

# JWT Configuration
JWT_ACCESS_TOKEN_SECRET=use_openssl_to_generate_32_byte_hex
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_SECRET=use_openssl_to_generate_32_byte_hex
JWT_REFRESH_TOKEN_EXPIRY=168h

# Redis Cache Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 🔒 SSL & Cloudflare Configuration (Full Strict)

We use a free 15-year **Cloudflare Origin Certificate** to secure the connection between Cloudflare and your VM.

### Step 1: Generate Cloudflare Certificate
1. Log into your **Cloudflare Dashboard** and select your domain.
2. Go to **SSL/TLS** ➔ **Origin Server**.
3. Click **Create Certificate**. Keep default settings (ECDSA, 15 years duration) and make sure `snippy.shamilkp.me` is listed in the hostnames. Click **Create**.
4. Cloudflare will display your **Origin Certificate** and **Private Key**.

### Step 2: Install Certificates on VM
Create the `certs/` directory inside `deploy/nginx/` and paste the certificates:
```bash
mkdir -p ~/url-shortener/deploy/nginx/certs

# Save the Origin Certificate
nano ~/url-shortener/deploy/nginx/certs/snippy.pem
# (Paste the Origin Certificate block here)

# Save the Private Key
nano ~/url-shortener/deploy/nginx/certs/snippy.key
# (Paste the Private Key block here)
```

### Step 3: Turn on Full (Strict) SSL
In your Cloudflare dashboard under **SSL/TLS** ➔ **Overview**, change the SSL/TLS encryption mode to **Full (strict)**.

---

## 🔄 Automated CI/CD Setup with GitHub Actions

Every time you push a version tag (`v*`), GitHub Actions automatically builds, compiles, tests, and deploys the stack to your VM.

### GitHub Repository Secrets Configuration
Go to your repository settings at **Settings** ➔ **Secrets and variables** ➔ **Actions**, and add the following repository secrets:

| Secret Name | Value |
| :--- | :--- |
| **`DOCKERHUB_USERNAME`** | Your Docker Hub account username. |
| **`DOCKERHUB_TOKEN`** | Your Docker Hub personal access token (PAT). |
| **`SSH_HOST`** | The VM's direct public IP address (do not use domain as it is proxied). |
| **`SSH_USERNAME`** | The SSH user to log in (e.g. `shamil`, `ubuntu`, `root`). |
| **`SSH_KEY`** | The private SSH key used to connect to your VM (`id_ed25519` or `id_rsa`). |
| **`SSH_PORT`** | SSH port of your server (defaults to `22` if not specified). |

### Deploying a New Release
Trigger a deployment by pushing a tag from your local terminal:
```bash
git tag v1.0.3
git push origin v1.0.3
```

The GHA pipeline will:
1. Fetch configurations via git checkout tags on the VM.
2. Pull updated production images from Docker Hub.
3. Reload Nginx and restart the stack automatically with zero downtime.

---

## 🛠️ Post-Deployment Administration

### Creating Admin Accounts in Production
Since Go is not installed on the VM host, run the script within a temporary Go container connected to the database bridge network:
```bash
docker run -it --rm \
  --network url-shortener_backend \
  --env-file .env \
  -e DB_HOST=postgres \
  -v $(pwd):/app \
  -w /app \
  golang:1.26-alpine \
  go run scripts/create_admin.go
```
Follow the interactive prompts to create your production administrator credentials.
