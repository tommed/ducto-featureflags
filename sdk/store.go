package sdk

// Store holds the loaded feature flags
type Store struct {
	flags map[string]Flag
}

func NewStore(flags map[string]Flag) *Store {
	return &Store{flags: flags}
}

// Evaluate supports rule-based matching
func (s *Store) Evaluate(key string, ctx EvalContext) bool {
	flag, ok := s.flags[key]
	if !ok {
		return false
	}
	return flag.Evaluate(ctx)
}

// IsEnabled returns true if the flag is defined and enabled
func (s *Store) IsEnabled(key string, ctx EvalContext) bool {
	return s.Evaluate(key, ctx)
}

// AllFlags returns the raw flag map
func (s *Store) AllFlags() map[string]Flag {
	return s.flags
}
