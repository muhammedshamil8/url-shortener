# URL Shortener Service

A production-style URL shortener built with **Go**, **Gin**, and **PostgreSQL**.

The goal of this project is to learn how production backend services are designed, tested, documented, monitored, and deployed—not just to build a CRUD application.

---

## Features

### Core API
- Create short URLs
- Redirect to original URLs
- List all URLs
- Delete URLs
- Automatic click counting
- URL validation

### Architecture
- Clean project structure
- Repository pattern
- Dependency Injection
- Interfaces
- Configuration package
- Standard API responses

### Database
- PostgreSQL
- Automatic database migrations

### Middleware
- Request ID middleware
- Structured request logging

### Observability
- Structured logging using `slog`
- Swagger / OpenAPI documentation

### Reliability
- Graceful shutdown
- Environment configuration

### Testing
- Unit tests
- Mock repository
- Handler tests
- Repository integration tests
- GitHub Actions CI

---

## Tech Stack

| Layer | Technology |
|--------|------------|
| Language | Go |
| Framework | Gin |
| Database | PostgreSQL |
| Documentation | Swagger / OpenAPI |
| Logging | slog |
| Testing | Go Testing |
| CI | GitHub Actions |

---

## Project Structure

```text
url-shortener/
├── docs/
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── logger/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── response/
│   └── utils/
├── .github/
│   └── workflows/
├── main.go
├── go.mod
└── README.md
```

---

## API Endpoints

| Method | Endpoint | Description |
|---------|----------|-------------|
| GET | `/health/api` | Health check |
| POST | `/shorten` | Create short URL |
| GET | `/{code}` | Redirect to original URL |
| GET | `/urls/all` | List all URLs |
| DELETE | `/{id}` | Delete URL |
| GET | `/swagger/index.html` | Swagger UI |

---

## Example

### Create Short URL

```http
POST /shorten
```

```json
{
  "url": "https://google.com"
}
```

Response

```json
{
  "status": "success",
  "data": {
    "id": 1,
    "original_url": "https://google.com",
    "short_code": "Ab3KdP",
    "short_url": "http://localhost:8080/Ab3KdP"
  },
  "request_id": "f8f1c2..."
}
```

---

## Development Roadmap

### Phase 1 — Backend Foundations ✅

- Gin REST API
- PostgreSQL
- Repository Pattern
- Dependency Injection
- Interfaces
- Unit & Integration Testing
- Swagger Documentation
- Structured Logging
- Request ID Middleware
- Standard API Responses
- Graceful Shutdown
- GitHub Actions CI

---

### Phase 2 — Production Readiness 🚧

- Dockerfile
- Docker Compose
- Makefile
- Rate Limiter
- CORS Middleware
- Request Validation Improvements
- Environment Validation
- Health Checks (`/live`, `/ready`)
- Better Error Handling
- Pagination
- Sorting
- Filtering

---

### Phase 3 — Scalability & Security

- Redis
- Response Caching
- User Management
- JWT Authentication
- Refresh Tokens
- Role-Based Access Control (RBAC)
- Prometheus Metrics
- OpenTelemetry Tracing
- CI/CD Improvements
- Production Deployment

---

## Learning Objectives

This project is being built to practice:

- Production Go project structure
- REST API development
- PostgreSQL
- Middleware design
- Repository pattern
- Dependency Injection
- Testing strategies
- API documentation
- Observability
- Production engineering
- Deployment workflows