package sdk

import (
	"context"
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

	// Assemble
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	// Write initial file
	writeTestFlags(t, file, `{
		"new_ui": { "enabled": false }
	}`)

	provider := NewFileProvider(file)
	store := NewDynamicStore(ctx, provider)
	err := store.Start()
	assert.NoError(t, err)

	// Wait for watcher to be fully registered
	time.Sleep(200 * time.Millisecond)

	assert.False(t, store.IsEnabled("new_ui", EvalContext{}))

	// Modify file to flip flag to true
	writeTestFlags(t, file, `{
		"new_ui": { "enabled": true }
	}`)

	// Wait for reload (polling version)
	time.Sleep(1 * time.Second)
	assert.True(t, store.IsEnabled("new_ui", EvalContext{}))
}

func TestWatchingStore_HandlesBadFileGracefully(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// Assemble
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	writeTestFlags(t, file, `{
		"test": { "enabled": true }
	}`)

	provider := NewFileProvider(file)
	store := NewDynamicStore(ctx, provider)
	err := store.Start()
	assert.NoError(t, err)

	// Break the file
	writeTestFlags(t, file, `{"flags": BROKEN_JSON`)

	// Did not replace the flags with broken ones
	assert.True(t, store.IsEnabled("test", EvalContext{}))
}
