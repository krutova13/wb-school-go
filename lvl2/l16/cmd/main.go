package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"wget/internal/app"
	"wget/internal/config"
)

func main() {
	var (
		url           = flag.String("url", "", "URL для загрузки (обязательный)")
		depth         = flag.Int("depth", 3, "Глубина рекурсии (по умолчанию: 3)")
		outputDir     = flag.String("output", "./downloaded", "Директория для сохранения файлов")
		concurrency   = flag.Int("concurrency", 5, "Количество одновременных загрузок")
		timeout       = flag.Duration("timeout", 30*time.Second, "Таймаут для HTTP запросов")
		respectRobots = flag.Bool("robots", true, "Соблюдать robots.txt")
		userAgent     = flag.String("user-agent", "Wget/1.0", "User-Agent для запросов")
	)
	flag.Parse()

	if *url == "" {
		fmt.Println("Ошибка: URL обязателен")
		flag.Usage()
		os.Exit(1)
	}

	cfg := &config.Config{
		URL:           *url,
		Depth:         *depth,
		OutputDir:     *outputDir,
		Concurrency:   *concurrency,
		Timeout:       *timeout,
		RespectRobots: *respectRobots,
		UserAgent:     *userAgent,
	}

	wget := app.NewWget(cfg)

	fmt.Printf("Начинаю загрузку %s с глубиной %d\n", cfg.URL, cfg.Depth)
	fmt.Printf("Файлы будут сохранены в: %s\n", cfg.OutputDir)

	if err := wget.Run(); err != nil {
		log.Fatalf("Ошибка при выполнении: %v", err)
	}

	fmt.Println("Загрузка завершена успешно!")
}
