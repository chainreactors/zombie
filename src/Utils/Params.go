package Utils

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

var Timeout int
var SSL bool

var (
	PortServer = map[int]string{
		21:    "FTP",
		22:    "SSH",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRESQL",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGO",
		5900:  "VNC",
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
		"VNC":           5900,
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
