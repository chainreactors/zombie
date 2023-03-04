package plugin

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	ldap "github.com/go-ldap/ldap/v3"
	"time"
)

type LdapService struct {
	*pkg.Task
	Input   string
	LdapCon *ldap.Conn
}

func (s *LdapService) Query() bool {
	panic("implement me")
}

func LdapConnect(info *pkg.Task) (con *ldap.Conn, err error) {
	var conn *ldap.Conn
	connectAddr := fmt.Sprintf(info.Address())

	ldap.DefaultTimeout = time.Duration(info.Timeout)
	conn, err = ldap.Dial("tcp", connectAddr)

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
