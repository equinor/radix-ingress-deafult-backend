package main_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
import main "github.com/equinor/radix-ingress-default-backend"

var testConfig = main.Config{Port: 65332, DefaultFormat: "text/html", ErrorFilesPath: "./www"}

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := main.Run(ctx, testConfig)
		require.NoError(t, err)
	}()
	time.Sleep(2 * time.Second)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:65332", nil)
	req.Header.Add("X-Namespace", "hello-world-dev")
	req.Header.Add("X-Code", "503")
	regularRequest, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer regularRequest.Body.Close()

	bytes, err := io.ReadAll(regularRequest.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bytes), "Server error - Radix")

	req, err = http.NewRequest(http.MethodGet, "http://localhost:65332", nil)
	req.Header.Add("X-Namespace", "equinor-web-sites-dev")
	req.Header.Add("X-Code", "503")
	equinorRequst, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer equinorRequst.Body.Close()

	bytes, err = io.ReadAll(equinorRequst.Body)
	require.NoError(t, err)
	assert.Contains(t, string(bytes), "Something went wrong - Equinor")
}
