package sdk

import (
	"github.com/tommed/ducto-featureflags/test"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var boolVariants = map[string]interface{}{
	"on":  true,
	"off": false,
}

func TestYAMLFlagFile_EvaluatesCorrectly(t *testing.T) {
	path := filepath.Join("..", "examples", "04-with_rules.yaml")

	store, err := NewStoreFromFile(path)
	assert.NoError(t, err)

	flag, found := store.Get("new_ui")
	assert.True(t, found)

	// Matching rule
	variant, val, ok, _ := flag.Evaluate(EvalContext{
		"env": "beta",
	})
	assert.Equal(t, "beta", variant)
	assert.Equal(t, 2, val)
	assert.True(t, ok, "did not evaluate flag")

	// Second rule match
	variant, val, ok, _ = flag.Evaluate(EvalContext{
		"env": "prod",
	})
	assert.True(t, ok)
	assert.Equal(t, "stable", variant)
	assert.Equal(t, 4, val)

	// No rule match, fallback to enabled
	variant, val, ok, _ = flag.Evaluate(EvalContext{
		"env": "dev",
	})
	assert.True(t, ok)
	assert.Equal(t, "dev", variant)
	assert.Equal(t, 0, val)

	// Missing flag
	_, found = store.Get("not_there")
	assert.False(t, found)
}

func TestFlagEvaluation_StaticOnly(t *testing.T) {
	f := Flag{
		Variants:       boolVariants,
		DefaultVariant: "on",
	}
	_, val, ok, _ := f.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, true, val)

	f.DefaultVariant = "off"
	_, val, ok, _ = f.Evaluate(EvalContext{})
	assert.True(t, ok)
	assert.Equal(t, false, val)
}

func TestFlagEvaluation_WithRules(t *testing.T) {
	f := Flag{
		Variants: boolVariants,
		Rules: []VariantRule{
			{If: map[string]string{"env": "prod"}, Variant: "on"},
			{If: map[string]string{"env": "dev"}, Variant: "off"},
		},
		DefaultVariant: "off",
	}

	_, val, ok, _ := f.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, ok)
	assert.Equal(t, true, val)

	_, val, ok, _ = f.Evaluate(EvalContext{"env": "dev"})
	assert.True(t, ok)
	assert.Equal(t, false, val)

	_, val, ok, _ = f.Evaluate(EvalContext{"env": "staging"})
	assert.True(t, ok)
	assert.Equal(t, false, val) // fallback
}

func TestStoreEvaluate(t *testing.T) {
	store, err := NewStoreFromBytesWithFormat([]byte(`{
		"new_ui": {
			"variants": `+test.BoolVariantsJSON()+`,
			"rules": [
				{ "if": { "env": "prod" }, "variant": "yes" }
			],
			"defaultVariant": "no"
		}
	}`), "json")
	assert.NoError(t, err)

	flag, found := store.Get("new_ui")
	assert.True(t, found)

	_, val, ok, _ := flag.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, ok)
	assert.Equal(t, true, val)

	_, val, ok, _ = flag.Evaluate(EvalContext{"env": "dev"})
	assert.True(t, ok)
	assert.Equal(t, false, val)

	_, found = store.Get("missing")
	assert.False(t, found)
}

func TestFlagEvaluation_FallbackToFalse(t *testing.T) {
	// No rules, no enabled
	f := Flag{}
	_, val, ok, _ := f.Evaluate(EvalContext{"env": "prod"})
	assert.False(t, ok)
	assert.Nil(t, val)

	// Rules don't match, and no enabled fallback
	f = Flag{
		Variants:       boolVariants,
		DefaultVariant: "off",
		Rules: []VariantRule{
			{If: map[string]string{"env": "qa"}, Variant: "on"},
		},
	}
	_, val, ok, _ = f.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, ok)
	assert.Equal(t, false, val)
}
