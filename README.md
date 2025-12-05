# Service 3: API Gateway

> HTTP REST API Gateway translating to gRPC backend services

**Protocol:** HTTP REST  
**Port:** 8080  
**Backend:** gRPC clients to User & Article Services  
**Dependencies:** User Service + Article Service

---

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Overview](#overview)
- [Setup Options](#setup-options)
  - [Option 1: Docker](#option-1-docker-recommended)
  - [Option 2: Terminal (Local)](#option-2-terminal-local-development)
- [Environment Configuration](#environment-configuration)
- [API Reference](#api-reference)
- [Response Format](#response-format)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

**⚠️ Service-3 requires Service-1 and Service-2 running first.**

### Standalone Mode (This Service Only)

```bash
# PREREQUISITE 1: Start Service-1 first
cd service-1-user
docker-compose up -d
sleep 15

# PREREQUISITE 2: Start Service-2
cd ../service-2-article
docker-compose up -d
sleep 10

# Then start Service-3 (Gateway)
cd ../service-3-gateway
git clone https://github.com/thatlq1812/service-3-gateway.git
cd service-3-gateway

# Configure (optional - defaults work fine)
cp .env.example .env

# Build and start (Gateway only, no database)
docker-compose build --no-cache
docker-compose up -d

# Wait for healthy status
sleep 5

# Verify
docker ps
curl http://localhost:8080/health

# Test API
curl http://localhost:8080/api/v1/users
```

**Exposed Ports:**
- `8080` - HTTP REST API Gateway

**External Dependencies (via host.docker.internal):**
- `host.docker.internal:50051` - User Service (required)
- `host.docker.internal:50052` - Article Service (required)

### Multi-Service Sequential Startup (Complete Platform)

**Order:** Service-1 → Service-2 → Service-3

```bash
# Step 1: Start User Service (Foundation)
cd service-1-user
docker-compose up -d
sleep 15

# Step 2: Start Article Service
cd ../service-2-article
docker-compose up -d
sleep 10

# Step 3: Start Gateway (this service)
cd ../service-3-gateway
docker-compose up -d
sleep 5

# Verify all services
docker ps
curl http://localhost:8080/health
```

**Important Notes:**
- Gateway is pure HTTP/gRPC proxy (no database)
- Requires both backend services healthy before starting
- All API requests translate to gRPC calls to Service-1 or Service-2
- If backend services restart, Gateway automatically reconnects

See [Setup Options](#setup-options) for detailed instructions.

---

## Prerequisites

### For Docker Setup (Recommended)

- **Docker** 20.10+ and **Docker Compose** 1.29+
- **Git** for cloning the repository
- 2GB RAM minimum
- Ports 8080, 50051, 50052 available

### For Local Development

- **Go** 1.21 or higher
- **curl** for testing HTTP endpoints
- Backend services running:
  - User Service (port 50051)
  - Article Service (port 50052)

### Install curl

**Windows (PowerShell):**
```powershell
# curl comes pre-installed with Windows 10+
curl --version
```

**Linux:**
```bash
sudo apt-get install curl
```

**macOS:**
```bash
# curl comes pre-installed
curl --version
```

---

## Overview

API Gateway provides HTTP REST interface for external clients, translating requests to gRPC calls to backend services.

**Features:**
- HTTP REST to gRPC translation
- Standardized JSON responses
- JWT authentication middleware
- Error code mapping (gRPC → REST)
- CORS support
- Request/response logging

**Technology Stack:**
- **Language:** Go 1.21+
- **Protocol:** HTTP REST (external) → gRPC (internal)
- **Router:** gorilla/mux
- **Integration:** User Service + Article Service gRPC clients

---

## Setup Options

### Option 1: Docker (Recommended)

**Prerequisites:**
- Docker 20.10+
- Docker Compose 1.29+
- Git

---

#### Option 1A: Run from Project Root (All Services)

Run complete platform with all services.

```bash
# 1. Clone repository
git clone https://github.com/thatlq1812/agrios.git
cd agrios

# 2. Start all services
# This will start:
#   - PostgreSQL (port 5432)
#   - Redis (port 6379)
#   - User Service (port 50051)
#   - Article Service (port 50052)
#   - Gateway (port 8080)
docker-compose up -d

# 3. Wait for initialization
sleep 15

# 4. Check service status
docker-compose ps

# Expected: All services Up (healthy)

# 5. Verify Gateway
curl http://localhost:8080/health

# Expected: {"status":"ok"}

# 6. View Gateway logs
docker logs agrios-gateway --tail 20
```

---

#### Option 1B: Run Gateway Only (Standalone)

Run Gateway with all backend services.

**Note:** Gateway needs both User Service and Article Service to function.

```bash
# 1. Clone repository
git clone https://github.com/thatlq1812/service-3-gateway.git
cd service-3-gateway

# 2. Configure environment (optional)
cp .env.example .env

# 3. Start Gateway with all dependencies
# This will start:
#   - PostgreSQL for users (port 5434)
#   - PostgreSQL for articles (port 5435)
#   - Redis (port 6381)
#   - User Service (port 50051)
#   - Article Service (port 50052)
#   - Gateway (port 8080)
docker-compose up -d

# 4. Wait for initialization
sleep 20

# 5. Check service status
docker-compose ps

# Expected output:
# gateway-postgres-user     Up (healthy)
# gateway-postgres-article  Up (healthy)
# gateway-redis             Up (healthy)
# gateway-user-service      Up (healthy)
# gateway-article-service   Up (healthy)
# gateway-app               Up (healthy)

# 6. Verify Gateway
curl http://localhost:8080/health

# Expected: {"status":"ok"}

# 7. View logs
docker logs gateway-app --tail 20
```

**Important Notes:**
- **All backend services included** - User Service + Article Service run automatically
- **PostgreSQL & Redis started automatically** - no manual installation needed
- Different ports used (5434, 5435, 6381) to avoid conflicts with root setup
- Database tables created automatically from migrations
- Gateway ready at http://localhost:8080

**Common Commands:**

```bash
# Rebuild after code changes (from root)
cd agrios
docker-compose up -d --build gateway

# Rebuild standalone (from service-3-gateway)
cd service-3-gateway
docker-compose up -d --build gateway

# Stop services
docker-compose down

# Remove volumes and clean data
docker-compose down -v

# Test API
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

---

### Option 2: Terminal (Local Development)

**Prerequisites:**
- Go 1.21+
- User Service running (port 50051)
- Article Service running (port 50052)

#### Step 1: Install Dependencies

```bash
cd service-3-gateway

# Download Go dependencies
go mod download

# Verify dependencies
go mod verify
```

#### Step 2: Start Backend Services

```bash
# Gateway requires both backend services to be running

# Option A: Start with Docker
docker-compose up -d user-service article-service

# Option B: Start manually in separate terminals
# Terminal 1: User Service
cd ../service-1-user
go run cmd/server/main.go

# Terminal 2: Article Service
cd ../service-2-article
go run cmd/server/main.go
```

#### Step 3: Configure Environment

```bash
# Copy example
cp .env.example .env

# Edit configuration
nano .env
```

**Required settings for local development:**
```env
# Backend Services
USER_SERVICE_HOST=localhost
USER_SERVICE_PORT=50051

ARTICLE_SERVICE_HOST=localhost
ARTICLE_SERVICE_PORT=50052

# Server
HTTP_PORT=8080

# CORS (for frontend development)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

#### Step 4: Build and Run

```bash
# Build
go build -o bin/gateway ./cmd/server

# Run
./bin/gateway

# Or run directly
go run cmd/server/main.go
```

**Expected output:**
```
2025/12/05 10:00:00 Connected to User Service at localhost:50051
2025/12/05 10:00:00 Connected to Article Service at localhost:50052
2025/12/05 10:00:00 API Gateway listening on :8080
```

#### Step 5: Verify Gateway

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected: {"status":"ok"}

# Test user registration
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"pass123"}'
```

---

## Environment Configuration

### Complete Environment Variables

```env
# Backend Services Configuration
USER_SERVICE_HOST=localhost         # User Service host (use 'user-service' for Docker)
USER_SERVICE_PORT=50051             # User Service port

ARTICLE_SERVICE_HOST=localhost      # Article Service host (use 'article-service' for Docker)
ARTICLE_SERVICE_PORT=50052          # Article Service port

# Gateway Server Configuration
HTTP_PORT=8080                      # HTTP server port

# CORS Configuration (for frontend apps)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With
CORS_ALLOW_CREDENTIALS=true

# Logging
LOG_LEVEL=info                      # debug, info, warn, error
LOG_FORMAT=json                     # json or text
```

### CORS Configuration

Enable CORS for frontend applications:

```env
# Development (allow all)
CORS_ALLOWED_ORIGINS=*

# Production (specific origins)
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

---

## API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check

```bash
GET /health

curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok"
}
```

---

### Authentication Endpoints

#### 1. Register User
```bash
POST /api/v1/users
Content-Type: application/json

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2025-12-05T10:00:00Z"
    }
  }
}
```

---

#### 2. Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Save tokens for subsequent requests:**
```bash
# Save to variable
export ACCESS_TOKEN="<your_access_token>"
export REFRESH_TOKEN="<your_refresh_token>"
```

---

#### 3. Refresh Token
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

#### 4. Logout
```bash
POST /api/v1/auth/logout
Authorization: Bearer <access_token>

curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:**
```json
{
  "code": "000",
  "message": "Logout successful"
}
```

---

### User Management Endpoints

#### 5. Get User
```bash
GET /api/v1/users/{id}

curl http://localhost:8080/api/v1/users/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2025-12-05T10:00:00Z",
      "updated_at": "2025-12-05T10:00:00Z"
    }
  }
}
```

---

#### 6. Update User
```bash
PUT /api/v1/users/{id}
Authorization: Bearer <access_token>
Content-Type: application/json

curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "johnsmith@example.com"
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "User updated successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Smith",
      "email": "johnsmith@example.com",
      "updated_at": "2025-12-05T11:00:00Z"
    }
  }
}
```

---

#### 7. Delete User
```bash
DELETE /api/v1/users/{id}
Authorization: Bearer <access_token>

curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:**
```json
{
  "code": "000",
  "message": "User deleted successfully"
}
```

---

#### 8. List Users
```bash
GET /api/v1/users?page=1&page_size=10

curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-12-05T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total_count": 25,
      "total_pages": 3
    }
  }
}
```

---

### Article Management Endpoints

#### 9. Create Article
```bash
POST /api/v1/articles
Authorization: Bearer <access_token>
Content-Type: application/json

curl -X POST http://localhost:8080/api/v1/articles \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Microservices",
    "content": "Microservices architecture is a design pattern..."
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "article": {
      "id": 1,
      "title": "Introduction to Microservices",
      "content": "Microservices architecture is...",
      "author": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
      },
      "created_at": "2025-12-05T10:00:00Z"
    }
  }
}
```

---

#### 10. Get Article
```bash
GET /api/v1/articles/{id}

curl http://localhost:8080/api/v1/articles/1
```

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "article": {
      "id": 1,
      "title": "Introduction to Microservices",
      "content": "Microservices architecture is...",
      "author": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
      },
      "created_at": "2025-12-05T10:00:00Z",
      "updated_at": "2025-12-05T10:00:00Z"
    }
  }
}
```

---

#### 11. Update Article
```bash
PUT /api/v1/articles/{id}
Authorization: Bearer <access_token>
Content-Type: application/json

