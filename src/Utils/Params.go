package Utils

import "time"

type IpInfo struct {
	Ip   string
	Port int
}

type ScanTask struct {
	Info     IpInfo
	Username string
	Password string
	Server   string
}

var Timeout time.Duration

var (
	PortServer = map[int]string{
		21:    "FTP",
		22:    "SSH",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRE",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGO",
	}
	ServerPort = map[string]int{
		"FTP":           21,
		"SSH":           22,
		"SMB":           445,
		"MSSQL":         1433,
		"MYSQL":         3306,
		"POSTGRESQL":    5432,
		"REDIS":         6379,
		"ELASTICSEARCH": 9200,
		"MONGO":         27017,
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
