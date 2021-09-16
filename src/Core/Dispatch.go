package Core

import (
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"fmt"
)

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
	//case "RDP":
	//	return &ExecAble.RdpService{
	//		Username: CurTask.Username,
	//		Password: CurTask.Password,
	//		IpInfo:   CurTask.Info,
	//	}
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
	case "MONGO":
		return &ExecAble.MongoService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.Info,
		}
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
