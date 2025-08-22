package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type ZerologAdapter struct{}

func (a *ZerologAdapter) NewLogEntry(r *http.Request) middleware.LogEntry {
	l := ForRequest(r)
	return &chiLogEntry{logger: l, req: r}
}

type chiLogEntry struct {
	logger zerolog.Logger
	req    *http.Request
}

func (e *chiLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.logger.Info().
		Str("method", e.req.Method).
		Str("path", e.req.URL.Path).
		Str("remote_addr", e.req.RemoteAddr).
		Int("status", status).
		Int("bytes", bytes).
		Dur("elapsed", elapsed).
		Str("user_agent", e.req.UserAgent()).
		Msgf("HTTP %s %s - %d (%s) - %d bytes in %v",
			e.req.Method, e.req.URL.Path, status, http.StatusText(status), bytes, elapsed)
}

func (e *chiLogEntry) Panic(v interface{}, stack []byte) {
	e.logger.Error().
		Interface("panic", v).
		Bytes("stack", stack).
		Msg("panic recovered")
}
