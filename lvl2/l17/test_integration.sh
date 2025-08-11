#!/bin/bash

# Интеграционный тест telnet-клиента

set -e

echo "=== Telnet Client Integration Test ==="
echo

# Сборка тестового сервера
echo "1. Сборка тестового сервера..."
go build -o test_server test_server.go
echo "✓ Тестовый сервер собран"
echo

# Запуск тестового сервера в фоне
echo "2. Запуск тестового сервера..."
./test_server &
SERVER_PID=$!
sleep 2
echo "✓ Тестовый сервер запущен (PID: $SERVER_PID)"
echo

# Тестирование подключения
echo "3. Тестирование подключения к серверу..."
echo "Отправка тестовых сообщений..."

# Создаем временный файл с тестовыми сообщениями
cat > /tmp/test_messages.txt << EOF
Hello, server!
This is a test message
Another line
EOF

# Запускаем telnet клиент с тестовыми сообщениями
echo "Подключение к localhost:8080..."
./telnet localhost 8080 < /tmp/test_messages.txt &
TELNET_PID=$!

# Ждем немного для обработки
sleep 3

# Завершаем процессы
echo "4. Завершение процессов..."
kill $TELNET_PID 2>/dev/null || true
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
wait $TELNET_PID 2>/dev/null || true

# Очистка
rm -f /tmp/test_messages.txt
echo "✓ Процессы завершены"
echo

echo "=== Интеграционный тест завершен ==="
echo
echo "Для интерактивного тестирования:"
echo "1. Запустите сервер: ./test_server"
echo "2. В другом терминале: ./telnet localhost 8080"
echo "3. Введите сообщения и нажмите Ctrl+D для завершения" 