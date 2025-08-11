package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// ReadInput читает входной файл построчно.
// Если указан аргумент командной строки, читает из файла,
// иначе читает из стандартного ввода
func ReadInput(filePath string) ([]string, error) {
	var lines []string
	var scanner *bufio.Scanner

	if filePath != "" {
		file, err := os.Open(filePath)
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
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// ProcessLines обрабатывает строки и находит совпадения.
// Подготавливает срез структур LineInfo для последующей обработки
func ProcessLines(lines []string, re *regexp.Regexp, opts *GrepOptions) []LineInfo {
	var linesInfos []LineInfo
	for i, line := range lines {
		match := re.MatchString(line)

		if opts.Invert {
			match = !match
		}

		linesInfos = append(linesInfos, LineInfo{
			Number: i + 1,
			Text:   line,
			Match:  match,
		})
	}
	return linesInfos
}

// OutputResult выводит результат
func OutputResult(lineInfos []LineInfo, opts *GrepOptions) error {
	if opts.Count {
		var count int
		for _, lineInfo := range lineInfos {
			if lineInfo.Match {
				count++
			}
		}
		fmt.Println(count)
		return nil
	}

	resultLines := collectOutputLines(lineInfos, opts)

	for _, line := range resultLines {
		if opts.Number {
			fmt.Printf("%d:", line.Number)
		}
		fmt.Println(line.Text)
	}

	return nil
}

func collectOutputLines(lineInfos []LineInfo, opts *GrepOptions) []LineInfo {
	var resultLines []LineInfo
	set := make(map[int]struct{})
	for i, line := range lineInfos {
		if line.Match {
			if opts.Before > 0 {
				start := i - opts.Before
				if start < 0 {
					start = 0
				}
				for j := start; j < i; j++ {
					resultLines = appendIfNotExist(set, resultLines, j, lineInfos[j])
				}
			}

			resultLines = appendIfNotExist(set, resultLines, i, lineInfos[i])

			if opts.After > 0 {
				end := i + opts.After + 1
				if end > len(lineInfos) {
					end = len(lineInfos)
				}
				for j := i + 1; j < end; j++ {
					resultLines = appendIfNotExist(set, resultLines, j, lineInfos[j])
				}
			}
		}
	}
	return resultLines
}

func appendIfNotExist(set map[int]struct{}, resultLines []LineInfo, iter int, line LineInfo) []LineInfo {
	if _, exists := set[iter]; !exists {
		resultLines = append(resultLines, line)
		set[iter] = struct{}{}
	}
	return resultLines
}
