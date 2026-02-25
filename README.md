# OnlineLeadership API

> RESTful API сервис для управления игровыми таблицами лидеров с рейтингом в реальном времени, отслеживанием очков, аутентификацией пользователей и полноценным стеком мониторинга.

---

## 🇷🇺 Русская версия

---

## Ключевые особенности

- 🏆 Таблицы лидеров в реальном времени с использованием Redis Sorted Sets
- 🔐 JWT аутентификация (access + refresh токены)
- 📊 **Реализован полноценный стек наблюдаемости на основе Prometheus и Grafana** с мониторингом задержки на уровне отдельных endpoint'ов и отслеживанием ресурсов контейнеров
- 🐳 Полная контейнеризация с Docker Compose
- 🧱 Чистая архитектура (Clean Architecture) с чётким разделением слоёв
- 📖 Документация API через Swagger UI

---

## Технологический стек

| Компонент | Технология |
|---|---|
| Язык | Go 1.25 |
| HTTP фреймворк | Gin |
| База данных | PostgreSQL 16 |
| Кэш / Лидерборд | Redis 7 (Sorted Sets) |
| Аутентификация | JWT (access + refresh) |
| Документация API | Swagger (swaggo/swag) |
| Мониторинг метрик | Prometheus |
| Дашборды | Grafana |
| Мониторинг контейнеров | cAdvisor |
| Конфигурация | Viper (config.yml + .env) |
| Логирование | slog (структурированные логи) |
| Контейнеризация | Docker & Docker Compose |

---

## Архитектура

Проект построен по принципам **Clean Architecture** с чётким разделением ответственности между слоями:

```
┌──────────────────────────────────────────────────────────┐
│  HTTP Layer  (Gin handlers, DTO, Swagger, Middleware)    │
│  internal/interfaces/http/handler/                       │
│  → string IDs в DTO, конвертация UUID на границе слоя    │
├──────────────────────────────────────────────────────────┤
│  Use Case Layer  (Бизнес-логика)                         │
│  internal/usecase/                                       │
│  → Работает только с uuid.UUID                          │
├──────────────────────────────────────────────────────────┤
│  Domain Layer  (Сущности, интерфейсы)                    │
│  internal/domain/                                        │
│  → Чистые доменные модели без зависимостей               │
├──────────────────────────────────────────────────────────┤
│  Infrastructure Layer  (PostgreSQL, Redis, JWT, Logger)  │
│  internal/infrastructure/                                │
│  → Реализации репозиториев, внешние зависимости          │
└──────────────────────────────────────────────────────────┘
```

### Правило типов ID

- **Domain / Service / Repository**: используют `uuid.UUID`
- **HTTP слой (DTO, запросы)**: используют `string`
- **Конвертация**: происходит **только** в handler'ах (`uuid.Parse`)

---

## Наблюдаемость и мониторинг

```
┌─────────────────────────────────────────────────────────────┐
│                      Docker Network                         │
│                                                             │
│  ┌──────────┐  /metrics   ┌────────────┐   PromQL          │
│  │ Backend  │ ──────────► │ Prometheus │ ◄────────────────┐ │
│  │  :8080   │             │   :9090    │                  │ │
│  └──────────┘             └────────────┘                  │ │
│                                    │                       │ │
│  ┌──────────┐  /metrics            ▼                       │ │
│  │ cAdvisor │ ──────────► ┌────────────┐  datasource:     │ │
│  │  :8081   │             │  Grafana   │ ◄── prometheus:9090│
│  └──────────┘             │   :3000    │                   │ │
│                           └────────────┘                   │ │
└─────────────────────────────────────────────────────────────┘
```

### Метрики приложения (Prometheus)

Бэкенд экспортирует метрики по адресу `backend:8080/metrics`.

