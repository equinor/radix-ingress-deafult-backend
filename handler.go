package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
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

func NewBackendController(path, defaultFormat string) RouteMapper {
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

			logDebugHeaders(r)
			format := r.Header.Get(FormatHeader)
			if format == "" {
				format = defaultFormat
				logger.Warn().Msgf("format not specified. Using %v", format)
			}

			cext, err := mime.ExtensionsByType(format)
			if err != nil {
				logger.Error().Msgf("unexpected error reading media type extension: %v. Using %v", err, ext)
				format = defaultFormat
			} else if len(cext) == 0 {
				logger.Warn().Msgf("couldn't get media type extension. Using %v", ext)
			} else {
				ext = cext[0]
			}
			w.Header().Set(ContentType, format)

			errCode := r.Header.Get(CodeHeader)
			code, err := strconv.Atoi(errCode)
			if err != nil {
				code = 404
				logger.Warn().Err(err).Str("code", errCode).Msg("unexpected error reading return code")
			}
			scode := strconv.Itoa(code)
			w.WriteHeader(code)

			appNs := r.Header.Get(Namespace)
			appNameParts := strings.Split(appNs, "-")
			appName := strings.Join(appNameParts[:len(appNameParts)-1], "-")

			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			// special case for compatibility
			if ext == ".htm" {
				ext = ".html"
			}

			file, err := findFile(
				fmt.Sprintf("%s/%s/%d%s", path, appName, code, ext),
				fmt.Sprintf("%s/%s/%cxx%s", path, appName, scode[0], ext),
				fmt.Sprintf("%s/%d%s", path, code, ext),
				fmt.Sprintf("%s/%cxx%s", path, scode[0], ext),
			)
			if err != nil {
				logger.Trace().Err(err).Msg("no matching file found")
				http.NotFound(w, r)
				return
			}

			f, err := os.Open(file)
			if err != nil {
				logger.Trace().Err(err).Msg("unexpected error opening file")
				http.NotFound(w, r)
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

func findFile(files ...string) (string, error) {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}
	return "", os.ErrNotExist
}

func logDebugHeaders(r *http.Request) {
	zerolog.Ctx(r.Context()).Trace().
		Str("FormatHeader", r.Header.Get(FormatHeader)).
		Str("CodeHeader", r.Header.Get(CodeHeader)).
		Str("ContentType", r.Header.Get(ContentType)).
		Str("OriginalURI", r.Header.Get(OriginalURI)).
		Str("Namespace", r.Header.Get(Namespace)).
		Str("IngressName", r.Header.Get(IngressName)).
		Str("ServiceName", r.Header.Get(ServiceName)).
		Str("ServicePort", r.Header.Get(ServicePort)).
		Str("RequestId", r.Header.Get(RequestId)).
		Msg("Request headers")
}
