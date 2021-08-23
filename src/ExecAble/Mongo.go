package ExecAble

//
//import (
//	"Zombie/src/Utils"
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"time"
//)
//
//type MongoService struct {
//	Utils.IpInfo
//	Username string `json:"username"`
//	Password string `json:"password"`
//	Input    string
//}
//
//func (s *MongoService) Query() bool {
//	return false
//}
//
//func (s *MongoService) GetInfo() bool {
//	return false
//}
//
//func (s *MongoService) Connect() bool {
//
//	err, res, client := MongoConnect(s.Username, s.Password, s.IpInfo)
//
//	if err == nil {
//		err = client.Ping(context.TODO(), nil)
//		//由于没有设计查询用途所以就直接关闭了
//		defer client.Disconnect(context.TODO())
//		if err == nil && res {
//
//			return true
//		}
//	}
//
//	return false
//
//}
//
//func (s *MongoService) DisConnect() bool {
//	return false
//}
//
//func (s *MongoService) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *MongoService) Output(res interface{}) {
//
//}
//
//func MongoConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, client *mongo.Client) {
//	var url string
//
//	if Password == "" {
//		url = fmt.Sprintf("mongodb://%v:%v", info.Ip, info.Port)
//	} else {
//		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", User, Password, info.Ip, info.Port)
//	}
//	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(Utils.Timeout) * time.Second)
//
//	// 连接到MongoDB
//	client, err = mongo.Connect(context.TODO(), clientOptions)
//	//defer client.Disconnect(context.TODO())
//	if err != nil {
//		result = false
//	}
//
//	return err, true, client
//}
