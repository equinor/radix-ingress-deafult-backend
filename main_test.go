package main_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
import main "github.com/equinor/radix-ingress-default-backend"

func TestRun(t *testing.T) {
	router := main.NewRouter(main.NewBackendController())
	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Add("X-Namespace", "hello-world-dev")
	req.Header.Add("X-Code", "503")
	regularRequest, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer regularRequest.Body.Close()

	bytes, err := io.ReadAll(regularRequest.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bytes), "Server error - Radix")

	req, err = http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Add("X-Namespace", "equinor-web-sites-dev")
	req.Header.Add("X-Code", "503")
	equinorRequst, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer equinorRequst.Body.Close()

	bytes, err = io.ReadAll(equinorRequst.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bytes), "Something went wrong - Equinor")
}
