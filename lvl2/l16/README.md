# Wget - Утилита для загрузки веб-страниц

Упрощенная реализация утилиты `wget` для загрузки веб-страниц вместе со всем вложенным контентом (ресурсы, ссылки), аналогичная `wget -m` (мирроринг сайта).

## Возможности

- ✅ Загрузка HTML-страниц с сохранением локально
- ✅ Рекурсивное скачивание ресурсов: CSS, JS, изображения
- ✅ Обработка относительных и абсолютных ссылок
- ✅ Предотвращение дублирования (не скачивает один ресурс несколько раз)
- ✅ Корректное формирование локальных путей для сохранения
- ✅ Избежание зацикливания по ссылкам
- ✅ Параллельное скачивание с ограничением количества одновременных загрузок
- ✅ Управление robots.txt (базовая поддержка)
- ✅ Обработка ошибок (сетевых, файловых)
- ✅ Таймауты на запросы
- ✅ Детальная статистика загрузки

### Структура проекта:
```
lvl2/l16/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── app/
│   │   └── wget.go          # Основная логика приложения
│   ├── config/
│   │   ├── config.go        # Конфигурация
│   │   └── errors.go        # Ошибки конфигурации
│   ├── downloader/
│   │   ├── downloader.go    # HTTP загрузчик
│   │   └── downloader_test.go
│   ├── filter/
│   │   └── filter.go        # Фильтр URL
│   ├── parser/
│   │   ├── parser.go        # Парсер HTML/CSS
│   │   └── parser_test.go
│   ├── reporter/
│   │   └── reporter.go      # Репортер прогресса
│   ├── storage/
│   │   ├── storage.go       # Файловое хранилище
│   │   └── storage_test.go
│   └── types/
│       └── types.go         # Интерфейсы и типы
├── go.mod
├── go.sum
├── demo.sh
└── README.md
```

## Установка и запуск

### Требования
- Go 1.21+

### Установка зависимостей
```bash
go mod tidy
```

### Сборка
```bash
go build -o wget cmd/main.go
```

### Использование
```bash
# Базовое использование
./wget -url https://example.com

# С настройками
./wget -url https://example.com -depth 2 -output ./downloaded -concurrency 10

# Полный список параметров
./wget -help
```

## Параметры командной строки

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `-url` | URL для загрузки (обязательный) | - |
| `-depth` | Глубина рекурсии | 3 |
| `-output` | Директория для сохранения файлов | ./downloaded |
| `-concurrency` | Количество одновременных загрузок | 5 |
| `-timeout` | Таймаут для HTTP запросов | 30s |
| `-robots` | Соблюдать robots.txt | true |
| `-user-agent` | User-Agent для запросов | Wget/1.0 |

## Примеры использования

### Загрузка простой страницы
```bash
./wget -url https://example.com
```

### Загрузка с ограниченной глубиной
```bash
./wget -url https://example.com -depth 1
```

### Параллельная загрузка
```bash
./wget -url https://example.com -concurrency 20
```

### Кастомная директория
```bash
./wget -url https://example.com -output /path/to/downloads
```

## Интерфейсы и компоненты

### Downloader
Интерфейс для загрузки ресурсов:
```go
type Downloader interface {
    Download(ctx context.Context, resourceURL *url.URL) (*Resource, error)
}
```

### Parser
Интерфейс для парсинга HTML и извлечения ссылок:
```go
type Parser interface {
    ParseHTML(content []byte, baseURL *url.URL) ([]*url.URL, error)
    ParseCSS(content []byte, baseURL *url.URL) ([]*url.URL, error)
}
```

### Storage
Интерфейс для сохранения ресурсов:
```go
type Storage interface {
    Save(resource *Resource) error
    Exists(localPath string) bool
    GetLocalPath(resourceURL *url.URL) string
}
```

### URLFilter
Интерфейс для фильтрации URL:
```go
type URLFilter interface {
    ShouldDownload(url *url.URL, depth int) bool
    MarkVisited(url *url.URL)
}
```

### ProgressReporter
Интерфейс для отчета о прогрессе:
```go
type ProgressReporter interface {
    ReportDownloaded(resource *Resource)
    ReportSkipped(url *url.URL, reason string)
    ReportError(url *url.URL, err error)
    PrintSummary()
}
```

## Тестирование

Запуск тестов:
```bash
go test ./...
```

Запуск тестов с покрытием:
```bash
go test -cover ./...
```

## Особенности реализации

### Обработка ссылок
- Автоматическое разрешение относительных ссылок
- Фильтрация javascript:, data:, mailto: ссылок
- Предотвращение зацикливания через отслеживание посещенных URL

### Параллельная обработка
- Использование горутин для параллельной загрузки
- Ограничение количества одновременных запросов
- Потокобезопасные структуры данных

### Обработка ошибок
- Graceful handling сетевых ошибок
- Логирование всех ошибок с контекстом
- Продолжение работы при ошибках отдельных ресурсов


## Ограничения

- Базовая поддержка robots.txt (без полного парсинга)
- Ограниченная поддержка JavaScript-генерируемых ссылок
- Нет поддержки аутентификации
- Нет поддержки cookies и сессий
