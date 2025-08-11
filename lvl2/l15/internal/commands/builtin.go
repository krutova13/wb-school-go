package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

// BuiltinCommandHandler обрабатывает встроенные команды
type BuiltinCommandHandler struct {
	handlers map[string]func(ctx context.Context, args []string) error
}

// NewBuiltinCommandHandler создает новый BuiltinCommandHandler
func NewBuiltinCommandHandler() *BuiltinCommandHandler {
	handler := &BuiltinCommandHandler{
		handlers: make(map[string]func(ctx context.Context, args []string) error),
	}

	handler.handlers["cd"] = handler.cd
	handler.handlers["pwd"] = handler.pwd
	handler.handlers["echo"] = handler.echo
	handler.handlers["kill"] = handler.kill
	handler.handlers["ps"] = handler.ps

	return handler
}

// CanHandle проверяет, может ли обработчик обработать команду
func (h *BuiltinCommandHandler) CanHandle(command string) bool {
	_, exists := h.handlers[command]
	return exists
}

// Execute выполняет встроенную команду
func (h *BuiltinCommandHandler) Execute(ctx context.Context, command string, args []string) error {
	handler, exists := h.handlers[command]
	if !exists {
		return fmt.Errorf("неизвестная встроенная команда: %s", command)
	}

	return handler(ctx, args)
}

func (h *BuiltinCommandHandler) cd(ctx context.Context, args []string) error {
	var path string
	if len(args) == 0 {
		// Переход в домашнюю директорию
		currentUser, err := user.Current()
		if err != nil {
			return fmt.Errorf("ошибка получения текущего пользователя: %w", err)
		}
		path = currentUser.HomeDir
	} else {
		path = args[0]
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("ошибка смены директории: %w", err)
	}

	return nil
}

func (h *BuiltinCommandHandler) pwd(ctx context.Context, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ошибка получения текущей директории: %w", err)
	}
	fmt.Println(dir)
	return nil
}

func (h *BuiltinCommandHandler) echo(ctx context.Context, args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func (h *BuiltinCommandHandler) kill(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("использование: kill <pid>")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("неверный PID: %s", args[0])
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("ошибка поиска процесса: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("ошибка отправки сигнала: %w", err)
	}

	fmt.Printf("Сигнал отправлен процессу %d\n", pid)
	return nil
}

func (h *BuiltinCommandHandler) ps(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, "ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
