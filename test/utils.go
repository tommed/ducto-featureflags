package test

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

func BoolVariants() map[string]interface{} {
	return map[string]interface{}{
		"yes": true,
		"no":  false,
	}
}

func BoolVariantsJSON() string {
	return Encode(BoolVariants(), "json")
}

func Encode(input any, format string) string {
	switch format {
	case "yaml":
		data, _ := yaml.Marshal(input)
		return string(data)
	default:
		data, _ := json.Marshal(input)
		return string(data)
	}
}
