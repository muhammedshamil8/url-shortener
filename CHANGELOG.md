# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-07-15

Initial stable release of the Snippy URL Shortener full-stack service.

### Added
- **Core Backend**: Built Go/Gin web service with Clean Architecture, dependency injection, and slog JSON logging.
- **Database Migrations**: Automatic table migrations for PostgreSQL (`users` and `urls` tables).
- **URL Shortening & Redirection**: Shortcode generation, validation, collision handling, and click-count metrics.
- **Security & RBAC**: JWT Access & Refresh token flows, password hashing with bcrypt, role authorization (Admin vs. User).
- **Redis Cache**: High-performance redirect lookup caching with automatic cache invalidation.
- **REST Filters**: Pagination, sorting, case-insensitive search, and filtering bounds (click count and creation dates) on URLs.
- **Interactive UI**: React + Vite SPA built with Tailwind CSS, hash routing, authentication state management, and toast notifications.
- **Developer Experience**:
  - Swagger UI API documentation.
  - Multi-stage optimized Docker builds and Docker Compose configuration.
  - Automation scripts using a `Makefile`.
  - GitHub Actions automated testing workflow.
  - `LICENSE` (MIT) and `CONTRIBUTING.md` guides.
