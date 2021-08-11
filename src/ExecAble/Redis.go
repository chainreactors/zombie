package ExecAble

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/go-redis/redis"
	"regexp"
	"strings"
	"time"
)

func RedisConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {

	opt := redis.Options{Addr: fmt.Sprintf("%v:%v", info.Ip, info.Port),
		Password: Password, DB: 0, DialTimeout: time.Duration(Utils.Timeout) * time.Second}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err = client.Ping().Result()
	if err == nil {
		result.Result = true
		redisinfo := client.Info().String()
		osreg, _ := regexp.Compile("os:(.*)\r\n")
		osname := osreg.FindString(redisinfo)
		getos, _ := regexp.Compile("(?i)linux")
		if getos.FindString(osname) != "" {
			result.Additional = "os: linux\t"

			isroot := client.ConfigSet("dir", "/root/.ssh/").String()

			if strings.Contains(isroot, "Permission denied") {
				result.Additional = "role: not root\t"
			} else if strings.Contains(isroot, "OK") {
				result.Additional += "role: root\texsit /root/.ssh"
			} else {
				result.Additional += "role: root\tdont have /root/.ssh"
			}

		} else {
			result.Additional = "os: windows\t"
		}

	}
	return err, result
}
