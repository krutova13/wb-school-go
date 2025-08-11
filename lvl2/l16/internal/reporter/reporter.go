package reporter

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"wget/internal/types"
)

// ProgressReporter реализует интерфейс ProgressReporter для отчета о прогрессе
type ProgressReporter struct {
	downloaded int
	skipped    int
	errors     int
	startTime  time.Time
	mux        sync.RWMutex
}

// NewProgressReporter создает новый репортер прогресса
func NewProgressReporter() *ProgressReporter {
	return &ProgressReporter{
		startTime: time.Now(),
	}
}

// ReportDownloaded сообщает о успешно загруженном ресурсе
func (r *ProgressReporter) ReportDownloaded(resource *types.Resource) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.downloaded++
	fmt.Printf("✓ Загружен: %s (%s)\n", resource.URL.String(), resource.Type.String())
}

// ReportSkipped сообщает о пропущенном URL
func (r *ProgressReporter) ReportSkipped(url *url.URL, reason string) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.skipped++
	fmt.Printf("- Пропущен: %s (%s)\n", url.String(), reason)
}

// ReportError сообщает об ошибке загрузки
func (r *ProgressReporter) ReportError(url *url.URL, err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.errors++
	fmt.Printf("✗ Ошибка: %s - %v\n", url.String(), err)
}

// GetStats возвращает статистику загрузки
func (r *ProgressReporter) GetStats() Stats {
	r.mux.RLock()
	defer r.mux.RUnlock()

	duration := time.Since(r.startTime)

	return Stats{
		Downloaded: r.downloaded,
		Skipped:    r.skipped,
		Errors:     r.errors,
		Duration:   duration,
	}
}

// PrintSummary выводит итоговую статистику
func (r *ProgressReporter) PrintSummary() {
	stats := r.GetStats()

	fmt.Println("\n=== Итоговая статистика ===")
	fmt.Printf("Загружено: %d файлов\n", stats.Downloaded)
	fmt.Printf("Пропущено: %d URL\n", stats.Skipped)
	fmt.Printf("Ошибок: %d\n", stats.Errors)
	fmt.Printf("Время выполнения: %v\n", stats.Duration)

	if stats.Duration > 0 {
		rate := float64(stats.Downloaded) / stats.Duration.Seconds()
		fmt.Printf("Скорость: %.2f файлов/сек\n", rate)
	}
}

// Stats содержит статистику загрузки
type Stats struct {
	Downloaded int
	Skipped    int
	Errors     int
	Duration   time.Duration
}
