# OnlineLeadership API

RESTful API service for managing game leaderboards with real-time ranking, score tracking, and user authentication.

## Tech Stack

- **Go 1.25** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL 16** - Primary database (user data, score history)
- **Redis 7** - In-memory store (leaderboard rankings, sorted sets)
- **JWT** - Authentication (access & refresh tokens)
- **Swagger/OpenAPI** - API documentation (`swaggo/swag`)
- **Docker & Docker Compose** - Containerization
- **sqlx** - SQL query builder
- **Viper** - Configuration management
- **slog** - Structured logging

## Architecture

The project follows **Clean Architecture** principles with clear layer separation:

```
┌─────────────────────────────────────────┐
│  HTTP Layer (Gin handlers)              │  ← DTOs use string IDs
├─────────────────────────────────────────┤
│  Use Case Layer (Business logic)        │  ← Uses uuid.UUID
├─────────────────────────────────────────┤
│  Domain Layer (Entities)                │  ← Pure domain models
├─────────────────────────────────────────┤
│  Repository Layer (Data access)         │  ← PostgreSQL + Redis
└─────────────────────────────────────────┘
```

### Layers

- **`internal/interfaces/http/handler`** - HTTP handlers, request/response DTOs, Swagger annotations
- **`internal/usecase`** - Business logic services (auth, admin, leaderboard, score)
- **`internal/domain`** - Domain models (User, Game, LeaderboardUser)
- **`internal/infrastructure`** - External dependencies (PostgreSQL, Redis, JWT, logger)

### ID Type Consistency

- **Domain/Service/Repository layers**: Use `uuid.UUID`
- **HTTP layer (DTOs, requests)**: Use `string`
- **Conversion**: Happens only at HTTP boundary (handlers)

## Features

### Authentication
- User registration with email and password
- JWT-based authentication (access + refresh tokens)
- Access token TTL: 30 minutes
- Refresh token TTL: 7 days
- Password hashing with bcrypt

### Game Management
- Create new games
- List all available games

### Score Tracking
- Submit player scores for specific games
- Persistent score history (PostgreSQL)
- Automatic leaderboard updates (Redis sorted sets)

### Leaderboards
- Global leaderboard (all players across all games)
- Per-game leaderboards
- User rank retrieval
- Pagination support (offset/limit)

## API Documentation

Swagger UI is available at: **http://localhost:8080/swagger/index.html**

### Authentication

Protected endpoints require JWT token in `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Endpoints

#### Public Endpoints
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and receive tokens
- `POST /admin/create` - Create a new game
- `GET /admin/games` - List all games

#### Protected Endpoints (require JWT)
- `POST /api/score/submit` - Submit player score
- `GET /api/leaderboard/global` - Get global leaderboard
- `GET /api/leaderboard/my` - Get current user's rank
- `POST /api/leaderboard/top` - Get top players for a specific game

## Environment Variables

Create a `.env` file in the project root:

```bash
# JWT Secrets (required)
JWT_ACCESS_SECRET=your-secret-access-key-here
JWT_REFRESH_SECRET=your-secret-refresh-key-here

# Database (optional, defaults in config.yml)
DB_PASSWORD=postgres
```

### Configuration Files

**`config.yml`** - Application configuration:
```yaml
port: "8080"
db:
  username: "postgres"
  host: "postgres"      # Use "localhost" for local development
  port: 5432
  dbname: "leaderboard"
  sslmode: "disable"
```

## Project Structure

```
OnlineLeadership/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain models (User, Game, LeaderboardUser)
│   ├── usecase/                 # Business logic services
│   │   ├── auth/
│   │   ├── admin/
│   │   ├── leaderboard/
│   │   └── score_history/
│   ├── infrastructure/          # External dependencies
│   │   ├── auth/                # JWT token manager
│   │   ├── logger/              # Structured logging
│   │   ├── postgres/            # Database connection & repositories
│   │   └── redis/               # Redis client
│   └── interfaces/
│       └── http/
│           ├── handler/         # HTTP handlers & DTOs
│           └── middleware/      # Request ID, authentication
├── migrations/                  # SQL migrations
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── docs/                        # Auto-generated Swagger docs
├── config.yml                   # App configuration
├── docker-compose.yml           # Docker orchestration
├── Dockerfile                   # Container image
└── go.mod                       # Go dependencies
```

## Running the Application

### Prerequisites
- Go 1.25+
- PostgreSQL 16
- Redis 7
- Docker & Docker Compose (optional)

### Option 1: Using Docker Compose (Recommended)

1. **Create `.env` file** with JWT secrets:
```bash
JWT_ACCESS_SECRET=your-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-key-change-in-production
DB_PASSWORD=postgres
```

2. **Start all services**:
```bash
docker-compose up -d
```

This will start:
- Application server on `http://localhost:8080`
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`
- Automatic database migrations

3. **Access Swagger UI**:
```
http://localhost:8080/swagger/index.html
```

4. **Stop services**:
```bash
docker-compose down
```

### Option 2: Local Development

1. **Install dependencies**:
```bash
go mod download
```

2. **Start PostgreSQL and Redis** (manually or via Docker):
```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=leaderboard postgres:16
docker run -d -p 6379:6379 redis:7
```

3. **Run migrations**:
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/leaderboard?sslmode=disable" up
```

