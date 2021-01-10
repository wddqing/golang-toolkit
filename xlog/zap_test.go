package xlog

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger(false)
	log.Debug("11111111", zap.Int("1", 1), zap.String("str", "name"))
}
