package sdk

type AnyStore interface {
	IsEnabled(key string, ctx EvalContext) bool
	AllFlags() map[string]Flag
}

// Store is a AnyStore which holds the provided flags and never updates.
type Store struct {
	flags map[string]Flag
}

func NewStore(flags map[string]Flag) AnyStore {
	return &Store{flags: flags}
}

// IsEnabled returns true if the flag is defined and enabled
func (s *Store) IsEnabled(key string, ctx EvalContext) bool {
	flag, ok := s.flags[key]
	if !ok {
		return false
	}
	return flag.Evaluate(ctx)
}

// AllFlags returns the raw flag map
func (s *Store) AllFlags() map[string]Flag {
	return s.flags
}
