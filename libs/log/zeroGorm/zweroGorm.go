package zeroGorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"time"
)

// NewZero is an utility function to generate a new ZeroLogger with some defaults
func NewZero() *zerolog.Logger {
	l := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	return &l
}

// Logger is the Logger that implements the Gorm.logger interface but delegates logging to Zerolog
type Logger struct {
	log *zerolog.Logger
	Cfg
}

type Cfg struct {
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

// New returns a new instance of ZeroGorm
func New(zero *zerolog.Logger, cfg Cfg) *Logger {
	if cfg.SlowThreshold == 0 {
		cfg.SlowThreshold = 300 * time.Millisecond
	}
	return &Logger{
		log: zero,
		Cfg: cfg,
	}
}

// LogMode returns a new logger with a different log level, pert of the gorm interface
func (z *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {

	var desLevel zerolog.Level
	switch level {
	case gormLogger.Error:
		desLevel = zerolog.ErrorLevel
	case gormLogger.Warn:
		desLevel = zerolog.WarnLevel
	case gormLogger.Info:
		desLevel = zerolog.InfoLevel
	case gormLogger.Silent:
		desLevel = zerolog.Disabled
	}

	subLogger := z.log.Level(desLevel)
	sub := Logger{
		log: &subLogger,
	}
	return &sub
}

func (z Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	z.log.Info().Msg(fmt.Sprintf(msg, data...))
}
func (z Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	z.log.Warn().Msg(fmt.Sprintf(msg, data...))
}
func (z Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	z.log.Error().Msg(fmt.Sprintf(msg, data...))
}

func (z Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	elapsed := time.Since(begin)
	switch {
	case err != nil && z.log.GetLevel() <= zerolog.ErrorLevel && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !z.IgnoreRecordNotFoundError):
		query, rows := fc()
		entry := z.log.Error()
		entry.Str("query", query).Float64("query_duration", float64(elapsed.Nanoseconds())/1e6)
		if rows != -1 {
			entry.Int64("rows", rows)
		}
		entry.Msg(err.Error())

	case elapsed > z.SlowThreshold && z.SlowThreshold != 0 && z.log.GetLevel() <= zerolog.WarnLevel:
		query, rows := fc()
		entry := z.log.Warn()

		entry.Str("query", query).Float64("query_duration", float64(elapsed.Nanoseconds())/1e6)
		if rows != -1 {
			entry.Int64("rows", rows)
		}
		entry.Int64("query_threshold", int64(z.SlowThreshold)).Msg("slow query")

	case z.log.GetLevel() <= zerolog.DebugLevel:
		query, rows := fc()
		entry := z.log.Info()

		entry.Str("query", query).Float64("query_duration", float64(elapsed.Nanoseconds())/1e6)
		if rows != -1 {
			entry.Int64("rows", rows)
		}
		entry.Msg("")
	}
}
