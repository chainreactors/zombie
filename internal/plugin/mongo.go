package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoService struct {
	*pkg.Task
	Input string
	conn  *mongo.Client
}

func (s *MongoService) Query() bool {
	return false
}

func (s *MongoService) GetInfo() bool {
	return false
}

func (s *MongoService) Connect() error {
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

func (s *MongoService) Close() error {
	if s.conn != nil {
		return s.conn.Disconnect(s.Context)
	}
	return nil
}

func (s *MongoService) SetQuery(query string) {
	s.Input = query
}

func (s *MongoService) Output(res interface{}) {

}
