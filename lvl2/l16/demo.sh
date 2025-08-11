#!/bin/bash

# Демо скрипт для тестирования утилиты wget

set -e

echo "=== Демонстрация утилиты Wget ==="
echo

if ! command -v go &> /dev/null; then
    echo "❌ Ошибка: Go не установлен"
    exit 1
fi

cd "$(dirname "$0")"

echo "📦 Устанавливаем зависимости..."
go mod tidy

echo "🔨 Собираем приложение..."
go build -o wget cmd/main.go

echo "✅ Приложение собрано успешно!"
echo

echo "🌐 Запускаем тестовый HTTP сервер..."
python3 -m http.server 8080 &
SERVER_PID=$!

# Ждем запуска сервера
sleep 2

echo "📥 Тестируем загрузку локальной страницы..."
echo "URL: http://localhost:8080"
echo "Глубина: 1"
echo "Директория: ./test-download"
echo

./wget -url http://localhost:8080 -depth 1 -output ./test-download -concurrency 3

echo
echo "📁 Проверяем загруженные файлы..."
if [ -d "./test-download" ]; then
    echo "Содержимое директории ./test-download:"
    find ./test-download -type f | head -10
    echo
    echo "Всего файлов: $(find ./test-download -type f | wc -l)"
else
    echo "❌ Директория ./test-download не найдена"
fi

echo
echo "🧹 Очистка..."

kill $SERVER_PID 2>/dev/null || true

rm -rf ./test-download
rm -f ./wget

echo "✅ Демонстрация завершена!"
echo
echo "Для использования с реальными сайтами:"
echo "  ./wget -url https://example.com -depth 2 -output ./downloaded"
echo
echo "Для просмотра всех параметров:"
echo "  ./wget -help" 