| Метрика | Тип | Описание |
|---|---|---|
| `http_request_duration_seconds` | Histogram | Время обработки запроса по endpoint'у и методу |
| `http_requests_total` | Counter | Общее количество запросов (метод, путь, статус) |
| `jwt_auth_failures_total` | Counter | Количество ошибок JWT аутентификации |
| `db_query_duration_seconds` | Histogram | Время выполнения запросов к PostgreSQL |
| `redis_operation_duration_seconds` | Histogram | Время операций с Redis |

### Мониторинг задержки (P95)

Для анализа производительности используйте P95 латентность:

```promql
# P95 задержка для всех HTTP endpoint'ов
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path)
)

# P95 задержка запросов к PostgreSQL
histogram_quantile(0.95,
  sum(rate(db_query_duration_seconds_bucket[5m])) by (le)
)

# P95 задержка операций Redis
histogram_quantile(0.95,
  sum(rate(redis_operation_duration_seconds_bucket[5m])) by (le)
)
```

### Мониторинг аутентификации JWT

```promql
# Скорость ошибок JWT за последние 5 минут
rate(jwt_auth_failures_total[5m])

# Суммарное количество ошибок JWT
sum(jwt_auth_failures_total)
```

### Мониторинг endpoint'ов

```promql
# RPS (запросов в секунду) по endpoint'ам
sum(rate(http_requests_total[1m])) by (path, method)

# Количество ошибок 5xx
sum(rate(http_requests_total{status=~"5.."}[5m])) by (path)

# Количество ошибок 4xx
sum(rate(http_requests_total{status=~"4.."}[5m])) by (path)
```

### Мониторинг контейнеров (cAdvisor)

cAdvisor собирает метрики Docker-контейнеров и экспортирует их в Prometheus.

```promql
# CPU utilization по контейнерам
sum(rate(container_cpu_usage_seconds_total{name!=""}[1m])) by (name)

# Использование памяти по контейнерам
sum(container_memory_usage_bytes{name!=""}) by (name)

# Сетевой трафик (входящий)
sum(rate(container_network_receive_bytes_total{name!=""}[1m])) by (name)
```

### Grafana дашборды

Grafana доступна по адресу: **http://localhost:3000**

- **Логин по умолчанию**: `admin` / `admin`
- **Data source**: Prometheus (`http://prometheus:9090`)

Рекомендуемые дашборды для импорта:
- **Gin HTTP metrics** — RPS, задержка, коды ответов
- **Go Runtime** — goroutines, GC, память (ID: `13240`)
- **PostgreSQL** — время запросов, соединения
- **Redis** — операции, hit rate
- **cAdvisor** — CPU, RAM контейнеров (ID: `14282`)

---

## Сервисы Docker Compose

| Сервис | Порт | Описание |
|---|---|---|
| `backend` | `8080` | Go приложение + `/metrics` endpoint |
| `postgres` | `5432` | PostgreSQL 16 |
| `redis` | `6379` | Redis 7 |
| `prometheus` | `9090` | Prometheus (scrape: `backend:8080/metrics`, `cadvisor:8081`) |
| `grafana` | `3000` | Grafana (datasource: `http://prometheus:9090`) |
| `cadvisor` | `8081` | Container CPU/RAM мониторинг |

---

## Документация API

Swagger UI: **http://localhost:8080/swagger/index.html**

### Публичные endpoint'ы

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/auth/register` | Регистрация пользователя |
| `POST` | `/auth/login` | Вход, получение токенов |
| `POST` | `/admin/create` | Создание новой игры |
| `GET` | `/admin/games` | Список всех игр |

### Защищённые endpoint'ы (требуют `Authorization: Bearer <token>`)

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/api/score/submit` | Отправка очков |
| `GET` | `/api/leaderboard/global` | Глобальный лидерборд |
| `GET` | `/api/leaderboard/my` | Ранг текущего пользователя |
| `POST` | `/api/leaderboard/top` | Топ игроков по игре |

---

## Переменные окружения

Создайте файл `.env` в корне проекта:

```env
JWT_ACCESS_SECRET=your-access-secret-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-change-in-production
DB_PASSWORD=postgres
```

