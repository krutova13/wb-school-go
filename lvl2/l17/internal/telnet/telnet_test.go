package telnet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"telnet/internal/types"
)

func TestNewClient(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}

	client := NewClient(config)

	assert.NotNil(t, client)
	assert.Equal(t, config, client.config)
	assert.NotNil(t, client.connManager)
	assert.NotNil(t, client.logger)
	assert.NotNil(t, client.errorHandler)
}

func TestClient_Connect_WithoutConnection(t *testing.T) {
	config := &types.Config{
		Host:    "invalid-host",
		Port:    "99999",
		Timeout: 1 * time.Millisecond,
	}

	client := NewClient(config)

	err := client.Connect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
}

func TestClient_Start_WithoutConnection(t *testing.T) {
	config := &types.Config{
		Host:    "localhost",
		Port:    "8080",
		Timeout: 10 * time.Second,
	}

	client := NewClient(config)

	err := client.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection is not established")
}
