package io

import (
	"bufio"
	"os"
	"strings"
)

// StdinReader реализует InputReader для чтения из stdin
type StdinReader struct {
	reader *bufio.Reader
}

// NewStdinReader создает новый StdinReader
func NewStdinReader() *StdinReader {
	return &StdinReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadPrompt читает строку из stdin
func (r *StdinReader) ReadPrompt() (string, error) {
	input, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}
