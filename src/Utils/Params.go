package Utils

import "context"

type IpInfo struct {
	Ip   string
	Port int
	SSL  bool
}

type ScanTask struct {
	Info     IpInfo
	Username string
	Password string
	Server   string
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
var SSL, IsAuto bool
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
		"SNMP":       161,
	}

	ExecPort = map[string]int{
		"MSSQL":      1433,
		"MYSQL":      3306,
		"POSTGRESQL": 5432,
	}

	ExecServer = map[int]string{
		1433: "MSSQL",
		3306: "MYSQL",
		5432: "POSTGRE",
	}
)
