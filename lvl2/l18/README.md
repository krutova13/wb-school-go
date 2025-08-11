# HTTP-сервер "Календарь"

HTTP-сервер для работы с небольшим календарем событий, реализованный на Go.

## Функциональность

### CRUD операции для событий

- **POST /create_event** — создание нового события
- **POST /update_event** — обновление существующего события
- **POST /delete_event** — удаление события
- **GET /events_for_day** — получение всех событий на день
- **GET /events_for_week** — получение событий на неделю
- **GET /events_for_month** — получение событий на месяц

## Формат запросов

### POST запросы

Данные могут передаваться в двух форматах:

1. **JSON** (Content-Type: application/json):
```json
{
  "user_id": "user123",
  "date": "2023-12-31",
  "text": "Новый год"
}
```

2. **Form-data** (application/x-www-form-urlencoded):
```
user_id=user123&date=2023-12-31&text=Новый год
```

### GET запросы

Параметры передаются через query string:
```
GET /events_for_day?user_id=user123&date=2023-12-31
```

## Формат ответов

### Успешный ответ
```json
{
  "result": "данные или сообщение"
}
```

### Ответ с ошибкой
```json
{
  "error": "описание ошибки"
}
```

## HTTP статус-коды

- **200 OK** — для успешных запросов
- **400 Bad Request** — для ошибок ввода (некорректные данные)
- **503 Service Unavailable** — для ошибок бизнес-логики
- **500 Internal Server Error** — для прочих ошибок

## Использование

### Требования

- Go 1.21 или выше
- curl (для тестирования)
- jq (для демонстрационного скрипта)

### Запуск

1. Клонируйте репозиторий и перейдите в папку проекта:
```bash
cd lvl2/l18
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Запустите сервер:
```bash
go run cmd/main.go
```

## Конфигурация

По умолчанию сервер запускается на порту 8080. Для изменения порта используйте переменную окружения:

```bash
export CALENDAR_PORT=3000
go run cmd/main.go
```


## Демонстрация

После запуска сервера выполните демонстрационный скрипт:
```bash
./demo.sh
```

**Примечание:** Для работы скрипта требуется установленный `jq` для форматирования JSON.



## Ручное тестирование

Запустите unit-тесты:
```bash
go test ./...
```

Запустите тесты с покрытием:
```bash
go test -cover ./...
```

## Примеры использования

### Создание события
```bash
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "date": "2023-12-31", "text": "Новый год"}'
```

### Получение событий на день
```bash
curl "http://localhost:8080/events_for_day?user_id=user123&date=2023-12-31"
```

### Обновление события
```bash
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{"id": "event_id", "user_id": "user123", "date": "2024-01-01", "text": "Обновленный текст"}'
```

### Удаление события
```bash
curl -X POST http://localhost:8080/delete_event \
  -H "Content-Type: application/json" \
  -d '{"id": "event_id", "user_id": "user123"}'
```

## Структура проекта

```
lvl2/l18/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── calendar/            # Бизнес-логика календаря
│   │   ├── calendar.go
│   │   └── calendar_test.go
│   ├── config/              # Конфигурация
│   │   └── config.go
│   ├── handlers/            # HTTP-обработчики
│   │   └── handlers.go
│   ├── middleware/          # Middleware
│   │   └── logger.go
│   └── types/               # Типы данных
│       └── types.go
├── go.mod
└── README.md
```

## Логирование

Сервер автоматически логирует все HTTP-запросы с информацией о:
- Методе запроса
- URL
- Времени выполнения

Логи выводятся в stdout.

Пример лога:
```
2025/08/11 15:05:41 Запрос: POST /create_event
2025/08/11 15:05:41 Запрос завершен: POST /create_event - 187.042µs
``` 