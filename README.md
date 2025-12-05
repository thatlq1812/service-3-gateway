# API Gateway

HTTP REST gateway for gRPC microservices.

**Port:** `8080` | **Protocol:** HTTP/REST

## Quick Start

```bash
# Run (User Service and Article Service must be running first!)
cp .env.example .env  # Edit with your config
go run cmd/server/main.go
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `USER_SERVICE_ADDR` | localhost:50051 | User Service gRPC address |
| `ARTICLE_SERVICE_ADDR` | localhost:50052 | Article Service gRPC address |
| `GATEWAY_PORT` | 8080 | HTTP server port |

---

# API Reference

## User APIs

### 1. Create User

`POST /api/v1/users`

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 2. Get User

`GET /api/v1/users/{id}`

```bash
curl http://localhost:8080/api/v1/users/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 3. Update User

`PUT /api/v1/users/{id}`

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john@example.com",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z"
  }
}
```

---

### 4. Delete User

`DELETE /api/v1/users/{id}`

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "success": true
  }
}
```

---

### 5. List Users

`GET /api/v1/users?page=1&page_size=10`

```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "users": [
      {"id": 1, "name": "John Doe", "email": "john@example.com"},
      {"id": 2, "name": "Jane Smith", "email": "jane@example.com"}
    ],
    "total": 25,
    "page": 1,
    "size": 10,
    "has_more": true
  }
}
```

---

## Auth APIs

### 6. Login

`POST /api/v1/auth/login`

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

---

### 7. Logout

`POST /api/v1/auth/logout`

**Requires:** `Authorization: Bearer <token>`

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "success": true
  }
}
```

---

## Article APIs

### 8. Create Article

`POST /api/v1/articles`

**Requires:** `Authorization: Bearer <token>`

```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "My Article",
    "content": "Article content here...",
    "user_id": 1
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "title": "My Article",
    "content": "Article content here...",
    "user_id": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 9. Get Article

`GET /api/v1/articles/{id}`

```bash
curl http://localhost:8080/api/v1/articles/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "title": "My Article",
    "content": "Article content here...",
    "user_id": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

---

### 10. Update Article

`PUT /api/v1/articles/{id}`

```bash
curl -X PUT http://localhost:8080/api/v1/articles/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content..."
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "id": 1,
    "title": "Updated Title",
    "content": "Updated content...",
    "user_id": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z"
  }
}
```

---

### 11. Delete Article

`DELETE /api/v1/articles/{id}`

```bash
curl -X DELETE http://localhost:8080/api/v1/articles/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "success": true
  }
}
```

---

### 12. List Articles

`GET /api/v1/articles?page=1&page_size=10&user_id=1`

```bash
# All articles
curl "http://localhost:8080/api/v1/articles?page=1&page_size=10"

# Filter by author
curl "http://localhost:8080/api/v1/articles?user_id=1"
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Article 1",
        "content": "...",
        "user_id": 1,
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com"
        }
      }
    ],
    "total": 25,
    "page": 1,
    "total_pages": 3
  }
}
```

---

## Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| `000` | 200 | Success |
| `001` | 409 | Already exists |
| `003` | 400 | Invalid argument |
| `005` | 404 | Not found |
| `016` | 401 | Unauthenticated |
| `013` | 500 | Internal error |

**Error Response:**
```json
{
  "code": "005",
  "message": "user not found"
}
```

---

## Quick Test Flow

```bash
# 1. Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"pass123"}'

# 2. Login (save token)
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}' | jq -r '.data.access_token')

# 3. Create article with token
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"My Article","content":"Hello!","user_id":1}'

# 4. Get article
curl http://localhost:8080/api/v1/articles/1
```
