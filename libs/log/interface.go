package log

type LeveledLogger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type SilentLog struct{}

func (l SilentLog) Debug(msg string) {}
func (l SilentLog) Info(msg string)  {}
func (l SilentLog) Warn(msg string)  {}
func (l SilentLog) Error(msg string) {}

func (l SilentLog) Debugf(template string, args ...interface{}) {}
func (l SilentLog) Infof(template string, args ...interface{})  {}
func (l SilentLog) Warnf(template string, args ...interface{})  {}
func (l SilentLog) Errorf(template string, args ...interface{}) {}
