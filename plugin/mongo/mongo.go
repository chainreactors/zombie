package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chainreactors/zombie/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSession struct {
	service string
	client  *mongo.Client
	ctx     context.Context
}

func (s *mongoSession) Service() string  { return s.service }
func (s *mongoSession) Raw() interface{} { return s.client }

func (s *mongoSession) Close() error {
	if s.client != nil {
		return s.client.Disconnect(s.ctx)
	}
	return nil
}

func (s *mongoSession) Query(query string, args ...any) ([][]string, error) {
	query = strings.TrimSpace(query)

	parts := strings.SplitN(query, " ", 2)
	cmd := parts[0]

	switch strings.ToLower(cmd) {
	case "show":
		if len(parts) > 1 && strings.HasPrefix(strings.ToLower(parts[1]), "db") {
			dbs, err := s.client.ListDatabaseNames(s.ctx, bson.D{})
			if err != nil {
				return nil, err
			}
			rows := [][]string{{"database"}}
			for _, db := range dbs {
				rows = append(rows, []string{db})
			}
			return rows, nil
		}
		if len(parts) > 1 && strings.HasPrefix(strings.ToLower(parts[1]), "collection") {
			dbAndRest := strings.TrimPrefix(strings.ToLower(parts[1]), "collections ")
			colls, err := s.client.Database(dbAndRest).ListCollectionNames(s.ctx, bson.D{})
			if err != nil {
				return nil, err
			}
			rows := [][]string{{"collection"}}
			for _, c := range colls {
				rows = append(rows, []string{c})
			}
			return rows, nil
		}
	}

	cmdDoc := bson.D{{Key: cmd, Value: 1}}
	if len(parts) > 1 {
		cmdDoc = append(cmdDoc, bson.E{Key: "arg", Value: parts[1]})
	}
	result := s.client.Database("admin").RunCommand(s.ctx, cmdDoc)
	if result.Err() != nil {
		return nil, result.Err()
	}
	raw, err := result.DecodeBytes()
	if err != nil {
		return nil, err
	}
	return [][]string{{"result"}, {raw.String()}}, nil
}

func init() {
	pkg.RegisterPlugin("mongo", &MongoPlugin{})
	pkg.Services.Register(&pkg.Service{Name: "mongo", DefaultPort: "27017", Alias: []string{"mongodb"}, Source: pkg.PluginSource})
}

type MongoPlugin struct{}

func (p *MongoPlugin) Name() string { return "mongo" }

func (p *MongoPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	var url string
	if task.Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", task.IP, task.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", task.Username, task.Password, task.IP, task.Port)
	}
	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(task.Timeout) * time.Second)

	client, err := mongo.Connect(task.Context, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(task.Context, nil); err != nil {
		client.Disconnect(task.Context)
		return nil, err
	}
	return &mongoSession{service: task.Service, client: client, ctx: task.Context}, nil
}

func (p *MongoPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
