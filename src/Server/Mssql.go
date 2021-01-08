package Server

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func MssqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool) {
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v", info.Ip,
		info.Port, User, Password, "master")

	db, err := sql.Open("mssql", dataSourceName)

	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result = true
		}
	}
	return err, result
}
