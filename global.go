package logging

import (
	"sync/atomic"

	"go.uber.org/zap"
)

var instance atomic.Pointer[zap.Logger]

func init() {
	instance.Store(zap.NewNop())
}

// SetLogger sets the current Logger instance for the logging system and returns the previous Logger instance.
// Parameter logger: A pointer to the zap.Logger instance to be set as the current Logger.
// Return value: Returns the previous zap.Logger instance.
func SetLogger(logger *zap.Logger) *zap.Logger {
	return instance.Swap(logger)
}

// GetLogger returns the global logger instance.
// This function takes no parameters.
// Return value: Returns a pointer to zap.Logger for subsequent logging operations.
func GetLogger() *zap.Logger {
	return instance.Load()
}