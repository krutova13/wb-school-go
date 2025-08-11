package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Command представляет разобранную команду
type Command struct {
	Name      string
	Args      []string
	Redirects []Redirect
}

// Redirect представляет перенаправление ввода/вывода
type Redirect struct {
	Type   RedirectType
	File   string
	Append bool // для >>
}

// RedirectType определяет тип перенаправления
type RedirectType int

const (
	// RedirectInput перенаправляет ввод из файла (<)
	RedirectInput RedirectType = iota // <
	// RedirectOutput перенаправляет вывод в файл (>)
	RedirectOutput // >
	// RedirectOutputAppend перенаправляет вывод в файл с добавлением (>>)
	RedirectOutputAppend // >>
	// RedirectError перенаправляет ошибки в файл (2>)
	RedirectError // 2>
)

// ConditionalCommand представляет условное выполнение команд
type ConditionalCommand struct {
	Left     *Command
	Operator ConditionalOperator
	Right    *Command
}

// ConditionalOperator определяет тип условного оператора
type ConditionalOperator int

const (
	// OperatorAnd логический оператор И (&&)
	OperatorAnd ConditionalOperator = iota // &&
	// OperatorOr логический оператор ИЛИ (||)
	OperatorOr // ||
)

// Pipeline представляет конвейер команд с условным выполнением
type Pipeline struct {
	Commands []interface{} // Command или ConditionalCommand
}

// Parser интерфейс для парсинга команд
type Parser interface {
	Parse(input string) (*Pipeline, error)
}

// DefaultParser реализует Parser для разбора команд shell
type DefaultParser struct{}

// NewDefaultParser создает новый DefaultParser
func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

// Parse разбирает строку ввода в конвейер команд
func (p *DefaultParser) Parse(input string) (*Pipeline, error) {
	// Сначала обрабатываем условные операторы, потом конвейеры
	// Это нужно для правильного приоритета операторов

	// Разбиваем по конвейерам, но сохраняем условные операторы внутри
	pipelineStrs := p.splitPipeline(input)
	commands := make([]interface{}, 0, len(pipelineStrs))

	for _, cmdStr := range pipelineStrs {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue
		}

		// Парсим условные команды
		parsed, err := p.parseConditionalCommand(cmdStr)
		if err != nil {
			return nil, err
		}
		commands = append(commands, parsed)
	}

	return &Pipeline{Commands: commands}, nil
}

