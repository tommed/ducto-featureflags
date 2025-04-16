package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentRollout(t *testing.T) {
	ten := 10

	rule := VariantRule{
		Percent: &ten,
		Seed:    "user_id",
		Variant: "on",
	}

	f := Flag{
		Variants:       boolVariants,
		Rules:          []VariantRule{rule},
		DefaultVariant: "off",
	}

	match := 0
	total := 1000
	for i := 0; i < total; i++ {
		uid := fmt.Sprintf("user-%d", i)
		result := f.Evaluate(EvalContext{"user_id": uid})
		if result.OK && result.Value == true {
			match++
		}
	}

	t.Logf("Matched %d out of %d (~%.1f%%)", match, total, float64(match)/float64(total)*100)
	assert.Greater(t, match, 50) // Expect some rollout
	assert.Less(t, match, 150)   // Should stay near 10%
}

func TestPercentWithSeedHashSHA256(t *testing.T) {
	percent := 100
	flag := Flag{
		Variants: boolVariants,
		Rules: []VariantRule{{
			Percent:  &percent,
			Seed:     "user_id",
			SeedHash: "sha256",
			Variant:  "on",
		}},
		DefaultVariant: "off",
	}

	ctx := EvalContext{"user_id": "abc123"}

	result := flag.Evaluate(ctx)
	assert.True(t, result.OK)
	assert.Equal(t, true, result.Value) // 100% should always be true

	// Now test 50% rollout
	percent = 50
	flag.Rules[0].Percent = &percent

	result = flag.Evaluate(ctx)
	t.Logf("Result for abc123 (sha256): %v", result.Value)
	assert.True(t, result.OK) // still should resolve, just might be false depending on hash
}

func TestPercentFallbackToHostname(t *testing.T) {
	percent := 100
	rule := VariantRule{
		Percent: &percent,
		Seed:    "HOSTNAME",
		Variant: "on",
	}

	flag := Flag{
		Variants:       boolVariants,
		Rules:          []VariantRule{rule},
		DefaultVariant: "off",
	}

	// No HOSTNAME in context — should fall back to env-hostname
	ctx := EvalContext{}
	result := flag.Evaluate(ctx)

	assert.True(t, result.OK)
	assert.Equal(t, true, result.Value)
}
