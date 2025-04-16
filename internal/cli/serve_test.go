package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-featureflags/test"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func TestRunRoot_Serve_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests in short mode")
	}
	type args struct {
		extension string
		decoder   func(data []byte, v any) error
	}
	var tests = []struct {
		name string
		args args
	}{
		{
			name: "no extension",
			args: args{
				extension: "",
				decoder:   json.Unmarshal,
			},
		},
		{
			name: ".json extension",
			args: args{
				extension: ".json",
				decoder:   json.Unmarshal,
			},
		},
		{
			name: ".yaml extension",
			args: args{
				extension: ".yaml",
				decoder:   yaml.Unmarshal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			file := filepath.Join(dir, "flags.json")

			err := os.WriteFile(file, []byte(`{
				"my_flag": {
					"variants": `+test.BoolVariantsJSON()+`,
					"rules": [
						{ "if": { "env": "prod" }, "variant": "yes" }
					],
					"defaultVariant": "no"
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
			time.Sleep(100 * time.Millisecond)

			// Make a request
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/flags%s?key=my_flag&env=prod",
				port,
				tt.args.extension))
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]bool
			_ = tt.args.decoder(body, &result)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, true, result["value"])
		})
	}
}

//goland:noinspection GoUnhandledErrorResult
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