curl -X PUT http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Advanced Microservices Patterns",
    "content": "Updated content..."
  }'
```

**Response:**
```json
{
  "code": "000",
  "message": "Article updated successfully",
  "data": {
    "article": {
      "id": 1,
      "title": "Advanced Microservices Patterns",
      "content": "Updated content...",
      "author": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
      },
      "updated_at": "2025-12-05T11:00:00Z"
    }
  }
}
```

---

#### 12. Delete Article
```bash
DELETE /api/v1/articles/{id}
Authorization: Bearer <access_token>

curl -X DELETE http://localhost:8080/api/v1/articles/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:**
```json
{
  "code": "000",
  "message": "Article deleted successfully"
}
```

---

#### 13. List Articles
```bash
GET /api/v1/articles?page=1&page_size=10&user_id=1

curl "http://localhost:8080/api/v1/articles?page=1&page_size=10&user_id=1"
```

**Query Parameters:**
- `page` (default: 1)
- `page_size` (default: 10)
- `user_id` (optional) - Filter by author

**Response:**
```json
{
  "code": "000",
  "message": "success",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "Introduction to Microservices",
        "content": "Microservices architecture is...",
        "author": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com"
        },
        "created_at": "2025-12-05T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 50
    }
  }
}
```

---

## Response Format

