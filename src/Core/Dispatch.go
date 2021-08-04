package Core

import (
	"Zombie/src/Database"
	"Zombie/src/Protocol"
	"Zombie/src/Utils"
	"Zombie/src/Web"
	"fmt"
)

var ScanSum = 0

func BruteDispatch(CurTask Utils.ScanTask) (err error, result Utils.BruteRes) {

	CurTask = Utils.UpdatePass(CurTask)

	switch CurTask.Server {

	case "POSTGRESQL":
		err, result = Database.PostgresConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MYSQL":
		err, result = Database.MysqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "REDIS":
		err, result = Database.RedisConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "SSH":
		err, result = Protocol.SSHConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MONGO":
		err, result = Database.MongoConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "MSSQL":
		err, result = Database.MssqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "VNC":
		err, result = Database.VNCConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "SMB":
		err, result = Protocol.SMBConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "ES":
		err, result = Web.EsConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "FTP":
		err, result = Protocol.FtpConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	case "SNMP":
		err, result = Database.SnmpConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
	case "TOMCAT":
		err, result = Web.TomcatConnect(CurTask.Username, CurTask.Password, CurTask.Info)
	default:
		fmt.Println("The Database isn't supported")
	}

	ScanSum += 1

	return err, result
}

func ExecDispatch(CurTask Utils.ScanTask) Database.SqlHandle {
	switch CurTask.Server {
	case "POSTGRESQL":
		return &Database.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &Database.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "MYSQL":
		return &Database.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "SNMP":
		return &Database.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	default:
		fmt.Println("The Database isn't supported")
	}

	return nil
}
