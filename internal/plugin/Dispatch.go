package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
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
	service pkg.Service
}

func (e NilConnError) Error() string {
	return e.service.String() + " has nil conn"
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
	switch task.Service {
	case pkg.POSTGRESQLService:
		return &PostgresService{
			Task:   task,
			Dbname: "postgres",
		}
	case pkg.MSSQLService:
		return &MssqlService{
			Task: task,
		}
	case pkg.MYSQLService:
		return &MysqlService{
			Task: task,
		}
	case pkg.ORACLEService:
		return &OracleService{
			Task: task,
		}
	case pkg.SNMPService:
		return &SnmpService{
			Task: task,
		}
	case pkg.SSHService:
		return &SshService{
			Task: task,
		}
	case pkg.RDPService:
		return &RdpService{
			Task: task,
		}
	case pkg.SMBService:
		return &SmbService{
			Task: task,
		}
	case pkg.FTPService:
		return &FtpService{
			Task: task,
		}
	case pkg.MONGOService:
		return &MongoService{
			Task: task,
		}
	case pkg.VNCService:
		return &VNCService{
			Task: task,
		}
	case pkg.REDISService:
		return &RedisService{
			Task: task,
		}
	case pkg.LDAPService:
		return &LdapService{
			Task: task,
		}
	case pkg.HTTPService:
		return &HttpService{
			Task: task,
			HttpInf: HttpInf{
				Path: task.Param["path"],
			},
		}
	case pkg.HTTPSService:
		return &HttpService{
			Task: task,
			HttpInf: HttpInf{
				Path: task.Param["path"],
			},
		}
	case pkg.SOCKS5Service:
		task.Timeout = 10
		return &Socks5Service{
			Task: task,
			//Socks5Inf: Socks5Inf{
			//	Url: task.Param["url"],
			//},
		}
	case pkg.TELNETService:
		task.Timeout = 10
		return &TelnetService{
			Task: task,
		}
	default:
		return nil
	}

}
