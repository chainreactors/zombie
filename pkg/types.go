package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers"
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
	FTPService        = &Service{Name: "ftp", DefaultPort: "21", Source: "plugin"}
	SSHService        = &Service{Name: "ssh", DefaultPort: "22", Source: "plugin"}
	SMBService        = &Service{Name: "smb", DefaultPort: "445", Source: "plugin"}
	MSSQLService      = &Service{Name: "mssql", DefaultPort: "1433", Source: "plugin"}
	MYSQLService      = &Service{Name: "mysql", DefaultPort: "3306", Source: "plugin"}
	POSTGRESQLService = &Service{Name: "postgresql", DefaultPort: "5432", Source: "plugin"}
	REDISService      = &Service{Name: "redis", DefaultPort: "6379", Source: "plugin"}
	MONGOService      = &Service{Name: "mongo", DefaultPort: "27017", Source: "plugin"}
	VNCService        = &Service{Name: "vnc", DefaultPort: "5900", Source: "plugin"}
	RDPService        = &Service{Name: "rdp", DefaultPort: "3389", Source: "plugin"}
	SNMPService       = &Service{Name: "snmp", DefaultPort: "161", Source: "plugin"}
	ORACLEService     = &Service{Name: "oracle", DefaultPort: "1521", Source: "plugin"}
	HTTPService       = &Service{Name: "http", DefaultPort: "80", Source: "plugin"}
	HTTPSService      = &Service{Name: "https", DefaultPort: "443", Source: "plugin"}
	GETService        = &Service{Name: "get", DefaultPort: "80", Source: "plugin"}
	PostService       = &Service{Name: "post", DefaultPort: "80", Source: "plugin"}
	LDAPService       = &Service{Name: "ldap", DefaultPort: "389", Source: "plugin"}
	SOCKS5Service     = &Service{Name: "socks5", DefaultPort: "1080", Source: "plugin"}
	TELNETService     = &Service{Name: "telnet", DefaultPort: "23", Source: "plugin"}
	POP3Service       = &Service{Name: "pop3", DefaultPort: "110", Source: "plugin"}
	RSYNCService      = &Service{Name: "rsync", DefaultPort: "873", Source: "plugin"}
	ZookeeperService  = &Service{Name: "zookeeper", DefaultPort: "2181", Source: "plugin"}
	AmqpService       = &Service{Name: "amqp", DefaultPort: "5672", Source: "plugin"}
	MqttService       = &Service{Name: "mqtt", DefaultPort: "1883", Source: "plugin"}
	MemcachedService  = &Service{Name: "memcached", DefaultPort: "11211", Source: "plugin"}
)

var Services = services{}

type services map[string]*Service

func (ss services) Register(s *Service) bool {
	if _, ok := ss[s.Name]; ok {
		return false
	}
	ss[s.Name] = s
	return true
}

func (ss services) DefaultPort(service string) string {
	if s, ok := ss[service]; ok {
		return s.DefaultPort
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

	// alias service
	Services.Register(&Service{Name: "tomcat", DefaultPort: "8080", Source: "plugin"})
}

const (
	PluginSource  = "plugin"
	NeutronSource = "neutron"
)

type Service struct {
	Name        string
	DefaultPort string
	Source      string
}

func (s Service) String() string {
	return s.Name
}

func GetDefault(port string) string {
	for _, s := range Services {
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

	s.WriteString(", " + r.Service + " login successfully\n")
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
	case "json":
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
