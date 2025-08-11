package io

import (
	"fmt"
	"os"
)

// StdoutWriter реализует OutputWriter для записи в stdout/stderr
type StdoutWriter struct{}

// NewStdoutWriter создает новый StdoutWriter
func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

// WritePrompt выводит приглашение командной строки
func (w *StdoutWriter) WritePrompt() {
	fmt.Print("minishell> ")
}

// WriteLine выводит строку в stdout
func (w *StdoutWriter) WriteLine(line string) {
	fmt.Println(line)
}

// WriteError выводит ошибку в stderr
func (w *StdoutWriter) WriteError(err error) {
	fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
}
