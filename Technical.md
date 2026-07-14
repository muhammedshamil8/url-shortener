# URL Shortener - Learning Roadmap

This project is not just a CRUD application.

The goal is to learn how production Go backend services are designed, tested, documented, deployed, and maintained.

---

# Phase 1 — Core Backend Development ✅

## Backend Framework
- [x] Initialize Go module
- [x] Gin Web Framework
- [x] Project structure

## Database
- [x] PostgreSQL
- [x] Database connection
- [x] Table migration
- [x] Repository pattern

## REST API
- [x] Health endpoint
- [x] Create Short URL
- [x] Redirect URL
- [x] List URLs
- [x] Delete URL

## Architecture
- [x] Packages
- [x] Models
- [x] Repository Layer
- [x] Handler Layer
- [x] Dependency Injection
- [x] Interfaces

## Validation
- [x] URL validation
- [x] Input validation

## Testing
- [x] Unit Testing
- [x] Mock Repository
- [x] Handler Tests
- [x] Repository Integration Tests
- [x] Test Coverage

## Configuration
- [x] Environment Variables
- [x] Config Package

## Logging
- [x] Structured Logging (slog)

## Middleware
- [x] Request ID Middleware
- [x] Request Logging Middleware

## API Responses
- [x] Standard Response Package

## API Documentation
- [x] Swagger / OpenAPI

## Reliability
- [x] Graceful Shutdown

## CI
- [x] GitHub Actions
- [x] Automatic Test Execution

---

# Phase 1 — Concepts Learned

## Go
- Packages
- Modules
- Interfaces
- Structs
- Methods
- Constructors
- Dependency Injection
- Error Handling
- Context

## HTTP
- REST APIs
- HTTP Status Codes
- Routing
- Middleware

## Database
- PostgreSQL
- SQL
- Migrations
- Repository Pattern

## Testing
- Unit Tests
- Integration Tests
- Mocking
- Dependency Injection for Tests

## Software Design
- Layered Architecture
- Separation of Concerns
- Clean Code
- Configuration Management

## Production Engineering
- Structured Logging
- Request IDs
- API Documentation
- Graceful Shutdown
- CI Pipelines

---

# Phase 2 — Production Features ✅

## Infrastructure
- [x] Dockerfile
- [x] Docker Compose
- [x] Makefile

## Middleware
- [x] Rate Limiter
- [x] CORS
- [x] Security Headers

## Validation
- [x] Request Validation Improvements
- [x] Environment Validation

## Health
- [x] /live endpoint
- [x] /ready endpoint

## API
- [x] Pagination
- [x] Sorting
- [x] Filtering

# Phase 3 — Production Features 🚧 

## Performance
- [ ] Redis
- [ ] Response Caching

## Authentication
- [ ] User Registration
- [ ] Login
- [ ] JWT Authentication
- [ ] Refresh Tokens
- [ ] Role-Based Access Control (RBAC)

## Observability
- [ ] Prometheus Metrics
- [ ] OpenTelemetry Tracing

## CI/CD
- [ ] Linting
- [ ] Security Scanning
- [ ] Docker Image Build
- [ ] Automatic Deployment

## Deployment
- [ ] Docker Deployment
- [ ] Render / Railway / Fly.io
- [ ] VPS Deployment

---

# Long-Term Learning Goals

- High-performance Go APIs
- Production Engineering
- Cloud-native Development
- Scalable Backend Systems
- Observability
- Distributed Systems