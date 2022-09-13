package core

import (
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

func PluginDispatch(task *utils.Task) plugin.Plugin {
	task.Service = strings.ToUpper(task.Service)
	switch task.Service {
	case "POSTGRESQL":
		return &plugin.PostgresService{
			Task:   task,
			Dbname: "postgres",
		}
	case "MSSQL":
		return &plugin.MssqlService{
			Task: task,
		}
	case "MYSQL":
		return &plugin.MysqlService{
			Task: task,
		}
	case "ORACLE":
		return &plugin.OracleService{
			Task: task,
		}
	case "SNMP":
		return &plugin.SnmpService{
			Task: task,
		}
	case "SSH":
		return &plugin.SshService{
			Task: task,
		}
	case "RDP":
		return &plugin.RdpService{
			Task: task,
		}
	case "SMB":
		return &plugin.SmbService{
			Task: task,
		}
	case "FTP":
		return &plugin.FtpService{
			Task: task,
		}

	//case "MONGO":
	//	return &Plugin.MongoService{
	//		Username: task.Username,
	//		Password: task.Password,
	//		Target:   task.Info,
	//	}
	case "VNC":
		return &plugin.VNCService{
			Task: task,
		}
	case "REDIS":
		return &plugin.RedisService{
			Task: task,
		}
	case "LDAP":
		return &plugin.LdapService{
			Task: task,
		}

	default:
		return nil
	}

}
