package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "successful logger creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			assert.NotNil(t, logger)
		})
	}
}

func TestWrapError(t *testing.T) {
	_ = NewLogger()

	tests := []struct {
		name        string
		inputErr    error
		inputMsg    string
		inputFields []zap.Field
		expectedLog string
		expectedErr string
	}{
		{
			name:        "basic error wrapping",
			inputErr:    errors.New("connection failed"),
			inputMsg:    "failed to connect",
			inputFields: []zap.Field{zap.String("db", "test")},
			expectedErr: "failed to connect: connection failed",
		},
		{
			name:        "nil error",
			inputErr:    nil,
			inputMsg:    "nil",
			inputFields: []zap.Field{},
			expectedErr: "nil: %!w(<nil>)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapError(tt.inputErr, tt.inputMsg, tt.inputFields...)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}