### Success Response

```json
{
  "code": "000",
  "message": "success",
  "data": {
    // Response data here
  }
}
```

### Error Response

```json
{
  "code": "001",
  "message": "Invalid email format"
}
```

### Error Codes

| Code | Meaning | HTTP Status | Example |
|------|---------|-------------|---------|
| 000 | Success | 200 | Operation completed |
| 001 | Invalid argument | 400 | Missing required field |
| 002 | Not found | 404 | User/Article not found |
| 003 | Already exists | 409 | Email already registered |
| 004 | Unauthorized | 401 | Invalid token |
| 005 | Permission denied | 403 | Cannot modify others' data |
| 006 | Internal error | 500 | Database connection failed |

---

## Testing

### Complete API Test Script

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=== 1. Register User ==="
curl -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"pass123"}'

echo -e "\n\n=== 2. Login ==="
LOGIN_RESP=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}')

echo $LOGIN_RESP | jq '.'

# Extract tokens
ACCESS_TOKEN=$(echo $LOGIN_RESP | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo $LOGIN_RESP | jq -r '.data.refresh_token')

echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"

echo -e "\n\n=== 3. Create Article ==="
curl -X POST $BASE_URL/articles \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Article","content":"This is a test article content."}'

echo -e "\n\n=== 4. Get Article ==="
curl $BASE_URL/articles/1

