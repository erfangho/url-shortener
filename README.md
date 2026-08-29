# URL Shortener API

A URL shortener built with Go, Gin, and GORM — featuring concurrent analytics processing, gRPC microservice communication, and Docker containerization.

> **Companion Project:** [url-shortener-analytics](https://github.com/erfangho/url-shortener-analytics) — gRPC analytics dashboard that receives click events in real-time.

## Architecture

```
┌─────────────────┐     gRPC      ┌──────────────────────┐
│  URL Shortener  │ ───────────→  │  Analytics Dashboard │
│  (Gin + GORM)   │  RecordClick  │  (gRPC Server)       │
│  Port: 8080     │               │  Port: 50051         │
└────────┬────────┘               └──────────────────────┘
         │
    ┌────┴────┐
    │  Redis  │  ← Cache layer
    └─────────┘
```

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go 1.27 |
| HTTP Framework | Gin |
| ORM | GORM (pure-Go SQLite) |
| Authentication | JWT (golang-jwt) |
| Caching | In-memory (RWMutex) + Redis |
| Inter-service Communication | gRPC + Protocol Buffers |
| Containerization | Docker (multi-stage build) + Docker Compose |
| API Documentation | Swagger (swaggo) |
| Password Hashing | bcrypt |

## Project Structure

```
url-shortener/
├── cmd/server/          → Application entrypoint
├── internal/
│   ├── config/          → Database, Redis, JWT configuration
│   ├── grpc/            → gRPC client for analytics service
│   ├── handler/         → HTTP request handlers (Gin)
│   ├── middleware/       → Auth middleware (JWT validation)
│   ├── model/           → Domain models (GORM)
│   ├── repository/      → Database access layer
│   ├── routes/          → Route registration
│   └── service/         → Business logic layer
├── pkg/cache/           → In-memory cache with TTL + Redis cache
├── proto/               → Protocol Buffer definitions
├── docs/                → Swagger documentation
├── Dockerfile           → Multi-stage Docker build
└── docker-compose.yml   → Container orchestration
```

## Features

### Core
- **URL Shortening** — Generate short codes with collision detection
- **Redirect** — 301 redirect to original URL
- **User Authentication** — Register, login with JWT tokens
- **Protected Routes** — Middleware-based JWT validation

### Concurrency Patterns
- **Worker Pool** — Background goroutines process analytics events from a buffered channel
- **Batch Processing** — Events are batched and saved to DB in single transactions
- **Ticker-based Flushing** — Time-based flush ensures events are processed even under low traffic
- **Graceful Shutdown** — Channel draining + WaitGroup for clean process exit
- **Signal Handling** — OS signal (SIGINT) triggers context cancellation

### Caching
- **In-memory Cache** — Thread-safe (RWMutex) with TTL and background cleanup
- **Redis Cache** — Distributed caching for user data
- **Cache-aside Pattern** — Check cache first, fallback to DB, update cache

### Microservices (gRPC)
- **Two-project Architecture** — URL Shortener + Analytics Dashboard as separate services
- **Protocol Buffers** — Service contract definition with code generation
- **Unary RPC** — RecordClick sends click events to dashboard
- **Client Wrapper Pattern** — Hides proto types behind domain models

### Production
- **Context Propagation** — `context.Context` flows through handler → service → repository
- **Structured Logging** — slog with file output
- **Environment Configuration** — `.env` files, no hardcoded secrets
- **Docker Multi-stage Build** — ~15MB production images with BuildKit cache mounts

## API Endpoints

| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/urls` | Create short URL | No |
| GET | `/urls` | List all URLs (paginated) | Yes |
| GET | `/urls/:shortCode` | Get URL info | No |
| GET | `/u/:shortCode` | Redirect to original URL | No |
| POST | `/users` | Register new user | No |
| GET | `/users` | List all users (paginated) | Yes |
| POST | `/login` | Login and get JWT token | No |

### API Testing (Bruno)

A [Bruno](https://www.usebruno.com/) API collection is included in `docs/bruno/` with pre-configured requests for all endpoints:

```
docs/bruno/url-shortener/
├── urls/           → Create, Get, List, Redirect
├── users/          → Register, List users
├── auth/           → Login
└── environments/   → Development environment config
```

Import the collection in Bruno to test all endpoints without manual cURL commands.

## Concurrency Deep Dive

### Analytics Pipeline

```
HTTP Request → Publish(event) → Buffered Channel → Worker Pool → Batch → DB Transaction
                                                    ↓
                                              gRPC Client → Analytics Dashboard
```

- **Buffered Channel** (size 100) — Thread-safe producer-consumer queue
- **3 Worker Goroutines** — Concurrent consumers reading from channel
- **Batch Size 2** — Events grouped before DB write
- **3-Second Ticker** — Flushes incomplete batches on timer
- **sync.WaitGroup** — Tracks active workers for graceful shutdown
- **Channel Close** — Signals workers to flush remaining events and exit

### Cache Concurrency

- **sync.RWMutex** — Multiple readers or single writer on in-memory cache
- **Background Cleanup** — Ticker-based goroutine removes expired entries under write lock
- **Redis** — Thread-safe by design, handles concurrent access natively

## Running Locally

```bash
# Without Docker
go run cmd/server/main.go

# With Docker Compose
docker compose up --build
```

Services:
- URL Shortener: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Analytics Dashboard: `localhost:50051` (gRPC)

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `SECRET_KEY` | JWT signing secret | (required) |
| `ANALYTICS_ADDR` | gRPC analytics server address | `localhost:50051` |
| `REDIS_ADDR` | Redis server address | `localhost:6379` |

## Analytics Dashboard

A separate gRPC service that receives click events in real-time.

```bash
cd url-shortener-analytics
go run cmd/server/main.go
```

Test with grpcurl:
```bash
grpcurl -plaintext -d '{"url_id": 1, "user_agent": "test", "ip_address": "127.0.0.1"}' \
  localhost:50051 analytics.AnalyticsService/RecordClick
```

## Key Go Concepts Demonstrated

| Concept | Implementation |
|---|---|
| Goroutines | Worker pool, background cleanup, graceful shutdown |
| Channels | Buffered event queue, signal handling |
| select | Timer-based batch flushing with channel reads |
| sync.WaitGroup | Worker lifecycle management |
| sync.RWMutex | Thread-safe in-memory cache |
| context.Context | Request cancellation propagation |
| interfaces | Repository pattern, dependency injection |
| struct embedding | UnimplementedAnalyticsServiceServer |
| error wrapping | Sentinel errors, gorm.ErrRecordNotFound |
| struct tags | JSON serialization, GORM column mapping, binding validation |
| Protocol Buffers | Service contract, message definitions |
| gRPC | Unary RPC, client wrapper pattern |

## License

MIT
