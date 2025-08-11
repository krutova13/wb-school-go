package app

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"wget/internal/config"
	"wget/internal/downloader"
	"wget/internal/filter"
	"wget/internal/parser"
	"wget/internal/reporter"
	"wget/internal/storage"
	"wget/internal/types"
)

// Wget основное приложение для загрузки веб-страниц
type Wget struct {
	config     *config.Config
	downloader types.Downloader
	parser     types.Parser
	storage    types.Storage
	filter     types.URLFilter
	reporter   types.ProgressReporter

	// Очередь URL для обработки
	urlChan chan *url.URL

	// Счетчик активных задач
	activeTasks int
	tasksMux    sync.Mutex

	// Флаг завершения
	done chan struct{}
}

// NewWget создает новый экземпляр Wget
func NewWget(cfg *config.Config) *Wget {
	return &Wget{
		config: cfg,
		done:   make(chan struct{}),
	}
}

// Run запускает процесс загрузки
func (w *Wget) Run() error {
	if err := w.config.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации: %w", err)
	}

	if err := w.initialize(); err != nil {
		return fmt.Errorf("ошибка инициализации: %w", err)
	}

	startURL, err := url.Parse(w.config.URL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга URL: %w", err)
	}

	w.urlChan = make(chan *url.URL, w.config.Concurrency*10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < w.config.Concurrency; i++ {
		wg.Add(1)
		go w.worker(ctx, &wg)
	}

	w.urlChan <- startURL

	go w.monitor(ctx, &wg)

	wg.Wait()

	w.reporter.PrintSummary()

	return nil
}

func (w *Wget) initialize() error {
	w.downloader = downloader.NewHTTPDownloader(w.config.Timeout, w.config.UserAgent)

	w.parser = parser.NewHTMLParser()

	fileStorage, err := storage.NewFileStorage(w.config.OutputDir)
	if err != nil {
		return fmt.Errorf("ошибка создания хранилища: %w", err)
	}
	w.storage = fileStorage

	urlFilter, err := filter.NewURLFilter(w.config)
	if err != nil {
		return fmt.Errorf("ошибка создания фильтра: %w", err)
	}
	w.filter = urlFilter

	w.reporter = reporter.NewProgressReporter()

	return nil
}

func (w *Wget) monitor(ctx context.Context, wg *sync.WaitGroup) {
	_ = wg
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tasksMux.Lock()
			active := w.activeTasks
			w.tasksMux.Unlock()

			if active == 0 && len(w.urlChan) == 0 {
				close(w.urlChan)
				return
			}
		}
	}
}

func (w *Wget) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case resourceURL, ok := <-w.urlChan:
			if !ok {
				// Канал закрыт, завершаем работу
				return
			}

			w.tasksMux.Lock()
			w.activeTasks++
			w.tasksMux.Unlock()

			w.processURL(ctx, resourceURL, 0)

			w.tasksMux.Lock()
			w.activeTasks--
			w.tasksMux.Unlock()
		}
	}
}

func (w *Wget) processURL(ctx context.Context, resourceURL *url.URL, depth int) {
	if !w.filter.ShouldDownload(resourceURL, depth) {
		w.reporter.ReportSkipped(resourceURL, "фильтр")
		return
	}

	w.filter.MarkVisited(resourceURL)

	localPath := w.storage.GetLocalPath(resourceURL)
	if w.storage.Exists(localPath) {
		w.reporter.ReportSkipped(resourceURL, "уже существует")
		return
	}

	resource, err := w.downloader.Download(ctx, resourceURL)
	if err != nil {
		w.reporter.ReportError(resourceURL, err)
		return
	}

	if err := w.storage.Save(resource); err != nil {
		w.reporter.ReportError(resourceURL, err)
		return
	}

	w.reporter.ReportDownloaded(resource)

	if resource.Type == types.ResourceTypeHTML && depth < w.config.Depth {
		w.extractAndQueueLinks(resource, depth+1)
	}

	if resource.Type == types.ResourceTypeCSS {
		w.extractCSSLinks(resource, depth+1)
	}
}

func (w *Wget) extractAndQueueLinks(resource *types.Resource, depth int) {
	urls, err := w.parser.ParseHTML(resource.Content, resource.URL)
	if err != nil {
		w.reporter.ReportError(resource.URL, fmt.Errorf("ошибка парсинга HTML: %w", err))
		return
	}

	for _, resourceURL := range urls {
		if w.filter.ShouldDownload(resourceURL, depth) {
			select {
			case w.urlChan <- resourceURL:
			default:
				// Очередь переполнена, пропускаем
				w.reporter.ReportSkipped(resourceURL, "очередь переполнена")
			}
		}
	}
}

func (w *Wget) extractCSSLinks(resource *types.Resource, depth int) {
	urls, err := w.parser.ParseCSS(resource.Content, resource.URL)
	if err != nil {
		w.reporter.ReportError(resource.URL, fmt.Errorf("ошибка парсинга CSS: %w", err))
		return
	}

	for _, resourceURL := range urls {
		if w.filter.ShouldDownload(resourceURL, depth) {
			select {
			case w.urlChan <- resourceURL:
			default:
				// Очередь переполнена, пропускаем
				w.reporter.ReportSkipped(resourceURL, "очередь переполнена")
			}
		}
	}
}
