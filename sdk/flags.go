package sdk

// Flag represents a single feature flag definition
type Flag struct {
	Disabled       bool                   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	DefaultVariant string                 `json:"defaultVariant" yaml:"defaultVariant"`
	Variants       map[string]interface{} `json:"variants" yaml:"variants"`
	Rules          []VariantRule          `json:"rules,omitempty" yaml:"rules,omitempty"`
}

// VariantRule is our v2 rule which is OpenFeature compatible and uses 'variants'
type VariantRule struct {
	If       map[string]string `json:"if,omitempty" yaml:"if,omitempty"`
	Variant  string            `json:"variant" yaml:"variant"` // name of the variant to use
	Percent  *int              `json:"percent,omitempty" yaml:"percent,omitempty"`
	Seed     string            `json:"seed,omitempty" yaml:"seed,omitempty"`
	SeedHash string            `json:"seed_hash,omitempty" yaml:"seed_hash,omitempty"` // optional: "sha256"
}

type EvalContext map[string]string

// Evaluate performs rule-based or fallback evaluation
// Evaluate resolves the flag to its chosen variant value
// File: sdk/flag.go or sdk/eval.go (your call)
// Flag.Evaluate now returns (variant, value, ok, matched)
func (f Flag) Evaluate(ctx EvalContext) (string, interface{}, bool, bool) {
	for _, rule := range f.Rules {
		if ruleMatches(rule, ctx) {
			// resolve variant name and value
			if rule.Variant != "" {
				val, ok := f.Variants[rule.Variant]
				if !ok {
					return rule.Variant, nil, false, true // matched rule, but invalid variant
				}
				return rule.Variant, val, true, true // success, from rule
			}
			// fallback to simple value (no named variant)
			return "", rule.Variant, true, true
		}
	}

	// fallback to default variant
	val, ok := f.Variants[f.DefaultVariant]
	if !ok {
		return f.DefaultVariant, nil, false, false // fallback failed
	}
	return f.DefaultVariant, val, true, false // success, from default
}

func ruleMatches(rule VariantRule, ctx EvalContext) bool {
	// Match conditions
	for k, v := range rule.If {
		if ctx[k] != v {
			return false
		}
	}

	// Match percent rollout (optional)
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
			// Fallback to hostname if seed is "HOSTNAME"
			if seedKey == "HOSTNAME" {
				seedVal = getHostname()
				if seedVal == "" {
					return false
				}
			} else {
				return false
			}
		}

		percent := hashToPercent(seedVal, rule.SeedHash)
		return percent < *rule.Percent
	}

	return true
}
