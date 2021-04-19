package Server

import "C"
import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type MysqlService struct {
	Username string
	Password string
	Utils.IpInfo
	Utils.MysqlInf
	Input  string
	SqlCon *sql.DB
}

func MysqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8", User,
		Password, info.Ip, info.Port, Utils.Timeout, Utils.Timeout, Utils.Timeout)
	db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		result = false
	}

	db.SetMaxOpenConns(60)
	db.SetMaxIdleConns(60)

	err = db.Ping()

	if err == nil {
		result = true
	}

	return err, result, db
}

func MysqlConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, db := MysqlConnect(User, Password, info)

	if err == nil {
		result.Result = res
	}

	_ = db.Close()

	return err, result
}

func (s *MysqlService) Connect() bool {
	err, _, db := MysqlConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SqlCon = db
		return true
	}
	return false
}

func MysqlQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {

	SqlCon.SetMaxOpenConns(60)
	SqlCon.SetMaxIdleConns(60)

	err = SqlCon.Ping()
	if err == nil {
		rows, err := SqlCon.Query(Query)
		if err == nil {
			Qresult, Columns = DoRowsMapper(rows)

		} else {
			fmt.Println("please check your query.")
			return err, Qresult, Columns
		}
	} else {
		fmt.Println("connect failed,please check your input.")
		return err, Qresult, Columns
	}

	return err, Qresult, Columns
}

func (s *MysqlService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetBaseInfo(s.SqlCon)

	res.Count = GetSummary(s.SqlCon)

	fmt.Printf("IP:%v\tServer:%v\nVersion:%v\tOS:%v\nSummary:%v", s.Ip, "Mysql", res.Version, res.OS, res.Count)
	return true
}

func (s *MysqlService) Query() bool {
	err, Qresult, Columns := MysqlQuery(s.SqlCon, s.Input)

	if err != nil {
		fmt.Println("something wrong")
		os.Exit(0)
	} else {
		Utils.OutPutQuery(Qresult, Columns)
	}
	return true
}

func (s *MysqlService) SetQuery(query string) {
	s.Input = query
}

func GetBaseInfo(SqlCon *sql.DB) Utils.MysqlInf {
	err, Qresult, Columns := MysqlQuery(SqlCon, "select VERSION(),@@version_compile_os")
	if err != nil {
		fmt.Println("something wrong")
	}
	res := Utils.GetBaseInfo(Qresult, Columns)
	return res
}

func GetSummary(SqlCon *sql.DB) string {
	err, Qresult, Columns := MysqlQuery(SqlCon, "select sum(table_rows) from  information_schema.tables where table_rows is not null")

	if err == nil {
		Count := Utils.GetSummary(Qresult, Columns)
		return Count
	}
	return ""
}
