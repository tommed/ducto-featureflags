package sdk

// Flag represents a single feature flag definition
type Flag struct {
	Enabled *bool  `json:"enabled,omitempty"` // Static fallback/default
	Rules   []Rule `json:"rules,omitempty"`   // Optional targeting logic
	// Future: rollout %, conditions, etc.
}

type Rule struct {
	If    map[string]string `json:"if,omitempty" yaml:"if,omitempty"`
	Value bool              `json:"value"`
	// Optional targeting
	Percent *int   `json:"percent,omitempty" yaml:"percent,omitempty"` // 0–100
	Seed    string `json:"seed,omitempty" yaml:"seed,omitempty"`       // key in context
}

type EvalContext map[string]string

// Evaluate performs rule-based or fallback evaluation
func (f Flag) Evaluate(ctx EvalContext) bool {
	// Check rules first
	for _, rule := range f.Rules {
		if ruleMatches(rule.If, ctx, rule) {
			return rule.Value
		}
	}
	// Fallback to Enabled
	if f.Enabled != nil {
		return *f.Enabled
	}
	return false
}

func ruleMatches(conditions map[string]string, ctx EvalContext, rule Rule) bool {
	for k, v := range conditions {
		if ctx[k] != v {
			return false
		}
	}
	// Handle percent rule (optional)
	if rule.Percent != nil {
		if *rule.Percent <= 0 {
			return false
		}
		seedKey := rule.Seed
		if seedKey == "" {
			return false
		}
		seedVal, ok := ctx[seedKey]
		if !ok {
			return false
		}
		percent := hashToPercent(seedVal)
		return percent < *rule.Percent
	}
	return true
}
