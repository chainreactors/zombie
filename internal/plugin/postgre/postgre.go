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
	Dbname string `json:"Dbname"`
	//PostgreInf
	Input string
	conn  *sql.DB
}

//type PostgreInf struct {
//	Version string
//	Count   int
//	OS      string
//}
//
//var PostgresCollectInfo string

//func (s *PostgresPlugin) GetInfo() bool {
//	//res := GetPostBaseInfo(s.conn)
//	//res.Count = GetPostgresSummary(s)
//	//s.PostgreInf = *res
//	////将结果放入管道
//	//s.Output(*s)
//	return true
//}
//
//func (s *PostgresPlugin) SetQuery(query string) {
//	s.Input = query
//}
//
//func (s *PostgresPlugin) SetDbname(db string) {
//	s.Dbname = db
//}

//func (s *PostgresPlugin) Output(res interface{}) {
//	//finres := res.(PostgresService)
//	//PostCollectInfo := ""
//	//PostCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", finres.IP, utils.OutputType, finres.Version, finres.OS, finres.Count)
//	//PostCollectInfo += "\n"
//	//fmt.Println(PostCollectInfo)
//	//switch utils.FileFormat {
//	//case "raw":
//	//	utils.TDatach <- PostCollectInfo
//	//case "json":
//	//	jsons, errs := json.Marshal(res)
//	//	if errs != nil {
//	//		fmt.Println(errs.Error())
//	//		return
//	//	}
//	//	utils.TDatach <- jsons
//	//}
//}

//func PostgresQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {
//	err = SqlCon.Ping()
//	if err == nil {
//		rows, err := SqlCon.Query(Query)
//		if err == nil {
//			Qresult, Columns = DoRowsMapper(rows)
//
//		} else {
//			if !utils.IsAuto {
//				fmt.Println("please check your query.")
//			}
//
//			return err, Qresult, Columns
//		}
//	} else {
//		fmt.Println("connect failed,please check your input.")
//		return err, Qresult, Columns
//	}
//
//	return err, Qresult, Columns
//}

//func (s *PostgresPlugin) Query() bool {
//
//	//defer s.conn.Close()
//	//err, Qresult, Columns := PostgresQuery(s.conn, s.Input)
//	//
//	//if err != nil {
//	//	fmt.Println("something wrong")
//	//	os.Exit(0)
//	//} else {
//	//	OutPutQuery(Qresult, Columns, true)
//	//}
//
//	return true
//}

func (s *PostgresPlugin) Login() error {
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

func (s *PostgresPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
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
