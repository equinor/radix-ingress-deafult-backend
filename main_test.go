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

	body, err := sendRequest(server.URL, "hello-world-dev", "503")
	require.NoError(t, err)
	assert.Contains(t, body, "Server error - Radix")

	body, err = sendRequest(server.URL, "equinor-web-sites-dev", "503")
	require.NoError(t, err)
	assert.Contains(t, body, "Something went wrong - Equinor")

	body, err = sendRequest(server.URL, "", "")
	require.NoError(t, err)
	assert.Contains(t, body, "Server error - Radix")

	body, err = sendRequest(server.URL+"/hello-world", "", "")
	require.NoError(t, err)
	assert.Contains(t, body, "Server error - Radix")

	body, err = sendRequest(server.URL+"/hello-world", "default", "503")
	require.NoError(t, err)
	assert.Contains(t, body, "Server error - Radix")
}

func sendRequest(url, namespace, code string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Namespace", namespace)
	req.Header.Add("X-Code", code)
	regularRequest, err := http.DefaultClient.Do(req)
	defer regularRequest.Body.Close()

	bytes, err := io.ReadAll(regularRequest.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
