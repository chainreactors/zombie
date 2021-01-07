package Server

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func PostgresConnect(User string, Password string,info Utils.IpInfo)(err error,result bool){
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", User,
		Password, info.Ip, info.Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)

	if err == nil {
		sqlStatement := "SELECT datname FROM pg_database"
		rows, err := db.Query(sqlStatement)
		if err != nil{
			return err ,result
		}
		result = true
		defer rows.Close()

		defer db.Close()
		err = db.Ping()
		if err == nil {
			result = true
		}
	}
	return err, result
}