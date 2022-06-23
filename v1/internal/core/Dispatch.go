package core

import (
	"Zombie/v1/internal/exec"
	"Zombie/v1/pkg/utils"
	"strings"
)

func ExecDispatch(CurTask utils.ScanTask) exec.ExecAble {
	CurTask.Server = strings.ToUpper(CurTask.Server)
	switch CurTask.Server {
	case "POSTGRESQL":
		return &exec.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &exec.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "MYSQL":
		return &exec.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "ORACLE":
		return &exec.OracleService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SNMP":
		return &exec.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SSH":
		return &exec.SshService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "RDP":
		return &exec.RdpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SMB":
		return &exec.SmbService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "FTP":
		return &exec.FtpService{
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
		return &exec.VNCService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "REDIS":
		return &exec.RedisService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "LDAP":
		return &exec.LdapService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}

	default:
		return nil
	}

}
