package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/httpx"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	InterruptError      = errors.New("interrupt")
	ErrorWrongUserOrPwd = errors.New("wrong username or password")
	NotImplUnauthorized = errors.New("not implemented unauthorized")
	RunOpt              = &runOpt{}
)

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }

func init() {
	RegisterServices()
}

var (
	UnknownService    = &Service{Name: "unknown", DefaultPort: "", Source: "unknown"}
	FTPService        = &Service{Name: "ftp", DefaultPort: "21", Source: PluginSource}
	SSHService        = &Service{Name: "ssh", DefaultPort: "22", Source: PluginSource}
	SMBService        = &Service{Name: "smb", DefaultPort: "445", Source: PluginSource}
	MSSQLService      = &Service{Name: "mssql", DefaultPort: "1433", Source: PluginSource}
	MYSQLService      = &Service{Name: "mysql", DefaultPort: "3306", Source: PluginSource}
	POSTGRESQLService = &Service{Name: "postgresql", DefaultPort: "5432", Alias: []string{"postgre"}, Source: PluginSource}
	REDISService      = &Service{Name: "redis", DefaultPort: "6379", Source: PluginSource}
	MONGOService      = &Service{Name: "mongo", DefaultPort: "27017", Alias: []string{"mongodb"}, Source: PluginSource}
	VNCService        = &Service{Name: "vnc", DefaultPort: "5900", Source: PluginSource}
	RDPService        = &Service{Name: "rdp", DefaultPort: "3389", Source: PluginSource}
	SNMPService       = &Service{Name: "snmp", DefaultPort: "161", Source: PluginSource}
	ORACLEService     = &Service{Name: "oracle", DefaultPort: "1521", Source: PluginSource}
	HTTPService       = &Service{Name: "http", DefaultPort: "80", Source: PluginSource}
	HTTPSService      = &Service{Name: "https", DefaultPort: "443", Source: PluginSource}
	GETService        = &Service{Name: "get", DefaultPort: "80", Source: PluginSource}
	PostService       = &Service{Name: "post", DefaultPort: "80", Source: PluginSource}
	LDAPService       = &Service{Name: "ldap", DefaultPort: "389", Source: PluginSource}
	SOCKS5Service     = &Service{Name: "socks5", DefaultPort: "1080", Source: PluginSource}
	TELNETService     = &Service{Name: "telnet", DefaultPort: "23", Source: PluginSource}
	POP3Service       = &Service{Name: "pop3", DefaultPort: "110", Alias: []string{"pop"}, Source: PluginSource}
	RSYNCService      = &Service{Name: "rsync", DefaultPort: "873", Source: PluginSource}
	ZookeeperService  = &Service{Name: "zookeeper", DefaultPort: "2181", Source: PluginSource}
	AmqpService       = &Service{Name: "amqp", DefaultPort: "5672", Source: PluginSource}
	MqttService       = &Service{Name: "mqtt", DefaultPort: "1883", Source: PluginSource}
	MemcachedService  = &Service{Name: "memcached", DefaultPort: "11211", Source: PluginSource}
	HTTPProxyService  = &Service{Name: "http_proxy", DefaultPort: "8080", Source: PluginSource}
	HTTPDigestService = &Service{Name: "digest", DefaultPort: "80", Source: PluginSource}
)

var Services = services{
	Plugins: map[string]*Service{},
	Aliases: map[string]*Service{},
}

type services struct {
	Plugins map[string]*Service
	Aliases map[string]*Service
}

func (ss *services) Get(name string) (*Service, bool) {
	if s, ok := ss.Plugins[name]; ok {
		return s, true
	}
	if s, ok := ss.Aliases[name]; ok {
		return s, true
	}
	return UnknownService, false
}

func (ss *services) Register(s *Service) bool {
	if _, ok := ss.Plugins[s.Name]; !ok {
		ss.Plugins[s.Name] = s
	}
	for _, a := range s.Alias {
		if _, ok := ss.Aliases[a]; !ok {
			ss.Aliases[a] = s
		}
	}
	return true
}

func (ss *services) DefaultPort(service string) string {
	if s, ok := ss.Get(service); ok {
		return s.DefaultPort
	} else if s := utils.ParsePortsString(service); len(s) > 0 {
		return s[0]
	}
	return ""
}

func RegisterServices() {
	Services.Register(FTPService)
	Services.Register(SSHService)
	Services.Register(SMBService)
	Services.Register(MSSQLService)
	Services.Register(MYSQLService)
	Services.Register(POSTGRESQLService)
	Services.Register(REDISService)
	Services.Register(MONGOService)
	Services.Register(VNCService)
	Services.Register(RDPService)
	Services.Register(SNMPService)
	Services.Register(ORACLEService)
	Services.Register(HTTPService)
	Services.Register(HTTPSService)
	Services.Register(GETService)
	Services.Register(PostService)
	Services.Register(LDAPService)
	Services.Register(SOCKS5Service)
	Services.Register(TELNETService)
	Services.Register(POP3Service)
	Services.Register(RSYNCService)
	Services.Register(ZookeeperService)
	Services.Register(AmqpService)
	Services.Register(MqttService)
	Services.Register(MemcachedService)
	Services.Register(HTTPProxyService)
	Services.Register(HTTPDigestService)
	// alias service
	//Services.Register(&Service{Name: "tomcat", DefaultPort: "8080", Source: PluginSource})
}

