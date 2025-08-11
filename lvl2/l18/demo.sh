#!/bin/bash

echo "🚀 Запуск демонстрации HTTP-сервера 'Календарь'"
echo "================================================"

echo "📡 Проверка доступности сервера..."
if ! curl -s http://localhost:8080/events_for_day > /dev/null; then
    echo "❌ Сервер не доступен на порту 8080"
    echo "Запустите сервер командой: go run cmd/main.go"
    exit 1
fi

echo "✅ Сервер доступен"
echo ""

echo "📅 Создание событий..."
echo ""

echo "1. Создание события 'Новый год'"
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "date": "2023-12-31", "text": "Новый год"}' \
  | jq '.'
echo ""

echo "2. Создание события 'Встреча с друзьями'"
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "date": "2023-12-31", "text": "Встреча с друзьями"}' \
  | jq '.'
echo ""

echo "3. Создание события 'Рабочая встреча'"
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "date": "2024-01-02", "text": "Рабочая встреча"}' \
  | jq '.'
echo ""

echo "4. Создание события для другого пользователя"
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user456", "date": "2023-12-31", "text": "Другой пользователь"}' \
  | jq '.'
echo ""

# Получение событий
echo "📋 Получение событий..."
echo ""

echo "5. События на день 2023-12-31 для user123"
curl -s "http://localhost:8080/events_for_day?user_id=user123&date=2023-12-31" | jq '.'
echo ""

echo "6. События на день 2024-01-02 для user123"
curl -s "http://localhost:8080/events_for_day?user_id=user123&date=2024-01-02" | jq '.'
echo ""

echo "7. События на неделю (начиная с 2023-12-25) для user123"
curl -s "http://localhost:8080/events_for_week?user_id=user123&date=2023-12-25" | jq '.'
echo ""

echo "8. События на месяц (декабрь 2023) для user123"
curl -s "http://localhost:8080/events_for_month?user_id=user123&date=2023-12-15" | jq '.'
echo ""

echo "✏️ Обновление события..."
echo ""

EVENT_ID=$(curl -s "http://localhost:8080/events_for_day?user_id=user123&date=2023-12-31" | jq -r '.result[0].id')

if [ "$EVENT_ID" != "null" ] && [ "$EVENT_ID" != "" ]; then
    echo "9. Обновление события с ID: $EVENT_ID"
    curl -X POST http://localhost:8080/update_event \
      -H "Content-Type: application/json" \
      -d "{\"id\": \"$EVENT_ID\", \"user_id\": \"user123\", \"date\": \"2024-01-01\", \"text\": \"Обновленный Новый год\"}" \
      | jq '.'
    echo ""
else
    echo "❌ Не удалось получить ID события для обновления"
fi

echo "🗑️ Удаление события..."
echo ""

EVENT_ID_TO_DELETE=$(curl -s "http://localhost:8080/events_for_day?user_id=user123&date=2023-12-31" | jq -r '.result[0].id')

if [ "$EVENT_ID_TO_DELETE" != "null" ] && [ "$EVENT_ID_TO_DELETE" != "" ]; then
    echo "10. Удаление события с ID: $EVENT_ID_TO_DELETE"
    curl -X POST http://localhost:8080/delete_event \
      -H "Content-Type: application/json" \
      -d "{\"id\": \"$EVENT_ID_TO_DELETE\", \"user_id\": \"user123\"}" \
      | jq '.'
    echo ""
else
    echo "❌ Не удалось получить ID события для удаления"
fi

echo "⚠️ Тестирование обработки ошибок..."
echo ""

echo "11. Попытка создания события с некорректной датой"
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "date": "2023-13-31", "text": "Некорректная дата"}' \
  | jq '.'
echo ""

echo "12. Попытка получения событий без user_id"
curl -s "http://localhost:8080/events_for_day?date=2023-12-31" | jq '.'
echo ""

echo "13. Попытка обновления несуществующего события"
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{"id": "несуществующий_id", "user_id": "user123", "date": "2024-01-01", "text": "Тест"}' \
  | jq '.'
echo ""

echo "🎉 Демонстрация завершена!"
echo "================================================"
echo "Сервер продолжает работать на http://localhost:8080"
echo "Для остановки сервера нажмите Ctrl+C в терминале с сервером" 