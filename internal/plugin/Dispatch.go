package plugin

import (
	"errors"
	"github.com/chainreactors/zombie/internal/plugin/ftp"
	"github.com/chainreactors/zombie/internal/plugin/mysql"
	"github.com/chainreactors/zombie/pkg"
)

var (
	ErrKnownPlugin = errors.New("not found plugin")
)

type Plugin interface {
	Name() string
	Unauth() (bool, error)
	Login() error
	Close() error
	GetBasic() *pkg.Basic
}

func Dispatch(task *pkg.Task) (Plugin, error) {
	switch task.Service {
	//case pkg.POSTGRESQLService:
	//	return &PostgresService{
	//		Task:   task,
	//		Dbname: "postgres",
	//	}
	//case pkg.MSSQLService:
	//	return &MssqlService{
	//		Task: task,
	//	}
	case pkg.MYSQLService:
		return &mysql.MysqlPlugin{
			Task: task,
		}, nil
	//case pkg.ORACLEService:
	//	return &OracleService{
	//		Task: task,
	//	}
	//case pkg.SNMPService:
	//	return &SnmpService{
	//		Task: task,
	//	}
	//case pkg.SSHService:
	//	return &SshService{
	//		Task: task,
	//	}
	//case pkg.RDPService:
	//	return &RdpService{
	//		Task: task,
	//	}
	//case pkg.SMBService:
	//	return &SmbService{
	//		Task: task,
	//	}
	case pkg.FTPService:
		return &ftp.FtpPlugin{
			Task: task,
		}, nil
	//case pkg.MONGOService:
	//	return &MongoService{
	//		Task: task,
	//	}
	//case pkg.VNCService:
	//	return &VNCService{
	//		Task: task,
	//	}
	//case pkg.REDISService:
	//	return &RedisService{
	//		Task: task,
	//	}
	//case pkg.LDAPService:
	//	return &LdapService{
	//		Task: task,
	//	}
	//case pkg.HTTPService:
	//	return &HttpService{
	//		Task: task,
	//		HttpInf: HttpInf{
	//			Path: task.Param["path"],
	//		},
	//	}
	//case pkg.HTTPSService:
	//	return &HttpService{
	//		Task: task,
	//		HttpInf: HttpInf{
	//			Path: task.Param["path"],
	//		},
	//	}
	//case pkg.SOCKS5Service:
	//	task.Timeout = 10
	//	return &Socks5Service{
	//		Task: task,
	//		//Socks5Inf: Socks5Inf{
	//		//	Url: task.Param["url"],
	//		//},
	//	}
	//case pkg.TELNETService:
	//	return &TelnetService{
	//		Task: task,
	//	}
	//case pkg.POP3Service:
	//	return &Pop3Service{
	//		Task: task,
	//	}
	//case pkg.RSYNCService:
	//	return &RsyncService{
	//		Task: task,
	//	}
	default:
		return nil, ErrKnownPlugin
	}

}
