package plugin

import (
	"github.com/chainreactors/zombie/pkg"
	ldap "github.com/go-ldap/ldap/v3"
)

type LdapService struct {
	*pkg.Task
	Input string
	conn  *ldap.Conn
}

func (s *LdapService) Query() bool {
	panic("implement me")
}

func (s *LdapService) Connect() error {
	var conn *ldap.Conn
	ldap.DefaultTimeout = s.Duration()
	conn, err := ldap.Dial("tcp", s.Address())

	if err != nil {
		return err
	}

	err = conn.Bind(s.Username, s.Password)
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *LdapService) Close() error {
	if s.conn != nil {
		s.conn.Close()
		return nil
	}
	return pkg.NilConnError{s.Service}
}

func (s *LdapService) SetQuery(query string) {
	s.Input = query
}

func (s *LdapService) Output(res interface{}) {

}

func (s *LdapService) GetInfo() bool {
	s.conn.Close()
	return true
}
