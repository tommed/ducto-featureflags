package sdk

import (
	"github.com/tommed/ducto-featureflags/test"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStoreFromBytes(t *testing.T) {
	input := []byte(`{
		"feature_a": { "variants": ` + test.BoolVariantsJSON() + `, "defaultVariant": "yes" },
		"feature_b": { "variants": ` + test.BoolVariantsJSON() + `, "defaultVariant": "no" }
	}`)

	store, err := NewStoreFromBytesWithFormat(input, "json")
	assert.NoError(t, err)
	assert.NotNil(t, store)

	ctx := EvalContext{}

	// feature_a
	flag, ok := store.Get("feature_a")
	assert.True(t, ok)
	_, val, ok, _ := flag.Evaluate(ctx)
	assert.True(t, ok)
	assert.Equal(t, true, val)

	// feature_b
	flag, ok = store.Get("feature_b")
	assert.True(t, ok)
	_, val, ok, _ = flag.Evaluate(ctx)
	assert.True(t, ok)
	assert.Equal(t, false, val)

	// nonexistent
	_, ok = store.Get("nonexistent")
	assert.False(t, ok)
}

func TestNewStoreFromBytes_Invalid(t *testing.T) {
	bad := []byte(`{ "flags": "not-an-object" }`)
	_, err := NewStoreFromBytesWithFormat(bad, "json")
	assert.Error(t, err)
}

func TestNewStoreFromFile(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "flags.json")

	err := os.WriteFile(file, []byte(`{
		"screen_mode": {
			"defaultVariant": "dark",
			"variants": {
				"dark": true,
				"light": false
			}
		}
	}`), 0644)
	assert.NoError(t, err)

	store, err := NewStoreFromFile(file)
	assert.NoError(t, err)

	flag, ok := store.Get("screen_mode")
	assert.True(t, ok)
	_, val, ok, _ := flag.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, true, val)
}

func TestNewStoreFromFile_BadFile(t *testing.T) {
	_, err := NewStoreFromFile("invalid-file.json")
	assert.Error(t, err)
}

func TestAllFlags(t *testing.T) {
	store, err := NewStoreFromBytesWithFormat([]byte(`{
		"x": { "variants": `+test.BoolVariantsJSON()+`, "defaultVariant": "yes" },
		"y": { "variants": `+test.BoolVariantsJSON()+`, "defaultVariant": "no" }
	}`), "json")
	require.NoError(t, err)

	flags := store.AllFlags()
	require.Len(t, flags, 2)

	// Check x
	xFlag, xOk := flags["x"]
	assert.True(t, xOk)
	_, xVal, xOk, _ := xFlag.Evaluate(EvalContext{})
	assert.True(t, xOk)
	assert.Equal(t, true, xVal.(bool))

	// Check y
	yFlag, yOk := flags["y"]
	assert.True(t, yOk)
	_, yVal, yOk, _ := yFlag.Evaluate(EvalContext{})
	assert.True(t, yOk)
	assert.Equal(t, false, yVal.(bool))
}
