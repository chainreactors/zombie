package plugin

import (
	"errors"
	"github.com/chainreactors/zombie/internal/plugin/ftp"
	"github.com/chainreactors/zombie/internal/plugin/http"
	"github.com/chainreactors/zombie/internal/plugin/ldap"
	"github.com/chainreactors/zombie/internal/plugin/mongo"
	"github.com/chainreactors/zombie/internal/plugin/mssql"
	"github.com/chainreactors/zombie/internal/plugin/mysql"
	"github.com/chainreactors/zombie/internal/plugin/oracle"
	"github.com/chainreactors/zombie/internal/plugin/pop3"
	"github.com/chainreactors/zombie/internal/plugin/postgre"
	"github.com/chainreactors/zombie/internal/plugin/rdp"
	"github.com/chainreactors/zombie/internal/plugin/redis"
	"github.com/chainreactors/zombie/internal/plugin/smb"
	"github.com/chainreactors/zombie/internal/plugin/snmp"
	"github.com/chainreactors/zombie/internal/plugin/socks5"
	"github.com/chainreactors/zombie/internal/plugin/ssh"
	"github.com/chainreactors/zombie/internal/plugin/telnet"
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
	//Execute() ([]byte, error)
	Close() error
	GetBasic() *pkg.Basic
}

func Dispatch(task *pkg.Task) (Plugin, error) {
	switch task.Service {
	case pkg.POSTGRESQLService:
		return &postgre.PostgresPlugin{
			Task:   task,
			Dbname: task.Param["dbname"],
		}, nil
	case pkg.MSSQLService:
		return &mssql.MssqlPlugin{
			Task:     task,
			Instance: task.Param["instance"],
		}, nil
	case pkg.MYSQLService:
		return &mysql.MysqlPlugin{Task: task}, nil
	case pkg.ORACLEService:
		return &oracle.OraclePlugin{
			Task:        task,
			SID:         task.Param["sid"],
			ServiceName: task.Param["service_name"],
		}, nil
	case pkg.SNMPService:
		return &snmp.SnmpPlugin{Task: task}, nil
	case pkg.SSHService:
		return &ssh.SshPlugin{
			Task: task,
		}, nil
	case pkg.RDPService:
		return &rdp.RdpPlugin{Task: task}, nil
	case pkg.SMBService:
		return &smb.SmbPlugin{Task: task}, nil
	case pkg.FTPService:
		return &ftp.FtpPlugin{Task: task}, nil
	case pkg.MONGOService:
		return &mongo.MongoPlugin{Task: task}, nil
	case pkg.VNCService:
		return &vnc.VNCPlugin{Task: task}, nil
	case pkg.REDISService:
		return &redis.RedisPlugin{Task: task}, nil
	case pkg.LDAPService:
		return &ldap.LdapPlugin{Task: task}, nil
	case pkg.HTTPService:
		return &http.HttpPlugin{
			Task: task,
			Path: task.Param["path"],
		}, nil
	case pkg.HTTPSService:
		return &http.HttpPlugin{
			Task: task,
			Path: task.Param["path"],
		}, nil
	case pkg.TomcatService:
		return &http.HttpPlugin{
			Task: task,
			Path: "manager",
		}, nil
	case pkg.KibanaService:
		return &http.HttpPlugin{Task: task}, nil
	case pkg.SOCKS5Service:
		task.Timeout = 10
		return &socks5.Socks5Plugin{Task: task}, nil
	case pkg.TELNETService:
		return &telnet.TelnetPlugin{Task: task}, nil
	case pkg.POP3Service:
		return &pop3.Pop3Plugin{Task: task}, nil
	//case pkg.RSYNCService:
	//	return &rsync.RsyncPlugin{Task: task}, nil
	default:
		return nil, ErrKnownPlugin
	}
}
