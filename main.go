package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// RequestId is a unique ID that identifies the request - same as for backend service
	RequestId = "X-Request-ID"
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	defer cancel()

	config := MustParseConfig()
	metricsController := NewMetricsController()
	backendController := NewBackendController(config.ErrorFilesPath, config.DefaultFormat, config.Debug)
	healthzController := NewHealthzController()

	router := NewRouter(backendController, metricsController, healthzController)

	log.Ctx(ctx).Info().Msgf("Starting server on http://localhost:%d", config.Port)
	err := Serve(ctx, config.Port, router)
	log.Err(err).Msg("Terminated")
}

func NewMetricsController() RouteMapper {
	return func(mux *http.ServeMux) {
		mux.Handle("GET /metrics", promhttp.Handler())
	}
}
func NewHealthzController() RouteMapper {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}
}

func NewBackendController(path, defaultFormat string, debug bool) RouteMapper {
	defaultExts, err := mime.ExtensionsByType(defaultFormat)
	if err != nil || len(defaultExts) == 0 {
		panic("couldn't get file extension for default format")
	}
	defaultExt := defaultExts[0]

	return func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(r.Context())
			start := time.Now()
			ext := defaultExt

			if debug {
				w.Header().Set(FormatHeader, r.Header.Get(FormatHeader))
				w.Header().Set(CodeHeader, r.Header.Get(CodeHeader))
				w.Header().Set(ContentType, r.Header.Get(ContentType))
				w.Header().Set(OriginalURI, r.Header.Get(OriginalURI))
				w.Header().Set(Namespace, r.Header.Get(Namespace))
				w.Header().Set(IngressName, r.Header.Get(IngressName))
				w.Header().Set(ServiceName, r.Header.Get(ServiceName))
				w.Header().Set(ServicePort, r.Header.Get(ServicePort))
				w.Header().Set(RequestId, r.Header.Get(RequestId))
			}

			format := r.Header.Get(FormatHeader)
			if format == "" {
				format = defaultFormat
				logger.Printf("format not specified. Using %v", format)
			}

			cext, err := mime.ExtensionsByType(format)
			if err != nil {
				logger.Printf("unexpected error reading media type extension: %v. Using %v", err, ext)
				format = defaultFormat
			} else if len(cext) == 0 {
				logger.Printf("couldn't get media type extension. Using %v", ext)
			} else {
				ext = cext[0]
			}
			w.Header().Set(ContentType, format)

			errCode := r.Header.Get(CodeHeader)
			code, err := strconv.Atoi(errCode)
			if err != nil {
				code = 404
				logger.Error().Err(err).Str("code", errCode).Msg("unexpected error reading return code")
			}
			w.WriteHeader(code)

			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			// special case for compatibility
			if ext == ".htm" {
				ext = ".html"
			}
			file := fmt.Sprintf("%v/%v%v", path, code, ext)
			f, err := os.Open(file)
			if err != nil {
				logger.Trace().Err(err).Msg("unexpected error opening file")
				scode := strconv.Itoa(code)
				file := fmt.Sprintf("%v/%cxx%v", path, scode[0], ext)
				f, err := os.Open(file)
				if err != nil {
					logger.Error().Err(err).Msg("unexpected error opening file")
					http.NotFound(w, r)
					return
				}
				defer f.Close()
				logger.Trace().Str("file", file).Int("code", code).Str("format", format).Msg("serving custom error response")
				_, _ = io.Copy(w, f)
				return
			}
			defer f.Close()
			logger.Trace().Str("file", file).Int("code", code).Str("format", format).Msg("serving custom error response")
			_, _ = io.Copy(w, f)

			duration := time.Now().Sub(start).Seconds()

			proto := strconv.Itoa(r.ProtoMajor)
			proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

			requestCount.WithLabelValues(proto).Inc()
			requestDuration.WithLabelValues(proto).Observe(duration)
		})
	}
}
