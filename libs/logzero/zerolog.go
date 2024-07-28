package logzero

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
	"time"
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

func GetLogLevel(in string) Level {
	in = strings.ToUpper(in)
	switch in {
	case "DEBUG":
		return DebugLevel
	case "WARN":
		return WarnLevel
	case "ERROR", "ERR":
		return ErrorLevel
	case "DISABLED", "OFF":
		return Disabled
	default:
		return InfoLevel

	}
}

func DefaultLogger(lev Level, w io.Writer) *zerolog.Logger {
	if w == nil {
		w = os.Stdout
	}
	l := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.Level(lev))
	return &l

}

func humanDur(i interface{}) string {
	switch d := i.(type) {
	case time.Duration:
		return d.String()
	case float64:
		return time.Duration(d * float64(time.Second)).String()
	default:
		return i.(string)
	}
}

func ConsoleFileOutput(file string) (zerolog.LevelWriter, error) {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	if file != "" {
		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		return zerolog.MultiLevelWriter(consoleWriter, f), nil
	}
	return zerolog.MultiLevelWriter(consoleWriter), nil
}

// SilentLogger returns a Zerologger that does not write any output
func SilentLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}
