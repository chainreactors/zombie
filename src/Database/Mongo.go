package Database

import (
	"Zombie/src/Utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func MongoConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, client *mongo.Client) {
	var url string

	if Password == "" {
		url = fmt.Sprintf("mongodb://%v:%v", info.Ip, info.Port)
	} else {
		url = fmt.Sprintf("mongodb://%v:%v@%v:%v", User, Password, info.Ip, info.Port)
	}
	clientOptions := options.Client().ApplyURI(url).SetConnectTimeout(time.Duration(Utils.Timeout) * time.Second)

	// 连接到MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	//defer client.Disconnect(context.TODO())
	if err != nil {
		result = false
	}

	return err, result, client
}

func MongoConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, client := MongoConnect(User, Password, info)
	defer client.Disconnect(context.TODO())
	if err == nil {
		err = client.Ping(context.TODO(), nil)
		if err == nil {
			res = true
			result.Result = res
		}
	}

	return err, result
}
