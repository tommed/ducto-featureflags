package sdk

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStoreFromBytes(t *testing.T) {
	input := []byte(`{
		"feature_a": { "enabled": true },
		"feature_b": { "enabled": false }
	}`)

	store, err := NewStoreFromBytesWithFormat(input, "json")
	assert.NoError(t, err)
	assert.NotNil(t, store)

	ctx := EvalContext{}
	assert.True(t, store.IsEnabled("feature_a", ctx))
	assert.False(t, store.IsEnabled("feature_b", ctx))
	assert.False(t, store.IsEnabled("nonexistent", ctx))
}

func TestNewStoreFromBytes_Invalid(t *testing.T) {
	bad := []byte(`{ "flags": "not-an-object" }`)
	_, err := NewStoreFromBytesWithFormat(bad, "json")
	assert.Error(t, err)
}

func TestNewStoreFromFile(t *testing.T) {
	tmp := t.TempDir()
	file := tmp + "/flags.json"

	err := os.WriteFile(file, []byte(`{
		"dark_mode": { "enabled": true }
	}`), 0644)
	assert.NoError(t, err)

	store, err := NewStoreFromFile(file)
	assert.NoError(t, err)
	assert.True(t, store.IsEnabled("dark_mode", EvalContext{}))
}

func TestNewStoreFromFile_BadFile(t *testing.T) {
	_, err := NewStoreFromFile("invalid-file.json")
	assert.Error(t, err)
}

func TestAllFlags(t *testing.T) {
	store, err := NewStoreFromBytesWithFormat([]byte(`{
		"x": { "enabled": true },
		"y": { "enabled": false }
	}`), "json")
	require.NoError(t, err)

	flags := store.AllFlags()
	require.Len(t, flags, 2)

	// Check x
	xFlag, xOk := flags["x"]
	assert.True(t, xOk)
	assert.NotNil(t, xFlag.Enabled)
	assert.True(t, *xFlag.Enabled)

	// Check y
	yFlag, yOk := flags["y"]
	assert.True(t, yOk)
	assert.NotNil(t, yFlag.Enabled)
	assert.False(t, *yFlag.Enabled)
}
