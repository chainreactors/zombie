package mongo

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoPlugin struct {
	*pkg.Task
	Input string
	conn  *mongo.Client
}

func (s *MongoPlugin) Unauth() (bool, error) {
	var err error
	var url string

	if s.Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", s.IP, s.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", "mongodbuser", s.Password, s.IP, s.Port)
	}
	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(s.Timeout) * time.Second)

	// 连接到MongoDB
	client, err := mongo.Connect(s.Context, clientOptions)
	if err != nil {
		return false, err
	}
	s.conn = client
	err = s.conn.Ping(s.Context, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *MongoPlugin) Name() string {
	return s.Service.String()
}

func (s *MongoPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *MongoPlugin) Login() error {
	var err error
	var url string

	if s.Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", s.IP, s.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", s.Username, s.Password, s.IP, s.Port)
	}
	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(s.Timeout) * time.Second)

	// 连接到MongoDB
	client, err := mongo.Connect(s.Context, clientOptions)
	if err != nil {
		return err
	}
	s.conn = client
	err = s.conn.Ping(s.Context, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Disconnect(s.Context)
	}
	return nil
}

//func (s *MongoPlugin) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *MongoPlugin) Output(res interface{}) {
//
//}
