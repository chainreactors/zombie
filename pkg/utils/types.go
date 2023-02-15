package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/logs"
	"strconv"
	"strings"
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

	DefaultUsernames = map[string][]string{
		"FTP":        {"ftp", "admin", "www", "wwwroot"},
		"MYSQL":      {"root", "mysql"},
		"MSSQL":      {"sa", "sql"},
		"SMB":        {"administrator", "admin"},
		"RDP":        {"administrator", "admin"},
		"POSTGRESQL": {"postgres", "admin"},
		"SSH":        {"root", "admin"},
		"MONGO":      {"root", "admin"},
		"REDIS":      {"root"},
	}

	DefaultPasswords = map[string][]string{
		"FTP":        {"123456", "admin", "admin123", "root", "q1w2e3r4", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"MYSQL":      {"123456", "admin", "admin123", "root", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"MSSQL":      {"123456", "admin", "admin123", "root", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"SMB":        {"123456", "admin", "admin123", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "sa123456"},
		"RDP":        {"123456", "admin", "admin123", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"POSTGRESQL": {"123456", "admin", "admin123", "root", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"SSH":        {"123456", "admin", "admin123", "root", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "admin@huawei.com"},
		"MONGO":      {"123456", "admin", "admin123", "root", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
		"REDIS":      {"123456", "admin", "admin123", "root", "", "password", "123", "1", "admin@123", "Admin@123", "P@ssw0rd", "Passw0rd", "12345678", "test", "test123", "123qwe!@#", "123456789", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "1qaz@WSX", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#"},
	}
)

func UseDefaultPassword(service string, top int) []string {
	if pwds, ok := DefaultPasswords[service]; ok {
		if top > len(pwds) {
			return pwds
		} else {
			return pwds[:top]
		}
	} else {
		return []string{"admin"}
	}
}

func UseDefaultUser(service string) []string {
	if users, ok := DefaultUsernames[service]; ok {
		return users
	} else {
		return []string{"admin"}
	}
}
