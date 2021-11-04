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
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	l := Zapper{
		Zap:   z,
		sugar: z.Sugar(),
	}
	return &l, nil
}

func (z Zapper) Debugf(template string, args ...interface{}) {
	z.sugar.Debugf(template, args)
}
func (z Zapper) Debug(msg string) {
	z.sugar.Debug(msg)
}

func (z Zapper) Infof(template string, args ...interface{}) {
	z.sugar.Infof(template, args)
}
func (z Zapper) Info(msg string) {
	z.sugar.Info(msg)
}

func (z Zapper) Warnf(template string, args ...interface{}) {
	z.sugar.Warnf(template, args)
}
func (z Zapper) Warn(msg string) {
	z.sugar.Warn(msg)
}

func (z Zapper) Errorf(template string, args ...interface{}) {
	z.sugar.Errorf(template, args)
}
func (z Zapper) Error(msg string) {
	z.sugar.Error(msg)
}
