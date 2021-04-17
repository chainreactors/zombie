package Core

import (
	"Zombie/src/Server"
	"Zombie/src/Utils"
	"fmt"
)

var ScanSum = 0

func BruteDispatch(CurTask Utils.ScanTask) (err error, result Utils.BruteRes) {

	CurTask = Utils.UpdatePass(CurTask)

	switch CurTask.Server {
	case "MYSQL":
		err, result = Server.MysqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "POSTGRESQL":
		err, result = Server.PostgresConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "REDIS":
		err, result = Server.RedisConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "SSH":
		err, result = Server.SSHConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MONGO":
		err, result = Server.MongoConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MSSQL":
		err, result = Server.MssqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "VNC":
		err, result = Server.VNCConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "SMB":
		err, result = Server.SMBConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "ES":
		err, result = Server.EsConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "FTP":
		err, result = Server.FtpConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "TOMCAT":
		err, result = Server.TomcatConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	default:
		fmt.Println("The Server isn't supported")
	}

	ScanSum += 1

	return err, result
}

func ExecDispatch(CurTask Utils.ScanTask, Query string) (err error, Qresult []map[string]string, Columns []string) {
	switch CurTask.Server {
	case "MYSQL":
		err, Qresult, Columns = Server.MysqlQuery(CurTask.Username, CurTask.Password, CurTask.Info, Query)
	case "POSTGRESQL":
		err, Qresult, Columns = Server.PostgresQuery(CurTask.Username, CurTask.Password, CurTask.Info, Query)
	case "MSSQL":
		err, Qresult, Columns = Server.MssqlQuery(CurTask.Username, CurTask.Password, CurTask.Info, Query)
	default:
		fmt.Println("The Server isn't supported")
	}

	return err, Qresult, Columns
}
