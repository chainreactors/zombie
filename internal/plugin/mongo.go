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
	s.conn, err = MongoConnect(s.Task)
	if err != nil {
		return err
	}
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

func MongoConnect(info *pkg.Task) (client *mongo.Client, err error) {
	var url string

	if info.Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", info.IP, info.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", info.Username, info.Password, info.IP, info.Port)
	}
	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(info.Timeout) * time.Second)

	// 连接到MongoDB
	client, err = mongo.Connect(info.Context, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
