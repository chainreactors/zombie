package plugin

import (
	"github.com/chainreactors/zombie/pkg"
	"strings"
)

func Dispatch(task *pkg.Task) Plugin {
	task.Service = strings.ToUpper(task.Service)
	switch task.Service {
	case "POSTGRESQL":
		return &PostgresService{
			Task:   task,
			Dbname: "postgres",
		}
	case "MSSQL":
		return &MssqlService{
			Task: task,
		}
	case "MYSQL":
		return &MysqlService{
			Task: task,
		}
	case "ORACLE":
		return &OracleService{
			Task: task,
		}
	case "SNMP":
		return &SnmpService{
			Task: task,
		}
	case "SSH":
		return &SshService{
			Task: task,
		}
	//case "RDP":
	//	return &RdpService{
	//		Task: task,
	//	}
	case "SMB":
		return &SmbService{
			Task: task,
		}
	case "FTP":
		return &FtpService{
			Task: task,
		}
	case "MONGO":
		return &MongoService{
			Task: task,
		}
	case "VNC":
		return &VNCService{
			Task: task,
		}
	case "REDIS":
		return &RedisService{
			Task: task,
		}
	case "LDAP":
		return &LdapService{
			Task: task,
		}

	default:
		return nil
	}

}
