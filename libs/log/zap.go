package log

import (
	"go.uber.org/zap"
)

// Zapper is a utility wrapper to Zap logger
type Zapper struct {
	Zap   *zap.Logger
	sugar *zap.SugaredLogger
}

// TODO set log level
// todo, do I want to recreate a logger for every library, or use a package variable
func NewZapper() (*Zapper, error) {
	z, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	l := Zapper{
		Zap:   z,
		sugar: z.Sugar(),
	}
	return &l, nil
}

//---  debug

func (z Zapper) Debug(msg string) {
	z.sugar.Debug(msg)
}
func (z Zapper) DebugF(template string, args ...interface{}) {
	z.sugar.Debugf(template, args...)
}
func (z Zapper) DebugW(msg string, kvs ...interface{}) {
	z.sugar.Debugw(msg, kvs...)
}

//---  info

func (z Zapper) Info(msg string) {
	z.sugar.Info(msg)
}
func (z Zapper) InfoF(template string, args ...interface{}) {
	z.sugar.Infof(template, args...)
}
func (z Zapper) InfoW(msg string, kvs ...interface{}) {
	z.sugar.Infow(msg, kvs...)
}

//---  warn

func (z Zapper) Warn(msg string) {
	z.sugar.Warn(msg)
}
func (z Zapper) WarnF(template string, args ...interface{}) {
	z.sugar.Warnf(template, args...)
}
func (z Zapper) WarnW(msg string, kvs ...interface{}) {
	z.sugar.Warnw(msg, kvs...)
}

//---  error

func (z Zapper) Error(msg string) {
	z.sugar.Error(msg)
}
func (z Zapper) ErrorF(template string, args ...interface{}) {
	z.sugar.Errorf(template, args...)
}
func (z Zapper) ErrorW(msg string, kvs ...interface{}) {
	z.sugar.Errorw(msg, kvs...)
}
