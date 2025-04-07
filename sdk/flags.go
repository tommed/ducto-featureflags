package sdk

import (
	"encoding/json"
	"fmt"
	"os"
)

// Flag represents a single feature flag definition
type Flag struct {
	Enabled bool `json:"enabled"`
	// Future: rollout %, conditions, etc.
}

// Store holds the loaded feature flags
type Store struct {
	flags map[string]Flag
}

// NewStoreFromFile loads flags from a JSON file into memory
func NewStoreFromFile(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read flag file: %w", err)
	}
	return NewStoreFromBytes(data)
}

// NewStoreFromBytes allows loading from embedded JSON or remote fetch
func NewStoreFromBytes(data []byte) (*Store, error) {
	var parsed struct {
		Flags map[string]Flag `json:"flags"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("parse flag JSON: %w", err)
	}
	return &Store{flags: parsed.Flags}, nil
}

// IsEnabled returns true if the flag is defined and enabled
func (s *Store) IsEnabled(key string) bool {
	flag, ok := s.flags[key]
	return ok && flag.Enabled
}

// AllFlags returns the raw flag map
func (s *Store) AllFlags() map[string]Flag {
	return s.flags
}
