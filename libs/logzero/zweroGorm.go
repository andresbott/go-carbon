package logzero

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

// ZeroGorm is the Logger that implements the Gorm.logger interface but delegates logging to Zerolog
type ZeroGorm struct {
	log zerolog.Logger
	Cfg
}

type Cfg struct {
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

// NewZeroGorm returns a new instance of ZeroGorm
func NewZeroGorm(zero zerolog.Logger, cfg Cfg) ZeroGorm {
	if cfg.SlowThreshold == 0 {
		cfg.SlowThreshold = 300 * time.Millisecond
	}
	return ZeroGorm{
		log: zero,
		Cfg: cfg,
	}
}

// LogMode returns a new logger with a different log level, pert of the gorm interface
func (z ZeroGorm) LogMode(level gormLogger.LogLevel) gormLogger.Interface {

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
	sub := ZeroGorm{
		log: subLogger,
	}
	return &sub
}

func (z ZeroGorm) Info(ctx context.Context, msg string, data ...interface{}) {
	z.log.Info().Msg(fmt.Sprintf(msg, data...))
}
func (z ZeroGorm) Warn(ctx context.Context, msg string, data ...interface{}) {
	z.log.Warn().Msg(fmt.Sprintf(msg, data...))
}
func (z ZeroGorm) Error(ctx context.Context, msg string, data ...interface{}) {
	z.log.Error().Msg(fmt.Sprintf(msg, data...))
}

func (z ZeroGorm) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

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
