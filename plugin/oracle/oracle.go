package oracle

import (
	"database/sql"
	"fmt"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin/internal/sqlsess"
	_ "github.com/sijms/go-ora/v2"
)

// OraclePlugin is a stateless factory that satisfies the Plugin interface.
type OraclePlugin struct{}

func (OraclePlugin) Name() string { return "oracle" }

// Open authenticates with the credentials from task and returns a SQLSession.
// It supports two modes: service_name (if task.Param["service_name"] is set)
// or SID (task.Param["sid"], defaulting to "orcl").
func (OraclePlugin) Open(task *pkg.Task) (pkg.Session, error) {
	if sn := task.Param["service_name"]; sn != "" {
		return dialServiceName(task, sn)
	}
	sid := task.Param["sid"]
	if sid == "" {
		sid = "orcl"
	}
	return dialSID(task, sid)
}

// Unauth is not implemented for Oracle.
func (OraclePlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return nil, pkg.NotImplUnauthorized
}

func dialSID(task *pkg.Task, sid string) (pkg.Session, error) {
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s?connection_timeout=%d&connection_pool_timeout=%d",
		task.Username, task.Password, task.IP, task.Port, sid, task.Timeout, task.Timeout)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &sqlsess.Session{DB: db, SvcName: "oracle"}, nil
}

func dialServiceName(task *pkg.Task, serviceName string) (pkg.Session, error) {
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/?service_name=%s&connection_timeout=%d&connection_pool_timeout=%d",
		task.Username, task.Password, task.IP, task.Port, serviceName, task.Timeout, task.Timeout)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &sqlsess.Session{DB: db, SvcName: "oracle"}, nil
}
