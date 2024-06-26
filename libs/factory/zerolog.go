package factory

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

// Level defines log levels.
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	NoLevel
	// Disabled disables the logger.
	Disabled
	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

func DefaultLogger(lev Level, w io.Writer) *zerolog.Logger {
	if w == nil {
		w = os.Stdout
	}

	l := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.Level(lev))

	return &l

}

// SilentLogger returns a Zerologger that does not write any output
func SilentLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}
