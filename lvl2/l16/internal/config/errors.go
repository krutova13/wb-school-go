package config

import "errors"

var (
	// ErrEmptyURL возвращается когда URL не может быть пустым
	ErrEmptyURL = errors.New("URL не может быть пустым")
	// ErrInvalidDepth возвращается когда глубина должна быть неотрицательной
	ErrInvalidDepth = errors.New("глубина должна быть неотрицательной")
	// ErrInvalidConcurrency возвращается когда количество одновременных загрузок должно быть положительным
	ErrInvalidConcurrency = errors.New("количество одновременных загрузок должно быть положительным")
	// ErrInvalidTimeout возвращается когда таймаут должен быть положительным
	ErrInvalidTimeout = errors.New("таймаут должен быть положительным")
)
