# 🔗 URL Shortener Service

[![Go Version](https://img.shields.io/github/go-mod/go-version/muhammedshamil8/url-shortener?color=00ADD8&style=flat-square)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/muhammedshamil8/url-shortener/go-tests.yml?branch=dev&style=flat-square)](https://github.com/muhammedshamil8/url-shortener/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhammedshamil8/url-shortener?style=flat-square)](https://goreportcard.com/report/github.com/muhammedshamil8/url-shortener)

A production-ready, full-stack high-performance URL shortener built with **Go (Gin)**, **PostgreSQL**, **Redis**, and **React (Vite + TypeScript + Tailwind CSS)**.

This service showcases clean-architecture development featuring dependency injection, JWT-based Authentication, Role-Based Access Control (RBAC), Redis-based response caching, comprehensive test suites, dynamic pagination/filtering, containerization, and automated CI/CD pipelines.

## 📐 System Architecture

```text
                                  Browser Client (HTTPS)
                                            │
                                            ▼
                                     Cloudflare Edge
                                            │
                                            ▼ (Port 443 SSL)
                                    Nginx Reverse Proxy
                                     ┌──────┴──────┐
                                     ▼             ▼
                            React Frontend      Go API (Port 8080)
                             (Static Nginx)   ┌──────┼──────┐
                                              ▼      ▼      ▼
                                         Middleware Handlers Swagger
                                                     │
                                                     ▼
                                              Repository Layer
                                            ┌────────┴────────┐
                                            ▼                 ▼
                                       PostgreSQL         Redis Cache
```

---

## 🚀 Key Features

* **⚡ Core API**: URL shortening, click tracking, custom error responses, and redirects.
* **🛡️ Security & Authentication**: JWT authentication with Access and Refresh token flow. Role-based route authorization (Admin & Users).
* **💻 Interactive Frontend**: Modern client dashboard built with React and Tailwind CSS, featuring hash-routing, Toast notifications, session/refresh handling, and Admin/User portals.
* **⚡ Cache Layer**: High-performance response caching using Redis to reduce database load and speed up redirects.
* **🔍 Search, Sort & Filter**: Case-insensitive filtering on original URLs or short codes. Dynamic sorting, pagination, and range bounds (`min_clicks`, `max_clicks`, `min_date`, `max_date`).
* **📦 Architecture**: Clean architecture with the Repository pattern, interfaces, and strict Dependency Injection.
* **🛡️ Reliability & Protection**: IP-based rate limiting, CORS configurations, environment validation, and graceful server shutdown.
* **🩺 Health Probes**: Live (`/live`) and database-connected Ready (`/ready`) endpoints.
* **📖 Interactive Docs**: API documentation generated with Swagger UI.
* **🧪 Testing**: 75%+ statement test coverage with mock repositories, controller unit tests, and database integration tests.
* **🔄 Automated CI/CD**: Full automated workflows to compile/test code, build multi-stage Docker images, publish to Docker Hub, generate GitHub Releases, and auto-deploy directly to VPS over SSH.


---

## 💎 Engineering Highlights

* **Clean Architecture**: Separation of concerns using handler, service, and database layers.
* **Repository Pattern**: Data layer abstraction to support clean database mocking during unit tests.
* **Dependency Injection**: Explicit constructors for dependencies, avoiding global states and singletons.
* **JWT Authentication**: Full authentication system utilizing cryptographically signed tokens.
* **Refresh Tokens**: Token rotation and refresh flow to allow continuous user sessions securely.
* **RBAC**: Role-Based Access Control allowing secure separation of permissions for Users and Admins.
* **Redis Caching**: Near-zero latency URL redirects using an in-memory Redis cache.
* **GitHub Actions**: Automated pipeline for building, testing, linting, and deploying releases.
* **Docker Multi-stage Builds**: Statically compiled lightweight Alpine binaries for low resource consumption and tiny image footprints.
* **Docker Compose**: Orchestration configurations tailored for local development and production.
* **Cloudflare**: Full SSL validation proxying and DNS routing setup.
* **Nginx Reverse Proxy**: Reverse proxying web client requests to React assets and backend API.
* **VPS Deployment**: Low-cost, high-performance hosting on virtual private server instances.
* **Automatic Releases**: Automatic tagged version releases on GitHub.
* **Docker Hub Images**: Container builds published to Docker Registry for simple distribution.

---
## 🛠️ Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Backend Language** | Go 1.26 | High-performance compiled backend language |
| **Web Framework** | Gin | Minimalist, fast HTTP router and middleware engine |
| **Frontend Framework** | React 18 & TypeScript 5 | Component-based interactive UI with static typing |
| **Build Tooling (Web)** | Vite 5 | Fast local development server and optimized bundles |
| **Styling** | Tailwind CSS 3 | Utility-first CSS framework for modern design |
| **Database** | PostgreSQL 17 | Relational database storage with auto migrations |
| **Caching** | Redis 7 | In-memory key-value store for redirect response caching |
| **Authentication** | JWT (v5) & bcrypt | Token-based secure user sessions and password hashing |
| **API Docs** | Swagger | Auto-generated OpenAPI 2.0 specifications |
| **Logging** | slog | Structured, leveled logging in JSON format |
| **CI/CD** | GitHub Actions | Automated lint, test, docker build, release, and SSH deploy |

---

## 📁 Project Structure

```text
url-shortener/
├── .github/
│   └── workflows/          # GitHub Actions CI/CD workflows
│       ├── go-tests.yml        # PR and push test checker
│       └── release-deploy.yml  # Tag-based build, release, and deploy pipeline
├── deploy/                 # Production deployment assets
│   └── nginx/              # Nginx proxy routing configurations
│       └── default.conf
├── docs/                   # Architectural & API documentation
│   ├── api.md              # REST API reference manual
│   ├── architecture.md     # System architecture overview
│   ├── deployment.md       # Production deployment instructions
│   └── monitoring.md       # Observability, logs, and metrics reference
├── frontend/               # React + Vite + TypeScript Frontend
│   ├── src/
│   │   ├── components/     # Reusable layout and notification components
│   │   ├── views/          # Pages (Landing, Login, Register, User/Admin Dashboards)
│   │   ├── api.ts          # Axios client with interceptors for refresh token flow
│   │   ├── router.tsx      # Hash router for SPA navigation
│   │   └── App.tsx         # Root view and state manager
│   ├── Dockerfile          # Multi-stage production Nginx & Dev builder
│   ├── package.json        # Frontend scripts and package dependencies
│   └── vite.config.ts      # Vite configuration (includes API dev proxy)
├── internal/
│   ├── auth/               # Password hashing and JWT generation/validation
│   ├── cache/              # Redis response caching abstractions
│   ├── config/             # Environment validation and parsing
│   ├── database/           # Postgres initialization and migrations
│   ├── handlers/           # HTTP controllers and routing handlers
│   ├── logger/             # Structured slog logger integration
│   ├── middleware/         # Rate limiter, CORS, request logging, JWT Auth
│   ├── models/             # Shared entities, request/response models, and Claims
│   ├── redis/              # Redis client initialization
│   ├── repository/         # Postgres queries and database layer
│   ├── response/           # Consistent, unified JSON response payloads
│   └── utils/              # Helper functions (e.g., URL validation, shortcode generators)
├── .env.example            # Environment template configuration
├── Dockerfile              # Multi-stage optimized builder image for backend API
├── docker-compose.dev.yml  # Multi-container orchestration config for local development
├── docker-compose.prod.yml # Production multi-container orchestration config using Docker Hub images
├── Makefile                # Shorthand CLI automation for development
├── LICENSE                 # MIT License file
├── CONTRIBUTING.md         # Guide on how to contribute
├── CHANGELOG.md            # Release changelog history
├── CODE_OF_CONDUCT.md      # Contributor Covenant Code of Conduct
├── SECURITY.md             # Security policy and disclosure process
└── main.go                 # Application bootstrap entrypoint
```

---

## ⚙️ Getting Started

### Prerequisites
* [Go 1.26+](https://golang.org/dl/)
* [Node.js 18+](https://nodejs.org/) & `npm`
* [Docker & Compose](https://www.docker.com/)

---

### Running the Services

#### 1. Running in Development (Local Docker Compose)
You can start the backend service alongside Postgres, Redis, and the hot-reloaded Frontend development server:
```bash
docker compose -f docker-compose.dev.yml up --build
```
* The backend API will be available at `http://localhost:8080`.
* The frontend will be available at `http://localhost:5173`.

#### 2. Running Backend Locally (No Docker)
1. Clone the repository and configure your environment:
   ```bash
   cp .env.example .env
   ```
   *(Update your DB connection details, Redis settings, and JWT secrets in `.env`)*

2. Run the application (using `air` for hot-reload):
   ```bash
   make dev
   ```
   *Alternatively, run `make run` for a standard start.*

3. Run the test suite:
   ```bash
   make test
   ```

#### 3. Running Frontend Locally (No Docker)
1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the local development server:
   ```bash
   npm run dev
   ```
   The frontend will be served at `http://localhost:5173`.

---

### 🛡️ Admin Account Setup

#### In Development
To access the Admin Dashboard features, you can create a user with the `admin` role using the interactive CLI utility locally:
```bash
make create-admin
```

#### In Production (Docker VM)
To run the admin creation script on your production VM without installing Go locally on the host, run this temporary Docker container command inside your project folder (`~/url-shortener`):
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

---

## 🚀 CI/CD & Deployments

The project features a fully automated Git Tag-based CI/CD pipeline using GitHub Actions. 

### Triggering a Release
To release and deploy a new version of the app (e.g. `v1.0.3`):
```bash
git tag v1.0.3
git push origin v1.0.3
```

This triggers the `.github/workflows/release-deploy.yml` pipeline which:
1. Runs all Go tests and validates compilation.
2. Builds optimized Docker images for the API backend and Frontend (Nginx static proxy).
3. Pushes the built images to Docker Hub under the version tag and `latest`.
4. Creates a GitHub Release with build logs.
5. Connects to the VM via SSH, checks out the tag configuration, pulls the new images, and safely restarts the stack.

For detailed VPS setup instructions, see the [Deployment Guide](docs/deployment.md).

---

## 📡 API Reference

### Health Check & Swagger Docs
* `GET /api/v1/live` — Liveness probe (always returns 200)
* `GET /api/v1/ready` — Readiness probe (pings database connection)
* `GET /swagger/index.html` — Swagger UI documentation

### Authentication Endpoints
* `POST /api/v1/auth/register` — Register a new user account.
* `POST /api/v1/auth/login` — Authenticate and retrieve Access + Refresh Tokens.
* `POST /api/v1/auth/refresh` — Refresh expired access token using the refresh token.

### Public Shortener Endpoints
* `POST /api/v1/shorten` — Create a shortened URL.
* `GET /{code}` — Redirects to original URL with a `302 Found` status and increments the click count.

### User Endpoints (Requires Access Token)
* `GET /api/v1/me` — Retrieve profile details of the authenticated user.
* `GET /api/v1/my/urls` — Get all shortened URLs belonging to the authenticated user.
* `PUT /api/v1/my/urls/:id` — Update the destination target original URL for a shortened URL.
* `DELETE /api/v1/my/urls/:id` — Delete a shortened URL owned by the authenticated user.

### Admin Endpoints (Requires Admin Access Token)
* `GET /api/v1/admin/urls` — List and filter all shortened URLs across all users.
* `DELETE /api/v1/admin/urls/:id` — Deletes any shortened URL by ID.
* `GET /api/v1/admin/users` — List all registered user accounts.
* `DELETE /api/v1/admin/users/:id` — Remove a user account (and cascade deletes their URLs).

---

## 📈 Development Roadmap

### Phase 1 — Backend Foundations ✅
* Gin HTTP routing engine
* PostgreSQL Repository pattern with migrations
* Dependency Injection & Interfaces
* Graceful server shutdown & CORS middleware
* OpenAPI/Swagger docs UI

### Phase 2 — Production Readiness ✅
* Multi-stage Docker containerization
* Makefile scripting for local workflows
* IP Rate limiting middleware
* Dynamic pagination, sorting, and filter bounds
* Liveness `/live` and readiness `/ready` probes

### Phase 3 — Scalability & Security ✅
* JWT Session-based Authentication & Token Refresh ✅
* Role-based Access Control (RBAC) (Admin/User separation) ✅
* Redis response caching layer ✅

### Phase 4 — Frontend Client ✅
* Vite + React + TypeScript setup with Tailwind CSS ✅
* Custom client Hash Routing (Landing, Login, Register, User/Admin panels) ✅
* Axios API integration with automatic token refreshing ✅
* Dynamic pagination, query searches, and filters UI ✅

### Phase 5 — Cloud Deployment ✅
* Cloudflare DNS, SSL & Proxy setup ✅
* Automated CI/CD deployments directly to VPS over SSH ✅

### Phase 6 — Advanced Features 🚧
* Malicious URL checking/protection
* Custom shortcode aliases
* WebSocket connections / Real-time statistics dashboard
