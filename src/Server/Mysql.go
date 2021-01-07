package Server

import (
	"Zombie/src/Utils"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func MysqlConnect(User string, Password string,info Utils.IpInfo)(err error,result bool){
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", User,
		Password, info.Ip, info.Port, "mysql")

	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		var bgCtx = context.Background()
		var ctx2SecondTimeout, cancelFunc2SecondTimeout = context.WithTimeout(bgCtx, time.Second*2)
		defer cancelFunc2SecondTimeout()
		defer db.Close()
		err = db.PingContext(ctx2SecondTimeout)
		if err == nil {
			result = true
		}
	}
	return err, result
}