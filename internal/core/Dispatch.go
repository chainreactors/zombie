package core

import (
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

func ExecDispatch(CurTask utils.ScanTask) plugin.ExecAble {
	CurTask.Server = strings.ToUpper(CurTask.Server)
	switch CurTask.Server {
	case "POSTGRESQL":
		return &plugin.PostgresService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
			Dbname:   "postgres",
		}
	case "MSSQL":
		return &plugin.MssqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "MYSQL":
		return &plugin.MysqlService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "ORACLE":
		return &plugin.OracleService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SNMP":
		return &plugin.SnmpService{
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SSH":
		return &plugin.SshService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "RDP":
		return &plugin.RdpService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "SMB":
		return &plugin.SmbService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "FTP":
		return &plugin.FtpService{
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
		return &plugin.VNCService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "REDIS":
		return &plugin.RedisService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}
	case "LDAP":
		return &plugin.LdapService{
			Username: CurTask.Username,
			Password: CurTask.Password,
			IpInfo:   CurTask.IpInfo,
		}

	default:
		return nil
	}

}
