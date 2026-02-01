# OnlineLeadership API

RESTful API сервис для управления игровыми таблицами лидеров с рейтингом в реальном времени, отслеживанием очков и аутентификацией пользователей.

## Технологический стек

- **Go 1.25** - Язык программирования
- **Gin** - HTTP веб-фреймворк
- **PostgreSQL 16** - Основная база данных (данные пользователей, история очков)
- **Redis 7** - In-memory хранилище (рейтинги лидерборда, sorted sets)
- **JWT** - Аутентификация (access и refresh токены)
- **Swagger/OpenAPI** - Документация API (`swaggo/swag`)
- **Docker & Docker Compose** - Контейнеризация
- **sqlx** - SQL query builder
- **Viper** - Управление конфигурацией
- **slog** - Структурированное логирование

## Архитектура

Проект следует принципам **Clean Architecture** с чётким разделением слоёв:

```
┌─────────────────────────────────────────┐
│  HTTP Layer (Gin handlers)              │  ← DTO используют string для ID
├─────────────────────────────────────────┤
│  Use Case Layer (Бизнес-логика)         │  ← Использует uuid.UUID
├─────────────────────────────────────────┤
│  Domain Layer (Сущности)                │  ← Чистые доменные модели
├─────────────────────────────────────────┤
│  Repository Layer (Доступ к данным)     │  ← PostgreSQL + Redis
└─────────────────────────────────────────┘
```

### Слои

- **`internal/interfaces/http/handler`** - HTTP обработчики, request/response DTO, Swagger аннотации
- **`internal/usecase`** - Сервисы бизнес-логики (auth, admin, leaderboard, score)
- **`internal/domain`** - Доменные модели (User, Game, LeaderboardUser)
- **`internal/infrastructure`** - Внешние зависимости (PostgreSQL, Redis, JWT, logger)

### Согласованность типов ID

- **Domain/Service/Repository слои**: Используют `uuid.UUID`
- **HTTP слой (DTO, запросы)**: Используют `string`
- **Конвертация**: Происходит только на границе HTTP (в handlers)

## Функциональность

### Аутентификация
- Регистрация пользователей с email и паролем
- JWT-аутентификация (access + refresh токены)
- Время жизни access токена: 30 минут
- Время жизни refresh токена: 7 дней
- Хеширование паролей с bcrypt

### Управление играми
- Создание новых игр
- Получение списка всех доступных игр

### Отслеживание очков
- Отправка очков игрока для конкретных игр
- Постоянное хранение истории очков (PostgreSQL)
- Автоматическое обновление лидерборда (Redis sorted sets)

### Таблицы лидеров
- Глобальный лидерборд (все игроки по всем играм)
- Лидерборды по конкретным играм
- Получение ранга пользователя
- Поддержка пагинации (offset/limit)

## Документация API

Swagger UI доступен по адресу: **http://localhost:8080/swagger/index.html**

### Аутентификация

Защищённые endpoint'ы требуют JWT токен в заголовке `Authorization`:

```
Authorization: Bearer <access_token>
```

### Endpoints

#### Публичные endpoints
- `POST /auth/register` - Регистрация нового пользователя
- `POST /auth/login` - Вход и получение токенов
- `POST /admin/create` - Создание новой игры
- `GET /admin/games` - Список всех игр

#### Защищённые endpoints (требуют JWT)
- `POST /api/score/submit` - Отправка очков игрока
- `GET /api/leaderboard/global` - Получение глобального лидерборда
- `GET /api/leaderboard/my` - Получение ранга текущего пользователя
- `POST /api/leaderboard/top` - Получение топ игроков для конкретной игры

## Переменные окружения

Создайте файл `.env` в корне проекта:

```bash
# JWT секреты (обязательно)
JWT_ACCESS_SECRET=ваш-секретный-ключ-доступа
JWT_REFRESH_SECRET=ваш-секретный-ключ-обновления

# База данных (опционально, значения по умолчанию в config.yml)
DB_PASSWORD=postgres
```

### Конфигурационные файлы

**`config.yml`** - Конфигурация приложения:
```yaml
port: "8080"
db:
  username: "postgres"
  host: "postgres"      # Используйте "localhost" для локальной разработки
  port: 5432
  dbname: "leaderboard"
  sslmode: "disable"
```

## Структура проекта

```
OnlineLeadership/
├── cmd/
│   └── app/
│       └── main.go              # Точка входа в приложение
├── internal/
│   ├── domain/                  # Доменные модели (User, Game, LeaderboardUser)
│   ├── usecase/                 # Сервисы бизнес-логики
│   │   ├── auth/
│   │   ├── admin/
│   │   ├── leaderboard/
│   │   └── score_history/
│   ├── infrastructure/          # Внешние зависимости
│   │   ├── auth/                # JWT менеджер токенов
│   │   ├── logger/              # Структурированное логирование
│   │   ├── postgres/            # Подключение к БД и репозитории
│   │   └── redis/               # Redis клиент
│   └── interfaces/
│       └── http/
│           ├── handler/         # HTTP обработчики и DTO
│           └── middleware/      # Request ID, аутентификация
├── migrations/                  # SQL миграции
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── docs/                        # Автогенерированная Swagger документация
├── config.yml                   # Конфигурация приложения
├── docker-compose.yml           # Docker оркестрация
├── Dockerfile                   # Образ контейнера
└── go.mod                       # Go зависимости
```

## Запуск приложения

### Требования
- Go 1.25+
- PostgreSQL 16
- Redis 7
- Docker & Docker Compose (опционально)

### Вариант 1: Использование Docker Compose (Рекомендуется)

