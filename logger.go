package convoy_go

import (
	"io"

	log "github.com/frain-dev/convoy/pkg/log"
)

// Level represents a log level.
type Level int32

const (
	// FatalLevel is used for undesired and unexpected events that
	// the program cannot recover from.
	FatalLevel Level = iota

	// ErrorLevel is used for undesired and unexpected events that
	// the program can recover from.
	ErrorLevel

	// WarnLevel is used for undesired but relatively expected events,
	// which may indicate a problem.
	WarnLevel

	// InfoLevel is used for general informational log messages.
	InfoLevel

	// DebugLevel is the lowest level of logging.
	// Debug logs are intended for debugging and development purposes.
	DebugLevel
)

// iLogger is a stripped down logging interface that
// clients can implement.
type iLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

func NewLogger(out io.Writer, lvl Level) *log.Logger {
	logger := log.NewLogger(out)
	logger.SetLevel(log.Level(lvl))

	return logger
}
