package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func initLogger(level string, pretty bool) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
		log.Warn().Msgf("Invalid log level '%s', fallback to '%s'", level, logLevel.String())
	}

	if logLevel == zerolog.NoLevel {
		logLevel = zerolog.InfoLevel
	}

	var logWriter io.Writer = os.Stderr
	if pretty {
		logWriter = &zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}
	}

	zerolog.DurationFieldUnit = time.Millisecond
	logger := zerolog.New(logWriter).Level(logLevel).With().Timestamp().Logger()

	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
}

func NewLoggingMiddleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		metrics := httpsnoop.CaptureMetrics(next, writer, request)
		log.Ctx(request.Context()).Info().
			Str("path", request.URL.Path).
			Str("referer", request.Referer()).
			Dur("duration", metrics.Duration).
			Int("status_code", metrics.Code).
			Int64("response_size", metrics.Written).
			Msg("Handled request")

		proto := fmt.Sprintf("%d.%d", request.ProtoMajor, request.ProtoMinor)
		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(metrics.Duration.Seconds())
	}
}

func NewZerologRequestIdMiddleware(useNginxRequestId bool) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var requestId string
		if useNginxRequestId {
			requestId = r.Header.Get(RequestId)
		} else {
			requestId = xid.New().String()
		}

		logger := log.Ctx(r.Context()).With().Str("request_id", requestId).Logger()
		r = r.WithContext(logger.WithContext(r.Context()))

		next(w, r)
	}
}
