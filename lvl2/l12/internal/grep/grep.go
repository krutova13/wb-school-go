package grep

import (
	"flag"
	"fmt"
	"greputil/internal"
	"regexp"
)

// ProcessGrepping выполняет основную логику поиска по паттерну
func ProcessGrepping() error {
	opts := internal.GetOpts()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("не указан паттерн для поиска")
	}
	pattern := args[0]

	var filePath string
	if len(args) > 1 {
		filePath = args[1]
	}

	lines, err := internal.ReadInput(filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения: %v", err)
	}

	if opts.Circle > 0 {
		opts.Before = opts.Circle
		opts.After = opts.Circle
	}

	var re *regexp.Regexp
	if opts.Fix {
		escapedPattern := regexp.QuoteMeta(pattern)
		re = compileRegex(opts, escapedPattern)
	} else {
		re = compileRegex(opts, pattern)
	}

	lineInfos := internal.ProcessLines(lines, re, opts)

	return internal.OutputResult(lineInfos, opts)
}

func compileRegex(opts *internal.GrepOptions, pattern string) *regexp.Regexp {
	if opts.Ignore {
		return regexp.MustCompile("(?i)" + pattern)
	}
	return regexp.MustCompile(pattern)
}