4. **Update `config.yml`** - change `host` to `localhost`:
```yaml
db:
  host: "localhost"
```

5. **Set environment variables**:
```bash
export JWT_ACCESS_SECRET="your-access-secret"
export JWT_REFRESH_SECRET="your-refresh-secret"
export DB_PASSWORD="postgres"
```

6. **Generate Swagger docs** (if not present):
```bash
swag init -g cmd/app/main.go -o ./docs
```

7. **Run the application**:
```bash
go run cmd/app/main.go
```

Server will start on `http://localhost:8080`

## Usage Examples

### 1. Register a User
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### 2. Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "password": "password123"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Create a Game
```bash
curl -X POST http://localhost:8080/admin/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chess"
  }'
```

Response:
```json
{
  "game_id": "987fcdeb-51a2-43f7-9876-543210fedcba"
}
```

### 4. Submit Score (Protected)
```bash
curl -X POST http://localhost:8080/api/score/submit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "game_id": "987fcdeb-51a2-43f7-9876-543210fedcba",
    "score": 1500
  }'
```

Response:
```json
{
  "status": "ok"
}
```

### 5. Get Global Leaderboard (Protected)
```bash
curl -X GET "http://localhost:8080/api/leaderboard/global?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Response:
```json
{
  "data": [
    {
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "score": 1500,
      "rank": 1
    }
  ]
}
```

### 6. Get My Rank (Protected)
```bash
curl -X GET http://localhost:8080/api/leaderboard/my \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Response:
```json
{
  "rank": 1
}
```

## Common Errors & Troubleshooting

### 1. "JWT secrets are not set"
**Solution**: Ensure `.env` file contains `JWT_ACCESS_SECRET` and `JWT_REFRESH_SECRET`

### 2. "db connect failed"
**Causes**:
- PostgreSQL not running
- Wrong credentials
- Wrong host (use `localhost` for local dev, `postgres` for Docker)

**Solution**: Check PostgreSQL connection and update `config.yml` or environment variables

### 3. "401 Unauthorized" on protected endpoints
**Causes**:
- Missing `Authorization` header
- Invalid or expired token
- Wrong token format

**Solution**: 
- Include header: `Authorization: Bearer <access_token>`
- Login again to get a fresh token
- Ensure "Bearer " prefix is included

### 4. "400 Bad Request - invalid game_id format"
**Cause**: Invalid UUID format in request body

**Solution**: Use valid UUID format (e.g., `123e4567-e89b-12d3-a456-426614174000`)

### 5. Swagger authentication not working
**Solution**: 
1. Click "Authorize" button in Swagger UI
2. Enter: `Bearer <your_access_token>` (include "Bearer " prefix)
3. Click "Authorize" then "Close"

### 6. "redis connection error"
**Cause**: Redis server not running

**Solution**: 
```bash
docker run -d -p 6379:6379 redis:7
```

## Database Schema

### Tables

**`users`**
- `id` (UUID, PK)
- `username` (TEXT, UNIQUE)
- `password_hash` (TEXT)
- `email` (TEXT, UNIQUE)
- `created_at` (TIMESTAMP)

**`games`**
- `id` (UUID, PK)
- `name` (TEXT, UNIQUE)

**`score_history`**
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `game_id` (UUID, FK → games)
- `score` (INT)
- `created_at` (TIMESTAMP)

### Redis Data Structures

- **Global leaderboard**: Sorted set `leaderboard:global`
  - Members: user IDs
  - Scores: total points

- **Game leaderboards**: Sorted set `leaderboard:game:{game_id}`
  - Members: user IDs
  - Scores: game-specific points

## Development

### Regenerate Swagger Documentation
```bash
swag init -g cmd/app/main.go -o ./docs
```

### Run Tests
```bash
go test ./...
```

### Build Binary
```bash
go build -o bin/app cmd/app/main.go
```

## Security Considerations

- **Production**: Change JWT secrets to strong, random values
- **HTTPS**: Use HTTPS in production (configure reverse proxy)
- **Rate Limiting**: Implement rate limiting for public endpoints
- **CORS**: Configure CORS if serving frontend from different origin
- **Admin Endpoints**: Add authentication/authorization to `/admin/*` routes

## License

MIT

## Support

For issues and questions, open an issue on the project repository.

