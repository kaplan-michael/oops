package logger

import (
	"go.uber.org/zap"
)

// Init initializes a global Zap logger based on the provided log level.
// It replaces the global logger so that zap.L() and zap.S() can be used directly.
func Init(logLevel string) error {
	var logger *zap.Logger
	var err error
	if logLevel == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
