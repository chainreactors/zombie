package oracle

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	_ "github.com/sijms/go-ora/v2"
)

type OraclePlugin struct {
	*pkg.Task
	//Input string
	SID         string
	ServiceName string
	conn        *sql.DB
}

func (s *OraclePlugin) Unauth() (bool, error) {
	return false, nil
}

func (s *OraclePlugin) Login() error {
	var err error
	if s.ServiceName != "" {
		s.conn, err = serviceNameLogin(s.Task, s.ServiceName)
	} else {
		s.conn, err = sidLogin(s.Task, s.SID)
	}

	err = s.conn.Ping()
	if err != nil {
		return err
	}

	return err
}

func (s *OraclePlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

func (s *OraclePlugin) Name() string {
	return s.Service.String()
}

func (s *OraclePlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

func sidLogin(task *pkg.Task, sid string) (*sql.DB, error) {
	if sid == "" {
		sid = "orcl"
	}
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s?connection_timeout=%d&connection_pool_timeout=%d", task.Username,
		task.Password, task.IP, task.Port, sid, task.Timeout, task.Timeout)

	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func serviceNameLogin(task *pkg.Task, serviceName string) (*sql.DB, error) {
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/?service_name=%s&connection_timeout=%d&connection_pool_timeout=%d", task.Username,
		task.Password, task.IP, task.Port, serviceName, task.Timeout, task.Timeout)

	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
