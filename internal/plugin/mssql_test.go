package plugin

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestScanMssql(t *testing.T) {
	//task := pkg.Task{"192.168.91.129", "1433", "MSSQL", "sa", "admin@123"}
	//s := MssqlService{task}
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;connection timeout=%v;encrypt=disable", "192.168.91.129",
		"1433", "sa", "admin@123", "master", 5)

	conn, err := sql.Open("mssql", dataSourceName)
	err = conn.Ping()
	fmt.Println(err)

}
