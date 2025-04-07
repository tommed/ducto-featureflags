package sdk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStoreFromBytes(t *testing.T) {
	input := []byte(`{
		"flags": {
			"feature_a": { "enabled": true },
			"feature_b": { "enabled": false }
		}
	}`)

	store, err := NewStoreFromBytes(input)
	assert.NoError(t, err)
	assert.NotNil(t, store)

	assert.True(t, store.IsEnabled("feature_a"))
	assert.False(t, store.IsEnabled("feature_b"))
	assert.False(t, store.IsEnabled("nonexistent"))
}

func TestNewStoreFromBytes_Invalid(t *testing.T) {
	bad := []byte(`{ "flags": "not-an-object" }`)
	_, err := NewStoreFromBytes(bad)
	assert.Error(t, err)
}

func TestNewStoreFromFile(t *testing.T) {
	tmp := t.TempDir()
	file := tmp + "/flags.json"

	err := os.WriteFile(file, []byte(`{
		"flags": {
			"dark_mode": { "enabled": true }
		}
	}`), 0644)
	assert.NoError(t, err)

	store, err := NewStoreFromFile(file)
	assert.NoError(t, err)
	assert.True(t, store.IsEnabled("dark_mode"))
}

func TestNewStoreFromFile_BadFile(t *testing.T) {
	_, err := NewStoreFromFile("invalid-file.json")
	assert.Error(t, err)
}

func TestAllFlags(t *testing.T) {
	store, err := NewStoreFromBytes([]byte(`{
		"flags": {
			"x": { "enabled": true },
			"y": { "enabled": false }
		}
	}`))
	assert.NoError(t, err)

	flags := store.AllFlags()
	assert.Len(t, flags, 2)
	assert.True(t, flags["x"].Enabled)
	assert.False(t, flags["y"].Enabled)
}
