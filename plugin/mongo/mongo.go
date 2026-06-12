package mongo

import (
	"context"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// mongoSession implements pkg.Session over a MongoDB client.
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

// MongoPlugin is stateless; all connection state lives in mongoSession.
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
	err = client.Ping(task.Context, nil)
	if err != nil {
		client.Disconnect(task.Context)
		return nil, err
	}

	return &mongoSession{service: task.Service, client: client, ctx: task.Context}, nil
}

func (p *MongoPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}
