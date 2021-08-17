package ExecAble

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/go-redis/redis"
	"regexp"
	"strings"
	"time"
)

type RedisService struct {
	Utils.IpInfo
	Username   string `json:"username"`
	Password   string `json:"password"`
	Additional string
	Input      string
}

func (s *RedisService) Query() bool {
	return false
}

func (s *RedisService) GetInfo() bool {
	return false
}

func (s *RedisService) Connect() bool {
	err, additional, res := RedisConnect(s.Username, s.Password, s.IpInfo)
	if err == nil && res {
		s.Additional = additional
		return true
	}
	return false

}

func (s *RedisService) DisConnect() bool {
	return false
}

func (s *RedisService) SetQuery(query string) {
	s.Input = query
}

func (s *RedisService) Output(res interface{}) {

}

func RedisConnect(User string, Password string, info Utils.IpInfo) (err error, additional string, result bool) {

	opt := redis.Options{Addr: fmt.Sprintf("%v:%v", info.Ip, info.Port),
		Password: Password, DB: 0, DialTimeout: time.Duration(Utils.Timeout) * time.Second}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err = client.Ping().Result()
	if err == nil {
		result = true
		redisinfo := client.Info().String()
		osreg, _ := regexp.Compile("os:(.*)\r\n")
		osname := osreg.FindString(redisinfo)
		getos, _ := regexp.Compile("(?i)linux")
		if getos.FindString(osname) != "" {
			additional = "os: linux\t"

			isroot := client.ConfigSet("dir", "/root/.ssh/").String()

			if strings.Contains(isroot, "Permission denied") {
				additional = "role: not root\t"
			} else if strings.Contains(isroot, "OK") {
				additional += "role: root\texsit /root/.ssh"
			} else {
				additional += "role: root\tdont have /root/.ssh"
			}

		} else {
			additional = "os: windows\t"
		}

	}
	return err, additional, result
}
