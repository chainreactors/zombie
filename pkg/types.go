package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/logs"
	"strconv"
	"strings"
)

var (
	Rules map[string]string = make(map[string]string)
)

type Task struct {
	IP         string             `json:"ip"`
	Port       string             `json:"port"`
	Service    string             `json:"service"`
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
	return strings.ToLower(t.Service) + "://" + t.Address()
}

func (t *Task) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s", t.Service, t.Username, t.Password, t.IP, t.Port)
}

func (t *Task) UintPort() uint16 {
	p, _ := strconv.Atoi(t.Port)
	return uint16(p)
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

var (
	ServicePortMap = map[string]string{
		"FTP":        "21",
		"SSH":        "22",
		"SMB":        "445",
		"MSSQL":      "1433",
		"MYSQL":      "3306",
		"POSTGRESQL": "5432",
		"REDIS":      "6379",
		"ES":         "9200",
		"MONGO":      "27017",
		"VNC":        "5900",
		"TOMCAT":     "8080",
		//"RDP":        "3389",
		"SNMP":   "161",
		"ORACLE": "1521",
		"21":     "FTP",
		"22":     "SSH",
		"445":    "SMB",
		"1433":   "MSSQL",
		"3306":   "MYSQL",
		"5432":   "POSTGRESQL",
		"6379":   "REDIS",
		"9200":   "ES",
		"27017":  "MONGO",
		"5900":   "VNC",
		"8080":   "TOMCAT",
		"161":    "SNMP",
		"3389":   "RDP",
		"1521":   "ORACLE",
	}
)

func UseDefaultPassword(service string, top int) []string {
	if pwds, ok := Keywords[strings.ToLower(service)+"_pwd"]; ok {
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
	if users, ok := Keywords[strings.ToLower(service)+"_user"]; ok {
		return users
	} else {
		return []string{"admin"}
	}
}
