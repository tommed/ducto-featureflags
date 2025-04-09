package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Store holds the loaded feature flags
type Store struct {
	flags map[string]Flag
}

func NewStore(flags map[string]Flag) *Store {
	return &Store{flags: flags}
}

// NewStoreFromFile loads flags from a JSON file into memory
func NewStoreFromFile(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read flag file: %w", err)
	}
	return NewStoreFromBytesWithFormat(data, detectFormat(path))
}

func detectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "json"
	}
}

// NewStoreFromBytesWithFormat allows loading from embedded YAML or JSON or remote fetch
func NewStoreFromBytesWithFormat(data []byte, format string) (*Store, error) {
	var parsed struct {
		Flags map[string]Flag `json:"flags"`
	}
	switch format {
	case "yaml":
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("parse YAML: %w", err)
		}
	default:
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}
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
