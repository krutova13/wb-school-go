package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Port int
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	portStr := os.Getenv("CALENDAR_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный порт: %v", err)
	}

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("порт должен быть в диапазоне 1-65535")
	}

	return &Config{
		Port: port,
	}, nil
}
