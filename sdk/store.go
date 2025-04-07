package sdk

import (
	"encoding/json"
	"fmt"
	"os"
)

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