### config.yml (основная конфигурация)

```yaml
port: "8080"
db:
  username: "postgres"
  host: "postgres"       # "localhost" для локальной разработки
  port: 5432
  dbname: "leaderboard"
  sslmode: "disable"
```

---

## Запуск

### Требования

- Go 1.25+
- Docker & Docker Compose

### Вариант 1: Docker Compose (рекомендуется)

```bash
# 1. Создайте .env файл
cp .env.example .env   # или создайте вручную

# 2. Запустите все сервисы
docker-compose up -d

# 3. Проверьте статус
docker-compose ps
```

**Доступные сервисы:**
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- cAdvisor: http://localhost:8081

```bash
# Остановить сервисы
docker-compose down

# Остановить и удалить volumes
docker-compose down -v
```

### Вариант 2: Локальная разработка

```bash
# 1. Запустите PostgreSQL и Redis
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=leaderboard postgres:16
docker run -d -p 6379:6379 redis:7

# 2. Установите зависимости
go mod download

# 3. Обновите config.yml: db.host = "localhost"

# 4. Установите переменные окружения
set JWT_ACCESS_SECRET=your-access-secret
set JWT_REFRESH_SECRET=your-refresh-secret
set DB_PASSWORD=postgres

# 5. Сгенерируйте Swagger (при необходимости)
swag init -g cmd/app/main.go -o ./docs

# 6. Запустите приложение
go run cmd/app/main.go
```

---

## Структура проекта

```
OnlineLeadership/
├── cmd/app/main.go                    # Точка входа
├── config/config.go                   # Конфигурация (Viper)
├── config.yml                         # Настройки приложения
├── internal/
│   ├── domain/                        # Доменные модели и интерфейсы
│   ├── usecase/                       # Бизнес-логика
│   │   ├── auth/
│   │   ├── admin/
│   │   ├── leaderboard/
│   │   └── score_history/
│   ├── infrastructure/
│   │   ├── auth/                      # JWT менеджер
│   │   ├── logger/                    # slog
│   │   ├── monitoring/                # Prometheus метрики
│   │   ├── postgres/                  # БД + репозитории
│   │   └── redis/                     # Redis клиент
│   └── interfaces/http/
│       ├── handler/                   # Gin обработчики, DTO, Swagger
│       └── middleware/                # Request ID, Auth, Metrics
├── migrations/                        # SQL миграции
├── monitoring/
│   └── prometheus.yml                 # Конфигурация scrape
├── docs/                              # Автогенерированная Swagger документация
├── docker-compose.yml
└── Dockerfile
```

---

## Схема базы данных

**`users`** — `id (UUID)`, `username`, `email`, `password_hash`, `created_at`

**`games`** — `id (UUID)`, `name`

**`score_history`** — `id (UUID)`, `user_id (FK)`, `game_id (FK)`, `score`, `created_at`

**Redis Sorted Sets:**
- `leaderboard:global` — глобальный рейтинг (member: user_id, score: сумма очков)
- `leaderboard:game:{game_id}` — рейтинг по игре

---

## Распространённые ошибки

| Ошибка | Причина | Решение |
|---|---|---|
| `JWT secrets are not set` | Не задан `.env` | Создайте `.env` с JWT секретами |
| `db connect failed` | PostgreSQL недоступен | Проверьте хост в `config.yml` |
| `401 Unauthorized` | Нет/истёк токен | Войдите заново, добавьте `Bearer ` |
| `400 invalid UUID` | Неверный формат ID | Используйте UUID формат |
| Swagger auth не работает | Не введён `Bearer ` | В Swagger UI: `Bearer <token>` |
| `redis connection error` | Redis не запущен | `docker run -d -p 6379:6379 redis:7` |
| Prometheus не видит метрики | Неверный scrape target | Проверьте `monitoring/prometheus.yml` |

---

## Рекомендации по безопасности

