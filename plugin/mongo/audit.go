package mongo

import (
	"fmt"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoSystemDBs = map[string]struct{}{
	"admin": {}, "local": {}, "config": {},
}

func (s *mongoSession) Audit(patterns []string, limit int) (map[string]string, error) {
	if limit <= 0 {
		limit = 100
	}

	dbs, err := s.client.ListDatabaseNames(s.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	remaining := limit
	done := false
	for _, db := range dbs {
		if done {
			break
		}
		if _, sys := mongoSystemDBs[db]; sys {
			continue
		}
		colls, err := s.client.Database(db).ListCollectionNames(s.ctx, bson.D{})
		if err != nil {
			continue
		}
		for _, coll := range colls {
			if remaining == 0 {
				done = true
				break
			}
			location := db + "." + coll
			opts := options.Find().SetLimit(int64(remaining))
			cursor, err := s.client.Database(db).Collection(coll).Find(s.ctx, bson.D{}, opts)
			if err != nil {
				continue
			}

			var docs []string
			for cursor.Next(s.ctx) {
				fields, err := matchingDocumentFields(cursor.Current, patterns)
				if err != nil || len(fields) == 0 {
					continue
				}
				docs = append(docs, strings.Join(fields, "\n"))
				remaining--
				if remaining == 0 {
					break
				}
			}
			_ = cursor.Close(s.ctx)

			if len(docs) > 0 {
				results[location] = strings.Join(docs, "\n")
			}
		}
	}
	return results, nil
}

func matchingDocumentFields(raw bson.Raw, patterns []string) ([]string, error) {
	normalized := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		if pattern = strings.ToLower(strings.TrimSpace(pattern)); pattern != "" {
			normalized = append(normalized, pattern)
		}
	}
	if len(normalized) == 0 {
		return nil, nil
	}

	var document bson.M
	if err := bson.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	var fields []string
	collectMatchingFields(document, "", normalized, &fields)
	return fields, nil
}

func collectMatchingFields(value interface{}, prefix string, patterns []string, fields *[]string) {
	switch typed := value.(type) {
	case bson.M:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			if matchesFieldName(key, patterns) {
				*fields = append(*fields, fmt.Sprintf("%s: %v", path, typed[key]))
				continue
			}
			collectMatchingFields(typed[key], path, patterns, fields)
		}
	case bson.A:
		for _, item := range typed {
			collectMatchingFields(item, prefix, patterns, fields)
		}
	case []interface{}:
		for _, item := range typed {
			collectMatchingFields(item, prefix, patterns, fields)
		}
	}
}

func matchesFieldName(field string, patterns []string) bool {
	field = strings.ToLower(field)
	for _, pattern := range patterns {
		if strings.Contains(field, pattern) {
			return true
		}
	}
	return false
}
