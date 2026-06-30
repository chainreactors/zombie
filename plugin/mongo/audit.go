package mongo

import (
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
	for _, db := range dbs {
		if _, sys := mongoSystemDBs[db]; sys {
			continue
		}
		colls, err := s.client.Database(db).ListCollectionNames(s.ctx, bson.D{})
		if err != nil {
			continue
		}
		for _, coll := range colls {
			location := db + "." + coll
			opts := options.Find().SetLimit(int64(limit))
			cursor, err := s.client.Database(db).Collection(coll).Find(s.ctx, bson.D{}, opts)
			if err != nil {
				continue
			}

			var docs []string
			for cursor.Next(s.ctx) {
				docs = append(docs, cursor.Current.String())
			}
			cursor.Close(s.ctx)

			if len(docs) > 0 {
				results[location] = strings.Join(docs, "\n")
			}
		}
	}
	return results, nil
}
