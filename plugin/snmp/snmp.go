package snmp

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/gosnmp/gosnmp"
	"time"
)

// snmpSession implements pkg.Session over an SNMP connection.
type snmpSession struct {
	service string
	conn    *gosnmp.GoSNMP
}

func (s *snmpSession) Service() string  { return s.service }
func (s *snmpSession) Raw() interface{} { return s.conn }

func (s *snmpSession) Close() error {
	if s.conn != nil {
		return s.conn.Conn.Close()
	}
	return nil
}

// SnmpPlugin is stateless; all connection state lives in snmpSession.
type SnmpPlugin struct{}

func (p *SnmpPlugin) Name() string { return "snmp" }

func (p *SnmpPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	return dial(task, task.Password)
}

func (p *SnmpPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "")
}

func dial(task *pkg.Task, community string) (pkg.Session, error) {
	conn := &gosnmp.GoSNMP{
		Target:             task.IP,
		Port:               task.UintPort(),
		Community:          community,
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(task.Timeout) * time.Second,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	return &snmpSession{service: task.Service, conn: conn}, nil
}
