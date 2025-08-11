package executor

import (
	"context"
	"fmt"
	"minishell/internal/parser"
	"os"
)

// CommandHandler интерфейс для обработки команд
type CommandHandler interface {
	CanHandle(command string) bool
	Execute(ctx context.Context, command string, args []string) error
}

// Executor реализует CommandExecutor
type Executor struct {
	parser   parser.Parser
	handlers []CommandHandler
}

// NewExecutor создает новый Executor
func NewExecutor(parser parser.Parser, handlers ...CommandHandler) *Executor {
	return &Executor{
		parser:   parser,
		handlers: handlers,
	}
}

// Execute выполняет команду или конвейер команд
func (e *Executor) Execute(ctx context.Context, input string) error {
	pipeline, err := e.parser.Parse(input)
	if err != nil {
		return fmt.Errorf("ошибка парсинга: %w", err)
	}

	if len(pipeline.Commands) == 0 {
		return nil
	}

	if len(pipeline.Commands) == 1 {
		// Одна команда
		return e.executeCommand(ctx, pipeline.Commands[0])
	}

	return e.executePipeline(ctx, pipeline.Commands)
}

// executeCommand выполняет одну команду (простую или условную)
func (e *Executor) executeCommand(ctx context.Context, cmd interface{}) error {
	switch c := cmd.(type) {
	case *parser.Command:
		return e.executeSimpleCommand(ctx, c)
	case *parser.ConditionalCommand:
		return e.executeConditionalCommand(ctx, c)
	default:
		return fmt.Errorf("неизвестный тип команды")
	}
}

// executeSimpleCommand выполняет простую команду с редиректами
func (e *Executor) executeSimpleCommand(ctx context.Context, cmd *parser.Command) error {
	// Сохраняем оригинальные потоки
	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Настраиваем редиректы
	if err := e.setupRedirects(cmd.Redirects); err != nil {
		return err
	}

	// Восстанавливаем потоки после выполнения команды
	defer func() {
		os.Stdin = originalStdin
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	// Выполняем команду
	for _, handler := range e.handlers {
		if handler.CanHandle(cmd.Name) {
			return handler.Execute(ctx, cmd.Name, cmd.Args)
		}
	}

	return e.executeExternalCommand(ctx, cmd.Name, cmd.Args)
}

// executeConditionalCommand выполняет условную команду
func (e *Executor) executeConditionalCommand(ctx context.Context, cmd *parser.ConditionalCommand) error {
	// Выполняем левую команду
	err := e.executeSimpleCommand(ctx, cmd.Left)

	switch cmd.Operator {
	case parser.OperatorAnd:
		// && - выполняем правую команду только если левая успешна
		if err != nil {
			return err
		}
		return e.executeCommand(ctx, cmd.Right)

	case parser.OperatorOr:
		// || - выполняем правую команду только если левая неуспешна
		if err == nil {
			return nil
		}
		return e.executeCommand(ctx, cmd.Right)

	default:
		return fmt.Errorf("неизвестный условный оператор")
	}
}

// setupRedirects настраивает перенаправления ввода/вывода
func (e *Executor) setupRedirects(redirects []parser.Redirect) error {
	for _, redirect := range redirects {
		switch redirect.Type {
		case parser.RedirectInput:
			file, err := os.Open(redirect.File)
			if err != nil {
				return fmt.Errorf("ошибка открытия файла для чтения %s: %w", redirect.File, err)
			}
			os.Stdin = file

		case parser.RedirectOutput:
			flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			if redirect.Append {
				flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
			}
			file, err := os.OpenFile(redirect.File, flags, 0644)
			if err != nil {
				return fmt.Errorf("ошибка открытия файла для записи %s: %w", redirect.File, err)
			}
			os.Stdout = file

		case parser.RedirectError:
			file, err := os.OpenFile(redirect.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("ошибка открытия файла для stderr %s: %w", redirect.File, err)
			}
			os.Stderr = file
		}
	}
	return nil
}

func (e *Executor) executePipeline(ctx context.Context, commands []interface{}) error {
	for _, cmd := range commands {
		if err := e.executeCommand(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) executeExternalCommand(ctx context.Context, command string, args []string) error {
	for _, handler := range e.handlers {
		if handler.CanHandle(command) {
			return handler.Execute(ctx, command, args)
		}
	}
	return fmt.Errorf("не найден обработчик для команды: %s", command)
}
