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
