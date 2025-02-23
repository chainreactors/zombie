package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
	"strconv"
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

type TaskMod int

const (
	TaskModBrute TaskMod = 0 + iota
	TaskModUnauth
	TaskModCheck
	TaskModSniper
	TaskModPitchfork
)

func (m TaskMod) String() string {
	switch m {
	case TaskModBrute:
		return "brute"
	case TaskModUnauth:
		return "unauth"
	case TaskModCheck:
		return "check"
	case TaskModSniper:
		return "sniper"
	case TaskModPitchfork:
		return "pitchfork"
	default:
		return "unknown"
	}
}

type Task struct {
	IP       string             `json:"ip"`
	Port     string             `json:"port"`
	Service  string             `json:"service"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Scheme   string             `json:"scheme"`
	Param    map[string]string  `json:"-"`
	Mod      TaskMod            `json:"-"`
	Timeout  int                `json:"-"`
	Context  context.Context    `json:"-"`
	Canceler context.CancelFunc `json:"-"`
	Locker   *sync.Mutex        `json:"-"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
}

func (t *Task) Address() string {
	return t.IP + ":" + t.Port
}

func (t *Task) URI() string {
	if t.Scheme != "" {
		return t.Scheme + "://" + t.Address()
	} else {
		return t.Service + "://" + t.Address()
	}
}

func (t *Task) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s", t.Scheme, t.Username, t.Password, t.IP, t.Port)
}

func (t *Task) UintPort() uint16 {
	p, _ := strconv.Atoi(t.Port)
	return uint16(p)
}

func (t *Task) Duration() time.Duration {
	return time.Duration(t.Timeout) * time.Second
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
	*Task
	Vulns      common.Vulns
	Extracteds parsers.Extracteds
	OK         bool
	Err        error
}

func (r *Result) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("[%s] ", r.Mod.String()))
	s.WriteString(r.URI())
	if r.Username != "" {
		s.WriteString(" " + r.Username)
	}
	if r.Password != "" {
		s.WriteString(" " + r.Password)
	}
	if len(r.Param) != 0 {
		s.WriteString(" " + fmt.Sprintf("%v", r.Param))
	}
	if r.Mod == TaskModCheck {
		s.WriteString(", " + r.Service + " maybe honeypot or unauth!!!\n")
	} else {
		s.WriteString(", " + r.Service + " login successfully\n")
	}

	return s.String()
}

func (r *Result) Json() string {
	bs, err := json.Marshal(r)
	if err != nil {
		logs.Log.Error(err.Error())
		return ""
	}
	return string(bs) + "\n"
}

func (r *Result) Format(form string) string {
	switch form {
	case "json", "jl":
		return r.Json()
	case "csv":
		return ""
	default:
		return r.String()
	}
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
