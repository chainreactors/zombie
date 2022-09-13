package plugin

import (
	"encoding/hex"
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/hirochachacha/go-smb2"

	//"github.com/hirochachacha/go-smb2"
	"net"
	"strings"
	"time"
)

type SmbService struct {
	*utils.Task
	Session *smb2.Session
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
	conn, err := SMBConnect(s.Task)
	if err != nil {
		return err
	}
	s.Session = conn
	return nil
}

func (s *SmbService) Close() error {
	if s.Session != nil {
		return s.Session.Logoff()
	}
	return NilConnError{s.Service}
}

func (s *SmbService) SetQuery(query string) {
	s.Input = query
}

func (s *SmbService) Output(res interface{}) {

}

func SMBConnect(info *utils.Task) (sess *smb2.Session, err error) {
	var user, domain string

	if strings.Contains(info.Username, "/") {
		user = strings.Split(info.Username, "/")[1]
		domain = strings.Split(info.Username, "/")[0]
	} else {
		user = info.Username
	}

	conn, err := net.DialTimeout("tcp", info.Address(), time.Duration(info.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	d := &smb2.Dialer{}
	if strings.HasPrefix(info.Password, "hash:") {
		hash := info.Password[5:]
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
			Password: info.Password,
		}
	}

	_ = conn.SetDeadline(time.Now().Add(time.Duration(utils.Timeout) * time.Second))

	s, _, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	//share, err := s.Mount("C$")
	//
	//fmt.Println(err.Error())
	return s, nil
}
