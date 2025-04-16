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
