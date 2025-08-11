package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSimpleCommand(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("ls -la")
	assert.NoError(t, err)
	assert.Len(t, result.Commands, 1)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, "ls", cmd.Name)
	assert.Equal(t, []string{"-la"}, cmd.Args)
}

func TestParseCommandWithRedirects(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("echo hello > output.txt")
	assert.NoError(t, err)
	assert.Len(t, result.Commands, 1)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, "echo", cmd.Name)
	assert.Equal(t, []string{"hello"}, cmd.Args)
	assert.Len(t, cmd.Redirects, 1)
	assert.Equal(t, RedirectOutput, cmd.Redirects[0].Type)
	assert.Equal(t, "output.txt", cmd.Redirects[0].File)
}

func TestParseCommandWithAppendRedirect(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("echo world >> output.txt")
	assert.NoError(t, err)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, RedirectOutputAppend, cmd.Redirects[0].Type)
	assert.True(t, cmd.Redirects[0].Append)
}

func TestParseCommandWithInputRedirect(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("cat < input.txt")
	assert.NoError(t, err)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, RedirectInput, cmd.Redirects[0].Type)
	assert.Equal(t, "input.txt", cmd.Redirects[0].File)
}

func TestParseConditionalAnd(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("echo success && echo done")
	assert.NoError(t, err)
	assert.Len(t, result.Commands, 1)

	cond := result.Commands[0].(*ConditionalCommand)
	assert.Equal(t, OperatorAnd, cond.Operator)
	assert.Equal(t, "echo", cond.Left.Name)
	assert.Equal(t, []string{"success"}, cond.Left.Args)
	assert.Equal(t, "echo", cond.Right.Name)
	assert.Equal(t, []string{"done"}, cond.Right.Args)
}

func TestParseConditionalOr(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("false || echo fallback")
	assert.NoError(t, err)

	cond := result.Commands[0].(*ConditionalCommand)
	assert.Equal(t, OperatorOr, cond.Operator)
	assert.Equal(t, "false", cond.Left.Name)
	assert.Equal(t, "echo", cond.Right.Name)
}

func TestParseEnvironmentVariables(t *testing.T) {
	// Устанавливаем тестовую переменную
	err := os.Setenv("TEST_VAR", "test_value")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer func() {
		err := os.Unsetenv("TEST_VAR")
		if err != nil {
			t.Errorf("Failed to unset environment variable: %v", err)
		}
	}()

	parser := NewDefaultParser()

	result, err := parser.Parse("echo $TEST_VAR")
	assert.NoError(t, err)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, []string{"test_value"}, cmd.Args)
}

func TestParseQuotedStrings(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse(`echo "hello world" 'test string'`)
	assert.NoError(t, err)

	cmd := result.Commands[0].(*Command)
	assert.Equal(t, []string{"hello world", "test string"}, cmd.Args)
}

func TestParseComplexCommand(t *testing.T) {
	err := os.Setenv("HOME", "/home/user")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
	defer func() {
		err := os.Unsetenv("HOME")
		if err != nil {
			t.Errorf("Failed to unset environment variable: %v", err)
		}
	}()

	parser := NewDefaultParser()

	result, err := parser.Parse(`ls $HOME > output.txt && echo "done" >> log.txt`)
	assert.NoError(t, err)

	cond := result.Commands[0].(*ConditionalCommand)
	assert.Equal(t, OperatorAnd, cond.Operator)

	// Проверяем левую команду
	left := cond.Left
	assert.Equal(t, "ls", left.Name)
	assert.Equal(t, []string{"/home/user"}, left.Args)
	assert.Len(t, left.Redirects, 1)
	assert.Equal(t, RedirectOutput, left.Redirects[0].Type)

	// Проверяем правую команду
	right := cond.Right
	assert.Equal(t, "echo", right.Name)
	assert.Equal(t, []string{"done"}, right.Args)
	assert.Len(t, right.Redirects, 1)
	assert.Equal(t, RedirectOutputAppend, right.Redirects[0].Type)
}

func TestParsePipeline(t *testing.T) {
	parser := NewDefaultParser()

	result, err := parser.Parse("ls | grep go | wc -l")
	assert.NoError(t, err)
	assert.Len(t, result.Commands, 3)

	// Проверяем первую команду
	cmd1 := result.Commands[0].(*Command)
	assert.Equal(t, "ls", cmd1.Name)

	// Проверяем вторую команду
	cmd2 := result.Commands[1].(*Command)
	assert.Equal(t, "grep", cmd2.Name)
	assert.Equal(t, []string{"go"}, cmd2.Args)

	// Проверяем третью команду
	cmd3 := result.Commands[2].(*Command)
	assert.Equal(t, "wc", cmd3.Name)
	assert.Equal(t, []string{"-l"}, cmd3.Args)
}

func TestParseErrorCases(t *testing.T) {
	parser := NewDefaultParser()

	// Тест с пустой командой
	_, err := parser.Parse("")
	assert.NoError(t, err) // Пустая строка не должна вызывать ошибку

	// Тест с неверным редиректом
	_, err = parser.Parse("echo >")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ожидается файл после редиректа")

	// Тест с неполным условным оператором (теперь это валидная команда "echo")
	_, err = parser.Parse("echo &&")
	assert.NoError(t, err) // Теперь это валидная команда
}