- **JWT Secrets**: Используйте сильные случайные секреты (мин. 64 символа), никогда не коммитте в VCS
- **HTTPS**: Используйте обратный прокси (Nginx/Traefik) с TLS в продакшене
- **Rate Limiting**: Добавьте ограничение частоты запросов для публичных endpoints
- **CORS**: Настройте CORS, если фронтенд обслуживается с другого домена
- **Admin Routes**: Защитите маршруты `/admin/*` с помощью аутентификации
- **Grafana**: Измените пароль администратора по умолчанию при первом входе
- **Prometheus**: Ограничьте доступ к `/metrics` в продакшене (например, только для внутренней сети)
- **Alerting**: Настройте Prometheus Alertmanager для оповещений об ошибках JWT и всплесках задержки

---

## Лицензия

MIT

## Поддержка

По вопросам и проблемам создавайте issue в репозитории проекта.

---

---

## 🇬🇧 English Version

---

# OnlineLeadership API

> A production-grade RESTful API service for real-time gaming leaderboards with score tracking, JWT authentication, and a full observability stack.

---

## Key Features

- 🏆 Real-time leaderboards powered by Redis Sorted Sets
- 🔐 JWT authentication with access + refresh tokens
- 📊 **Implemented a real-time observability stack using Prometheus and Grafana** with endpoint-level latency monitoring and container resource tracking via cAdvisor
- 🐳 Full containerization with Docker Compose
- 🧱 Clean Architecture with strict layer separation
- 📖 API documentation via Swagger UI

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | Gin |
| Database | PostgreSQL 16 |
| Cache / Leaderboard | Redis 7 (Sorted Sets) |
| Authentication | JWT (access + refresh) |
| API Docs | Swagger (swaggo/swag) |
| Metrics | Prometheus |
| Dashboards | Grafana |
| Container Monitoring | cAdvisor |
| Configuration | Viper (config.yml + .env) |
| Logging | slog (structured logs) |
| Containerization | Docker & Docker Compose |

---

## Architecture

The project follows **Clean Architecture** principles with strict separation of concerns:

```
┌──────────────────────────────────────────────────────────┐
│  HTTP Layer  (Gin handlers, DTO, Swagger, Middleware)    │
│  internal/interfaces/http/handler/                       │
│  → string IDs in DTOs, UUID conversion at layer boundary │
├──────────────────────────────────────────────────────────┤
│  Use Case Layer  (Business Logic)                        │
│  internal/usecase/                                       │
│  → Works exclusively with uuid.UUID                     │
├──────────────────────────────────────────────────────────┤
│  Domain Layer  (Entities, Interfaces)                    │
│  internal/domain/                                        │
│  → Pure domain models with no external dependencies      │
├──────────────────────────────────────────────────────────┤
│  Infrastructure Layer  (PostgreSQL, Redis, JWT, Logger)  │
│  internal/infrastructure/                                │
│  → Repository implementations, external adapters         │
└──────────────────────────────────────────────────────────┘
```

### ID Type Consistency Rule

- **Domain / Service / Repository**: use `uuid.UUID`
- **HTTP layer (DTOs, request params)**: use `string`
- **Conversion**: happens **only** in handlers via `uuid.Parse`

---

## Observability & Monitoring

```
┌─────────────────────────────────────────────────────────────┐
│                      Docker Network                         │
│                                                             │
│  ┌──────────┐  /metrics   ┌────────────┐   PromQL          │
│  │ Backend  │ ──────────► │ Prometheus │ ◄────────────────┐ │
│  │  :8080   │             │   :9090    │                  │ │
│  └──────────┘             └────────────┘                  │ │
│                                    │                       │ │
│  ┌──────────┐  /metrics            ▼                       │ │
│  │ cAdvisor │ ──────────► ┌────────────┐  datasource:     │ │
│  │  :8081   │             │  Grafana   │ ◄── prometheus:9090│
│  └──────────┘             │   :3000    │                   │ │
│                           └────────────┘                   │ │
└─────────────────────────────────────────────────────────────┘
```

