package Server

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func RedisConnect(User string, Password string,info Utils.IpInfo)(err error,result bool){

	opt := redis.Options{Addr: fmt.Sprintf("%v:%v", info.Ip, info.Port),
		Password: Password, DB: 0, DialTimeout: 2 * time.Second}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err = client.Ping().Result()
	if err == nil {
		result = true
	}
	return err, result
}