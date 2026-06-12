package mysql

import (
	"database/sql"
	"fmt"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin/internal/sqlsess"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type nilLog struct{}

func (nilLog) Print(v ...interface{}) {}

// MysqlPlugin is a stateless factory that satisfies the Plugin interface.
type MysqlPlugin struct{}

func (MysqlPlugin) Name() string { return "mysql" }

// Open authenticates with the credentials from task and returns a SQLSession.
func (MysqlPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	return dial(task, task.Username, task.Password)
}

// Unauth attempts an unauthenticated connection (root with empty password).
func (MysqlPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "root", "")
}

// dial builds a MySQL DSN, connects, pings, and wraps the *sql.DB in a
// sqlsess.Session so it satisfies pkg.SQLSession.
func dial(task *pkg.Task, user, pass string) (pkg.Session, error) {
	mysql.SetLogger(nilLog{})

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8",
		user, pass, task.IP, task.Port, task.Timeout, task.Timeout, task.Timeout)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &sqlsess.Session{
		DB:      db,
		SvcName: "mysql",
	}, nil
}
