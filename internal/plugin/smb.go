package plugin

import (
	"encoding/hex"
	"github.com/chainreactors/zombie/pkg"
	"github.com/hirochachacha/go-smb2"
	//"github.com/hirochachacha/go-smb2"
	"net"
	"strings"
	"time"
)

type SmbService struct {
	*pkg.Task
	conn    *smb2.Session
	Version string
	Input   string
}

func (s *SmbService) Query() bool {
	return false
}

func (s *SmbService) GetInfo() bool {
	return false
}

func (s *SmbService) Connect() error {
	var user, domain string

	if strings.Contains(s.Username, "/") {
		user = strings.Split(s.Username, "/")[1]
		domain = strings.Split(s.Username, "/")[0]
	} else {
		user = s.Username
	}

	c, err := net.DialTimeout("tcp", s.Address(), time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return err
	}

	d := &smb2.Dialer{}
	if strings.HasPrefix(s.Password, "hash:") {
		hash := s.Password[5:]
		buf := make([]byte, len(hash)/2)
		hex.Decode(buf, []byte(hash))
		d.Initiator = &smb2.NTLMInitiator{
			User:   user,
			Domain: domain,
			Hash:   buf,
		}
	} else {
		d.Initiator = &smb2.NTLMInitiator{
			User:   user,
			Domain: domain,
			//Hash: buf,
			Password: s.Password,
		}
	}

	_ = c.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second))

	conn, _, err := d.Dial(c)
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

func (s *SmbService) Close() error {
	if s.conn != nil {
		return s.conn.Logoff()
	}
	return NilConnError{s.Service}
}

func (s *SmbService) SetQuery(query string) {
	s.Input = query
}

func (s *SmbService) Output(res interface{}) {

}
