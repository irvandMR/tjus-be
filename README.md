# TJUS Backend API (PajangAset)

A clean and fast REST API for an asset consignment catalog (cars, heavy equipment, property, etc.) built with **Go (Golang)**, **Gin**, **PostgreSQL**, and **Redis**.

---

## Tech Stack

- **Language & Framework**: Go 1.26, Gin Web Framework
- **Database**: PostgreSQL 16 (`pgx/v5` raw SQL connection pool, no ORM)
- **Database Migration**: Goose (auto-runs on server startup)
- **Caching & Security**: Redis 7 (`go-redis/v9`)
- **Authentication**: JWT (`golang-jwt/jwt/v5`) + `bcrypt` password hashing
- **Logging**: Go standard `log/slog` (Structured JSON logs in production)
- **Testing**: Go `testing`, `testify` (`assert`, `require`, `mock`), `httptest`
- **DevOps**: Multi-stage Docker, Docker Compose, GitHub Actions CI
- **API Docs**: OpenAPI 3.0 (Swagger UI)

---

## Key Features

- **Modular Feature-Based Architecture**: Code is organized cleanly by domain (`asset`, `category`, `auth`, `setting`, `health`).
- **Public vs Admin Data Separation**: Public visitors can browse active assets without seeing private owner phone numbers or internal notes.
- **Dynamic Attributes per Category**: Each category (e.g., Car, Property) can define its own custom attributes (`Year`, `Transmission`, `Land Size`, etc.).
- **Search, Filter & Dynamic Sorting**: Supports keyword search, category filter, min/max price, pagination, and safe sorting (`newest`, `oldest`, `price_asc`, `price_desc`, `popular`).
- **Redis Caching (Cache-Aside)**: Asset detail pages are cached in Redis for fast response times and automatically cleared when an asset is updated or deleted.
- **JWT Logout Blacklist & Rate Limiting**: Logged-out tokens are stored in Redis until expiry, and sensitive endpoints (`/auth/login` and media upload) are protected against brute-force/spam with Redis rate limiting (`429 Too Many Requests`).
- **Safe Media Upload**: Checks file size, file extension, and real file content (`MIME` type sniffing) before saving images.
- **Graceful Shutdown**: Waits for active requests to finish cleanly and closes PostgreSQL & Redis connections safely when the server stops.

---

## Project Structure

```text
tjus-be/
├── .github/workflows/   # GitHub Actions CI pipeline (vet, test, build)
├── cmd/tjus/            # Main application entrypoint (main.go)
├── docs/                # OpenAPI spec (swagger.yaml) & Swagger UI
├── internal/
│   ├── app/             # Server setup, router wiring & graceful shutdown
│   ├── config/          # Env loader, PostgreSQL pool, Redis client & Goose migration
│   ├── features/        # Business modules (asset, auth, category, health, setting)
│   ├── middleware/      # JWT Auth, Redis Rate Limit & Structured Logger (slog)
│   └── pkg/             # Shared generic utilities (httputil, redisutil) & base entity
├── migrations/          # SQL schema migrations
├── uploads/             # Local folder for uploaded asset images
├── Dockerfile           # Multi-stage Docker build
└── docker-compose.yml   # Runs App + PostgreSQL + Redis together
```

---

## Quick Start (With Docker Compose)

The easiest way to run the full stack (API + PostgreSQL + Redis) is using Docker.

### 1. Copy the environment file
```bash
cp .env.example .env
```

### 2. Start all containers
```bash
docker compose up -d --build
```
Database tables are created automatically by **Goose** when the app starts.

### 3. Check if the server is running
- **Health Check**: `http://localhost:8080/api/v1/health`
- **Swagger API Docs**: `http://localhost:8080/swagger`

---

## Running Locally (Without Docker for App)

If you want to run the Go app directly on your machine:

1. Make sure PostgreSQL and Redis are running, and update `DB_URL` and `REDIS_ADDR` in `.env` (for example `localhost:5433` and `localhost:6380`).
2. Download dependencies and start the server:
   ```bash
   go mod download
   go run ./cmd/tjus/main.go
   ```

---

## Running Tests & Code Coverage

Unit tests use Mock repositories, so they run fast **without needing a live database**.

```bash
# Run linter check
go vet ./...

# Run all unit tests
go test -v ./...

# Generate and open visual HTML code coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Main API Endpoints (`/api/v1`)

Full request and response examples are available at **`http://localhost:8080/swagger`**.

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` & `/health/db` | Public | Check server and database status |
| `POST` | `/auth/login` | Public | Admin login (Rate limited: 5x/min) |
| `POST` | `/auth/logout` | Admin | Revoke JWT token in Redis blacklist |
| `GET` | `/assets` | Public | List active assets (search, filter, sort, pagination) |
| `GET` | `/assets/slug/:slug` | Public | Get asset detail by slug (cached in Redis) |
| `GET` | `/assets/admin` | Admin | List all assets for Admin dashboard (includes owner info) |
| `GET` | `/assets/:id` | Admin | Get full asset detail by ID |
| `POST` | `/assets` | Admin | Create asset + owner + attributes in one DB transaction |
| `POST` | `/assets/update` | Admin | Modular asset update (general, owner, attributes, status) |
| `POST` | `/assets/:id/media` | Admin | Upload asset images (Rate limited: 10x/min) |
| `DELETE` | `/assets/media/:media_id` | Admin | Delete an asset image |
| `DELETE` | `/assets/:id` | Admin | Soft-delete an asset |
| `GET` | `/category` | Public | List all categories |
| `GET` | `/category/:id/attributes` | Public | List attributes for a category |
| `POST` | `/category` & `/category/update` | Admin | Create or update a category |
| `DELETE` | `/category/:id` | Admin | Delete a category |
| `POST` | `/category/attribute` & `/category/attribute/update` | Admin | Create or update category attributes |
| `DELETE` | `/category/attribute/:attr_id` | Admin | Delete a category attribute |
| `GET` / `POST` | `/setting` & `/setting/update` | Admin | View, create, or update site settings |
