package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunRoot_Serve_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests in short mode")
	}
	dir := t.TempDir()
	file := filepath.Join(dir, "flags.json")

	err := os.WriteFile(file, []byte(`{
		"my_flag": {
			"rules": [
				{ "if": { "env": "prod" }, "value": true }
			],
			"enabled": false
		}
	}`), 0644)
	assert.NoError(t, err)

	// Use a unique port to avoid collisions
	port := "9173"
	go func() {
		stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
		_ = RunRoot([]string{"serve", "-file", file, "-addr", ":" + port}, stdout, stderr)
	}()

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)

	// Make a request
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/flags?key=my_flag&env=prod", port))
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]bool
	_ = json.Unmarshal(body, &result)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, result["enabled"])
}

func TestServe_WithToken_EnforcesAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests in short mode")
	}
	// Write a basic flag file
	dir := t.TempDir()
	file := filepath.Join(dir, "flags.json")
	err := os.WriteFile(file, []byte(`{
		"flags": {
			"secure_feature": { "enabled": true }
		}
	}`), 0644)
	assert.NoError(t, err)

	port := "9191"
	token := "secret-token"

	// Start server in background
	go Serve([]string{
		"-file", file,
		"-addr", ":" + port,
		"-token", token,
	}, io.Discard, io.Discard)

	time.Sleep(300 * time.Millisecond) // wait for server to bind

	url := "http://localhost:" + port + "/api/flags?key=secure_feature"

	// No token should 401
	resp, err := http.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// With token should 200
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp2, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}
