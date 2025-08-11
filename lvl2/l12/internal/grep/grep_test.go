package grep

import (
	"greputil/internal"
	"regexp"
	"testing"
)

func withOptions(options internal.GrepOptions, testFunc func()) {
	originalOpts := *internal.GetOpts()

	defer func() {
		opts := internal.GetOpts()
		*opts = originalOpts
	}()

	opts := internal.GetOpts()
	*opts = options

	testFunc()
}

func TestProcessLines(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		pattern string
		opts    internal.GrepOptions
		want    []internal.LineInfo
	}{
		{
			name:    "флаг -i (игнорирование регистра)",
			lines:   []string{"Hello", "WORLD", "hello", "World"},
			pattern: "hello",
			opts:    internal.GrepOptions{Ignore: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "Hello", Match: true},
				{Number: 2, Text: "WORLD", Match: false},
				{Number: 3, Text: "hello", Match: true},
				{Number: 4, Text: "World", Match: false},
			},
		},
		{
			name:    "флаг -v (инвертированный поиск)",
			lines:   []string{"hello", "world", "hello world", "test"},
			pattern: "hello",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello", Match: false},
				{Number: 2, Text: "world", Match: true},
				{Number: 3, Text: "hello world", Match: false},
				{Number: 4, Text: "test", Match: true},
			},
		},
		{
			name:    "флаг -i и -v вместе",
			lines:   []string{"Hello", "WORLD", "hello", "World"},
			pattern: "hello",
			opts:    internal.GrepOptions{Ignore: true, Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "Hello", Match: false},
				{Number: 2, Text: "WORLD", Match: true},
				{Number: 3, Text: "hello", Match: false},
				{Number: 4, Text: "World", Match: true},
			},
		},
		{
			name:    "флаг -F (фиксированная строка) - точка как символ",
			lines:   []string{"hello.world", "hello world", "hello"},
			pattern: "hello.world",
			opts:    internal.GrepOptions{Fix: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello.world", Match: true},
				{Number: 2, Text: "hello world", Match: false},
				{Number: 3, Text: "hello", Match: false},
			},
		},
		{
			name:    "флаг -F с игнорированием регистра",
			lines:   []string{"Hello.World", "HELLO.WORLD", "hello.world"},
			pattern: "hello.world",
			opts:    internal.GrepOptions{Fix: true, Ignore: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "Hello.World", Match: true},
				{Number: 2, Text: "HELLO.WORLD", Match: true},
				{Number: 3, Text: "hello.world", Match: true},
			},
		},
		{
			name:    "флаг -F с инвертированным поиском",
			lines:   []string{"hello.world", "test", "hello world"},
			pattern: "hello.world",
			opts:    internal.GrepOptions{Fix: true, Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello.world", Match: false},
				{Number: 2, Text: "test", Match: true},
				{Number: 3, Text: "hello world", Match: true},
			},
		},
		{
			name:    "все флаги вместе (-i -v -F)",
			lines:   []string{"Hello.World", "test", "HELLO.WORLD"},
			pattern: "hello.world",
			opts:    internal.GrepOptions{Fix: true, Ignore: true, Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "Hello.World", Match: false},
				{Number: 2, Text: "test", Match: true},
				{Number: 3, Text: "HELLO.WORLD", Match: false},
			},
		},
		{
			name:    "поиск с якорем начала строки",
			lines:   []string{"hello world", "world hello", "hello"},
			pattern: "^hello",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello world", Match: true},
				{Number: 2, Text: "world hello", Match: false},
				{Number: 3, Text: "hello", Match: true},
			},
		},
		{
			name:    "поиск с якорем конца строки",
			lines:   []string{"hello world", "world hello", "hello"},
			pattern: "hello$",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello world", Match: false},
				{Number: 2, Text: "world hello", Match: true},
				{Number: 3, Text: "hello", Match: true},
			},
		},
		{
			name:    "поиск с якорем начала и конца",
			lines:   []string{"hello", "hello world", "world hello"},
			pattern: "^hello$",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello", Match: true},
				{Number: 2, Text: "hello world", Match: false},
				{Number: 3, Text: "world hello", Match: false},
			},
		},
		{
			name:    "поиск цифр с инвертированием",
			lines:   []string{"abc123", "def", "456ghi", "test"},
			pattern: "[0-9]+",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "abc123", Match: false},
				{Number: 2, Text: "def", Match: true},
				{Number: 3, Text: "456ghi", Match: false},
				{Number: 4, Text: "test", Match: true},
			},
		},
		{
			name:    "поиск с квантификаторами",
			lines:   []string{"a", "aa", "aaa", "b", "aaaa"},
			pattern: "a+",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "a", Match: true},
				{Number: 2, Text: "aa", Match: true},
				{Number: 3, Text: "aaa", Match: true},
				{Number: 4, Text: "b", Match: false},
				{Number: 5, Text: "aaaa", Match: true},
			},
		},
		{
			name:    "поиск с квантификаторами и инвертированием",
			lines:   []string{"a", "aa", "b", "c"},
			pattern: "a+",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "a", Match: false},
				{Number: 2, Text: "aa", Match: false},
				{Number: 3, Text: "b", Match: true},
				{Number: 4, Text: "c", Match: true},
			},
		},
		{
			name:    "поиск альтернатив",
			lines:   []string{"hello", "world", "hello world", "test"},
			pattern: "hello|world",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello", Match: true},
				{Number: 2, Text: "world", Match: true},
				{Number: 3, Text: "hello world", Match: true},
				{Number: 4, Text: "test", Match: false},
			},
		},
		{
			name:    "поиск альтернатив с игнорированием регистра",
			lines:   []string{"Hello", "WORLD", "hello world", "test"},
			pattern: "hello|world",
			opts:    internal.GrepOptions{Ignore: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "Hello", Match: true},
				{Number: 2, Text: "WORLD", Match: true},
				{Number: 3, Text: "hello world", Match: true},
				{Number: 4, Text: "test", Match: false},
			},
		},
		{
			name:    "поиск групп",
			lines:   []string{"hello123", "world456", "test", "hello789"},
			pattern: "(hello|world)[0-9]+",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello123", Match: true},
				{Number: 2, Text: "world456", Match: true},
				{Number: 3, Text: "test", Match: false},
				{Number: 4, Text: "hello789", Match: true},
			},
		},
		{
			name:    "поиск групп с инвертированием",
			lines:   []string{"hello123", "world456", "test", "hello789"},
			pattern: "(hello|world)[0-9]+",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello123", Match: false},
				{Number: 2, Text: "world456", Match: false},
				{Number: 3, Text: "test", Match: true},
				{Number: 4, Text: "hello789", Match: false},
			},
		},
		{
			name:    "поиск специальных символов",
			lines:   []string{"hello.world", "hello+world", "hello*world", "hello.world"},
			pattern: "hello\\.world",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello.world", Match: true},
				{Number: 2, Text: "hello+world", Match: false},
				{Number: 3, Text: "hello*world", Match: false},
				{Number: 4, Text: "hello.world", Match: true},
			},
		},
		{
			name:    "поиск специальных символов с фиксированной строкой",
			lines:   []string{"hello.world", "hello+world", "hello*world"},
			pattern: "hello.world",
			opts:    internal.GrepOptions{Fix: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello.world", Match: true},
				{Number: 2, Text: "hello+world", Match: false},
				{Number: 3, Text: "hello*world", Match: false},
			},
		},
		{
			name:    "поиск пустой строки",
			lines:   []string{"", "hello", "", "world"},
			pattern: "^$",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "", Match: true},
				{Number: 2, Text: "hello", Match: false},
				{Number: 3, Text: "", Match: true},
				{Number: 4, Text: "world", Match: false},
			},
		},
		{
			name:    "поиск пустой строки с инвертированием",
			lines:   []string{"", "hello", "", "world"},
			pattern: "^$",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "", Match: false},
				{Number: 2, Text: "hello", Match: true},
				{Number: 3, Text: "", Match: false},
				{Number: 4, Text: "world", Match: true},
			},
		},
		{
			name:    "поиск с множественными совпадениями в строке",
			lines:   []string{"hello hello", "world", "hello world hello"},
			pattern: "hello",
			opts:    internal.GrepOptions{},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello hello", Match: true},
				{Number: 2, Text: "world", Match: false},
				{Number: 3, Text: "hello world hello", Match: true},
			},
		},
		{
			name:    "поиск с множественными совпадениями и инвертированием",
			lines:   []string{"hello hello", "world", "hello world hello"},
			pattern: "hello",
			opts:    internal.GrepOptions{Invert: true},
			want: []internal.LineInfo{
				{Number: 1, Text: "hello hello", Match: false},
				{Number: 2, Text: "world", Match: true},
				{Number: 3, Text: "hello world hello", Match: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withOptions(tt.opts, func() {
				var re *regexp.Regexp
				pattern := tt.pattern

				if tt.opts.Fix {
					pattern = regexp.QuoteMeta(tt.pattern)
				}

				if tt.opts.Ignore {
					re = regexp.MustCompile("(?i)" + pattern)
				} else {
					re = regexp.MustCompile(pattern)
				}

				opts := internal.GetOpts()

				result := internal.ProcessLines(tt.lines, re, opts)

				if len(result) != len(tt.want) {
					t.Errorf("ProcessLines() вернул %d строк, ожидалось %d", len(result), len(tt.want))
					return
				}

				for i, line := range result {
					expected := tt.want[i]

					if line.Number != expected.Number {
						t.Errorf("строка %d: номер строки = %d, ожидалось %d", i+1, line.Number, expected.Number)
					}

					if line.Text != expected.Text {
						t.Errorf("строка %d: текст = %q, ожидалось %q", i+1, line.Text, expected.Text)
					}

					if line.Match != expected.Match {
						t.Errorf("строка %d: совпадение = %v, ожидалось %v", i+1, line.Match, expected.Match)
					}
				}
			})
		})
	}
}
