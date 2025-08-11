package internal

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// ReadInput читает входной файл построчно.
// Если указан аргумент командной строки, читает из файла,
// иначе читает из стандартного ввода
func ReadInput() ([]string, error) {
	var lines []string
	var scanner *bufio.Scanner

	args := flag.Args()
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		line := scanner.Text()
		opts := GetOpts()
		if opts.TrailingBlanks {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// PrepareLines преобразует строки в срез структур для последующей сортировки
func PrepareLines(lines []string) SortableLines {
	var sortableLines SortableLines

	for _, line := range lines {
		fields := strings.Split(line, "\t")
		key := line

		opts := GetOpts()
		if opts.Column > 0 && opts.Column <= len(fields) {
			key = fields[opts.Column-1]
		}

		sortableLines = append(sortableLines, Line{
			Original: line,
			Fields:   fields,
			Key:      key,
		})
	}

	return sortableLines
}

// OutputLines выводит отсортированные строки.
// Если включена опция Unique, выводит только уникальные строки
func OutputLines(sortableLines SortableLines) {
	seen := make(map[string]bool)
	opts := GetOpts()

	for _, line := range sortableLines {
		if opts.Unique {
			if seen[line.Original] {
				continue
			}
			seen[line.Original] = true
		}
		fmt.Println(line.Original)
	}
}