### Application Metrics (Prometheus)

The backend exposes metrics at `backend:8080/metrics`.

| Metric | Type | Description |
|---|---|---|
| `http_request_duration_seconds` | Histogram | Request processing time per endpoint and method |
| `http_requests_total` | Counter | Total request count (method, path, status) |
| `jwt_auth_failures_total` | Counter | JWT authentication failure count |
| `db_query_duration_seconds` | Histogram | PostgreSQL query execution time |
| `redis_operation_duration_seconds` | Histogram | Redis operation latency |

### Latency Monitoring (P95)

```promql
# P95 latency for all HTTP endpoints
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path)
)

# P95 PostgreSQL query latency
histogram_quantile(0.95,
  sum(rate(db_query_duration_seconds_bucket[5m])) by (le)
)

# P95 Redis operation latency
histogram_quantile(0.95,
  sum(rate(redis_operation_duration_seconds_bucket[5m])) by (le)
)
```

### JWT Authentication Monitoring

```promql
# JWT failure rate over last 5 minutes
rate(jwt_auth_failures_total[5m])

# Total JWT failure count
sum(jwt_auth_failures_total)
```

### Endpoint Monitoring

```promql
# RPS (requests per second) per endpoint
sum(rate(http_requests_total[1m])) by (path, method)

# 5xx error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) by (path)

# 4xx error rate
sum(rate(http_requests_total{status=~"4.."}[5m])) by (path)
```

### Container Monitoring (cAdvisor)

cAdvisor collects Docker container metrics and exposes them to Prometheus.

```promql
# CPU utilization per container
sum(rate(container_cpu_usage_seconds_total{name!=""}[1m])) by (name)

# Memory usage per container
sum(container_memory_usage_bytes{name!=""}) by (name)

# Inbound network traffic per container
sum(rate(container_network_receive_bytes_total{name!=""}[1m])) by (name)
```

### Grafana Dashboards

Grafana is available at: **http://localhost:3000**

- **Default credentials**: `admin` / `admin`
- **Data source**: Prometheus (`http://prometheus:9090`)

Recommended dashboards to import:
- **Gin HTTP Metrics** — RPS, latency, response codes
- **Go Runtime** — goroutines, GC, memory (ID: `13240`)
- **PostgreSQL** — query time, connections
- **Redis** — operations, hit rate
- **cAdvisor Docker** — CPU, RAM per container (ID: `14282`)

---

## Docker Compose Services

| Service | Port | Description |
|---|---|---|
| `backend` | `8080` | Go application + `/metrics` endpoint |
| `postgres` | `5432` | PostgreSQL 16 |
| `redis` | `6379` | Redis 7 |
| `prometheus` | `9090` | Prometheus (scrapes `backend:8080/metrics`, `cadvisor:8081`) |
| `grafana` | `3000` | Grafana (datasource: `http://prometheus:9090`) |
| `cadvisor` | `8081` | Container CPU/RAM monitoring |

---

## API Documentation

Swagger UI: **http://localhost:8080/swagger/index.html**

### Public Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/auth/register` | Register a new user |
| `POST` | `/auth/login` | Login and obtain tokens |
| `POST` | `/admin/create` | Create a new game |
| `GET` | `/admin/games` | List all games |

### Protected Endpoints (require `Authorization: Bearer <token>`)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/score/submit` | Submit player score |
| `GET` | `/api/leaderboard/global` | Get global leaderboard |
| `GET` | `/api/leaderboard/my` | Get current user rank |
| `POST` | `/api/leaderboard/top` | Get top players for a game |

---

## Environment Variables

Create a `.env` file in the project root:

```env
JWT_ACCESS_SECRET=your-access-secret-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-change-in-production
DB_PASSWORD=postgres
```

### config.yml

```yaml
port: "8080"
db:
  username: "postgres"
  host: "postgres"       # Use "localhost" for local development
  port: 5432
  dbname: "leaderboard"
  sslmode: "disable"
```

