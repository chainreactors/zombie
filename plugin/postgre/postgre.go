package postgre

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/chainreactors/zombie/pkg"
	"github.com/chainreactors/zombie/plugin/internal/sqlsess"
	_ "github.com/lib/pq"
)

// PostgresPlugin is stateless; all connection state lives in sqlsess.Session.
type PostgresPlugin struct{}

func (PostgresPlugin) Name() string { return "postgresql" }

// Open authenticates with the credentials from task and returns a SQLSession.
func (PostgresPlugin) Open(task *pkg.Task) (pkg.Session, error) {
	return dial(task, task.Username, task.Password)
}

// Unauth attempts an unauthenticated connection (empty user and password).
func (PostgresPlugin) Unauth(task *pkg.Task) (pkg.Session, error) {
	return dial(task, "", "")
}

// dial builds a lib/pq DSN, opens and pings the database, then wraps it in a
// sqlsess.Session with SvcName "postgresql".
func dial(task *pkg.Task, user, password string) (pkg.Session, error) {
	dbname := task.Param["dbname"]
	if dbname == "" {
		dbname = "postgres"
	}

	dsn := strings.Join([]string{
		fmt.Sprintf("host=%v", task.IP),
		fmt.Sprintf("port=%v", task.Port),
		fmt.Sprintf("user=%v", user),
		fmt.Sprintf("password=%v", password),
		fmt.Sprintf("dbname=%s", dbname),
		"sslmode=disable",
		fmt.Sprintf("connect_timeout=%d", task.Timeout),
	}, " ")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &sqlsess.Session{
		DB:      db,
		SvcName: "postgresql",
	}, nil
}
