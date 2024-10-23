package main

import (
	"context"
	"net/http"
	"os/signal"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	defer cancel()

	config := MustParseConfig()
	initLogger(config.LogLevel, config.LogPretty)

	log.Info().Interface("config", config).Msg("Starting")

	router := NewRouter(
		NewBackendController(config.ErrorFilesPath, config.DefaultFormat),
		NewMetricsController(),
		NewHealthzController(),
	)

	log.Ctx(ctx).Info().Msgf("Listening on http://localhost:%d", config.Port)

	err := Serve(ctx, config.Port, router)
	log.Err(err).Msg("Terminated")
}

func NewHealthzController() RouteMapper {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
	}
}
