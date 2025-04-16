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
	result := flag.Evaluate(EvalContext{
		"env": "beta",
	})
	assert.Equal(t, "beta", result.Variant)
	assert.Equal(t, 2, result.Value)
	assert.True(t, result.OK, "did not evaluate flag")

	// Second rule match
	result = flag.Evaluate(EvalContext{
		"env": "prod",
	})
	assert.True(t, result.OK)
	assert.Equal(t, "stable", result.Variant)
	assert.Equal(t, 4, result.Value)

	// No rule match, fallback to enabled
	result = flag.Evaluate(EvalContext{
		"env": "dev",
	})
	assert.True(t, result.OK)
	assert.Equal(t, "dev", result.Variant)
	assert.Equal(t, 0, result.Value)

	// Missing flag
	_, found = store.Get("not_there")
	assert.False(t, found)
}

func TestFlagEvaluation_StaticOnly(t *testing.T) {
	f := Flag{
		Variants:       boolVariants,
		DefaultVariant: "on",
	}
	result := f.Evaluate(EvalContext{})
	assert.True(t, result.OK)
	assert.Equal(t, true, result.Value)

	f.DefaultVariant = "off"
	result = f.Evaluate(EvalContext{})
	assert.True(t, result.OK)
	assert.Equal(t, false, result.Value)
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

	result := f.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, result.OK)
	assert.Equal(t, true, result.Value)

	result = f.Evaluate(EvalContext{"env": "dev"})
	assert.True(t, result.OK)
	assert.Equal(t, false, result.Value)

	result = f.Evaluate(EvalContext{"env": "staging"})
	assert.True(t, result.OK)
	assert.Equal(t, false, result.Value) // fallback
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

	result := flag.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, result.OK)
	assert.Equal(t, true, result.Value)

	result = flag.Evaluate(EvalContext{"env": "dev"})
	assert.True(t, result.OK)
	assert.Equal(t, false, result.Value)

	_, found = store.Get("missing")
	assert.False(t, found)
}

func TestFlagEvaluation_FallbackToFalse(t *testing.T) {
	// No rules, no enabled
	f := Flag{}
	result := f.Evaluate(EvalContext{"env": "prod"})
	assert.False(t, result.OK)
	assert.Nil(t, result.Value)

	// Rules don't match, and no enabled fallback
	f = Flag{
		Variants:       boolVariants,
		DefaultVariant: "off",
		Rules: []VariantRule{
			{If: map[string]string{"env": "qa"}, Variant: "on"},
		},
	}
	result = f.Evaluate(EvalContext{"env": "prod"})
	assert.True(t, result.OK)
	assert.Equal(t, false, result.Value)
}