1. **Создайте файл `.env`** с JWT секретами:
```bash
JWT_ACCESS_SECRET=ваш-секретный-ключ-измените-в-продакшене
JWT_REFRESH_SECRET=ваш-ключ-обновления-измените-в-продакшене
DB_PASSWORD=postgres
```

2. **Запустите все сервисы**:
```bash
docker-compose up -d
```

Это запустит:
- Сервер приложения на `http://localhost:8080`
- PostgreSQL на `localhost:5432`
- Redis на `localhost:6379`
- Автоматические миграции базы данных

3. **Откройте Swagger UI**:
```
http://localhost:8080/swagger/index.html
```

4. **Остановите сервисы**:
```bash
docker-compose down
```

### Вариант 2: Локальная разработка

1. **Установите зависимости**:
```bash
go mod download
```

2. **Запустите PostgreSQL и Redis** (вручную или через Docker):
```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=leaderboard postgres:16
docker run -d -p 6379:6379 redis:7
```

3. **Выполните миграции**:
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/leaderboard?sslmode=disable" up
```

4. **Обновите `config.yml`** - измените `host` на `localhost`:
```yaml
db:
  host: "localhost"
```

5. **Установите переменные окружения**:
```bash
export JWT_ACCESS_SECRET="ваш-секрет-доступа"
export JWT_REFRESH_SECRET="ваш-секрет-обновления"
export DB_PASSWORD="postgres"
```

6. **Сгенерируйте Swagger документацию** (если отсутствует):
```bash
swag init -g cmd/app/main.go -o ./docs
```

7. **Запустите приложение**:
```bash
go run cmd/app/main.go
```

Сервер запустится на `http://localhost:8080`

## Примеры использования

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "password123"
  }'
```

Ответ:
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### 2. Вход в систему
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "password": "password123"
  }'
```

Ответ:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Создание игры
```bash
curl -X POST http://localhost:8080/admin/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Шахматы"
  }'
```

Ответ:
```json
{
  "game_id": "987fcdeb-51a2-43f7-9876-543210fedcba"
}
```

### 4. Отправка очков (Защищённый endpoint)
```bash
curl -X POST http://localhost:8080/api/score/submit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ВАШ_ACCESS_TOKEN" \
  -d '{
    "game_id": "987fcdeb-51a2-43f7-9876-543210fedcba",
    "score": 1500
  }'
```

Ответ:
```json
{
  "status": "ok"
}
```

### 5. Получение глобального лидерборда (Защищённый endpoint)
```bash
curl -X GET "http://localhost:8080/api/leaderboard/global?offset=0&limit=10" \
  -H "Authorization: Bearer ВАШ_ACCESS_TOKEN"
```

Ответ:
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

### 6. Получение своего ранга (Защищённый endpoint)
```bash
curl -X GET http://localhost:8080/api/leaderboard/my \
  -H "Authorization: Bearer ВАШ_ACCESS_TOKEN"
```

Ответ:
```json
{
  "rank": 1
}
```

## Распространённые ошибки и решения

### 1. "JWT secrets are not set"
**Решение**: Убедитесь, что файл `.env` содержит `JWT_ACCESS_SECRET` и `JWT_REFRESH_SECRET`

### 2. "db connect failed"
**Причины**:
- PostgreSQL не запущен
- Неверные учётные данные
- Неверный host (используйте `localhost` для локальной разработки, `postgres` для Docker)

**Решение**: Проверьте подключение к PostgreSQL и обновите `config.yml` или переменные окружения

### 3. "401 Unauthorized" на защищённых endpoints
**Причины**:
- Отсутствует заголовок `Authorization`
- Недействительный или просроченный токен
- Неверный формат токена

**Решение**: 
- Добавьте заголовок: `Authorization: Bearer <access_token>`
- Войдите заново для получения свежего токена
- Убедитесь, что префикс "Bearer " включён

### 4. "400 Bad Request - invalid game_id format"
**Причина**: Невалидный формат UUID в теле запроса

**Решение**: Используйте валидный формат UUID (например, `123e4567-e89b-12d3-a456-426614174000`)

### 5. Аутентификация в Swagger не работает
**Решение**: 
1. Нажмите кнопку "Authorize" в Swagger UI
2. Введите: `Bearer <ваш_access_token>` (включая префикс "Bearer ")
3. Нажмите "Authorize", затем "Close"

### 6. "redis connection error"
**Причина**: Redis сервер не запущен

**Решение**: 
```bash
docker run -d -p 6379:6379 redis:7
```

## Схема базы данных

### Таблицы

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

### Структуры данных Redis

- **Глобальный лидерборд**: Sorted set `leaderboard:global`
  - Элементы: ID пользователей
  - Очки: общее количество баллов

- **Лидерборды игр**: Sorted set `leaderboard:game:{game_id}`
  - Элементы: ID пользователей
  - Очки: баллы по конкретной игре

## Разработка

### Регенерация Swagger документации
```bash
swag init -g cmd/app/main.go -o ./docs
```

### Запуск тестов
```bash
go test ./...
```

### Сборка бинарного файла
```bash
go build -o bin/app cmd/app/main.go
```

## Рекомендации по безопасности

- **Production**: Измените JWT секреты на сильные случайные значения
- **HTTPS**: Используйте HTTPS в продакшене (настройте reverse proxy)
- **Rate Limiting**: Реализуйте ограничение частоты запросов для публичных endpoints
- **CORS**: Настройте CORS, если фронтенд обслуживается с другого домена
- **Admin Endpoints**: Добавьте аутентификацию/авторизацию для маршрутов `/admin/*`

## Лицензия

MIT

## Поддержка

По вопросам и проблемам создавайте issue в репозитории проекта.

