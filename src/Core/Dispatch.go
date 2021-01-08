package Core

import (
	"Zombie/src/Server"
	"Zombie/src/Utils"
	"fmt"
)

var ScanSum = 0

func Dispatch(CurTask Utils.ScanTask) (err error, result bool) {

	switch CurTask.Server {
	case "MYSQL":
		err, result = Server.MysqlConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "POSTGRESQL":
		err, result = Server.PostgresConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "REDIS":
		err, result = Server.RedisConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MONGO":
		err, result = Server.MongoConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MSSQL":
		err, result = Server.MssqlConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	default:
		fmt.Println("The Server isn't supported")
	}

	ScanSum += 1

	return err, result
}
