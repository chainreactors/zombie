package mysql

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type nilLog struct {
}

func (l nilLog) Print(v ...interface{}) {

}

type MysqlPlugin struct {
	*pkg.Task
	input string
	conn  *sql.DB
}

func (s *MysqlPlugin) Name() string {
	return s.Service
}

func (s *MysqlPlugin) Unauth() (bool, error) {
	// mysql none pass
	mysql.SetLogger(nilLog{})
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8", "root",
		"", s.IP, s.Port, s.Timeout, s.Timeout, s.Timeout)
	conn, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return false, err
	}

	//conn.SetMaxOpenConns(60)
	//conn.SetMaxIdleConns(60)

	err = conn.Ping()
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

func (s *MysqlPlugin) Login() error {
	mysql.SetLogger(nilLog{})
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8", s.Username,
		s.Password, s.IP, s.Port, s.Timeout, s.Timeout, s.Timeout)
	conn, err := sql.Open("mysql", dataSourceName)
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

func (s *MysqlPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *MysqlPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}
