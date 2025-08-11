package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExternalCommandHandler обрабатывает внешние команды
type ExternalCommandHandler struct{}

// NewExternalCommandHandler создает новый ExternalCommandHandler
func NewExternalCommandHandler() *ExternalCommandHandler {
	return &ExternalCommandHandler{}
}

// CanHandle проверяет, может ли обработчик обработать команду
func (h *ExternalCommandHandler) CanHandle(command string) bool {
	_ = command
	// Внешний обработчик обрабатывает все команды, которые не являются встроенными
	// Это будет проверяться в Executor после проверки встроенных команд
	return true
}

// Execute выполняет внешнюю команду
func (h *ExternalCommandHandler) Execute(ctx context.Context, command string, args []string) error {
	execPath, err := h.findExecutable(command)
	if err != nil {
		return fmt.Errorf("команда не найдена: %s", command)
	}

	cmd := exec.CommandContext(ctx, execPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (h *ExternalCommandHandler) findExecutable(command string) (string, error) {
	if filepath.IsAbs(command) || strings.Contains(command, "/") {
		return command, nil
	}

	path := os.Getenv("PATH")
	paths := strings.Split(path, ":")

	for _, p := range paths {
		execPath := filepath.Join(p, command)
		if _, err := os.Stat(execPath); err == nil {
			return execPath, nil
		}
	}

	return "", fmt.Errorf("команда не найдена")
}
