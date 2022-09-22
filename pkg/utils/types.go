package utils

import (
	"context"
	"fmt"
	"github.com/chainreactors/ipcs"
	"strconv"
)

type Task struct {
	*ipcs.Addr
	Service    string `json:"Server"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	ExecString string `json:"exec"`
	Instance   string
	Timeout    int
	Context    context.Context
	Canceler   context.CancelFunc
}

func (r Task) Address() string {
	return r.Addr.String()
}

func (r Task) UintPort() uint16 {
	p, _ := strconv.Atoi(r.Port)
	return uint16(p)
}

type Codebook struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

type Result struct {
	*Task
	OK         bool
	Err        error
	Additional string `json:"additional"`
}

func (r Result) String() string {
	return fmt.Sprintf("[+] %s\t%s\t%s\n", r.Address(), r.Username, r.Password)
}

var OutputType string
var IsAuto, More bool
var FileFormat string
var BrutedList []Result

var (
	ValueableSlice = []string{"PWD", "PASS", "PASSWORD", "CERT", "EMAIL", "MOBILE", "PAPER"}
)

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
		"RDP":        "3389",
		"SNMP":       "161",
		"ORACLE":     "1521",
		"21":         "FTP",
		"22":         "SSH",
		"445":        "SMB",
		"1433":       "MSSQL",
		"3306":       "MYSQL",
		"5432":       "POSTGRESQL",
		"6379":       "REDIS",
		"9200":       "ES",
		"27017":      "MONGO",
		"5900":       "VNC",
		"8080":       "TOMCAT",
		"161":        "SNMP",
		"3389":       "RDP",
		"1521":       "ORACLE",
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
		"FTP":        {"123456", "admin", "admin123", "root", "q1w2e3r4", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"MYSQL":      {"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"MSSQL":      {"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"SMB":        {"123456", "admin", "admin123", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456"},
		"RDP":        {"123456", "admin", "admin123", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"POSTGRESQL": {"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"SSH":        {"123456", "admin", "admin123", "root", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "admin@huawei.com"},
		"MONGO":      {"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
		"REDIS":      {"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "1QAZ2wsx", "#EDC2wsX", "We1c0me!", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!"},
	}
)