echo -e "\n\n=== 5. List Articles ==="
curl "$BASE_URL/articles?page=1&page_size=10"

echo -e "\n\n=== 6. Update Article ==="
curl -X PUT $BASE_URL/articles/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated content"}'

echo -e "\n\n=== 7. Refresh Token ==="
curl -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"

echo -e "\n\n=== 8. Logout ==="
curl -X POST $BASE_URL/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo -e "\n\n=== 9. Verify Token Invalid ==="
curl $BASE_URL/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

Save as `test-gateway.sh` and run:
```bash
chmod +x test-gateway.sh
./test-gateway.sh
```

---

## Troubleshooting

### curl not found

**Problem:** `curl: command not found` when testing APIs

**Solution:**

**Windows:**
```powershell
# curl is built-in on Windows 10+
# If missing, download from https://curl.se/windows/

# Alternative: Use PowerShell Invoke-WebRequest
Invoke-WebRequest -Uri http://localhost:8080/health
```

**Linux:**
```bash
sudo apt-get install curl
```

**macOS:**
```bash
# curl is pre-installed
# If missing: brew install curl
```

---

### Gateway Won't Start

**Problem:** Service fails to start

**Check logs:**
```bash
# Docker
docker-compose logs -f gateway

# Local
# Check terminal output
```

**Common causes:**
1. Backend services not running
2. Port 8080 already in use
3. Missing environment variables

**Solutions:**
```bash
# Check backend services
docker-compose ps user-service article-service

# Start backend services
docker-compose up -d user-service article-service

# Check port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                  # Linux/Mac

# Verify .env configuration
cat .env
```

---

### Cannot Connect to Backend Services

**Problem:** `failed to connect to user service` or `failed to connect to article service`

**Solutions:**
```bash
# 1. Check backend services are running
docker-compose ps user-service article-service

# 2. Test backend services directly
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50052 list

# 3. Verify HOST configuration
cat .env | grep SERVICE_HOST

# 4. Check Docker network (if using Docker)
docker network inspect agrios_default

# 5. Restart gateway after backend is ready
docker-compose restart gateway
```

---

### CORS Errors

**Problem:** Browser shows CORS errors

