package Database

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"strings"
)

type PostgresService struct {
	Username string
	Password string
	Dbname   string
	Utils.IpInfo
	Utils.PostgreInf
	Input  string
	SqlCon *sql.DB
}

var PostgresCollectInfo string

func (s *PostgresService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetBaseInfo(s.SqlCon)
	res.Count = GetPostgresSummary(s)
	PostgresCollectInfo = ""
	if res != nil {
		PostgresCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", s.Ip, "Postgres", res.Version, res.OS, res.Count)
		PostgresCollectInfo += "\n"
	}
	//将结果放入管道
	Utils.QDatach <- PostgresCollectInfo
	return true
}

func (s *PostgresService) SetQuery(query string) {
	s.Input = query
}

func (s *PostgresService) SetDbname(db string) {
	s.Dbname = db
}

func PostgresConnect(User string, Password string, info Utils.IpInfo, dbname string) (err error, result bool, db *sql.DB) {
	dataSourceName := strings.Join([]string{
		fmt.Sprintf("connect_timeout=%d", Utils.Timeout),
		fmt.Sprintf("dbname=%s", dbname),
		fmt.Sprintf("host=%v", info.Ip),
		fmt.Sprintf("password=%v", Password),
		fmt.Sprintf("port=%v", info.Port),
		"sslmode=disable",
		fmt.Sprintf("user=%v", User),
	}, " ")

	db, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		result = false
		return err, result, nil
	}

	err = db.Ping()
	if err == nil {
		result = true

	}

	return err, result, db

}

func PostgresConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, db := PostgresConnect(User, Password, info, "postgres")

	if err == nil {
		result.Result = res
		_ = db.Close()
	}

	return err, result
}

func PostgresQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {
	err = SqlCon.Ping()
	if err == nil {
		rows, err := SqlCon.Query(Query)
		if err == nil {
			Qresult, Columns = DoRowsMapper(rows)

		} else {
			if !Utils.IsAuto {
				fmt.Println("please check your query.")
			}

			return err, Qresult, Columns
		}
	} else {
		fmt.Println("connect failed,please check your input.")
		return err, Qresult, Columns
	}

	return err, Qresult, Columns
}

func (s *PostgresService) Query() bool {

	defer s.SqlCon.Close()
	err, Qresult, Columns := PostgresQuery(s.SqlCon, s.Input)

	if err != nil {
		fmt.Println("something wrong")
		os.Exit(0)
	} else {
		Utils.OutPutQuery(Qresult, Columns, true)
	}

	return true
}

func (s *PostgresService) Connect() bool {
	err, _, db := PostgresConnect(s.Username, s.Password, s.IpInfo, s.Dbname)
	if err == nil {
		s.SqlCon = db
		return true
	}
	return false
}

func GetBaseInfo(SqlCon *sql.DB) *Utils.PostgreInf {

	res := Utils.PostgreInf{}

	err, Qresult, Columns := PostgresQuery(SqlCon, "SHOW server_version;")

	if err != nil {
		fmt.Println("something wrong")
		return nil
	}

	VerOs := Utils.GetSummary(Qresult, Columns)

	VerOs = strings.Replace(VerOs, "(", "", 1)
	VerOs = strings.Replace(VerOs, ")", "", 1)

	VerOsList := strings.Split(VerOs, " ")

	if len(VerOsList) < 2 {
		fmt.Println("something wrong in split")

		return nil
	}

	res.Version = VerOsList[0]
	res.OS = VerOsList[1]

	return &res
}

func GetPostgresSummary(s *PostgresService) int {

	var db []string
	var sum int

	err, Qresult, Columns := PostgresQuery(s.SqlCon, "SELECT datname FROM pg_database")

	for _, items := range Qresult {
		for _, cname := range Columns {
			db = append(db, items[cname])
		}
	}

	if err != nil {
		fmt.Println("something wrong")
		return 0
	}

	_, Qresult, Columns = PostgresQuery(s.SqlCon, "SELECT sum(n_live_tup) FROM pg_stat_user_tables")

	CurIntSum := Utils.GetSummary(Qresult, Columns)

	CurSum, err := strconv.Atoi(CurIntSum)

	if err == nil {
		sum += CurSum
	}

	s.SqlCon.Close()

	for _, dbname := range db {
		if dbname == "postgres" {
			continue
		}

		s.SetDbname(dbname)
		succ := s.Connect()
		if succ {
			_, Qresult, Columns = PostgresQuery(s.SqlCon, "SELECT sum(n_live_tup) FROM pg_stat_user_tables")

			CurIntSum = Utils.GetSummary(Qresult, Columns)

			CurSum, err = strconv.Atoi(CurIntSum)

			if err == nil {
				sum += CurSum
			}
			s.SqlCon.Close()
		}
	}

	return sum

}
