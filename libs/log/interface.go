package log

type LeveledLogger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)

	LeveledFormattedLogger
	LeveledStructuredLogger
}

// LeveledFormattedLogger allows to send log messages with formatted fields in the template
type LeveledFormattedLogger interface {
	DebugF(template string, args ...interface{})
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
}

// LeveledStructuredLogger uses key value pairs as structured data to be logged
// pretty much inspired on the functionality of zapper
type LeveledStructuredLogger interface {
	DebugW(msg string, args ...interface{})
	InfoW(msg string, args ...interface{})
	WarnW(msg string, args ...interface{})
	ErrorW(msg string, args ...interface{})
}

// SilentLog is an empty implementation of LeveledLogger that will not produce any log output
type SilentLog struct{}

func (l SilentLog) Debug(msg string) {}
func (l SilentLog) Info(msg string)  {}
func (l SilentLog) Warn(msg string)  {}
func (l SilentLog) Error(msg string) {}

func (l SilentLog) DebugF(template string, args ...interface{}) {}
func (l SilentLog) InfoF(template string, args ...interface{})  {}
func (l SilentLog) WarnF(template string, args ...interface{})  {}
func (l SilentLog) ErrorF(template string, args ...interface{}) {}

func (l SilentLog) DebugW(template string, args ...interface{}) {}
func (l SilentLog) InfoW(template string, args ...interface{})  {}
func (l SilentLog) WarnW(template string, args ...interface{})  {}
func (l SilentLog) ErrorW(template string, args ...interface{}) {}
