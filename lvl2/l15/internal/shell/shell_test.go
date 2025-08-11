package shell

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandExecutor - мок для CommandExecutor
type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(ctx context.Context, input string) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// MockInputReader - мок для InputReader
type MockInputReader struct {
	mock.Mock
}

func (m *MockInputReader) ReadPrompt() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// MockOutputWriter - мок для OutputWriter
type MockOutputWriter struct {
	mock.Mock
}

func (m *MockOutputWriter) WritePrompt() {
	m.Called()
}

func (m *MockOutputWriter) WriteLine(line string) {
	m.Called(line)
}

func (m *MockOutputWriter) WriteError(err error) {
	m.Called(err)
}

func TestNewShell(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	sh := NewShell(executor, reader, writer)

	assert.NotNil(t, sh)
	assert.Equal(t, executor, sh.executor)
	assert.Equal(t, reader, sh.reader)
	assert.Equal(t, writer, sh.writer)
	assert.NotNil(t, sh.ctx)
	assert.NotNil(t, sh.cancel)
}

func TestShellClose(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	sh := NewShell(executor, reader, writer)

	assert.NotPanics(t, func() {
		sh.Close()
	})
}

func TestShellRunWithEOF(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", io.EOF).Once()
	writer.On("WriteLine", "\nДо свидания!").Once()

	sh := NewShell(executor, reader, writer)

	err := sh.Run()

	assert.NoError(t, err)
	writer.AssertExpectations(t)
	reader.AssertExpectations(t)
}

func TestShellRunWithEmptyInput(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", nil).Once()
	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", io.EOF).Once()
	writer.On("WriteLine", "\nДо свидания!").Once()

	sh := NewShell(executor, reader, writer)

	err := sh.Run()

	assert.NoError(t, err)
	writer.AssertExpectations(t)
	reader.AssertExpectations(t)
}

func TestShellRunWithCommand(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("pwd", nil).Once()
	executor.On("Execute", mock.Anything, "pwd").Return(nil).Once()
	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", io.EOF).Once()
	writer.On("WriteLine", "\nДо свидания!").Once()

	sh := NewShell(executor, reader, writer)

	err := sh.Run()

	assert.NoError(t, err)
	writer.AssertExpectations(t)
	reader.AssertExpectations(t)
	executor.AssertExpectations(t)
}

func TestShellRunWithCommandError(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	expectedErr := errors.New("command error")

	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("invalid", nil).Once()
	executor.On("Execute", mock.Anything, "invalid").Return(expectedErr).Once()
	writer.On("WriteError", expectedErr).Once()
	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", io.EOF).Once()
	writer.On("WriteLine", "\nДо свидания!").Once()

	sh := NewShell(executor, reader, writer)

	err := sh.Run()

	assert.NoError(t, err)
	writer.AssertExpectations(t)
	reader.AssertExpectations(t)
	executor.AssertExpectations(t)
}

func TestShellRunWithReadError(t *testing.T) {
	executor := &MockCommandExecutor{}
	reader := &MockInputReader{}
	writer := &MockOutputWriter{}

	readErr := errors.New("read error")

	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", readErr).Once()
	writer.On("WriteError", mock.AnythingOfType("*fmt.wrapError")).Once()
	writer.On("WritePrompt").Once()
	reader.On("ReadPrompt").Return("", io.EOF).Once()
	writer.On("WriteLine", "\nДо свидания!").Once()

	sh := NewShell(executor, reader, writer)

	err := sh.Run()

	assert.NoError(t, err)
	writer.AssertExpectations(t)
	reader.AssertExpectations(t)
}
