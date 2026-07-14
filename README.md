# 🔗 URL Shortener Service

[![Go Version](https://img.shields.io/github/go-mod/go-version/muhammedshamil8/url-shortener?color=00ADD8&style=flat-square)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/muhammedshamil8/url-shortener/go-tests.yml?branch=dev&style=flat-square)](https://github.com/muhammedshamil8/url-shortener/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/muhammedshamil8/url-shortener?style=flat-square)](https://goreportcard.com/report/github.com/muhammedshamil8/url-shortener)

A production-ready, high-performance URL shortener built with **Go**, **Gin**, and **PostgreSQL**.

This service demonstrates production-grade Go backend architecture, featuring dependency injection, robust middlewares, comprehensive unit/integration test suites, dynamic pagination/filtering, containerization, and automated CI pipelines.

---

## 🚀 Key Features

* **⚡ Core API**: URL shortening, click tracking, custom error responses, and redirects.
* **🔍 Search & Filter**: Case-insensitive filtering on original URLs or short codes.
* **📊 Click & Date Bounds**: Retrieve URLs filtered by click ranges (`min_clicks`, `max_clicks`) and creation dates (`min_date`, `max_date`).
* **📦 Architecture**: Clean architecture with the Repository pattern, interfaces, and strict Dependency Injection.
* **🛡️ Security & Reliability**: IP-based rate limiting, CORS configuration, environment validation, and graceful server shutdown.
* **🩺 Health Probes**: Live (`/live`) and database-connected Ready (`/ready`) endpoints.
* **📖 Documentation**: Interactive API documentation generated with Swagger UI.
* **🧪 Testing**: 70%+ test coverage with mock repositories, controller unit tests, and database integration tests.

---

## 🛠️ Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Language** | Go 1.26 | High-performance compiled backend language |
| **Web Framework** | Gin | Minimalist, fast HTTP router and middleware engine |
| **Database** | PostgreSQL 16 | Relational database storage with auto migrations |
| **API Docs** | Swagger | Auto-generated OpenAPI 2.0 specifications |
| **Logging** | slog | Structured, leveled logging in JSON format |
| **CI** | GitHub Actions | Automated lint, vet, and test checks |

---

## 📁 Project Structure

```text
url-shortener/
├── docs/                   # Auto-generated Swagger spec files
├── internal/
│   ├── config/             # Environment validation and parsing
│   ├── database/           # Postgres initialization and migrations
│   ├── handlers/           # HTTP controllers and routing handlers
│   ├── logger/             # Structured slog logger integration
│   ├── middleware/         # Rate limiter, CORS, request logging, and UUID tracking
│   ├── models/             # Shared entities and filtering ListOptions
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
   *(Update your DB connection details in `.env`)*

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

### Health check & Docs
* `GET /live` — Liveness probe (always returns 200)
* `GET /ready` — Readiness probe (pings database connection)
* `GET /swagger/index.html` — Swagger UI documentation

### Shorten URL
* `POST /shorten`
  * **Payload:**
    ```json
    {
      "url": "https://github.com/muhammedshamil8/url-shortener"
    }
    ```
  * **Response (201):**
    ```json
    {
      "status": "success",
      "data": {
        "id": 1,
        "original_url": "https://github.com/muhammedshamil8/url-shortener",
        "short_code": "xY7z9P",
        "short_url": "http://localhost:8080/xY7z9P"
      },
      "request_id": "8e3c1a-..."
    }
    ```

### Redirect
* `GET /{code}` — Redirects with a `302 Found` header to the original URL and increments click counts.

### Delete URL
* `DELETE /{id}` — Deletes URL matching the database primary ID.

### List All URLs (Paginated & Filtered)
* `GET /urls/all`
  * **Query Parameters:**
    | Parameter | Type | Default | Description |
    | :--- | :--- | :--- | :--- |
    | `page` | `int` | `1` | Page number |
    | `limit` | `int` | `20` | Results per page (Max: `100`) |
    | `sort` | `string` | `created_at` | Sort key (`created_at`, `click_count`, `short_code`) |
    | `order` | `string` | `DESC` | Sort direction (`ASC`, `DESC`) |
    | `search` | `string` | `""` | Case-insensitive search on URLs and short codes |
    | `min_clicks`| `int` | `""` | Filter URLs with clicks `>= min_clicks` |
    | `max_clicks`| `int` | `""` | Filter URLs with clicks `<= max_clicks` |
    | `min_date` | `string` | `""` | RFC3339 formatted start date threshold |
    | `max_date` | `string` | `""` | RFC3339 formatted end date threshold |

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
* Dynamic pagination, sorting, and filter bounds (click counts, dates, and search queries)
* Liveness `/live` and readiness `/ready` probes

### Phase 3 — Scalability & Security 🚧
* Redis response caching layer
* JWT session-based Authentication
* Prometheus metrics exports
* OpenTelemetry tracing instrumentation