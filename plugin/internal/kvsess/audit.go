package kvsess

import (
	"strings"
)

func (s *RedisSession) Audit(patterns []string, limit int) (map[string]string, error) {
	if limit <= 0 {
		limit = 100
	}

	var allKeys []string
	for _, p := range patterns {
		keys, err := s.Keys("*" + p + "*")
		if err != nil {
			continue
		}
		allKeys = append(allKeys, keys...)
	}

	seen := make(map[string]struct{})
	results := make(map[string]string)
	count := 0
	for _, key := range allKeys {
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		if count >= limit {
			break
		}
		val, err := s.Get(key)
		if err != nil || len(val) == 0 {
			continue
		}
		results[key] = strings.TrimSpace(string(val))
		count++
	}
	return results, nil
}
