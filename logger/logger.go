package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/ip812/blog/config"
	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func New(cfg *config.Config) Logger {
	var log zerolog.Logger
	if cfg.App.Env == "local" {
		log = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger().Level(zerolog.InfoLevel)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().UTC()
		}
		return &LoggerLocal{
			Log: &log,
		}
	}

	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	return &LoggerProd{
		Log: &log,
	}
}

// Logger implementation for local development (plain text)
type LoggerLocal struct {
	Log *zerolog.Logger
}

func (l *LoggerLocal) Debug(msg string, args ...interface{}) {
	l.Log.Debug().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerLocal) Info(msg string, args ...interface{}) {
	l.Log.Info().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerLocal) Warn(msg string, args ...interface{}) {
	l.Log.Warn().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerLocal) Error(msg string, args ...interface{}) {
	l.Log.Error().Msg(fmt.Sprintf(msg, args...))
}

// Logger implementation for prod (json)
type LoggerProd struct {
	Log *zerolog.Logger
}

func (l *LoggerProd) Debug(msg string, args ...interface{}) {
	l.Log.Debug().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerProd) Info(msg string, args ...interface{}) {
	l.Log.Info().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerProd) Warn(msg string, args ...interface{}) {
	l.Log.Warn().Msg(fmt.Sprintf(msg, args...))
}

func (l *LoggerProd) Error(msg string, args ...interface{}) {
	l.Log.Error().Msg(fmt.Sprintf(msg, args...))
}
