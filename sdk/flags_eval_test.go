package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	store, err := NewStoreFromBytes([]byte(`{
		"flags": {
			"new_ui": {
				"rules": [
					{ "if": { "env": "prod" }, "value": true }
				],
				"enabled": false
			}
		}
	}`))
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
