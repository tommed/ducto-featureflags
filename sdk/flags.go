package sdk

// Flag represents a single feature flag definition
type Flag struct {
	Enabled *bool  `json:"enabled,omitempty"` // Static fallback/default
	Rules   []Rule `json:"rules,omitempty"`   // Optional targeting logic
	// Future: rollout %, conditions, etc.
}

type Rule struct {
	If    map[string]string `json:"if,omitempty"` // Simple AND-matching conditions
	Value bool              `json:"value"`        // Result if rule matches
}

type EvalContext map[string]string

// Evaluate performs rule-based or fallback evaluation
func (f Flag) Evaluate(ctx EvalContext) bool {
	// Check rules first
	for _, rule := range f.Rules {
		if ruleMatches(rule.If, ctx) {
			return rule.Value
		}
	}
	// Fallback to Enabled
	if f.Enabled != nil {
		return *f.Enabled
	}
	return false
}

func ruleMatches(conditions map[string]string, ctx EvalContext) bool {
	for k, v := range conditions {
		if ctx[k] != v {
			return false
		}
	}
	return true
}
