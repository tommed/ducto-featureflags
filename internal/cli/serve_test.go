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
	dir := t.TempDir()
	file := filepath.Join(dir, "flags.json")

	err := os.WriteFile(file, []byte(`{
		"flags": {
			"my_flag": {
				"rules": [
					{ "if": { "env": "prod" }, "value": true }
				],
				"enabled": false
			}
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
