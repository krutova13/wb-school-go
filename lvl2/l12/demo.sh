#!/bin/bash

echo "=== Демонстрация Grep Utility ==="
echo

echo "1. Простой поиск слова 'hello':"
./grep "hello" data/test_data.txt
echo

echo "2. Поиск 'hello' с игнорированием регистра (-i):"
./grep -i "hello" data/test_data.txt
echo

echo "3. Инвертированный поиск 'hello' (-v):"
./grep -v "hello" data/test_data.txt
echo

echo "4. Поиск 'hello' с номерами строк (-n):"
./grep -n "hello" data/test_data.txt
echo

echo "5. Подсчет количества совпадений 'hello' (-c):"
./grep -c "hello" data/test_data.txt
echo

echo "6. Поиск 'world' с контекстом (-C 2):"
./grep -C 2 "world" data/test_data.txt
echo

echo "7. Поиск 'world' с контекстом до (-B 2):"
./grep -B 2 "world" data/test_data.txt
echo

echo "8. Поиск 'world' с контекстом после (-A 2):"
./grep -A 2 "world" data/test_data.txt
echo

echo "9. Поиск 'hello' с игнорированием регистра и инвертированием (-i -v):"
./grep -i -v "hello" data/test_data.txt
echo

echo "10. Поиск 'hello' с номерами строк и игнорированием регистра (-n -i):"
./grep -n -i "hello" data/test_data.txt
echo

echo "11. Поиск 'hello world' с фиксированной строкой (-F):"
./grep -F "hello world" data/test_data.txt
echo

echo "12. Поиск 'hello world' с фиксированной строкой и игнорированием регистра (-F -i):"
./grep -F -i "hello world" data/test_data.txt
echo

echo "13. Поиск с регулярным выражением (якорь начала строки):"
./grep "^Третья" data/test_data.txt
echo

echo "14. Поиск с регулярным выражением (якорь конца строки):"
./grep "строка$" data/test_data.txt
echo

echo "15. Поиск цифр с регулярным выражением:"
./grep "[0-9]+" data/test_data.txt
echo

echo "16. Поиск цифр с инвертированием:"
./grep -v "[0-9]+" data/test_data.txt
echo

echo "17. Поиск с альтернативами (hello|world):"
./grep "hello|world" data/test_data.txt
echo

echo "18. Поиск с альтернативами и игнорированием регистра:"
./grep -i "hello|world" data/test_data.txt
echo

echo "19. Поиск пустых строк:"
echo -e "Строка 1\n\nСтрока 2\n\nСтрока 3" | ./grep "^$"
echo

echo "20. Поиск пустых строк с инвертированием:"
echo -e "Строка 1\n\nСтрока 2\n\nСтрока 3" | ./grep -v "^$"
echo

echo "21. Поиск с множественными совпадениями в строке:"
echo -e "hello hello world\nworld\nhello world hello" | ./grep "hello"
echo

echo "22. Поиск с множественными совпадениями и инвертированием:"
echo -e "hello hello world\nworld\nhello world hello" | ./grep -v "hello"
echo

echo "23. Поиск специальных символов (точка как символ):"
echo -e "hello.world\nhello+world\nhello*world" | ./grep "hello\\.world"
echo

echo "24. Поиск специальных символов с фиксированной строкой:"
echo -e "hello.world\nhello+world\nhello*world" | ./grep -F "hello.world"
echo

echo "25. Комбинация всех флагов (-i -v -F -n):"
./grep -i -v -F -n "hello world" data/test_data.txt
echo

echo "=== Демонстрация завершена ==="