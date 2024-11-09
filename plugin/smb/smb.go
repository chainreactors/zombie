package smb

import (
	"github.com/chainreactors/utils/encode"
	"github.com/chainreactors/zombie/pkg"
	"github.com/hirochachacha/go-smb2"
	"net"
	"strings"
	"time"
)

type SmbPlugin struct {
	*pkg.Task
	conn    *smb2.Session
	Version string
	Input   string
}

func (s *SmbPlugin) Unauth() (bool, error) {
	user, domain := pkg.SplitUserDomain(s.Username)

	dialer := &smb2.Dialer{}
	dialer.Initiator = &smb2.NTLMInitiator{
		User:     user,
		Domain:   domain,
		Password: "",
	}

	c, err := net.DialTimeout("tcp", s.Address(), time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return false, err
	}

	conn, err := dialer.Dial(c)
	if err != nil {
		return false, err
	}
	// todo anon
	_, err = conn.ListSharenames()
	if err != nil {
		return false, err
	}
	s.conn = conn

	return true, nil
}

func (s *SmbPlugin) Login() error {
	var user, domain string

	if strings.Contains(s.Username, "/") {
		user = strings.Split(s.Username, "/")[1]
		domain = strings.Split(s.Username, "/")[0]
	} else {
		user = s.Username
	}

	dialer := &smb2.Dialer{}
	method, pwd := pkg.ParseMethod(s.Password)
	if method == "hash" {
		dialer.Initiator = &smb2.NTLMInitiator{
			User:   user,
			Domain: domain,
			Hash:   encode.HexDecode(pwd),
		}
	} else {
		dialer.Initiator = &smb2.NTLMInitiator{
			User:     user,
			Domain:   domain,
			Password: s.Password,
		}
	}

	c, err := net.DialTimeout("tcp", s.Address(), time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return err
	}

	conn, err := dialer.Dial(c)
	if err != nil {
		return err
	}
	// todo anon
	_, err = conn.ListSharenames()
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *SmbPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Logoff()
	}
	return nil
}

func (s *SmbPlugin) Name() string {
	return s.Service
}

func (s *SmbPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}
