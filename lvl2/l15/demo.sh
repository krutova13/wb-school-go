#!/bin/bash

echo "=== Демонстрация Minishell - Полная версия ==="
echo ""

echo "=== БАЗОВЫЕ КОМАНДЫ ==="
echo ""

echo "1. Команда pwd:"
echo "pwd" | ./minishell
echo ""

echo "2. Команда echo:"
echo "echo Hello from minishell!" | ./minishell
echo ""

echo "3. Команда ls:"
echo "ls -la" | ./minishell
echo ""

echo "4. Конвейер ls | grep go:"
echo "ls | grep go" | ./minishell
echo ""

echo "5. Конвейер ps | grep go | wc -l:"
echo "ps | grep go | wc -l" | ./minishell
echo ""

echo "6. Команда cd и pwd:"
echo -e "cd ..\npwd" | ./minishell
echo ""

echo "7. Команда ps (список процессов):"
echo "ps aux | head -5" | ./minishell
echo ""

echo "=== РАСШИРЕННЫЕ ВОЗМОЖНОСТИ ==="
echo ""

echo "8. Условное выполнение команд (&&):"
echo "echo 'Первая команда успешна' && echo 'Вторая команда выполнилась'"
echo "echo 'Первая команда успешна' && echo 'Вторая команда выполнилась'" | ./minishell
echo ""

echo "9. Условное выполнение команд (||):"
echo "false || echo 'Команда выполнилась после неудачи'"
echo "false || echo 'Команда выполнилась после неудачи'" | ./minishell
echo ""

echo "10. Переменные окружения:"
echo "echo 'Домашняя директория: $HOME'"
echo "echo 'Домашняя директория: $HOME'" | ./minishell
echo ""

echo "11. Редирект вывода (>):"
echo "echo 'Этот текст будет в файле' > test_output.txt"
echo "echo 'Этот текст будет в файле' > test_output.txt" | ./minishell
echo "Содержимое файла test_output.txt:"
cat test_output.txt
echo ""

echo "12. Редирект вывода с добавлением (>>):"
echo "echo 'Дополнительная строка' >> test_output.txt"
echo "echo 'Дополнительная строка' >> test_output.txt" | ./minishell
echo "Обновленное содержимое файла test_output.txt:"
cat test_output.txt
echo ""

echo "13. Редирект ввода (<):"
echo "cat < test_output.txt"
echo "cat < test_output.txt" | ./minishell
echo ""

echo "14. Комбинированные команды:"
echo "echo 'Создаем файл' > combined.txt && echo 'Добавляем строку' >> combined.txt && cat < combined.txt"
echo "echo 'Создаем файл' > combined.txt && echo 'Добавляем строку' >> combined.txt && cat < combined.txt" | ./minishell
echo ""

echo "15. Кавычки и переменные:"
echo "echo 'Текущий путь: $PWD'"
echo "echo 'Текущий путь: $PWD'" | ./minishell
echo ""

echo "16. Конвейеры с условным выполнением:"
echo "ls | grep go && echo 'Найдены Go файлы' || echo 'Go файлы не найдены'"
echo "ls | grep go && echo 'Найдены Go файлы' || echo 'Go файлы не найдены'" | ./minishell
echo ""

echo "17. Обработка ошибок:"
echo "nonexistent_command && echo 'Эта команда не должна выполниться'"
echo "nonexistent_command && echo 'Эта команда не должна выполниться'" | ./minishell
echo ""

echo "18. Сложная команда с переменными и редиректами:"
echo "echo 'Пользователь: $USER' > user_info.txt && echo 'Путь: $PWD' >> user_info.txt && cat user_info.txt"
echo "echo 'Пользователь: $USER' > user_info.txt && echo 'Путь: $PWD' >> user_info.txt && cat user_info.txt" | ./minishell
echo ""

echo "19. Очистка тестовых файлов:"
echo "rm test_output.txt combined.txt user_info.txt"
echo "rm test_output.txt combined.txt user_info.txt" | ./minishell
echo ""

echo "=== Демонстрация завершена ===" 