// Package plugins assembles the protocol drivers bundled with Zombie.
package plugins

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin"
	"github.com/chainreactors/zombie/plugin/ftp"
	httpplugin "github.com/chainreactors/zombie/plugin/http"
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
)

type builtin struct {
	service pkg.Service
	plugin  plugin.Plugin
}

var bundled = []builtin{
	{pkg.Service{Name: "ftp", DefaultPort: "21"}, &ftp.FtpPlugin{}},
	{pkg.Service{Name: "http", DefaultPort: "80"}, &httpplugin.HttpAuthPlugin{}},
	{pkg.Service{Name: "https", DefaultPort: "443"}, &httpplugin.HttpAuthPlugin{}},
	{pkg.Service{Name: "get", DefaultPort: "80"}, httpplugin.NewHTTPPlugin("GET")},
	{pkg.Service{Name: "post", DefaultPort: "80"}, httpplugin.NewHTTPPlugin("POST")},
	{pkg.Service{Name: "http_proxy", DefaultPort: "8080"}, &httpplugin.HTTPProxyPlugin{}},
	{pkg.Service{Name: "digest", DefaultPort: "80"}, &httpplugin.HTTPDigestPlugin{}},
	{pkg.Service{Name: "ldap", DefaultPort: "389"}, &ldap.LdapPlugin{}},
	{pkg.Service{Name: "memcached", DefaultPort: "11211"}, &memcache.MemcachePlugin{}},
	{pkg.Service{Name: "mongo", Alias: []string{"mongodb"}, DefaultPort: "27017"}, &mongo.MongoPlugin{}},
	{pkg.Service{Name: "amqp", DefaultPort: "5672"}, &mq.AMQPPlugin{}},
	{pkg.Service{Name: "mqtt", DefaultPort: "1883"}, &mq.MQTTPlugin{}},
	{pkg.Service{Name: "mssql", DefaultPort: "1433"}, &mssql.MssqlPlugin{}},
	{pkg.Service{Name: "mysql", DefaultPort: "3306"}, &mysql.MysqlPlugin{}},
	{pkg.Service{Name: "oracle", DefaultPort: "1521"}, &oracle.OraclePlugin{}},
	{pkg.Service{Name: "pop3", Alias: []string{"pop"}, DefaultPort: "110"}, &pop3.Pop3Plugin{}},
	{pkg.Service{Name: "postgresql", Alias: []string{"postgre"}, DefaultPort: "5432"}, &postgre.PostgresPlugin{}},
	{pkg.Service{Name: "rdp", DefaultPort: "3389"}, &rdp.RdpPlugin{}},
	{pkg.Service{Name: "redis", DefaultPort: "6379"}, &redis.RedisPlugin{}},
	{pkg.Service{Name: "rsync", DefaultPort: "873"}, &rsync.RsyncPlugin{}},
	{pkg.Service{Name: "smb", DefaultPort: "445"}, &smb.SmbPlugin{}},
	{pkg.Service{Name: "snmp", DefaultPort: "161"}, &snmp.SnmpPlugin{}},
	{pkg.Service{Name: "socks5", DefaultPort: "1080"}, &socks5.Socks5Plugin{}},
	{pkg.Service{Name: "ssh", DefaultPort: "22"}, &ssh.SshPlugin{}},
	{pkg.Service{Name: "vnc", DefaultPort: "5900"}, &vnc.VNCPlugin{}},
	{pkg.Service{Name: "zookeeper", DefaultPort: "2181"}, &zookeeper.ZookeeperPlugin{}},
}

func RegisterServices() {
	for _, entry := range bundled {
		service := entry.service
		if service.Source == "" {
			service.Source = pkg.PluginSource
		}
		pkg.Services.Register(&service)
	}
}

func Default() map[string]plugin.Plugin {
	RegisterServices()
	registry := make(map[string]plugin.Plugin, len(bundled))
	for _, entry := range bundled {
		registry[entry.service.Name] = entry.plugin
		for _, alias := range entry.service.Alias {
			registry[alias] = entry.plugin
		}
	}
	return registry
}

func Fallback() plugin.Plugin {
	return &neutron.NeutronPlugin{}
}
