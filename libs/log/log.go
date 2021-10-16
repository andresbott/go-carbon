package log

// boilerplate interface
type leveledLogger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}

type leveledLoggerf interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type silentLog struct{}

func (l silentLog) Debug(msg string) {}
func (l silentLog) Info(msg string)  {}
func (l silentLog) Warn(msg string)  {}
func (l silentLog) Error(msg string) {}

func (l silentLog) Debugf(template string, args ...interface{}) {}
func (l silentLog) Infof(template string, args ...interface{})  {}
func (l silentLog) Warnf(template string, args ...interface{})  {}
func (l silentLog) Errorf(template string, args ...interface{}) {}
