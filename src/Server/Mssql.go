package Server

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"time"
)

func MssqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;connection timeout=%v;encrypt=disable", info.Ip,
		info.Port, User, Password, "master", time.Duration(Utils.Timeout)*time.Second)

	db, err = sql.Open("mssql", dataSourceName)

	if err != nil {
		result = false
	}
	return err, result, db
}

func MssqlConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, db := MssqlConnect(User, Password, info)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result.Result = res
		}
	}

	return err, result
}

func MssqlQuery(User string, Password string, info Utils.IpInfo, Query string) (err error, Qresult []map[string]string, Columns []string) {
	err, _, db := MssqlConnect(User, Password, info)
	if err != nil {
		fmt.Println("connect failed,please check your input.")
	} else {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			rows, err := db.Query(Query)
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
	}
	return err, Qresult, Columns
}
