package sdk

type EvalContext map[string]string

type EvaluationResult struct {
	Variant string
	Value   interface{}
	OK      bool
	Matched bool
}

// Evaluate performs rule-based or fallback evaluation
// Evaluate resolves the flag to its chosen variant value
// File: sdk/flag.go or sdk/eval.go (your call)
// Flag.Evaluate now returns (variant, value, ok, matched)
func (f Flag) Evaluate(ctx EvalContext) EvaluationResult {
	for _, rule := range f.Rules {
		if ruleMatches(rule, ctx) {
			if rule.Variant == "" {
				return EvaluationResult{Variant: "", OK: false, Matched: true}
			}

			v, found := f.Variants[rule.Variant]
			if !found {
				return EvaluationResult{Variant: rule.Variant, OK: false, Matched: true}
			}

			return EvaluationResult{
				Variant: rule.Variant,
				Value:   v,
				OK:      true,
				Matched: true,
			}
		}
	}

	// fallback to default
	v, found := f.Variants[f.DefaultVariant]
	if !found {
		return EvaluationResult{Variant: f.DefaultVariant, OK: false}
	}

	return EvaluationResult{
		Variant: f.DefaultVariant,
		Value:   v,
		OK:      true,
		Matched: false,
	}
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
