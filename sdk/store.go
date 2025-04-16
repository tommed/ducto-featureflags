package sdk

type AnyStore interface {
	Get(key string) (Flag, bool) // For OpenFeature compatibility
	AllFlags() map[string]Flag
}

// Store is a AnyStore which holds the provided flags and never updates.
type Store struct {
	flags map[string]Flag
}

func NewStore(flags map[string]Flag) AnyStore {
	return &Store{flags: flags}
}

// Get just returns the Flag now, so doesn't need EvalContext at this stage
func (s *Store) Get(key string) (Flag, bool) {
	f, ok := s.flags[key]
	return f, ok
}

// AllFlags returns the raw flag map
func (s *Store) AllFlags() map[string]Flag {
	return s.flags
}
