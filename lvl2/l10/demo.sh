#!/bin/bash

echo "=== Демонстрация Sort Utility ==="
echo

echo "1. Базовая сортировка:"
echo -e "zebra\napple\nbanana" | ./sort
echo

echo "2. Числовая сортировка:"
echo -e "10\n2\n1\n20" | ./sort -n
echo

echo "3. Сортировка в обратном порядке:"
echo -e "apple\nbanana\ncherry" | ./sort -r
echo

echo "4. Сортировка по второму столбцу (числовая):"
echo -e "melon\t5\napple\t3\nbanana\t1" | ./sort -k 2 -n
echo

echo "5. Сортировка по месяцам:"
echo -e "Dec\nJan\nMar\nFeb" | ./sort -M
echo

echo "6. Сортировка человекочитаемых размеров:"
echo -e "1KB\n2MB\n500B\n1GB" | ./sort -h
echo

echo "7. Удаление дубликатов:"
echo -e "apple\nbanana\napple\ncherry\nbanana" | ./sort -u
echo

echo "8. Проверка сортированности (отсортировано):"
echo -e "apple\nbanana\ncherry" | ./sort -c
echo

echo "9. Проверка сортированности (не отсортировано):"
echo -e "cherry\napple\nbanana" | ./sort -c
echo

echo "10. Комбинация флагов (-k 2 -n -r):"
echo -e "melon\t5\napple\t3\nbanana\t1" | ./sort -k 2 -n -r
echo

echo "11. Сортировка файла test_data.txt по второму столбцу:"
./sort -k 2 -n data/test_data.txt
echo

echo "12. Сортировка файла test_data.txt по третьему столбцу (месяцы):"
./sort -k 3 -M data/test_data.txt
echo

echo "13. Сортировка numeric_data.txt по размерам файлов:"
./sort -k 2 -h data/numeric_data.txt
echo
