package core

import (
	"fmt"
	"strings"
)

func parseKeyValueArgs(values []string) (map[string]interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}
	parsed := make(map[string]interface{}, len(values))
	for _, value := range values {
		key, val, ok := strings.Cut(value, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid -V/-var value %q, expected key=value", value)
		}
		parsed[strings.TrimSpace(key)] = val
	}
	return parsed, nil
}

func parsePayloadArgs(values []string) (map[string]interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}
	grouped := make(map[string][]string, len(values))
	for _, value := range values {
		key, val, ok := strings.Cut(value, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid --payload value %q, expected key=value", value)
		}
		grouped[key] = append(grouped[key], val)
	}

	parsed := make(map[string]interface{}, len(grouped))
	for key, vals := range grouped {
		parsed[key] = vals
	}
	return parsed, nil
}