---

## Running the Application

### Option 1: Docker Compose (recommended)

```bash
# 1. Create .env file
# 2. Start all services
docker-compose up -d

# 3. Check service status
docker-compose ps
```

**Available services:**
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- cAdvisor: http://localhost:8081

```bash
# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Option 2: Local Development

```bash
# 1. Start PostgreSQL and Redis
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=leaderboard postgres:16
docker run -d -p 6379:6379 redis:7

# 2. Install dependencies
go mod download

# 3. Update config.yml: set db.host = "localhost"

# 4. Set environment variables (Windows)
set JWT_ACCESS_SECRET=your-access-secret
set JWT_REFRESH_SECRET=your-refresh-secret
set DB_PASSWORD=postgres

# 5. Regenerate Swagger docs (if needed)
swag init -g cmd/app/main.go -o ./docs

# 6. Run the application
go run cmd/app/main.go
```

---

## Project Structure

```
OnlineLeadership/
├── cmd/app/main.go                    # Entry point
├── config/config.go                   # Configuration (Viper)
├── config.yml                         # Application settings
├── internal/
│   ├── domain/                        # Domain models and interfaces
│   ├── usecase/                       # Business logic services
│   │   ├── auth/
│   │   ├── admin/
│   │   ├── leaderboard/
│   │   └── score_history/
│   ├── infrastructure/
│   │   ├── auth/                      # JWT token manager
│   │   ├── logger/                    # slog structured logger
│   │   ├── monitoring/                # Prometheus metric definitions
│   │   ├── postgres/                  # DB connection + repositories
│   │   └── redis/                     # Redis client
│   └── interfaces/http/
│       ├── handler/                   # Gin handlers, DTOs, Swagger
│       └── middleware/                # Request ID, Auth, Metrics
├── migrations/                        # SQL migration files
├── monitoring/
│   └── prometheus.yml                 # Prometheus scrape configuration
├── docs/                              # Auto-generated Swagger docs
├── docker-compose.yml
└── Dockerfile
```

---

## Database Schema

**`users`** — `id (UUID)`, `username`, `email`, `password_hash`, `created_at`

**`games`** — `id (UUID)`, `name`

**`score_history`** — `id (UUID)`, `user_id (FK)`, `game_id (FK)`, `score`, `created_at`

**Redis Sorted Sets:**
- `leaderboard:global` — global ranking (member: user_id, score: total points)
- `leaderboard:game:{game_id}` — per-game ranking

---

## Common Errors & Troubleshooting

| Error | Cause | Solution |
|---|---|---|
| `JWT secrets are not set` | Missing `.env` file | Create `.env` with JWT secrets |
| `db connect failed` | PostgreSQL unavailable | Check host in `config.yml` |
| `401 Unauthorized` | Missing/expired token | Re-login, add `Bearer ` prefix |
| `400 invalid UUID` | Malformed ID format | Use valid UUID format |
| Swagger auth fails | Missing `Bearer ` prefix | In Swagger UI: enter `Bearer <token>` |
| `redis connection error` | Redis not running | `docker run -d -p 6379:6379 redis:7` |
| Prometheus no metrics | Wrong scrape target | Check `monitoring/prometheus.yml` |

---

## Production Recommendations

- **JWT Secrets**: Use strong random secrets (min. 64 chars), never commit to VCS
- **HTTPS**: Use a reverse proxy (Nginx/Traefik) with TLS in production
- **Rate Limiting**: Add rate limiting middleware for public endpoints
- **CORS**: Configure CORS if frontend is served from a different domain
- **Admin Routes**: Protect `/admin/*` endpoints with authentication
- **Grafana**: Change default admin password on first login
- **Prometheus**: Restrict `/metrics` endpoint access in production (e.g., internal network only)
- **Alerting**: Set up Prometheus Alertmanager for JWT failure and latency spike alerts

---

## License

MIT

## Support

For questions and issues, please open an issue in the project repository.

