package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// CommandExecutor определяет интерфейс для выполнения команд
type CommandExecutor interface {
	Execute(ctx context.Context, input string) error
}

// InputReader определяет интерфейс для чтения ввода
type InputReader interface {
	ReadPrompt() (string, error)
}

// OutputWriter определяет интерфейс для записи вывода
type OutputWriter interface {
	WritePrompt()
	WriteLine(line string)
	WriteError(err error)
}

// Shell представляет основной интерпретатор командной строки
type Shell struct {
	executor CommandExecutor
	reader   InputReader
	writer   OutputWriter
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewShell создает новый экземпляр shell
func NewShell(executor CommandExecutor, reader InputReader, writer OutputWriter) *Shell {
	ctx, cancel := context.WithCancel(context.Background())
	return &Shell{
		executor: executor,
		reader:   reader,
		writer:   writer,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Run запускает основной цикл shell
func (s *Shell) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		s.cancel()
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.writer.WriteLine("\nДо свидания!")
			return nil
		default:
			s.writer.WritePrompt()

			input, err := s.reader.ReadPrompt()
			if err != nil {
				if err == io.EOF {
					s.writer.WriteLine("\nДо свидания!")
					return nil
				}
				s.writer.WriteError(fmt.Errorf("ошибка чтения: %w", err))
				continue
			}

			if input == "" {
				continue
			}

			if err := s.executor.Execute(s.ctx, input); err != nil {
				s.writer.WriteError(err)
			}
		}
	}
}

// Close освобождает ресурсы shell
func (s *Shell) Close() {
	s.cancel()
}