**Solutions:**
```bash
# 1. Add frontend origin to .env
CORS_ALLOWED_ORIGINS=http://localhost:3000

# 2. Restart gateway
docker-compose restart gateway

# 3. Verify CORS headers in response
curl -I http://localhost:8080/api/v1/users \
  -H "Origin: http://localhost:3000"

# Should see:
# Access-Control-Allow-Origin: http://localhost:3000
```

---

### Container Name Conflicts

**Problem:** `container name already in use`

**Cause:** Previous containers not removed

**Solutions:**
```bash
# List all containers
docker ps -a

# Remove specific containers
docker rm -f gateway-app
docker rm -f gateway-user-service
docker rm -f gateway-article-service
docker rm -f gateway-postgres-user
docker rm -f gateway-postgres-article
docker rm -f gateway-redis

# Or remove all stopped containers
docker container prune

# Then restart
docker-compose up -d
```

---

### Port Conflicts

**Problem:** `port is already allocated`

**Cause:** Another service using the same port

**Solutions:**
```bash
# Check what's using the ports
netstat -ano | findstr :8080    # Windows
netstat -ano | findstr :50051
netstat -ano | findstr :50052

lsof -i :8080                   # Linux/Mac
lsof -i :50051
lsof -i :50052

# Option 1: Stop conflicting service
# Find PID and kill it

# Option 2: Change ports in docker-compose.yml
# Edit ports section:
# - "8081:8080"  # Use 8081 instead of 8080
```

---

### 401 Unauthorized

**Problem:** API returns 401 even with token

**Possible causes:**
1. Token expired (access token: 15 minutes)
2. Token blacklisted (after logout)
3. Invalid token format
4. Missing Authorization header

**Solutions:**
```bash
# 1. Check token format
# Must be: Authorization: Bearer <token>

# 2. Use refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"

# 3. Login again
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'

# 4. Verify token is valid
# (Check User Service directly via grpcurl)
```

---

### 404 Not Found

**Problem:** Endpoint returns 404

**Solutions:**
```bash
# 1. Check endpoint path
# Correct: /api/v1/users
# Wrong: /users or /api/users

# 2. List available routes
curl http://localhost:8080/

# 3. Check method (GET/POST/PUT/DELETE)
```

---

## Project Structure

```
service-3-gateway/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── client/
│   │   ├── user_client.go       # User Service gRPC client
│   │   └── article_client.go    # Article Service gRPC client
│   ├── handler/
│   │   ├── user_handler.go      # User endpoints
│   │   ├── article_handler.go   # Article endpoints
│   │   └── health_handler.go    # Health check
│   ├── middleware/
│   │   ├── auth_middleware.go   # JWT authentication
│   │   ├── cors_middleware.go   # CORS headers
│   │   └── logging_middleware.go # Request logging
│   ├── response/
│   │   └── response.go          # Response formatting
│   └── config/
│       └── config.go            # Configuration
├── .env.example                 # Environment template
├── Dockerfile                   # Docker configuration
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

---

## Development Commands

```bash
# Install dependencies
go mod download

# Update dependencies
go mod tidy

# Build
go build -o bin/gateway ./cmd/server

# Run
./bin/gateway

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

---

## Middleware

### Authentication Middleware

Protects endpoints requiring authentication:

```go
// internal/middleware/auth_middleware.go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        
        // Validate token with User Service
        // If valid, add user_id to request context
        
        next.ServeHTTP(w, r)
    })
}
```

**Protected endpoints:**
- POST /api/v1/articles
- PUT /api/v1/articles/{id}
- DELETE /api/v1/articles/{id}
- PUT /api/v1/users/{id}
- DELETE /api/v1/users/{id}
- POST /api/v1/auth/logout

---

## Additional Resources

- **[Main Project README](../README.md)** - Complete platform documentation
- **[Deployment Guide](../DEPLOYMENT.md)** - Production deployment steps
- **[User Service](../service-1-user/README.md)** - Authentication & JWT documentation
- **[Article Service](../service-2-article/README.md)** - Content management & graceful degradation

---

**Service Version:** 1.0.0  
**Last Updated:** December 5, 2025  
**Maintainer:** thatlq1812@gmail.com
