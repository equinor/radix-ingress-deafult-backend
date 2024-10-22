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

	err := Run(ctx, config)

	log.Err(err).Msg("Terminated")
}

func Run(ctx context.Context, config Config) error {
	router := NewRouter(
		NewBackendController(config.ErrorFilesPath, config.DefaultFormat),
		NewMetricsController(),
		NewHealthzController(),
	)

	log.Ctx(ctx).Info().Msgf("Listening on http://localhost:%d", config.Port)
	return Serve(ctx, config.Port, router)
}

func NewHealthzController() RouteMapper {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
	}
}
