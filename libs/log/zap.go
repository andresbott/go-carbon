package log

import (
	"go.uber.org/zap"
)

// Zapper is a utility wrapper to Zap logger
type Zapper struct {
	Zap *zap.Logger
}

// TODO set log level
// todo, do I want to recreate a logger for every library, or use a package variable
func NewZapper() (*Zapper, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	l := Zapper{
		Zap: z,
	}
	return &l, nil
}

type level int8

const (
	debug level = iota
	info
	warn
	err
)

func (z Zapper) write(level level, template string, args ...interface{}) {

	sugar := z.Zap.Sugar()

	switch level {
	case info:
		sugar.Infof(template, args)
		break
	case warn:
		sugar.Warnf(template, args)
		break
	case err:
		sugar.Errorf(template, args)
	case debug:
		sugar.Debugf(template, args)
	}
	z.Zap.Sync() // flushes buffer, if any

}

func (z Zapper) Debugf(template string, args ...interface{}) {
	z.write(debug, template, args)
}

func (z Zapper) Debug(msg string) {
	z.write(debug, msg)
}

func (z Zapper) Infof(template string, args ...interface{}) {
	z.write(info, template, args)
}
func (z Zapper) Info(msg string) {
	z.write(info, msg)
}

func (z Zapper) Warnf(template string, args ...interface{}) {
	z.write(info, template, args)
}
func (z Zapper) Warn(msg string) {
	z.write(warn, msg)
}

func (z Zapper) Errorf(template string, args ...interface{}) {
	z.write(info, template, args)
}
func (z Zapper) Error(msg string) {
	z.write(err, msg)
}
