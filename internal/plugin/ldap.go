package plugin

import (
	"crypto/tls"
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	"github.com/go-ldap/ldap/v3"
)

type LdapService struct {
	*utils.Task
	Input   string
	LdapCon *ldap.Conn
}

func (s *LdapService) Query() bool {
	panic("implement me")
}

func LdapConnect(info *utils.Task) (con *ldap.Conn, err error) {
	enableTLS := true

	var conn *ldap.Conn
	connectAddr := fmt.Sprintf(info.Address())

	if enableTLS {
		conn, err = ldap.DialTLS("tcp", connectAddr, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", connectAddr)
	}
	if err != nil {
		return nil, err
	}

	err = conn.Bind(info.Username, info.Password)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *LdapService) Connect() error {
	conn, err := LdapConnect(s.Task)
	if err != nil {
		return err
	}
	s.LdapCon = conn
	return nil
}

func (s *LdapService) Close() error {
	if s.LdapCon != nil {
		s.LdapCon.Close()
		return nil
	}
	return NilConnError{s.Service}
}

func (s *LdapService) SetQuery(query string) {
	s.Input = query
}

func (s *LdapService) Output(res interface{}) {

}

func (s *LdapService) GetInfo() bool {
	defer s.LdapCon.Close()

	return true
}
