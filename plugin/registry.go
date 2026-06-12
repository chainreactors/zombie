package plugin

import (
	"github.com/chainreactors/zombie/plugin/ftp"
	"github.com/chainreactors/zombie/plugin/http"
	"github.com/chainreactors/zombie/plugin/ldap"
	"github.com/chainreactors/zombie/plugin/memcache"
	"github.com/chainreactors/zombie/plugin/mongo"
	"github.com/chainreactors/zombie/plugin/mq"
	"github.com/chainreactors/zombie/plugin/mssql"
	"github.com/chainreactors/zombie/plugin/mysql"
	"github.com/chainreactors/zombie/plugin/neutron"
	"github.com/chainreactors/zombie/plugin/oracle"
	"github.com/chainreactors/zombie/plugin/pop3"
	"github.com/chainreactors/zombie/plugin/postgre"
	"github.com/chainreactors/zombie/plugin/rdp"
	"github.com/chainreactors/zombie/plugin/redis"
	"github.com/chainreactors/zombie/plugin/rsync"
	"github.com/chainreactors/zombie/plugin/smb"
	"github.com/chainreactors/zombie/plugin/snmp"
	"github.com/chainreactors/zombie/plugin/socks5"
	"github.com/chainreactors/zombie/plugin/ssh"
	"github.com/chainreactors/zombie/plugin/vnc"
	"github.com/chainreactors/zombie/plugin/zookeeper"
	"github.com/chainreactors/zombie/pkg"
)

var registry = map[string]Plugin{}

func Register(name string, p Plugin) {
	registry[name] = p
}

func Get(service string) (Plugin, bool) {
	p, ok := registry[service]
	return p, ok
}

func DefaultRegistry() map[string]Plugin {
	m := make(map[string]Plugin, len(registry))
	for k, v := range registry {
		m[k] = v
	}
	return m
}

func init() {
	Register(pkg.SSHService.Name, &ssh.SshPlugin{})
	Register(pkg.MYSQLService.Name, &mysql.MysqlPlugin{})
	Register(pkg.POSTGRESQLService.Name, &postgre.PostgresPlugin{})
	Register(pkg.MSSQLService.Name, &mssql.MssqlPlugin{})
	Register(pkg.ORACLEService.Name, &oracle.OraclePlugin{})
	Register(pkg.REDISService.Name, &redis.RedisPlugin{})
	Register(pkg.MONGOService.Name, &mongo.MongoPlugin{})
	Register(pkg.SMBService.Name, &smb.SmbPlugin{})
	Register(pkg.FTPService.Name, &ftp.FtpPlugin{})
	Register(pkg.LDAPService.Name, &ldap.LdapPlugin{})
	Register(pkg.VNCService.Name, &vnc.VNCPlugin{})
	Register(pkg.RDPService.Name, &rdp.RdpPlugin{})
	Register(pkg.SNMPService.Name, &snmp.SnmpPlugin{})
	Register(pkg.POP3Service.Name, &pop3.Pop3Plugin{})
	Register(pkg.ZookeeperService.Name, &zookeeper.ZookeeperPlugin{})
	Register(pkg.MemcachedService.Name, &memcache.MemcachePlugin{})
	Register(pkg.AmqpService.Name, &mq.AMQPPlugin{})
	Register(pkg.MqttService.Name, &mq.MQTTPlugin{})
	Register(pkg.RSYNCService.Name, &rsync.RsyncPlugin{})
	Register(pkg.SOCKS5Service.Name, &socks5.Socks5Plugin{})
	Register(pkg.HTTPService.Name, &http.HttpAuthPlugin{})
	Register(pkg.HTTPSService.Name, &http.HttpAuthPlugin{})
	Register(pkg.HTTPProxyService.Name, &http.HTTPProxyPlugin{})
	Register(pkg.HTTPDigestService.Name, &http.HTTPDigestPlugin{})
	Register(pkg.GETService.Name, http.NewHTTPPlugin("GET"))
	Register(pkg.PostService.Name, http.NewHTTPPlugin("POST"))
	Register("neutron", &neutron.NeutronPlugin{})
}
