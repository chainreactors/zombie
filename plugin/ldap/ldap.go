package ldap

import (
	"errors"

	"github.com/chainreactors/zombie/pkg"
	ldap "github.com/go-ldap/ldap/v3"
)

// ldapSession implements pkg.DirectorySession over a bound LDAP connection.
type ldapSession struct {
	service string
	conn    *ldap.Conn
}

func (s *ldapSession) Service() string  { return s.service }
func (s *ldapSession) Raw() interface{} { return s.conn }

func (s *ldapSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *ldapSession) Search(baseDN, filter string, attrs []string) ([]map[string][]string, error) {
	req := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attrs,
		nil,
	)
	res, err := s.conn.Search(req)
	if err != nil {
		return nil, err
	}
	results := make([]map[string][]string, 0, len(res.Entries))
	for _, entry := range res.Entries {
		m := make(map[string][]string, len(entry.Attributes)+1)
		m["dn"] = []string{entry.DN}
		for _, attr := range entry.Attributes {
			m[attr.Name] = attr.Values
		}
		results = append(results, m)
	}
	return results, nil
}

// LdapPlugin is stateless; all connection state lives in ldapSession.
type LdapPlugin struct{}

func (p *LdapPlugin) Name() string { return "ldap" }

func (p *LdapPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	ldap.DefaultTimeout = task.Duration()
	conn, err := ldap.Dial("tcp", task.Address())
	if err != nil {
		return nil, err
	}
	if err := conn.Bind(task.Username, task.Password); err != nil {
		conn.Close()
		return nil, err
	}
	return &ldapSession{service: task.Service, conn: conn}, nil
}

func (p *LdapPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, errors.New("ldap: unauthenticated access not supported")
}
