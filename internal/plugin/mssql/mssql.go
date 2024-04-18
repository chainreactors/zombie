package mssql

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlPlugin struct {
	*pkg.Task
	//MssqlInf
	//Input string
	Instance string
	conn     *sql.DB
}

func (s *MssqlPlugin) Name() string {
	return s.Service
}

func (s *MssqlPlugin) Login() error {
	if s.Instance == "" {
		s.Instance = "master"
	}
	dataSourceName := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;connection timeout=%d;encrypt=disable", s.IP,
		s.Port, s.Username, s.Password, s.Instance, s.Timeout)

	//time.Duration(Utils.Timeout)*time.Second
	conn, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}

func (s *MssqlPlugin) Unauth() (bool, error) {
	dataSourceName := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%v;database=%v;connection timeout=%v;encrypt=disable", s.IP,
		s.Port, "sa", "", "master", s.Timeout)

	//time.Duration(Utils.Timeout)*time.Second
	conn, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return false, err
	}

	err = conn.Ping()
	if err != nil {
		return false, err
	}

	s.conn = conn
	return true, nil
}
func (s *MssqlPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *MssqlPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
