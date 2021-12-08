package log

// LeveledStructuredLogger uses key value pairs as structured data to be logged
// pretty much inspired on the functionality of zapper
type LeveledStructuredLogger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SilentLog is an empty implementation of LeveledLogger that will not produce any log output
type SilentLog struct{}

func (l SilentLog) Debug(msg string, args ...interface{}) {}
func (l SilentLog) Info(msg string, args ...interface{})  {}
func (l SilentLog) Warn(msg string, args ...interface{})  {}
func (l SilentLog) Error(msg string, args ...interface{}) {}
