# 🔗 URL Shortener Service

[![Go Version](https://img.shields.io/github/go-mod/go-version/muhammedshamil8/url-shortener?color=00ADD8&style=flat-square)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/muhammedshamil8/url-shortener/go-tests.yml?branch=dev&style=flat-square)](https://github.com/muhammedshamil8/url-shortener/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhammedshamil8/url-shortener?style=flat-square)](https://goreportcard.com/report/github.com/muhammedshamil8/url-shortener)

A production-ready, full-stack high-performance URL shortener built with **Go (Gin)**, **PostgreSQL**, **Redis**, and **React (Vite + TypeScript + Tailwind CSS)**.

This service showcases clean-architecture development featuring dependency injection, JWT-based Authentication, Role-Based Access Control (RBAC), Redis-based response caching, comprehensive test suites, dynamic pagination/filtering, containerization, and automated CI pipelines.

## 📐 System Architecture

```text
                React
                  │
             Axios Client
                  │
                  ▼
            Gin HTTP Server
      ┌──────────┼───────────┐
      ▼          ▼           ▼
 Middleware   Handlers   Swagger
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

---

## 🛠️ Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Backend Language** | Go 1.26.3 | High-performance compiled backend language |
| **Web Framework** | Gin | Minimalist, fast HTTP router and middleware engine |
| **Frontend Framework** | React 18 & TypeScript 5 | Component-based interactive UI with static typing |
| **Build Tooling (Web)** | Vite 5 | Fast local development server and optimized bundles |
| **Styling** | Tailwind CSS 3 | Utility-first CSS framework for modern design |
| **Database** | PostgreSQL 17 | Relational database storage with auto migrations |
| **Caching** | Redis 7 | In-memory key-value store for redirect response caching |
| **Authentication** | JWT (v5) & bcrypt | Token-based secure user sessions and password hashing |
| **API Docs** | Swagger | Auto-generated OpenAPI 2.0 specifications |
| **Logging** | slog | Structured, leveled logging in JSON format |
| **CI** | GitHub Actions | Automated lint, vet, and test checks |

---

## 📁 Project Structure

```text
url-shortener/
├── .github/
│   └── workflows/          # GitHub Actions CI workflow files
├── docs/                   # Architectural & API documentation
│   ├── api.md              # REST API reference manual
│   ├── architecture.md     # System architecture overview
│   ├── deployment.md       # Production deployment instructions
│   ├── monitoring.md       # Observability, logs, and metrics reference
│   ├── docs.go             # Auto-generated Swagger spec configurations
│   ├── swagger.json        
│   └── swagger.yaml        
├── frontend/               # React + Vite + TypeScript + Tailwind CSS Frontend
│   ├── src/
│   │   ├── components/     # Reusable layout and notification components
│   │   ├── views/          # Pages (Landing, Login, Register, User/Admin Dashboards)
│   │   ├── api.ts          # Axios client with interceptors for refresh token flow
│   │   ├── router.tsx      # Hash router for SPA navigation
│   │   └── App.tsx         # Root view and state manager
│   ├── package.json        # Frontend scripts and package dependencies
│   └── vite.config.ts      # Vite configuration (includes API dev proxy)
├── internal/
│   ├── auth/               # Password hashing and JWT generation/validation
│   ├── cache/              # Redis response caching abstractions
│   ├── config/             # Environment validation and parsing
│   ├── database/           # Postgres initialization and migrations
│   ├── handlers/           # HTTP controllers and routing handlers (User, Auth, Admin, URL)
│   ├── logger/             # Structured slog logger integration
│   ├── middleware/         # Rate limiter, CORS, request logging, JWT Auth, and Admin auth
│   ├── models/             # Shared entities, request/response models, and Claims
│   ├── redis/              # Redis client initialization
│   ├── repository/         # Postgres queries and database layer
│   ├── response/           # Consistent, unified JSON response payloads
│   └── utils/              # Helper functions (e.g., URL validation, shortcode generators)
├── .env.example            # Environment template configuration
├── Dockerfile              # Multi-stage optimized builder image
├── docker-compose.yml      # Multi-container orchestration config (App, Postgres, Redis)
├── Makefile                # Shorthand CLI automation for backend
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
* [PostgreSQL 16+](https://www.postgresql.org/) OR [Docker](https://www.docker.com/)

---

### Running the Services

#### 1. Spin up Backend & Databases with Docker (Recommended)
You can start the backend service alongside Postgres and Redis instantly:
```bash
docker-compose up --build
```
The backend API will be available at `http://localhost:8080`.

#### 2. Run Backend Locally
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

#### 3. Run Frontend Locally
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
   The frontend will be served at `http://localhost:3000`.

#### 4. Create an Admin User
To access the Admin Dashboard features, you can create a user with the `admin` role using the CLI utility:
* **Interactive Mode:**
  ```bash
  make create-admin
  ```
  *(You will be prompted to enter the username, email, and password securely)*

* **Non-Interactive Mode (with flags):**
  ```bash
  go run scripts/create_admin.go -username <username> -email <email> -password <password>
  ```

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
* `POST /api/v1/shorten` — Create a shortened URL (can optionally include authorization headers to link to user profiles).
* `GET /{code}` — Redirects to original URL with a `302 Found` status and increments the click count.

### User Endpoints (Requires Access Token)
* `GET /api/v1/me` — Retrieve profile details of the authenticated user.
* `GET /api/v1/my/urls` — Get all shortened URLs belonging to the authenticated user.
* `PUT /api/v1/my/urls/:id` — Update the destination target original URL for a shortened URL.
* `DELETE /api/v1/my/urls/:id` — Delete a shortened URL owned by the authenticated user.

### Admin Endpoints (Requires Admin Access Token)
* `GET /api/v1/admin/urls` — List and filter all shortened URLs across all users.
  * **Query Parameters:** `page`, `limit`, `sort`, `order`, `search`, `min_clicks`, `max_clicks`, `min_date`, `max_date`
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

### Phase 5 — Cloud Deployment 🚧
* Cloudflare DNS & SSL certificate configuration
* Production deployment to VPS or PaaS (Railway, Render, etc.)

### Phase 6 — Advanced Features 🚧
* Malicious URL checking/protection
* Custom shortcode aliases
* WebSocket connections / Real-time statistics dashboard
