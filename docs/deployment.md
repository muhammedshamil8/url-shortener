# 🚀 Deployment Guide

This document describes how to deploy the Snippy URL Shortener to development, staging, and production environments.

---

## 📦 Containerization with Docker

The application is containerized using a multi-stage `Dockerfile` to optimize the final image size and enforce security.

### Multi-Stage Dockerfile Strategy
1. **Builder Stage**: Uses the official `golang:alpine` image to install dependencies, run code verifications, and compile a statically linked binary.
2. **Production Stage**: Copies the compiled binary into a minimal `alpine:latest` runner image. No source code or compiler tools are included in the final image, reducing the attack surface.

To build the image manually:
```bash
docker build -t url-shortener:latest .
```

---

## 🛠️ Deployment Targets

### 1. Multi-Container Orchestration (Docker Compose)
For local testing and single-node VPS deployments, `docker-compose.yml` configures the app, PostgreSQL, and Redis databases to run together.

Run in the background:
```bash
docker compose up -d --build
```

### 2. Platform as a Service (Railway / Render)
To deploy on PaaS environments:
1. **Database provision**: Set up a PostgreSQL instance and a Redis instance using the platform's database services.
2. **Environment Variables**: Bind the environment variables from the PaaS service (e.g. `DATABASE_URL` or individual `DB_*` parameters) to the app container.
3. **Build Command**: Set the root directory as the build context. The platform will automatically detect the `Dockerfile` and execute the multi-stage build.

---

## 🔒 Security Best Practices for Production

### 1. Reverse Proxy & SSL (Nginx / Cloudflare)
Never expose the Go server (`8080`) directly to the public internet. Instead:
- Deploy **Nginx** or **Traefik** as a reverse proxy in front of the Go application.
- Terminate SSL/HTTPS at the proxy level.
- Configure Cloudflare in front of the reverse proxy for SSL validation, DDoS protection, and IP masking.

### 2. Database Protection
- Ensure your PostgreSQL and Redis instances are **not** accessible on public IP addresses. Keep them within a private VPC or private Docker bridge network.
- Change all default passwords (`DB_PASSWORD` and `REDIS_PASSWORD`) before deploying.
- Use an SSL connection for PostgreSQL by setting `DB_SSLMODE=require` (or `verify-full`).

### 3. JWT Secret Rotation
Ensure your access and refresh secrets are highly secure:
```bash
# Generate a secure 64-character hex secret
openssl rand -hex 32
```
