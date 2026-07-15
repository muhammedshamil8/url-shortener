# 📡 API Reference Manual

This manual documents the REST endpoints, authentication flows, and request/response payloads for the URL Shortener service.

---

## 🔑 Authentication Flow

All endpoints requiring authorization use **Bearer Token Authentication** in the request headers:
```text
Authorization: Bearer <your_access_token>
```

---

## 🌐 Public Endpoints

### 1. Create Shortened URL
* **Endpoint**: `POST /api/v1/shorten`
* **Headers**: `Authorization: Bearer <token>` (Optional. If provided, links URL to user account)
* **Request Body**:
  ```json
  {
    "url": "https://github.com/muhammedshamil8/url-shortener"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "success",
    "data": {
      "id": 12,
      "original_url": "https://github.com/muhammedshamil8/url-shortener",
      "short_code": "xY8z",
      "short_url": "http://localhost:8080/xY8z"
    },
    "request_id": "936d5258-29bf-407b-8d02-a726ea90cbbd"
  }
  ```

### 2. URL Redirection
* **Endpoint**: `GET /{code}`
* **Response (302 Found)**: Redirects client to the original URL target.

---

## 🔒 Authentication Endpoints

### 1. Register User
* **Endpoint**: `POST /api/v1/auth/register`
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "success",
    "message": "User registered successfully"
  }
  ```

### 2. Login User
* **Endpoint**: `POST /api/v1/auth/login`
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "access_token": "eyJhbG...",
      "refresh_token": "eyJhbG..."
    }
  }
  ```

### 3. Refresh Access Token
* **Endpoint**: `POST /api/v1/auth/refresh`
* **Request Body**:
  ```json
  {
    "refresh_token": "eyJhbG..."
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "access_token": "eyJhbG..."
    }
  }
  ```

---

## 👤 User Endpoints (Requires Access Token)

### 1. Get Profile
* **Endpoint**: `GET /api/v1/me`
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "user": {
        "id": 1,
        "email": "user@example.com",
        "name": "John Doe",
        "role": "user"
      }
    }
  }
  ```

### 2. List Personal URLs
* **Endpoint**: `GET /api/v1/my/urls`
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "urls": [
        {
          "id": 12,
          "original_url": "https://github.com/muhammedshamil8/url-shortener",
          "short_code": "xY8z",
          "clicks": 5,
          "created_at": "2026-07-15T10:00:00Z"
        }
      ]
    }
  }
  ```

### 3. Update Redirect Target
* **Endpoint**: `PUT /api/v1/my/urls/:id`
* **Request Body**:
  ```json
  {
    "url": "https://google.com"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "URL updated successfully"
  }
  ```

### 4. Delete Personal URL
* **Endpoint**: `DELETE /api/v1/my/urls/:id`
* **Response (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "URL deleted successfully"
  }
  ```

---

## 🛡️ Admin Endpoints (Requires Admin Access Token)

### 1. List and Filter All URLs
* **Endpoint**: `GET /api/v1/admin/urls`
* **Query Parameters**:
  - `page`: Page index (default: `1`)
  - `limit`: Page item count (default: `10`)
  - `sort`: Column name to sort by (e.g. `clicks`, `created_at`)
  - `order`: Sort order (`asc` or `desc`)
  - `search`: Search query string
  - `min_clicks`/`max_clicks`: Click ranges
  - `min_date`/`max_date`: Creation date ranges
* **Response (200 OK)**: Returns paginated JSON list of all matching shortened links.

### 2. Delete Any URL
* **Endpoint**: `DELETE /api/v1/admin/urls/:id`

### 3. List All Users
* **Endpoint**: `GET /api/v1/admin/users`

### 4. Delete Any User
* **Endpoint**: `DELETE /api/v1/admin/users/:id`
