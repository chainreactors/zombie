package service

import (
	"fmt"

	"github.com/chainreactors/neutron/operators"
	"github.com/chainreactors/neutron/protocols"
)

const ServiceProtocol protocols.ProtocolType = 6

var _ protocols.Request = &Request{}

// RawCommander is an optional interface for sessions that support
// arbitrary command execution beyond their typed interface (e.g., Redis CONFIG/SLAVEOF).
type RawCommander interface {
	Command(name string, args ...string) (interface{}, error)
}

// Request implements protocols.Request for the service protocol.
type Request struct {
	operators.Operators `json:",inline" yaml:",inline"`

	ID  string `json:"id,omitempty" yaml:"id,omitempty"`
	Ops []*Op  `json:"ops" yaml:"ops"`

	AttackType string                 `json:"attack,omitempty" yaml:"attack,omitempty"`
	Payloads   map[string]interface{} `json:"payloads,omitempty" yaml:"payloads,omitempty"`

	StopAtFirstMatch bool `json:"stop-at-first-match,omitempty" yaml:"stop-at-first-match,omitempty"`

	CompiledOperators *operators.Operators       `json:"-" yaml:"-"`
	options           *protocols.ExecuterOptions `json:"-" yaml:"-"`
}

// Op is a single operation against a session.
// Each Op sets exactly one session-type field — the type is inferred from which field is non-empty.
type Op struct {
	Shell string  `json:"shell,omitempty" yaml:"shell,omitempty"` // ShellSession: command to execute
	DB    string  `json:"db,omitempty" yaml:"db,omitempty"`       // SQLSession: SQL query to execute
	KV    string  `json:"kv,omitempty" yaml:"kv,omitempty"`       // KVSession: command expression (GET key / KEYS * / CONFIG SET ...)
	File  *FileOp `json:"file,omitempty" yaml:"file,omitempty"`   // FileSession: list or read
	LDAP  *LDAPOp `json:"ldap,omitempty" yaml:"ldap,omitempty"`   // DirectorySession: search

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Legacy aliases kept so existing service templates continue to load.
	Exec      string  `json:"exec,omitempty" yaml:"exec,omitempty"`
	Query     string  `json:"query,omitempty" yaml:"query,omitempty"`
	Databases bool    `json:"databases,omitempty" yaml:"databases,omitempty"`
	Get       string  `json:"get,omitempty" yaml:"get,omitempty"`
	Keys      string  `json:"keys,omitempty" yaml:"keys,omitempty"`
	Cmd       string  `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	List      string  `json:"list,omitempty" yaml:"list,omitempty"`
	Read      string  `json:"read,omitempty" yaml:"read,omitempty"`
	Search    *LDAPOp `json:"search,omitempty" yaml:"search,omitempty"`
}

// FileOp specifies a file session operation. Set exactly one field.
type FileOp struct {
	List string `json:"list,omitempty" yaml:"list,omitempty"`
	Read string `json:"read,omitempty" yaml:"read,omitempty"`
}

// LDAPOp specifies an LDAP search operation.
type LDAPOp struct {
	BaseDN string   `json:"base-dn" yaml:"base-dn"`
	Filter string   `json:"filter" yaml:"filter"`
	Attrs  []string `json:"attrs,omitempty" yaml:"attrs,omitempty"`
}

func (r *Request) Type() protocols.ProtocolType {
	return ServiceProtocol
}

func (r *Request) GetID() string {
	return r.ID
}

func (r *Request) Requests() int {
	return len(r.Ops)
}

func (r *Request) GetCompiledOperators() []*operators.Operators {
	return []*operators.Operators{r.CompiledOperators}
}

func (r *Request) Compile(options *protocols.ExecuterOptions) error {
	r.options = options
	if len(r.Matchers) > 0 || len(r.Extractors) > 0 {
		compiled := &r.Operators
		if err := compiled.Compile(); err != nil {
			return fmt.Errorf("could not compile operators: %w", err)
		}
		r.CompiledOperators = compiled
	}
	return nil
}
