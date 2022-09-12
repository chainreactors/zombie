package core

import (
	exec2 "github.com/chainreactors/zombie/internal/exec"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

func ExecDispatch(CurTask utils.ScanTask) exec2.ExecAble {
	CurTask.Server = strings.ToUpper(CurTask.Server)
	switch CurTask.Server {
	case "POSTGRESQL":
		return &exec2.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &exec2.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "MYSQL":
		return &exec2.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "ORACLE":
		return &exec2.OracleService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SNMP":
		return &exec2.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SSH":
		return &exec2.SshService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "RDP":
		return &exec2.RdpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SMB":
		return &exec2.SmbService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "FTP":
		return &exec2.FtpService{
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
		return &exec2.VNCService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "REDIS":
		return &exec2.RedisService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "LDAP":
		return &exec2.LdapService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}

	default:
		return nil
	}

}
