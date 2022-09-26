package plugin

import "fmt"

const (
	FTP = iota
	LDAP
	MSSQL
	MYSQL
	ORACLE
	POSTGRESQL
	RDP
	SMB
	SNMP
	SSH
	VNC
)

type Plugin interface {
	Query() bool
	GetInfo() bool
	Connect() error
	SetQuery(string)
	Output(interface{})
	Close() error
}

type NilConnError struct {
	service string
}

func (e NilConnError) Error() string {
	return e.service + "has nil conn"
}

type TimeoutError struct {
	err     error
	timeout int
	service string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("%s spended out of %ds, %s", e.service, e.timeout, e.err.Error())
}

func (e TimeoutError) Unwrap() error { return e.err }
