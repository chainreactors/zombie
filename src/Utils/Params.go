package Utils

import "context"

type IpInfo struct {
	Ip   string `json:"IP"`
	Port int    `json:"Port"`
	SSL  bool
}

type ScanTask struct {
	Info     IpInfo
	Username string
	Password string
	Server   string
	Input    string
}

type OutputRes struct {
	Type       string `json:"type"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Additional string `json:"additional"`
}

type BruteRes struct {
	Result     bool
	Additional string
}

var (
	Thread  int
	Simple  bool
	Timeout int
	Proc    int
)

var File string
var OutputType string
var SSL, IsAuto, More bool
var FileFormat string
var ChildContext context.Context
var ChildCancel context.CancelFunc

var (
	PortServer = map[int]string{
		21:    "FTP",
		22:    "SSH",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRESQL",
		6379:  "REDIS",
		9200:  "ES",
		27017: "MONGO",
		5900:  "VNC",
		8080:  "TOMCAT",
		161:   "SNMP",
		3389:  "RDP",
	}
	ServerPort = map[string]int{
		"FTP":        21,
		"SSH":        22,
		"SMB":        445,
		"MSSQL":      1433,
		"MYSQL":      3306,
		"POSTGRESQL": 5432,
		"REDIS":      6379,
		"ES":         9200,
		"MONGO":      27017,
		"VNC":        5900,
		"TOMCAT":     8080,
		"RDP":        3389,
		"SNMP":       161,
	}

	DefaultUserDict = map[string][]string{
		"FTP":        {"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
		"MYSQL":      {"root", "mysql"},
		"MSSQL":      {"sa", "sql"},
		"SMB":        {"administrator", "admin", "guest"},
		"POSTGRESQL": {"postgres", "admin"},
		"SSH":        {"root", "admin"},
		"MONGO":      {"root", "admin"},
	}

	DefaultPasswords = []string{"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "%user%", "%user%1", "%user%111", "%user%123", "%user%@123", "%user%_123", "%user%#123", "%user%@111", "%user%@2019", "%user%@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456"}

	//ExecPort = map[string]int{
	//	"MSSQL":      1433,
	//	"MYSQL":      3306,
	//	"POSTGRESQL": 5432,
	//	"SNMP":       161,
	//	"SSH":        22,
	//	"RDP":        3389,
	//}
	//
	//ExecServer = map[int]string{
	//	1433: "MSSQL",
	//	3306: "MYSQL",
	//	5432: "POSTGRE",
	//	161:  "SNMP",
	//	22:   "SSH",
	//	3389: "RDP",
	//}
)
