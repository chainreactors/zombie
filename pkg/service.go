package pkg

func init() {
	RegisterServices()
}

var (
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
	LDAPService       = &Service{Name: "ldap", DefaultPort: "389", Source: "plugin"}
	SOCKS5Service     = &Service{Name: "socks5", DefaultPort: "1080", Source: "plugin"}
	TELNETService     = &Service{Name: "telnet", DefaultPort: "23", Source: "plugin"}
	POP3Service       = &Service{Name: "pop3", DefaultPort: "110", Source: "plugin"}
	RSYNCService      = &Service{Name: "rsync", DefaultPort: "873", Source: "plugin"}
	UnknownService    = &Service{Name: "unknown", DefaultPort: "", Source: "unknown"}
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

func UseDefaultPassword(service string, top int) []string {
	if pwds, ok := Keywords[service+"_pwd"]; ok {
		if top == 0 || top > len(pwds) {
			return pwds
		} else {
			return pwds[:top]
		}
	} else {
		return Keywords["top10_pwd"]
	}
}

func UseDefaultUser(service string) []string {
	if users, ok := Keywords[service+"_user"]; ok {
		return users
	} else {
		return Keywords["top10_user"]
	}
}
