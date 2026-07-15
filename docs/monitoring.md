# 📊 Observability & Monitoring

This document details logging practices, health checking mechanisms, and plans for metrics and tracing integration.

---

## 🪵 Structured Logging (`slog`)

The application implements structured logging in JSON format using Go's standard library `slog` package. 

### Why Structured Logging?
Standard text logs are difficult to parse automatically. Structured JSON logs allow logging aggregators (e.g. Datadog, ELK Stack, Grafana Loki) to index log parameters for easy searching, filtering, and alert configuration.

### Log Example
```json
{
  "time": "2026-07-15T15:15:30Z",
  "level": "INFO",
  "msg": "Redirected URL successfully",
  "request_id": "8b7e28aa-3ef4-47c3-bd02-c94d6e902bbf",
  "short_code": "abc",
  "original_url": "https://google.com"
}
```

---

## 🩺 Health Probes (Liveness & Readiness)

The backend provides two HTTP health endpoints under the `/api/v1` namespace for orchestrators (like Kubernetes, Docker, or Render) to monitor container status:

### 1. Liveness Probe (`GET /api/v1/live`)
* **Purpose**: Verifies that the HTTP server is running and responsive.
* **Response**: Returns a `200 OK` status and a success message.
* **Orchestrator Action**: If this endpoint fails, the container is dead and should be restarted.

### 2. Readiness Probe (`GET /api/v1/ready`)
* **Purpose**: Verifies that the application is fully ready to receive traffic. It tests external network dependencies (in this case, pings the PostgreSQL database connection).
* **Response**:
  * `200 OK` if the database is reachable.
  * `503 Service Unavailable` if the database connection fails.
* **Orchestrator Action**: If this endpoint fails, the container should be removed from the load balancer until it reports ready.

---

## 📈 Future Metrics & Observability Roadmap

We plan to export application metrics and traces for real-time dashboards:

### 1. Prometheus Metrics
We will expose a `/metrics` scrape endpoint displaying:
- HTTP request counter (`http_requests_total`) partitioned by path, method, and status code.
- Request duration histograms (`http_request_duration_seconds`) to track API latency.
- Go runtime stats (heap size, active goroutines).
- Redis cache hit/miss rates.

### 2. Grafana Dashboards
A pre-configured Grafana dashboard JSON file will be provided to visualize:
- Average request throughput (RPS).
- API latency percentiles (p50, p95, p99).
- Cache efficiency ratios.
- Database connection pool utilization.
