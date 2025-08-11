package main

import (
	"log"

	"minishell/internal/commands"
	"minishell/internal/executor"
	"minishell/internal/io"
	"minishell/internal/parser"
	"minishell/internal/shell"
)

func main() {
	defaultParser := parser.NewDefaultParser()
	builtinHandler := commands.NewBuiltinCommandHandler()
	externalHandler := commands.NewExternalCommandHandler()

	exec := executor.NewExecutor(defaultParser, builtinHandler, externalHandler)

	reader := io.NewStdinReader()
	writer := io.NewStdoutWriter()

	sh := shell.NewShell(exec, reader, writer)
	defer sh.Close()

	if err := sh.Run(); err != nil {
		log.Fatalf("Ошибка запуска shell: %v", err)
	}
}
