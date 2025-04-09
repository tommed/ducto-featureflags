package sdk

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// NewStoreFromFile loads flags from a JSON file into memory
func NewStoreFromFile(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read flag file: %w", err)
	}
	return NewStoreFromBytesWithFormat(data, DetectFormat(path))
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

func DetectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "json"
	}
}
