# URL Shortener Service

A production-style URL shortener built with Go and Gin.

## Features

### Version 1
- [x] Create short URL
- [x] Redirect to original URL
- [x] List all URLs
- [x] Delete URL
- [x] Tests

### Version 2
- [ ] User registration
- [ ] Login
- [ ] JWT authentication
- [ ] User dashboard

### Version 3
- [ ] Custom aliases
- [ ] Click analytics
- [ ] URL expiration
- [ ] QR code generation

### Version 4
- [ ] Rate limiting
- [ ] Redis caching
- [ ] Admin panel
- [ ] Abuse reporting
- [ ] Safe Browsing integration

---

## Tech Stack

| Layer | Technology |
|--------|------------|
| Language | Go |
| Framework | Gin |
| Database | PostgreSQL |
| Cache | Redis (later) |
| Authentication | JWT |
| Frontend | React |
| Reverse Proxy | Nginx |
| Deployment | Docker + Linux VPS |

---

## Project Structure

```text
url-shortener/
├── main.go
├── handlers.go
├── database.go
├── models.go
├── helpers.go
├── .env
├── .gitignore
└── README.md
```

---

## API

### Create Short URL

```http
POST /shorten
```

Request

```json
{
    "url": "https://www.boot.dev/certificates/2d34e01f-b7f8-4228-a854-deab3119f51a"
}
```

Response

```json
{
    "short_url": "https://shamilkp.me/Ab3KdP"
}
```

---

### Redirect

```http
GET /:code
```

Example

```
GET /Ab3KdP
```

Returns

```
302 Redirect
```

---

## Roadmap

- [ ] PostgreSQL integration
- [ ] Redis caching
- [ ] Authentication
- [ ] Analytics
- [ ] Docker
- [ ] CI/CD
- [ ] Deployment to shamilkp.me
- [ ] HTTPS