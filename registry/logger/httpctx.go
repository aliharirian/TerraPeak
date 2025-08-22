package logger

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func ForRequest(r *http.Request) zerolog.Logger {
	l := With().Logger()
	if r == nil {
		return l
	}
	if id := middleware.GetReqID(r.Context()); id != "" {
		l = l.With().Str("request_id", id).Logger()
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		l = l.With().Str("remote_ip", rip).Logger()
	}
	return l
}
