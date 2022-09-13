package plugin

import (
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/go-redis/redis"
	"time"
)

type RedisService struct {
	*utils.Task
	conn       *redis.Client
	Additional string
	Input      string
}

func (s *RedisService) Query() bool {
	return false
}

func (s *RedisService) GetInfo() bool {
	return false
}

func (s *RedisService) Connect() error {
	conn, err := RedisConnect(s.Task)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil

}

func (s *RedisService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return NilConnError{s.Service}
}

func (s *RedisService) SetQuery(query string) {
	s.Input = query
}

func (s *RedisService) Output(res interface{}) {

}

func RedisConnect(task *utils.Task) (client *redis.Client, err error) {
	opt := redis.Options{Addr: task.Address(),
		Password: task.Password, DB: 0, DialTimeout: time.Duration(utils.Timeout) * time.Second}
	client = redis.NewClient(&opt)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	//if err == nil {
	//	result = true
	//	redisinfo := client.Info().String()
	//	osreg, _ := regexp.Compile("os:(.*)\r\n")
	//	osname := osreg.FindString(redisinfo)
	//	getos, _ := regexp.Compile("(?i)linux")
	//	if getos.FindString(osname) != "" {
	//		additional = "os: linux\t"
	//
	//		isroot := client.ConfigSet("dir", "/root/.ssh/").String()
	//
	//		if strings.Contains(isroot, "Permission denied") {
	//			additional = "role: not root\t"
	//		} else if strings.Contains(isroot, "OK") {
	//			additional += "role: root\texsit /root/.ssh"
	//		} else {
	//			additional += "role: root\tdont have /root/.ssh"
	//		}
	//
	//	} else {
	//		additional = "os: windows\t"
	//	}
	//}
	return client, nil
}
