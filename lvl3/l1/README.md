# Delayed Notifier

Простая система отложенных уведомлений с поддержкой Telegram и Email каналов, построенная на Go

## 🚀 Быстрый старт

### 1. Запуск сервисов
```bash
docker-compose up -d
```

### 2. Запуск приложения
```bash
go run ./cmd
```

### 3. Открыть веб-интерфейс
```
http://localhost:8080/web/
```

## 📁 Структура проекта

```
lvl3/
├── cmd/                    # Точка входа
│   └── main.go
├── internal/              # Внутренние пакеты
│   ├── app/              # Основное приложение и зависимости
│   ├── cache/            # Redis кэширование
│   ├── config/           # Конфигурация приложения
│   ├── domain/           # Доменные модели
│   ├── dto/              # Data Transfer Objects
│   ├── handlers/         # HTTP обработчики
│   ├── queue/            # Очереди сообщений (RabbitMQ)
│   ├── repository/       # Репозитории (PostgreSQL)
│   ├── sender/           # Отправители (Telegram, Email)
│   ├── service/          # Бизнес-логика и воркеры
│   └── validation/       # Валидация запросов
├── web/                  # Веб-интерфейс
├── migrations/           # Миграции БД
└── config.yaml          # Конфигурация
```

### Основные компоненты:

- **HTTP API** - RESTful API
- **RabbitMQ** - очереди для отложенных уведомлений с retry механизмом
- **Redis** - кэширование статусов уведомлений
- **PostgreSQL** - хранение данных
- **Workers** - фоновые воркеры для обработки очереди
- **Multiple Channels** - поддержка Email и Telegram

## 🔧 Конфигурация

Все настройки в файле `config.yaml`, чувствительные данные в .env.

Пример .env конфигурации представлен в .env.example файле

## 📧 Поддерживаемые каналы

### Telegram
- Простая настройка через Bot Token и Chat ID
- Поддержка HTML разметки
- Автоматические retry при ошибках

### Email
- Полная конфигурация SMTP
- Поддержка Gmail, Outlook, Yahoo и других провайдеров
- HTML и текстовые версии писем
- Кастомные темы и отправители

## 🐳 Docker

### Сервисы:
- **PostgreSQL** - база данных (порт 5432)
- **Redis** - кэширование (порт 6379)
- **RabbitMQ** - очереди сообщений (порты 5672, 15672)
- **Migrate** - автоматические миграции БД

### Запуск:
```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

## 📊 API Endpoints

### Создание уведомления
```bash
POST /api/v1/notify
Content-Type: application/json

{
  "payload": "Hello World",
  "notification_date": "2024-12-31T23:59:59Z",
  "recipient_id": "user@example.com",
  "channel": "email",
  "email_config": {
    "subject": "Notification",
    "from_name": "Service",
    "from_email": "sender@example.com",
    "smtp_host": "smtp.gmail.com",
    "smtp_port": 587,
    "username": "sender@gmail.com",
    "password": "app_password"
  }
}
```

**Ответ:**
```json
{
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "payload": "Hello World",
    "notification_date": "2024-12-31T23:59:59Z",
    "recipient_id": "user@example.com",
    "channel": "email"
  }
}
```

### Получение статуса
```bash
GET /api/v1/notify/{id}
```

**Ответ:**
```json
{
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "sent",
    "payload": "Hello World",
    "notification_date": "2024-12-31T23:59:59Z",
    "recipient_id": "user@example.com",
    "channel": "email",
    "retries": 0
  }
}
```

### Отмена уведомления
```bash
DELETE /api/v1/notify/{id}
```

**Ответ:**
```json
{
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "cancelled"
  }
}
```

## 🔄 Статусы уведомлений

- **pending** - ожидает отправки
- **sent** - успешно отправлено
- **failed** - ошибка отправки
- **cancelled** - отменено пользователем

## 🚀 Особенности

### Retry механизм
- Экспоненциальная задержка при повторных попытках
- Настраиваемое количество попыток
- Автоматическое логирование ошибок

### Кэширование
- Redis для быстрого доступа к статусам
- TTL для автоматической очистки
- Fallback на базу данных при недоступности кэша

### Graceful Shutdown
- Корректное завершение воркеров
- Закрытие соединений с БД и Redis
- Обработка сигналов системы

### Логирование
- Структурированные логи (JSON/Console)
- Настраиваемые уровни логирования
- Контекстная информация в логах

## 🔍 Отладка

### Логи
```bash
# Debug уровень
LOGGING_LEVEL=debug go run ./cmd

# Console формат (по умолчанию)
LOGGING_FORMAT=console go run ./cmd
```

### Переменные окружения
Все настройки можно переопределить через переменные окружения:
```bash
HTTP_SERVER_PORT=8080
POSTGRES_HOST=localhost
REDIS_URL=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

## 🧪 Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск конкретного теста
go test ./internal/service/
```

## 📝 Разработка

### Требования
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 13+
- Redis 6+
- RabbitMQ 3.8+

### Установка зависимостей
```bash
go mod download
```

### Сборка
```bash
go build ./cmd/main.go
```