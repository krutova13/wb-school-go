package parser

import (
	"flag"
	"fmt"
	"time"

	"telnet/internal/types"
)

// Parser интерфейс для парсинга аргументов
type Parser interface {
	Parse(args []string) (*types.Config, error)
}

// CommandLineParser реализует парсинг аргументов командной строки
type CommandLineParser struct{}

// NewCommandLineParser создает новый парсер
func NewCommandLineParser() *CommandLineParser {
	return &CommandLineParser{}
}

// Parse парсит аргументы командной строки и возвращает конфигурацию
func (p *CommandLineParser) Parse(args []string) (*types.Config, error) {
	fs := flag.NewFlagSet("telnet", flag.ExitOnError)

	var timeout time.Duration
	fs.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connection")

	err := fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	positionalArgs := fs.Args()
	if len(positionalArgs) < 2 {
		return nil, fmt.Errorf("usage: telnet [--timeout=10s] host port")
	}

	host := positionalArgs[0]
	port := positionalArgs[1]

	if host == "" || port == "" {
		return nil, fmt.Errorf("host and port are required")
	}

	return &types.Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}, nil
}
