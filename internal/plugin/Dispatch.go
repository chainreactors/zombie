package plugin

import (
	"errors"
	"github.com/chainreactors/zombie/internal/plugin/ftp"
	"github.com/chainreactors/zombie/internal/plugin/http"
	"github.com/chainreactors/zombie/internal/plugin/ldap"
	"github.com/chainreactors/zombie/internal/plugin/mongo"
	"github.com/chainreactors/zombie/internal/plugin/mssql"
	"github.com/chainreactors/zombie/internal/plugin/mysql"
	"github.com/chainreactors/zombie/internal/plugin/neutron"
	"github.com/chainreactors/zombie/internal/plugin/oracle"
	"github.com/chainreactors/zombie/internal/plugin/pop3"
	"github.com/chainreactors/zombie/internal/plugin/postgre"
	"github.com/chainreactors/zombie/internal/plugin/rdp"
	"github.com/chainreactors/zombie/internal/plugin/redis"
	"github.com/chainreactors/zombie/internal/plugin/rsync"
	"github.com/chainreactors/zombie/internal/plugin/smb"
	"github.com/chainreactors/zombie/internal/plugin/snmp"
	"github.com/chainreactors/zombie/internal/plugin/socks5"
	"github.com/chainreactors/zombie/internal/plugin/ssh"
	"github.com/chainreactors/zombie/internal/plugin/vnc"
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
	GetResult() *pkg.Result
}

func Dispatch(task *pkg.Task) Plugin {
	switch task.Service {
	case pkg.POSTGRESQLService:
		return &postgre.PostgresPlugin{
			Task:   task,
			Dbname: task.Param["dbname"],
		}
	case pkg.MSSQLService:
		return &mssql.MssqlPlugin{
			Task:     task,
			Instance: task.Param["instance"],
		}
	case pkg.MYSQLService:
		return &mysql.MysqlPlugin{Task: task}
	case pkg.ORACLEService:
		return &oracle.OraclePlugin{
			Task:        task,
			SID:         task.Param["sid"],
			ServiceName: task.Param["service_name"],
		}
	case pkg.SNMPService:
		return &snmp.SnmpPlugin{Task: task}
	case pkg.SSHService:
		return &ssh.SshPlugin{
			Task: task,
		}
	case pkg.RDPService:
		return &rdp.RdpPlugin{Task: task}
	case pkg.SMBService:
		return &smb.SmbPlugin{Task: task}
	case pkg.FTPService:
		return &ftp.FtpPlugin{Task: task}
	case pkg.MONGOService:
		return &mongo.MongoPlugin{Task: task}
	case pkg.VNCService:
		return &vnc.VNCPlugin{Task: task}
	case pkg.REDISService:
		return &redis.RedisPlugin{Task: task}
	case pkg.LDAPService:
		return &ldap.LdapPlugin{Task: task}
	case pkg.HTTPService:
		return &http.HttpPlugin{
			Task: task,
			Path: task.Param["path"],
		}
	case pkg.HTTPSService:
		return &http.HttpPlugin{
			Task: task,
			Path: task.Param["path"],
		}
	case pkg.SOCKS5Service:
		task.Timeout = 10
		return &socks5.Socks5Plugin{
			Task: task,
			Url:  task.Param["url"],
		}
	//case pkg.TELNETService:
	//	return &telnet.TelnetPlugin{Task: task}, nil
	case pkg.POP3Service:
		return &pop3.Pop3Plugin{Task: task}
	case pkg.RSYNCService:
		return &rsync.RsyncPlugin{Task: task}
	default:
		return &neutron.NeutronPlugin{
			Task:    task,
			Service: task.Service,
		}
	}
}
