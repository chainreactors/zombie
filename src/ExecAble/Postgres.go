package ExecAble

import (
	"Zombie/src/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"strings"
)

type PostgresService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"Dbname"`
	PostgreInf
	Input  string
	SqlCon *sql.DB
}

type PostgreInf struct {
	Version string
	Count   int
	OS      string
}

var PostgresCollectInfo string

func (s *PostgresService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetPostBaseInfo(s.SqlCon)
	res.Count = GetPostgresSummary(s)
	s.PostgreInf = *res
	//将结果放入管道
	s.Output(*s)
	return true
}

func (s *PostgresService) SetQuery(query string) {
	s.Input = query
}

func (s *PostgresService) SetDbname(db string) {
	s.Dbname = db
}

func (s *PostgresService) Output(res interface{}) {
	finres := res.(PostgresService)
	PostCollectInfo := ""
	PostCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
	PostCollectInfo += "\n"
	fmt.Println(PostCollectInfo)
	switch Utils.FileFormat {
	case "raw":
		Utils.TDatach <- PostCollectInfo
	case "json":
		jsons, errs := json.Marshal(res)
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}
		Utils.TDatach <- jsons
	}
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
		OutPutQuery(Qresult, Columns, true)
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

func (s *PostgresService) DisConnect() bool {
	s.SqlCon.Close()
	return false
}

func GetPostBaseInfo(SqlCon *sql.DB) *PostgreInf {

	res := PostgreInf{}

	err, Qresult, Columns := PostgresQuery(SqlCon, "SHOW server_version;")

	if err != nil {
		fmt.Println("something wrong")
		return nil
	}

	VerOs := GetSummary(Qresult, Columns)

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

	CurIntSum := GetSummary(Qresult, Columns)

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

			CurIntSum = GetSummary(Qresult, Columns)

			CurSum, err = strconv.Atoi(CurIntSum)

			if err == nil {
				sum += CurSum
			}
			s.SqlCon.Close()
		}
	}

	return sum

}