const (
	PluginSource  = "plugin"
	NeutronSource = "neutron"
)

type Service struct {
	Name        string
	Alias       []string
	DefaultPort string
	Source      string
}

func (s Service) String() string {
	return s.Name
}

func GetDefault(port string) string {
	for _, s := range Services.Plugins {
		if s.DefaultPort == port {
			return s.Name
		}
	}
	return UnknownService.Name
}

// DialFunc 是与 proxyclient.Dial 兼容的拨号函数签名。使用普通函数类型而非
// 直接依赖 proxyclient，避免给 zombie 引入更高的 Go 版本要求（proxyclient
// 需要 go1.24，而 zombie 仍为 go1.16）。SDK 层可直接把 proxyclient.Dial /
// dialer.DialContext 赋值给该字段。
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

// DialTimeoutFunc 与 NewSocketWithDialer / Task.DialTimeout 的签名一致，
// 供 socket 风格的插件（如 rsync）传递代理拨号器。
type DialTimeoutFunc func(network, address string, timeout time.Duration) (net.Conn, error)

type Task struct {
	*parsers.ZombieResult
	Timeout  int                `json:"-"`
	Context  context.Context    `json:"-"`
	Canceler context.CancelFunc `json:"-"`
	Locker   *sync.Mutex        `json:"-"`
	// ProxyDial 非 nil 时，插件应使用它建立连接而非直接 net.Dial。
	ProxyDial DialFunc `json:"-"`
}

func (t *Task) Duration() time.Duration {
	return time.Duration(t.Timeout) * time.Second
}

// DialTimeout 按 task 配置建立连接：设置了 ProxyDial 则走代理，否则直连。
// network 通常为 "tcp"。
func (t *Task) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	if t.ProxyDial != nil {
		ctx := t.Context
		if ctx == nil {
			ctx = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return t.ProxyDial(ctx, network, address)
	}
	return net.DialTimeout(network, address, timeout)
}

// HTTPClient 返回一个 per-task 的 *http.Client（统一经 utils/httpx 构造，零全局）。
// 设置了 ProxyDial 时连接走代理；否则直连。所有 http 系插件应使用它，
// 替代 http.DefaultClient，以保证代理生效且并发隔离。
func (t *Task) HTTPClient(followRedirects bool) *http.Client {
	cfg := httpx.ClientConfig{
		Timeout:            t.Duration(),
		FollowRedirects:    followRedirects,
		InsecureSkipVerify: true,
	}
	if t.ProxyDial != nil {
		cfg.DialContext = httpx.DialContextFunc(t.ProxyDial)
	}
	return httpx.NewHTTPClient(cfg)
}

func NewResult(task *Task, err error) *Result {
	if err != nil {
		return &Result{
			Task: task,
			OK:   false,
			Err:  err,
		}
	} else {
		return &Result{
			Task: task,
			OK:   true,
		}
	}
}

type Result struct {
	*Task         `json:",inline"`
	Vulns         common.Vulns       `json:"vulns,omitempty"`
	Extracteds    parsers.Extracteds `json:"extracteds,omitempty"`
	OK            bool               `json:"ok,omitempty"`
	Err           error              `json:"err,omitempty"`
	ActionResults []*ActionResult    `json:"action_results,omitempty"`
	Loot          map[string][]byte  `json:"loot,omitempty"`
}

func (r *Result) Merge(ar *ActionResult) {
	if ar == nil {
		return
	}
	r.Extracteds = append(r.Extracteds, ar.Extracteds...)
	for k, v := range ar.Vulns {
		if r.Vulns == nil {
			r.Vulns = make(common.Vulns)
		}
		r.Vulns[k] = v
	}
	for k, v := range ar.Loot {
		if r.Loot == nil {
			r.Loot = map[string][]byte{}
		}
		r.Loot[k] = v
	}
	r.ActionResults = append(r.ActionResults, ar)
}

type runOpt struct {
	Raw bool
}

func ParseMethod(input string) (string, string) {
	if RunOpt.Raw {
		return "", input
	}
	if strings.HasPrefix(input, "pk:") {
		return "pk", input[3:]
	} else if strings.HasPrefix(input, "hash:") {
		return "hash", input[5:]
	} else if strings.HasPrefix(input, "raw:") {
		return "raw", input[4:]
	} else {
		return "", input
	}
}
