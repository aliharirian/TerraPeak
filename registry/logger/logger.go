package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

var base zerolog.Logger

// Init configures the global logger once (call from main()).
func Init(app string, w io.Writer, level string, timeFormat string) {
	if w == nil {
		w = os.Stdout
	}
	if timeFormat == "" {
		timeFormat = "15:04:05.0000T2006-01-02"
	}
	zerolog.TimeFieldFormat = timeFormat

	lvl := zerolog.InfoLevel
	switch level {
	case "debug":
		lvl = zerolog.DebugLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	case "fatal":
		lvl = zerolog.FatalLevel
	}

	zerolog.SetGlobalLevel(lvl)

	base = zerolog.New(w).With().Timestamp().Str("app", app).Logger()
	log.Logger = base
}

// With returns a child logger with extra fields (structured logging).
// Usage: logger.With().Str("ns", ns).Str("name", name).Infof("provider requested")
func With() zerolog.Context {
	return base.With()
}

func Debugf(msg string, args ...any) { base.Debug().Msgf(msg, args...) }
func Infof(msg string, args ...any)  { base.Info().Msgf(msg, args...) }
func Warnf(msg string, args ...any)  { base.Warn().Msgf(msg, args...) }
func Errorf(msg string, args ...any) { base.Error().Msgf(msg, args...) }
func Fatalf(msg string, args ...any) { base.Fatal().Msgf(msg, args...) }
