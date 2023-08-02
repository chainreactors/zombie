package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"strings"
)

const (
	FTP = iota
	LDAP
	MSSQL
	MYSQL
	ORACLE
	POSTGRESQL
	RDP
	SMB
	SNMP
	SSH
	VNC
)

type Plugin interface {
	Query() bool
	GetInfo() bool
	Connect() error
	SetQuery(string)
	Output(interface{})
	Close() error
}

type NilConnError struct {
	service string
}

func (e NilConnError) Error() string {
	return e.service + "has nil conn"
}

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }

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
	case "RDP":
		return &RdpService{
			Task: task,
		}
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
