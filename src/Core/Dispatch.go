package Core

import (
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"strings"
)

func ExecDispatch(CurTask Utils.ScanTask) ExecAble.ExecAble {
	CurTask.Server = strings.ToUpper(CurTask.Server)
	switch CurTask.Server {
	case "POSTGRESQL":
		return &ExecAble.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &ExecAble.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "MYSQL":
		return &ExecAble.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "ORACLE":
		return &ExecAble.OracleService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SNMP":
		return &ExecAble.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SSH":
		return &ExecAble.SshService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "RDP":
		return &ExecAble.RdpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SMB":
		return &ExecAble.SmbService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "FTP":
		return &ExecAble.FtpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
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
			IpInfo:   CurTask.IpInfo,
		}
	case "REDIS":
		return &ExecAble.RedisService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "LDAP":
		return &ExecAble.LdapService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}

	default:
		return nil
	}

	return nil
}
