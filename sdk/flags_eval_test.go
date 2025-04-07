package sdk

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLFlagFile_EvaluatesCorrectly(t *testing.T) {
	path := filepath.Join("..", "examples", "with_rules.yaml")

	store, err := NewStoreFromFile(path)
	assert.NoError(t, err)

	// Matching rule
	result := store.IsEnabled("new_ui", EvalContext{
		"env":   "prod",
		"group": "beta",
	})
	assert.True(t, result)

	// Second rule match
	result = store.IsEnabled("new_ui", EvalContext{
		"env": "prod",
	})
	assert.False(t, result)

	// No rule match, fallback
	result = store.IsEnabled("new_ui", EvalContext{
		"env": "dev",
	})
	assert.True(t, result)

	// Missing flag
	result = store.IsEnabled("not_there", EvalContext{})
	assert.False(t, result)
}

func TestFlagEvaluation_StaticOnly(t *testing.T) {
	f := Flag{Enabled: boolPtr(true)}
	assert.True(t, f.Evaluate(EvalContext{}))

	f = Flag{Enabled: boolPtr(false)}
	assert.False(t, f.Evaluate(EvalContext{}))
}

func TestFlagEvaluation_WithRules(t *testing.T) {
	f := Flag{
		Rules: []Rule{
			{If: map[string]string{"env": "prod"}, Value: true},
			{If: map[string]string{"env": "dev"}, Value: false},
		},
		Enabled: boolPtr(false), // fallback
	}

	assert.True(t, f.Evaluate(EvalContext{"env": "prod"}))
	assert.False(t, f.Evaluate(EvalContext{"env": "dev"}))
	assert.False(t, f.Evaluate(EvalContext{"env": "staging"}))
}

func TestStoreEvaluate(t *testing.T) {
	store, err := NewStoreFromBytesWithFormat([]byte(`{
		"flags": {
			"new_ui": {
				"rules": [
					{ "if": { "env": "prod" }, "value": true }
				],
				"enabled": false
			}
		}
	}`), "json")
	assert.NoError(t, err)

	assert.True(t, store.Evaluate("new_ui", EvalContext{"env": "prod"}))
	assert.False(t, store.Evaluate("new_ui", EvalContext{"env": "dev"}))
	assert.False(t, store.Evaluate("missing", EvalContext{"env": "prod"}))
}

func TestFlagEvaluation_FallbackToFalse(t *testing.T) {
	// No rules, no enabled
	f := Flag{}
	result := f.Evaluate(EvalContext{"env": "prod"})
	assert.False(t, result)

	// Rules don't match, and no enabled fallback
	f = Flag{
		Rules: []Rule{
			{If: map[string]string{"env": "qa"}, Value: true},
		},
	}
	result = f.Evaluate(EvalContext{"env": "prod"})
	assert.False(t, result)
}

func boolPtr(b bool) *bool {
	return &b
}
