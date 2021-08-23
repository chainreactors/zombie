package Core

import (
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"fmt"
)

//var ScanSum = 0

//func BruteDispatch(CurTask Utils.ScanTask) (err error, result Utils.BruteRes) {
//
//	CurTask = Utils.UpdatePass(CurTask)
//
//	switch CurTask.Server {
//
//	case "POSTGRESQL":
//		err, result = ExecAble.PostgresConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "MYSQL":
//		err, result = ExecAble.MysqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "REDIS":
//		err, result = ExecAble.RedisConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "SSH":
//		err, result = ExecAble.SSHConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "MONGO":
//		err, result = ExecAble.MongoConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "MSSQL":
//		err, result = ExecAble.MssqlConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "VNC":
//		err, result = ExecAble.VNCConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "SMB":
//		err, result = ExecAble.SMBConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "ES":
//		err, result = Web.EsConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "FTP":
//		err, result = ExecAble.FtpConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "SNMP":
//		err, result = ExecAble.SnmpConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "TOMCAT":
//		err, result = Web.TomcatConnect(CurTask.Username, CurTask.Password, CurTask.Info)
//	case "RDP":
//		err, result = ExecAble.RdpConnectTest(CurTask.Username, CurTask.Password, CurTask.Info)
//
//	default:
//		fmt.Println("The ExecAble isn't supported")
//	}
//
//	ScanSum += 1
//
//	return err, result
//}

func ExecDispatch(CurTask Utils.ScanTask) ExecAble.ExecAble {
	switch CurTask.Server {
	case "POSTGRESQL":
		return &ExecAble.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &ExecAble.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "MYSQL":
		return &ExecAble.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "SNMP":
		return &ExecAble.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "SSH":
		return &ExecAble.SshService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "RDP":
		return &ExecAble.RdpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "SMB":
		return &ExecAble.SmbService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "FTP":
		return &ExecAble.FtpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	//case "MONGO":
	//	return &ExecAble.MongoService{
	//		Username: CurTask.Username,
	//		Password: CurTask.Password,
	//		IpInfo:   CurTask.Info,
	//	}
	case "VNC":
		return &ExecAble.VNCService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
	case "REDIS":
		return &ExecAble.RedisService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}

	default:
		fmt.Println("The ExecAble isn't supported")
	}

	return nil
}