// splitPipeline разбивает строку по конвейерам, но не внутри условных выражений
func (p *DefaultParser) splitPipeline(input string) []string {
	var result []string
	var current strings.Builder
	var parenCount int // для будущей поддержки скобок

	for i := 0; i < len(input); i++ {
		char := input[i]

		// Проверяем на двойной символ | (||)
		if char == '|' && i+1 < len(input) && input[i+1] == '|' {
			// Это ||, пропускаем
			current.WriteByte(char)
			current.WriteByte(input[i+1])
			i++ // Пропускаем следующий символ
		} else if char == '|' && parenCount == 0 {
			// Нашли одиночный конвейер
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// parseConditionalCommand парсит команды с условным выполнением
func (p *DefaultParser) parseConditionalCommand(input string) (interface{}, error) {
	// Ищем операторы && и ||
	andIndex := strings.Index(input, " && ")
	orIndex := strings.Index(input, " || ")

	// Если нет условных операторов, парсим как простую команду
	if andIndex == -1 && orIndex == -1 {
		return p.parseSimpleCommand(input)
	}

	// Определяем какой оператор использовать (приоритет слева направо)
	var operatorIndex int
	var operator ConditionalOperator

	if andIndex != -1 && (orIndex == -1 || andIndex < orIndex) {
		operatorIndex = andIndex
		operator = OperatorAnd
	} else {
		operatorIndex = orIndex
		operator = OperatorOr
	}

	// Разбиваем строку по оператору
	leftStr := strings.TrimSpace(input[:operatorIndex])
	rightStr := strings.TrimSpace(input[operatorIndex+4:]) // +4 для " && " или " || "

	// Парсим левую и правую части
	left, err := p.parseSimpleCommand(leftStr)
	if err != nil {
		return nil, err
	}

	// Рекурсивно парсим правую часть (может содержать еще условные операторы)
	right, err := p.parseConditionalCommand(rightStr)
	if err != nil {
		return nil, err
	}

	// Если правая часть - простая команда, создаем условную команду
	if cmd, ok := right.(*Command); ok {
		return &ConditionalCommand{
			Left:     left,
			Operator: operator,
			Right:    cmd,
		}, nil
	}

	// Если правая часть уже условная команда, создаем цепочку
	if cond, ok := right.(*ConditionalCommand); ok {
		// Создаем новую условную команду, где левая часть - наша левая команда,
		// а правая - вся цепочка условных команд
		return &ConditionalCommand{
			Left:     left,
			Operator: operator,
			Right:    cond.Left, // Берем только левую часть цепочки
		}, nil
	}

	return nil, fmt.Errorf("неизвестный тип команды")
}

// parseSimpleCommand парсит простую команду с редиректами
func (p *DefaultParser) parseSimpleCommand(input string) (*Command, error) {
	// Разбиваем на части, сохраняя кавычки
	parts := p.splitPreservingQuotes(input)
	if len(parts) == 0 {
		return nil, fmt.Errorf("пустая команда")
	}

	command := &Command{
		Name:      parts[0],
		Args:      make([]string, 0),
		Redirects: make([]Redirect, 0),
	}

	// Обрабатываем аргументы и редиректы
	for i := 1; i < len(parts); i++ {
		part := parts[i]

		// Проверяем на редирект
		if p.isRedirect(part) {
			if i+1 >= len(parts) {
				return nil, fmt.Errorf("ожидается файл после редиректа: %s", part)
			}

			redirect, err := p.parseRedirect(part, parts[i+1])
			if err != nil {
				return nil, err
			}
			command.Redirects = append(command.Redirects, redirect)
			i++ // Пропускаем следующий аргумент (имя файла)
		} else {
			// Подставляем переменные окружения
			expanded := p.expandEnvironmentVariables(part)
			command.Args = append(command.Args, expanded)
		}
	}

	return command, nil
}

// isRedirect проверяет, является ли токен редиректом
func (p *DefaultParser) isRedirect(token string) bool {
	return token == "<" || token == ">" || token == ">>" || token == "2>"
}

// parseRedirect парсит редирект
func (p *DefaultParser) parseRedirect(operator, file string) (Redirect, error) {
	var redirectType RedirectType
	var shouldAppend bool

	switch operator {
	case "<":
		redirectType = RedirectInput
	case ">":
		redirectType = RedirectOutput
	case ">>":
		redirectType = RedirectOutputAppend
		shouldAppend = true
	case "2>":
		redirectType = RedirectError
	default:
		return Redirect{}, fmt.Errorf("неизвестный тип редиректа: %s", operator)
	}

	return Redirect{
		Type:   redirectType,
		File:   p.expandEnvironmentVariables(file),
		Append: shouldAppend,
	}, nil
}

// splitPreservingQuotes разбивает строку, сохраняя кавычки
func (p *DefaultParser) splitPreservingQuotes(input string) []string {
	var result []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	for _, char := range input {
		if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
			} else {
				current.WriteRune(char)
			}
		} else if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// expandEnvironmentVariables подставляет переменные окружения
func (p *DefaultParser) expandEnvironmentVariables(input string) string {
	// Простая подстановка переменных вида $VAR
	re := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[1:] // Убираем $
		if value, exists := os.LookupEnv(varName); exists {
			return value
		}
		return "" // Если переменная не найдена, возвращаем пустую строку
	})
}
