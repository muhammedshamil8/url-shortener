# 🔗 URL Shortener Service

[![Go Version](https://img.shields.io/github/go-mod/go-version/muhammedshamil8/url-shortener?color=00ADD8&style=flat-square)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/muhammedshamil8/url-shortener/go-tests.yml?branch=dev&style=flat-square)](https://github.com/muhammedshamil8/url-shortener/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhammedshamil8/url-shortener?style=flat-square)](https://goreportcard.com/report/github.com/muhammedshamil8/url-shortener)

A production-ready, high-performance URL shortener built with **Go**, **Gin**, and **PostgreSQL**.

This service demonstrates production-grade Go backend architecture, featuring dependency injection, robust middlewares, JWT-based Authentication, Role-Based Access Control (RBAC), comprehensive unit/integration test suites, dynamic pagination/filtering, containerization, and automated CI pipelines.

---

## 🚀 Key Features

* **⚡ Core API**: URL shortening, click tracking, custom error responses, and redirects.
* **🛡️ Security & Authentication**: JWT authentication with Access and Refresh token flow. Role-based route authorization (Admin & Users).
* **🔍 Search & Filter**: Case-insensitive filtering on original URLs or short codes.
* **📊 Click & Date Bounds**: Retrieve URLs filtered by click ranges (`min_clicks`, `max_clicks`) and creation dates (`min_date`, `max_date`).
* **📦 Architecture**: Clean architecture with the Repository pattern, interfaces, and strict Dependency Injection.
* **🛡️ Reliability & Protection**: IP-based rate limiting, CORS configuration, environment validation, and graceful server shutdown.
* **🩺 Health Probes**: Live (`/live`) and database-connected Ready (`/ready`) endpoints.
* **📖 Documentation**: Interactive API documentation generated with Swagger UI.
* **🧪 Testing**: 75%+ statement test coverage with mock repositories, controller unit tests, and database integration tests.

---

## 🛠️ Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Language** | Go 1.26 | High-performance compiled backend language |
| **Web Framework** | Gin | Minimalist, fast HTTP router and middleware engine |
| **Database** | PostgreSQL 16 | Relational database storage with auto migrations |
| **Authentication** | JWT (v5) | Token-based secure user sessions and authorization |
| **API Docs** | Swagger | Auto-generated OpenAPI 2.0 specifications |
| **Logging** | slog | Structured, leveled logging in JSON format |
| **CI** | GitHub Actions | Automated lint, vet, and test checks |

---

## 📁 Project Structure

```text
url-shortener/
├── docs/                   # Auto-generated Swagger spec files
├── internal/
│   ├── auth/               # Password hashing and JWT generation/validation
│   ├── config/             # Environment validation and parsing
│   ├── database/           # Postgres initialization and migrations
│   ├── handlers/           # HTTP controllers and routing handlers (User, Auth, Admin, URL)
│   ├── logger/             # Structured slog logger integration
│   ├── middleware/         # Rate limiter, CORS, request logging, JWT Auth, and Admin authorization
│   ├── models/             # Shared entities, request/response models, and Claims
│   ├── repository/         # Postgres queries and database layer
│   ├── response/           # Consistent, unified JSON response payloads
│   └── utils/              # Helper functions (e.g., URL validation, shortcode generators)
├── .github/
│   └── workflows/          # GitHub Actions CI workflow files
├── Makefile                # Shorthand CLI automation
├── docker-compose.yml      # Multi-container orchestration config
├── Dockerfile              # Multi-stage optimized builder image
└── main.go                 # Application bootstrap entrypoint
```

---

## ⚙️ Getting Started

### Prerequisites
* [Go 1.26+](https://golang.org/dl/)
* [PostgreSQL 16](https://www.postgresql.org/) OR [Docker](https://www.docker.com/)

### Running with Docker (Recommended)
Spin up the service and a database instance instantly:
```bash
docker-compose up --build
```
The API will be available at `http://localhost:8080`.

### Running Locally
1. Clone the repository and configure your environment:
   ```bash
   cp .env.example .env
   ```
   *(Update your DB connection details and JWT secrets in `.env`)*

2. Run the application:
   ```bash
   make run
   ```

3. Run the test suite:
   ```bash
   make test
   ```

---

## 📡 API Reference

### Health Check & Swagger Docs
* `GET /api/v1/live` — Liveness probe (always returns 200)
* `GET /api/v1/ready` — Readiness probe (pings database connection)
* `GET /swagger/index.html` — Swagger UI documentation

### Authentication Endpoints
* `POST /api/v1/auth/register` — Register a new account.
* `POST /api/v1/auth/login` — Authenticate and retrieve Access + Refresh Tokens.
* `POST /api/v1/auth/refresh` — Refresh expired access token using a refresh token.

### Public Shortener Endpoints
* `POST /api/v1/shorten` — Create a shortened URL.
* `GET /{code}` — Redirects with a `302 Found` header to the original URL and increments click counts.

### User Endpoints (Requires Access Token)
* `GET /api/v1/me` — Retrieve profile details of authenticated user.
* `GET /api/v1/my/urls` — Get all shortened URLs belonging to the authenticated user.
* `DELETE /api/v1/my/urls/:id` — Delete a shortened URL owned by the authenticated user.

### Admin Endpoints (Requires Admin Access Token)
* `GET /api/v1/admin/urls` — List and filter all shortened URLs across all users.
  * **Query Parameters:** `page`, `limit`, `sort`, `order`, `search`, `min_clicks`, `max_clicks`, `min_date`, `max_date`
* `DELETE /api/v1/admin/urls/:id` — Deletes any shortened URL by ID.
* `GET /api/v1/admin/users` — List all registered user accounts.
* `DELETE /api/v1/admin/users/:id` — Remove a user account (and cascade delete their URLs).

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

### Phase 3 — Scalability & Security ✅ / 🚧
* JWT Session-based Authentication & Token Refresh ✅
* Role-based Access Control (RBAC) (Admin/User separation) ✅
* Redis response caching layer 🚧
* Prometheus metrics exports 🚧
* OpenTelemetry tracing instrumentation 🚧