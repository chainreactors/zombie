package ExecAble

import (
	"Zombie/src/Utils"
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
)

type LdapService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
	LdapCon  *ldap.Conn
}

func (s *LdapService) Query() bool {
	panic("implement me")
}

func LdapConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, con *ldap.Conn) {

	enableTLS := true

	var ldapCon *ldap.Conn
	connectAddr := fmt.Sprintf("%s:%d", info.Ip, info.Port)

	if enableTLS {
		ldapCon, err = ldap.DialTLS("tcp", connectAddr, &tls.Config{InsecureSkipVerify: true})
	} else {
		ldapCon, err = ldap.Dial("tcp", connectAddr)
	}

	err = ldapCon.Bind(User, Password)

	if err == nil {
		result = true
	}
	return err, result, ldapCon
}

func (s *LdapService) Connect() bool {
	err, _, ldapCon := LdapConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		return true
	}
	s.LdapCon = ldapCon
	return false
}

func (s *LdapService) DisConnect() bool {
	s.LdapCon.Close()
	return false
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
