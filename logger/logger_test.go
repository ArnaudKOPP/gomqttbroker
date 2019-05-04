/* Copyright (c) 2019, Arnaud KOPP
 */
package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGet(t *testing.T) {
	var l *zap.Logger
	logger := Get()

	assert.NotNil(t, logger)
	assert.IsType(t, l, logger)
}

func TestNewDevLogger(t *testing.T) {
	logger, err := NewDevLogger()

	assert.Nil(t, err)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))
}

func TestNewProdLogger(t *testing.T) {
	logger, err := NewProdLogger()

	assert.Nil(t, err)
	assert.False(t, logger.Core().Enabled(zap.DebugLevel))
}
