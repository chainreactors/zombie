package Core

import (
	"Zombie/src/Server"
	"Zombie/src/Utils"
	"fmt"
	"strings"
)

func Dispatch(CurTask Utils.ScanTask)(err error, result bool){
	ServerName := strings.ToUpper(CurTask.Server)




	switch ServerName {
	case "MYSQL":
		err, result = Server.MysqlConnect(CurTask.Username , CurTask.Password , CurTask.Info )
	case "POSTGRESQL":
		err, result = Server.PostgresConnect(CurTask.Username , CurTask.Password , CurTask.Info )
	case "REDIS":
		err, result = Server.RedisConnect(CurTask.Username , CurTask.Password , CurTask.Info )
	default:
		fmt.Println("The Server isn't supported")
	}
	return err, result
}

