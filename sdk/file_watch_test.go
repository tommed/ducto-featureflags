package sdk

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func writeTestFlags(t *testing.T, path string, flags string) {
	t.Helper()

	// Wait to ensure the watcher has started fully
	time.Sleep(200 * time.Millisecond)

	// Write to a temp file first
	tmp := path + ".tmp"
	err := os.WriteFile(tmp, []byte(flags), 0644)
	assert.NoError(t, err)

	// Rename atomically
	err = os.Rename(tmp, path)
	assert.NoError(t, err)

	// give fs layer time to register file (esp. on macOS)
	time.Sleep(100 * time.Millisecond)
}

func TestWatchingStore_ReloadsOnChange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	// Write initial file
	writeTestFlags(t, file, `{
		"flags": {
			"new_ui": { "enabled": false }
		}
	}`)

	store, err := NewFileWatchingStore(file)
	assert.NoError(t, err)
	defer store.Close()

	// Wait for watcher to be fully registered
	time.Sleep(200 * time.Millisecond)

	assert.False(t, store.IsEnabled("new_ui", EvalContext{}))

	// Modify file to flip flag to true
	writeTestFlags(t, file, `{
		"flags": {
			"new_ui": { "enabled": true }
		}
	}`)

	// Wait for reload (polling version)
	timeout := time.After(3 * time.Second)
	var seenChange bool

	for !seenChange {
		select {
		case <-store.Events():
			seenChange = true
		case <-timeout:
			t.Fatal("timeout waiting for fsnotify reload")
		}
	}

	assert.True(t, store.IsEnabled("new_ui", EvalContext{}))
}

func TestWatchingStore_HandlesBadFileGracefully(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	writeTestFlags(t, file, `{
		"flags": {
			"test": { "enabled": true }
		}
	}`)

	store, err := NewFileWatchingStore(file)
	assert.NoError(t, err)
	defer store.Close()

	// Break the file
	writeTestFlags(t, file, `{"flags": BROKEN_JSON`)

	select {
	case err := <-store.Errors():
		assert.Contains(t, err.Error(), "reload failed")
	case <-time.After(2 * time.Second):
		t.Fatal("expected reload error not received")
	}
}
