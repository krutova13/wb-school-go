#!/bin/bash

set -e

echo "=== Telnet Client Demo ==="
echo

echo "1. Сборка проекта..."
go build -o telnet cmd/main.go
echo "✓ Проект собран"
echo

echo "2. Тестирование парсера аргументов..."
echo "Тест с валидными аргументами:"
./telnet localhost 8080 2>&1 | head -5 || true
echo

echo "Тест с невалидными аргументами:"
./telnet 2>&1 | head -3 || true
echo

echo "Тест с таймаутом:"
./telnet --timeout=5s localhost 8080 2>&1 | head -5 || true
echo

echo "3. Запуск тестов..."
go test ./... -v
echo "✓ Тесты пройдены"
echo

echo "4. Демонстрация подключения к echo серверу..."
echo "Подключение к echo.websocket.org:80..."
echo "Отправка HTTP GET запроса..."
echo

cat > /tmp/http_request.txt << EOF
GET / HTTP/1.1
Host: echo.websocket.org
Connection: close

EOF

timeout 10s ./telnet echo.websocket.org 80 < /tmp/http_request.txt || true

echo
echo "✓ Демонстрация завершена"
echo

rm -f /tmp/http_request.txt
echo "5. Очистка временных файлов..."
echo "✓ Очистка завершена"
echo

echo "=== Демо завершено ==="
echo
echo "Для тестирования вручную:"
echo "1. Запустите netcat сервер: nc -l 8080"
echo "2. В другом терминале: ./telnet localhost 8080"
echo "3. Введите сообщения и нажмите Ctrl+D для завершения" 