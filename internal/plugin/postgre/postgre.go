package postgre

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	_ "github.com/lib/pq"
	"strings"
)

type PostgresPlugin struct {
	*pkg.Task
	Dbname string
	//PostgreInf
	//Input string
	conn *sql.DB
}

func (s *PostgresPlugin) Login() error {
	if s.Dbname == "" {
		s.Dbname = "postgres"
	}
	dataSourceName := strings.Join([]string{
		fmt.Sprintf("connect_timeout=%d", s.Timeout),
		fmt.Sprintf("dbname=%s", s.Dbname),
		fmt.Sprintf("host=%v", s.IP),
		fmt.Sprintf("password=%v", s.Password),
		fmt.Sprintf("port=%v", s.Port),
		"sslmode=disable",
		fmt.Sprintf("user=%v", s.Username),
	}, " ")

	conn, err := sql.Open("postgres", dataSourceName)
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

func (s *PostgresPlugin) Unauth() (bool, error) {
	dataSourceName := strings.Join([]string{
		fmt.Sprintf("connect_timeout=%d", s.Timeout),
		fmt.Sprintf("dbname=%s", s.Dbname),
		fmt.Sprintf("host=%v", s.IP),
		fmt.Sprintf("password=%v", ""),
		fmt.Sprintf("port=%v", s.Port),
		"sslmode=disable",
		fmt.Sprintf("user=%v", ""),
	}, " ")

	conn, err := sql.Open("postgres", dataSourceName)
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

func (s *PostgresPlugin) Name() string {
	return s.Service.String()
}

func (s *PostgresPlugin) GetResult() *pkg.Result {
	// todo list dbs
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *PostgresPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

//func GetPostBaseInfo(SqlCon *sql.DB) *PostgreInf {
//
//	res := PostgreInf{}
//
//	err, Qresult, Columns := PostgresQuery(SqlCon, "SHOW server_version;")
//
//	if err != nil {
//		fmt.Println("something wrong")
//		return nil
//	}
//
//	VerOs := GetSummary(Qresult, Columns)
//
//	VerOs = strings.Replace(VerOs, "(", "", 1)
//	VerOs = strings.Replace(VerOs, ")", "", 1)
//
//	VerOsList := strings.Split(VerOs, " ")
//
//	if len(VerOsList) < 2 {
//		fmt.Println("something wrong in split")
//
//		return nil
//	}
//
//	res.Version = VerOsList[0]
//	res.OS = VerOsList[1]
//
//	return &res
//}
//
//func GetPostgresSummary(s *PostgresService) int {
//	var db []string
//	var sum int
//
//	err, Qresult, Columns := PostgresQuery(s.conn, "SELECT datname FROM pg_database")
//
//	for _, items := range Qresult {
//		for _, cname := range Columns {
//			db = append(db, items[cname])
//		}
//	}
//
//	if err != nil {
//		fmt.Println("something wrong")
//		return 0
//	}
//
//	_, Qresult, Columns = PostgresQuery(s.conn, "SELECT sum(n_live_tup) FROM pg_stat_user_tables")
//	CurIntSum := GetSummary(Qresult, Columns)
//	CurSum, err := strconv.Atoi(CurIntSum)
//	if err == nil {
//		sum += CurSum
//	}
//
//	s.conn.Close()
//
//	for _, dbname := range db {
//		if dbname == "postgres" {
//			continue
//		}
//
//		s.SetDbname(dbname)
//		err := s.Connect()
//		if err == nil {
//			_, Qresult, Columns = PostgresQuery(s.conn, "SELECT sum(n_live_tup) FROM pg_stat_user_tables")
//			CurIntSum = GetSummary(Qresult, Columns)
//			CurSum, err = strconv.Atoi(CurIntSum)
//			if err == nil {
//				sum += CurSum
//			}
//			s.conn.Close()
//		}
//	}
//
//	return sum
//}
