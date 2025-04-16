package openfeature

import (
	"context"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-featureflags/sdk"
	"testing"
	"time"
)

func TestProviderIntegration(t *testing.T) {
	store := sdk.NewStore(map[string]sdk.Flag{
		"my_flag": {
			DefaultVariant: "greet",
			Variants:       map[string]interface{}{"greet": "hello world", "farewell": "goodbye world"},
		},
	})
	err := openfeature.SetProvider(NewProvider(store))
	assert.NoError(t, err)

	client := openfeature.NewClient("test")
	evalCtx := openfeature.NewEvaluationContext("user-123", map[string]interface{}{
		"env": "prod",
	})

	time.Sleep(1 * time.Millisecond)
	val, err := client.StringValue(context.Background(), "my_flag", "fallback", evalCtx)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", val)
}

func makeTestProvider(flags string) openfeature.FeatureProvider {
	store, err := sdk.NewStoreFromBytesWithFormat([]byte(flags), "json")
	if err != nil {
		panic(err)
	}
	return NewProvider(store)
}

func withContext(attrs map[string]interface{}) context.Context {
	evalCtx := openfeature.NewEvaluationContext("tester", attrs)
	return openfeature.WithTransactionContext(context.Background(), evalCtx)
}

func TestBooleanEvaluation(t *testing.T) {
	provider := makeTestProvider(`{
		"beta_enabled": {
			"defaultVariant": "on",
			"variants": {
				"on": true,
				"off": false
			},
			"rules": [
				{ "if": { "group": "beta" }, "variant": "on" }
			]
		}
	}`)

	ctx := map[string]interface{}{"group": "beta"}
	detail := provider.BooleanEvaluation(context.Background(), "beta_enabled", false, ctx)

	assert.Equal(t, true, detail.Value)
	assert.Equal(t, "on", detail.Variant)
	assert.Equal(t, openfeature.TargetingMatchReason, detail.Reason)
}

func TestBooleanEvaluation_DefaultFallback(t *testing.T) {
	provider := makeTestProvider(`{
		"dark_mode": {
			"defaultVariant": "off",
			"variants": {
				"on": true,
				"off": false
			}
		}
	}`)
	var ctx map[string]interface{} = nil
	detail := provider.BooleanEvaluation(context.Background(), "dark_mode", true, ctx)

	assert.Equal(t, false, detail.Value)
	assert.Equal(t, "off", detail.Variant)
	assert.Equal(t, openfeature.DefaultReason, detail.Reason)
}

func TestStringEvaluation(t *testing.T) {
	provider := makeTestProvider(`{
		"color": {
			"defaultVariant": "red",
			"variants": {
				"red": "red",
				"blue": "blue"
			},
			"rules": [
				{ "if": { "user": "a" }, "variant": "blue" }
			]
		}
	}`)

	ctx := map[string]interface{}{"user": "a"}
	detail := provider.StringEvaluation(context.Background(), "color", "fallback", ctx)
	assert.Equal(t, "blue", detail.Value)
	assert.Equal(t, openfeature.TargetingMatchReason, detail.Reason)
}

func TestIntEvaluation(t *testing.T) {
	provider := makeTestProvider(`{
		"max_limit": {
			"defaultVariant": "low",
			"variants": {
				"low": 5,
				"high": 10
			},
			"rules": [
				{ "if": { "env": "prod" }, "variant": "high" }
			]
		}
	}`)

	ctx := map[string]interface{}{"env": "prod"}
	detail := provider.IntEvaluation(context.Background(), "max_limit", 0, ctx)
	assert.Equal(t, int64(10), detail.Value)
	assert.Equal(t, openfeature.TargetingMatchReason, detail.Reason)
}

func TestFloatEvaluation(t *testing.T) {
	provider := makeTestProvider(`{
		"threshold": {
			"defaultVariant": "low",
			"variants": {
				"low": 0.5,
				"high": 1.0
			},
			"rules": [
				{ "if": { "env": "staging" }, "variant": "high" }
			]
		}
	}`)

	ctx := map[string]interface{}{"env": "staging"}
	detail := provider.FloatEvaluation(context.Background(), "threshold", 0.0, ctx)
	assert.Equal(t, 1.0, detail.Value)
	assert.Equal(t, openfeature.TargetingMatchReason, detail.Reason)
}

func TestObjectEvaluation(t *testing.T) {
	provider := makeTestProvider(`{
		"profile": {
			"defaultVariant": "basic",
			"variants": {
				"basic": { "mode": "limited", "features": ["read"] },
				"pro": { "mode": "full", "features": ["read", "write"] }
			},
			"rules": [
				{ "if": { "plan": "pro" }, "variant": "pro" }
			]
		}
	}`)

	ctx := map[string]interface{}{"plan": "pro"}
	detail := provider.ObjectEvaluation(context.Background(), "profile", nil, ctx)
	assert.Equal(t, map[string]interface{}{
		"mode":     "full",
		"features": []interface{}{"read", "write"},
	}, detail.Value)
	assert.Equal(t, openfeature.TargetingMatchReason, detail.Reason)
}
