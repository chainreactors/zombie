package mssql

import (
	"database/sql"
	"fmt"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin/internal/sqlsess"
	_ "github.com/denisenkom/go-mssqldb"
)

// MssqlPlugin is a stateless factory that satisfies the Plugin interface.
type MssqlPlugin struct{}

func (MssqlPlugin) Name() string { return "mssql" }

// Open authenticates with the credentials from task and returns a SQLSession.
func (MssqlPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	instance := task.Param["instance"]
	if instance == "" {
		instance = "master"
	}
	return dial(task, task.Username, task.Password, instance)
}

// Unauth attempts an unauthenticated connection using sa with an empty password.
func (MssqlPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "sa", "", "master")
}

func dial(task *pkg.Task, user, password, instance string) (pkg.Session, error) {
	dsn := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;connection timeout=%d;encrypt=disable",
		task.IP, task.Port, user, password, instance, task.Timeout)

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &sqlsess.Session{DB: db, SvcName: "mssql"}, nil
}
