package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentRollout(t *testing.T) {
	ten := 10

	rule := Rule{
		Percent: &ten,
		Seed:    "user_id",
		Value:   true,
	}

	f := Flag{
		Rules:   []Rule{rule},
		Enabled: boolPtr(false), // fallback
	}

	match := 0
	total := 1000
	for i := 0; i < total; i++ {
		uid := fmt.Sprintf("user-%d", i)
		if f.Evaluate(EvalContext{"user_id": uid}) {
			match++
		}
	}

	t.Logf("Matched %d out of %d (~%.1f%%)", match, total, float64(match)/float64(total)*100)
	assert.Greater(t, match, 50) // Expect some rollout
	assert.Less(t, match, 150)   // Should stay near 10%
}

func TestPercentFallbackToHostname(t *testing.T) {
	percent := 100
	rule := Rule{
		Percent: &percent,
		Seed:    "HOSTNAME",
		Value:   true,
	}
	flag := Flag{
		Rules: []Rule{rule},
	}

	// Remove HOSTNAME from context
	ctx := EvalContext{}
	assert.Equal(t, true, flag.Evaluate(ctx)) // Should fallback to env
}
