// File: ducto-featureflags/internal/cli/runner_test.go
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeTempFlags(t *testing.T, json string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "flags.json")
	err := os.WriteFile(path, []byte(json), 0644)
	assert.NoError(t, err)
	return path
}

func TestRun_QuerySingleFlag(t *testing.T) {
	flags := `{
		"flags": {
			"beta": { "enabled": true }
		}
	}`
	path := writeTempFlags(t, flags)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	code := Run([]string{"-file", path, "-key", "beta"}, stdout, stderr)

	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), `"beta" is true`)
}

func TestRun_MissingKeyAndList(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	code := Run([]string{"-file", "flags.json"}, stdout, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "must provide either -key or -list")
}

func TestRun_MissingFile(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	code := Run([]string{"-list", "-file", ""}, stdout, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "missing required flag: -file")
}

func TestRun_InvalidFlagFile(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	code := Run([]string{"-file", "nonexistent.json", "-key", "x"}, stdout, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "failed to load flags")
}

func TestRun_InvalidArgs(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Missing value for -key
	code := Run([]string{"-key"}, stdout, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "failed to parse args")
}

func TestRun_UnknownFlag(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	code := Run([]string{"-unknown"}, stdout, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "flag provided but not defined")
}

func TestRun_WithContextEvaluation(t *testing.T) {
	flags := `{
		"flags": {
			"canary_mode": {
				"rules": [
					{ "if": { "user_group": "beta" }, "value": true }
				],
				"enabled": false
			}
		}
	}`

	path := writeTempFlags(t, flags)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	code := Run([]string{
		"-file", path,
		"-key", "canary_mode",
		"--ctx", "user_group=beta",
	}, stdout, stderr)

	assert.Equal(t, 0, code)
	assert.Contains(t, stdout.String(), `"canary_mode" is true`)
}
