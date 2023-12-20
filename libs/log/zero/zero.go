package zero

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
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

type Zero struct {
	ZeroLog *zerolog.Logger
}

func NewZero(lev Level, w io.Writer) *Zero {
	if w == nil {
		w = os.Stdout
	}
	l := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.Level(lev))

	z := Zero{
		ZeroLog: &l,
	}
	return &z
}

func msgPayload(ev *zerolog.Event, msg string, kvs ...interface{}) {
	n := len(kvs)
	if n%2 == 0 {
		for i := 0; i < n; i = i + 2 {
			ev.Str(fmt.Sprintf("%v", kvs[i]), fmt.Sprintf("%v", kvs[i+1]))
		}
	} else {
		var payload []string
		for i := 0; i < n; i++ {
			payload = append(payload, fmt.Sprintf("%v", kvs[i]))
		}
		ev.Str("data", strings.Join(payload, ","))
	}
	ev.Msg(msg)
}

func (z Zero) Debug(msg string, kvs ...interface{}) {
	event := z.ZeroLog.Debug()
	msgPayload(event, msg, kvs...)
}

func (z Zero) Info(msg string, kvs ...interface{}) {
	event := z.ZeroLog.Info()
	msgPayload(event, msg, kvs...)
}

func (z Zero) Warn(msg string, kvs ...interface{}) {
	event := z.ZeroLog.Warn()
	msgPayload(event, msg, kvs...)
}

func (z Zero) Error(msg string, kvs ...interface{}) {
	event := z.ZeroLog.Error()
	msgPayload(event, msg, kvs...)
}
