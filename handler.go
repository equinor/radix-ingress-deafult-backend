package main

import (
	"embed"
	_ "embed"
	"net/http"
	"strconv"
	"strings"

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

//go:embed www/*.html
var www embed.FS

func NewBackendController() RouteMapper {

	return func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(r.Context())

			logDebugHeaders(r)

			errCode := r.Header.Get(CodeHeader)
			code, err := strconv.Atoi(errCode)
			if err != nil {
				code = 404
				logger.Warn().Err(err).Str("code", errCode).Msg("unexpected error reading return code")
			}
			scode := strconv.Itoa(code)
			w.WriteHeader(code)

			format := r.Header.Get(FormatHeader)
			accepts := r.Header.Get("Accept")
			acceptsHtml := strings.Contains(format, "text/html") || format == "" ||
				strings.Contains(accepts, "text/html") ||
				strings.Contains(accepts, "*/*")

			w.Header().Set("X-Content-Type-Options", "nosniff")

			if !acceptsHtml {
				w.Header().Set(ContentType, "text/plain")
				_, _ = w.Write([]byte(scode))
				logger.Warn().Msgf("format %s is unknown (not unspecified or html)", format)
				return
			}

			w.Header().Set(ContentType, "text/html; charset=utf-8")
			filename := getHtmlFile(r.Header.Get(Namespace))
			htmlFile, err := www.ReadFile("www/" + filename)
			if err != nil {
				filename = "default.html"
				htmlFile, err = www.ReadFile("www/" + filename)

				if err != nil {
					logger.Error().Err(err).Msg("unable to serve error page")
					http.NotFound(w, r)
					return
				}
			}

			logger.Trace().Str("file", filename).Msg("serving custom error response")
			_, _ = w.Write(htmlFile)
		})
	}
}

func getHtmlFile(namespace string) string {
	if namespace == "" {
		return "default.html"
	}

	appNameParts := strings.Split(namespace, "-")
	if len(appNameParts) < 2 {
		return "default.html"
	}

	appName := strings.Join(appNameParts[:len(appNameParts)-1], "-")
	return appName + ".html"
}

func logDebugHeaders(r *http.Request) {
	headers := r.Header.Clone()
	if headers.Get("Authorization") != "" {
		headers.Set("Authorization", "***removed***")
	}
	if headers.Get("Cookie") != "" {
		headers.Set("Cookie", "***removed***")
	}
	if headers.Get("Cookie2") != "" {
		headers.Set("Cookie2", "***removed***")
	}

	zerolog.Ctx(r.Context()).Trace().Interface("headers", headers).Msg("Request headers")
}
