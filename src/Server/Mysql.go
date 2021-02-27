package Server

import (
	"Zombie/src/Utils"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func MysqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", User,
		Password, info.Ip, info.Port, "mysql")

	db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		result = false
	}

	return err, result, db
}

func MysqlConnectTest(User string, Password string, info Utils.IpInfo) (err error, result bool) {
	err, result, db := MysqlConnect(User, Password, info)

	if err == nil {
		defer db.Close()
		var bgCtx = context.Background()
		var ctx2SecondTimeout, cancelFunc2SecondTimeout = context.WithTimeout(bgCtx, time.Second*2)
		defer cancelFunc2SecondTimeout()
		err = db.PingContext(ctx2SecondTimeout)
		if err == nil {
			result = true
		}
	}

	return err, result
}

func MysqlQuery(User string, Password string, info Utils.IpInfo, Query string) (err error, Qresult []map[string]string, Columns []string) {
	err, _, db := MysqlConnect(User, Password, info)

	if err != nil {
		fmt.Println("connect failed,please check your input.")
	} else {
		defer db.Close()
		var bgCtx = context.Background()
		var ctx2SecondTimeout, cancelFunc2SecondTimeout = context.WithTimeout(bgCtx, time.Second*2)
		defer cancelFunc2SecondTimeout()
		err = db.PingContext(ctx2SecondTimeout)
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
