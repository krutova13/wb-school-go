package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"telnet/internal/types"
)

func TestCommandLineParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *types.Config
		wantErr bool
	}{
		{
			name: "valid arguments with default timeout",
			args: []string{"localhost", "8080"},
			want: &types.Config{
				Host:    "localhost",
				Port:    "8080",
				Timeout: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid arguments with custom timeout",
			args: []string{"--timeout=5s", "localhost", "8080"},
			want: &types.Config{
				Host:    "localhost",
				Port:    "8080",
				Timeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "missing host and port",
			args:    []string{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing port",
			args:    []string{"localhost"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty host",
			args:    []string{"", "8080"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty port",
			args:    []string{"localhost", ""},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewCommandLineParser()
			got, err := p.Parse(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Host, got.Host)
				assert.Equal(t, tt.want.Port, got.Port)
				assert.Equal(t, tt.want.Timeout, got.Timeout)
			}
		})
	}
}
