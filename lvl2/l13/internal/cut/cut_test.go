package cut

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []FieldRange
		wantErr bool
	}{
		{
			name:  "отдельные поля",
			input: "1,3,5",
			want: []FieldRange{
				{Start: 1, End: 1},
				{Start: 3, End: 3},
				{Start: 5, End: 5},
			},
		},
		{
			name:  "диапазоны",
			input: "1-3,5-7",
			want: []FieldRange{
				{Start: 1, End: 3},
				{Start: 5, End: 7},
			},
		},
		{
			name:  "смешанные поля и диапазоны",
			input: "1,3-5,7,9-10",
			want: []FieldRange{
				{Start: 1, End: 1},
				{Start: 3, End: 5},
				{Start: 7, End: 7},
				{Start: 9, End: 10},
			},
		},
		{
			name:    "пустая строка",
			input:   "",
			wantErr: true,
		},
		{
			name:    "неверный диапазон",
			input:   "1-",
			wantErr: true,
		},
		{
			name:    "неверный номер поля",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "отрицательный номер",
			input:   "-1",
			wantErr: true,
		},
		{
			name:    "неверный диапазон (начало > конца)",
			input:   "5-3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCutProcessor_ProcessLine(t *testing.T) {
	tests := []struct {
		name         string
		fields       string
		delimiter    string
		separated    bool
		input        string
		want         string
		shouldOutput bool
	}{
		{
			name:         "выбор отдельных полей",
			fields:       "1,3",
			delimiter:    "\t",
			input:        "a\tb\tc\td",
			want:         "a\tc",
			shouldOutput: true,
		},
		{
			name:         "выбор диапазона полей",
			fields:       "2-4",
			delimiter:    "\t",
			input:        "a\tb\tc\td\te",
			want:         "b\tc\td",
			shouldOutput: true,
		},
		{
			name:         "смешанные поля и диапазоны",
			fields:       "1,3-4",
			delimiter:    "\t",
			input:        "a\tb\tc\td\te",
			want:         "a\tc\td",
			shouldOutput: true,
		},
		{
			name:         "поле за границами",
			fields:       "1,10",
			delimiter:    "\t",
			input:        "a\tb\tc",
			want:         "a",
			shouldOutput: true,
		},
		{
			name:         "разделитель запятая",
			fields:       "1,3",
			delimiter:    ",",
			input:        "a,b,c,d",
			want:         "a,c",
			shouldOutput: true,
		},
		{
			name:         "флаг -s, строка без разделителя",
			fields:       "1,2",
			delimiter:    "\t",
			separated:    true,
			input:        "строка без разделителя",
			want:         "",
			shouldOutput: false,
		},
		{
			name:         "флаг -s, строка с разделителем",
			fields:       "1,2",
			delimiter:    "\t",
			separated:    true,
			input:        "a\tb\tc",
			want:         "a\tb",
			shouldOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := NewProcessor(tt.fields, tt.delimiter, tt.separated)
			assert.NoError(t, err)

			got, shouldOutput := processor.ProcessLine(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.shouldOutput, shouldOutput)
		})
	}
}

func TestCutProcessor(t *testing.T) {
	tests := []struct {
		name      string
		fields    string
		delimiter string
		separated bool
		wantErr   bool
	}{
		{
			name:      "корректные параметры",
			fields:    "1,3-5",
			delimiter: "\t",
			separated: false,
			wantErr:   false,
		},
		{
			name:      "неверные поля",
			fields:    "abc",
			delimiter: "\t",
			separated: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := NewProcessor(tt.fields, tt.delimiter, tt.separated)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, processor)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, processor)
			assert.Equal(t, tt.delimiter, processor.delimiter)
			assert.Equal(t, tt.separated, processor.separated)
		})
	}
}
