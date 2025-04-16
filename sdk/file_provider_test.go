package sdk

import (
	"context"
	"github.com/tommed/ducto-featureflags/test"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	writeTestFlags(t, file, `{
		"new_ui": { "variants": `+test.BoolVariantsJSON()+`, "defaultVariant": "no" }
	}`)

	provider := NewFileProvider(file)
	store := NewDynamicStore(ctx, provider)
	err := store.Start()
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	flag, ok := store.Get("new_ui")
	assert.True(t, ok)

	_, val, ok, _ := flag.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, false, val)

	writeTestFlags(t, file, `{
		"new_ui": { "variants": `+test.BoolVariantsJSON()+`, "defaultVariant": "yes" }
	}`)

	time.Sleep(1 * time.Second)

	flag, ok = store.Get("new_ui")
	assert.True(t, ok)

	_, val, ok, _ = flag.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, true, val)
}

func TestWatchingStore_HandlesBadFileGracefully(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	writeTestFlags(t, file, `{
		"test": { "variants": `+test.BoolVariantsJSON()+`, "defaultVariant": "yes" }
	}`)

	provider := NewFileProvider(file)
	store := NewDynamicStore(ctx, provider)
	err := store.Start()
	assert.NoError(t, err)

	writeTestFlags(t, file, `{"flags": BROKEN_JSON`)

	// Should continue serving last good config
	flag, ok := store.Get("test")
	assert.True(t, ok)

	_, val, ok, _ := flag.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, true, val)
}
