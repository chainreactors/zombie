package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/logs"
	"strconv"
	"time"
)

type Task struct {
	IP         string             `json:"ip"`
	Port       string             `json:"port"`
	Service    Service            `json:"service"`
	Username   string             `json:"username"`
	Password   string             `json:"password"`
	ExecString string             `json:"exec"`
	Instance   string             `json:"-"`
	Timeout    int                `json:"-"`
	Context    context.Context    `json:"-"`
	Canceler   context.CancelFunc `json:"-"`
	OutputCh   chan *Result       `json:"-"`
}

func (t *Task) Address() string {
	return t.IP + ":" + t.Port
}

func (t *Task) URI() string {
	return t.Service.String() + "://" + t.Address()
}

func (t *Task) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s", t.Service, t.Username, t.Password, t.IP, t.Port)
}

func (t *Task) UintPort() uint16 {
	p, _ := strconv.Atoi(t.Port)
	return uint16(p)
}

func (t *Task) Duration() time.Duration {
	return time.Duration(t.Timeout) * time.Second
}

type Result struct {
	*Task
	OK         bool
	Err        error
	Additional string `json:"additional"`
}

func (r *Result) String() string {
	return fmt.Sprintf("[+] %s://%s\t%s\t%s\n", r.Service, r.Address(), r.Username, r.Password)
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

//var (
//	ValueableSlice = []string{"PWD", "PASS", "PASSWORD", "CERT", "EMAIL", "MOBILE", "PAPER"}
//)

type Service string

var (
	FTPService        Service = "ftp"
	SSHService        Service = "ssh"
	SMBService        Service = "smb"
	MSSQLService      Service = "mssql"
	MYSQLService      Service = "mysql"
	POSTGRESQLService Service = "postgresql"
	REDISService      Service = "redis"
	ESService         Service = "es"
	MONGOService      Service = "mongo"
	VNCService        Service = "vnc"
	RDPService        Service = "rdp"
	SNMPService       Service = "snmp"
	ORACLEService     Service = "oracle"
	HTTPService       Service = "http"
	HTTPSService      Service = "https"
	LDAPService       Service = "ldap"
	UnknownService    Service = ""
)

var Services = map[Service]string{
	FTPService:        "21",
	SSHService:        "22",
	SMBService:        "445",
	MSSQLService:      "1433",
	MYSQLService:      "3306",
	POSTGRESQLService: "5432",
	REDISService:      "6379",
	ESService:         "9200",
	MONGOService:      "27017",
	VNCService:        "5900",
	RDPService:        "3389",
	SNMPService:       "161",
	ORACLEService:     "1521",
	LDAPService:       "389",
	HTTPService:       "80",
	HTTPSService:      "443",
}

func (s Service) String() string {
	return string(s)
}
func (s Service) DefaultPort() string {
	if port, ok := Services[s]; ok {
		return port
	}
	return ""
}

func GetDefault(s string) Service {
	switch s {
	case "22":
		return SSHService
	case "21":
		return FTPService
	case "445":
		return SMBService
	case "1433":
		return MSSQLService
	case "3306":
		return MYSQLService
	case "5432":
		return POSTGRESQLService
	case "6379":
		return REDISService
	case "9200":
		return ESService
	case "27017":
		return MONGOService
	case "5900":
		return VNCService
	case "3389":
		return RDPService
	case "161":
		return SNMPService
	case "1521":
		return ORACLEService
	case "80":
		return HTTPService
	case "443":
		return HTTPSService
	case "389":
		return LDAPService
	default:
		return UnknownService
	}
}

func UseDefaultPassword(service string, top int) []string {
	if pwds, ok := Keywords[service+"_pwd"]; ok {
		if top == 0 || top > len(pwds) {
			return pwds
		} else {
			return pwds[:top]
		}
	} else {
		return []string{"admin"}
	}
}

func UseDefaultUser(service string) []string {
	if users, ok := Keywords[service+"_user"]; ok {
		return users
	} else {
		return []string{"admin"}
	}
}
