package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func ApplyConfig(logPath, logLevel string) error {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("logger cannot open log file: %w", err)
	}
	Logger = zerolog.New(file).With().Timestamp().Logger().Level(level(logLevel))

	return nil
}

func level(lev string) zerolog.Level {
	switch lev {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "nolevel":
		return zerolog.NoLevel
	case "disabled":
		return zerolog.Disabled
	case "trace":
		return zerolog.TraceLevel
	}

	return zerolog.InfoLevel
}
