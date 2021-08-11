package ExecAble

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
	"strings"
)

type MssqlService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	MssqlInf
	Input  string
	SqlCon *sql.DB
}

type MssqlInf struct {
	Version     string `json:"version"`
	Count       int    `json:"count"`
	OS          string `json:"os"`
	XpCmdShell  string `json:"xp_cmdshell"`
	SP_OACREATE string `json:"sp_oacreate"`
}

var MssqlCollectInfo string

func MssqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;connection timeout=%v;encrypt=disable", info.Ip,
		info.Port, User, Password, "master", Utils.Timeout)

	//time.Duration(Utils.Timeout)*time.Second
	db, err = sql.Open("mssql", dataSourceName)

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

func MssqlConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, db := MssqlConnect(User, Password, info)

	if err == nil {
		result.Result = res
		_ = db.Close()
	}

	return err, result

}

func MssqlQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {

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

func (s *MssqlService) Query() bool {

	defer s.SqlCon.Close()
	err, Qresult, Columns := MssqlQuery(s.SqlCon, s.Input)

	if err != nil {
		fmt.Println("something wrong")
		return false
	} else {
		OutPutQuery(Qresult, Columns, true)
	}

	return true
}

func (s *MssqlService) SetQuery(query string) {
	s.Input = query
}

func (s *MssqlService) Connect() bool {
	err, _, db := MssqlConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SqlCon = db
		return true
	}
	return false
}

func (s *MssqlService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetMssqlBaseInfo(s.SqlCon)

	if res == nil {
		return false
	}

	MssqlCollectInfo = ""

	res.Count = GetMssqlSummary(s.SqlCon)

	res = GetMssqlVulnableInfo(s.SqlCon, res)
	s.MssqlInf = *res
	//将结果放入管道
	Utils.TDatach <- *s

	return true
}

func GetMssqlSummary(SqlCon *sql.DB) int {

	var db []string
	var sum int

	err, Qresult, Columns := MssqlQuery(SqlCon, "select name from sysdatabases where dbid>4")

	if err != nil {
		fmt.Println("something wrong")
		return 0
	}

	for _, items := range Qresult {
		for _, cname := range Columns {
			db = append(db, items[cname])
		}
	}

	for _, dbname := range db {
		curinput := fmt.Sprintf("use %s;select SUM(i.rows )as [RowCount] from sys.tables as t, sysindexes as i where t.object_id = i.id and i.indid <=1;", dbname)

		_, Qresult, Columns = MssqlQuery(SqlCon, curinput)

		CurIntSum := GetSummary(Qresult, Columns)

		CurSum, err := strconv.Atoi(CurIntSum)

		if err != nil {
			continue
		}

		sum += CurSum

	}

	return sum

}

func GetMssqlBaseInfo(SqlCon *sql.DB) *MssqlInf {

	res := MssqlInf{}

	err, Qresult, Columns := MssqlQuery(SqlCon, "select @@version")
	if err != nil {
		fmt.Println("something wrong at get version")
		return nil
	}
	info := GetSummary(Qresult, Columns)

	infolist := strings.Split(info, "\n")

	for _, in := range infolist {
		if strings.Contains(in, "SQL") {
			res.Version = in
		} else if strings.Contains(in, "bit") {
			res.OS = in
		}
	}

	return &res
}

func GetMssqlVulnableInfo(SqlCon *sql.DB, res *MssqlInf) *MssqlInf {
	err, Qresult, Columns := MssqlQuery(SqlCon, "select count(*) from master.dbo.sysobjects where xtype='x' and name='xp_cmdshell'")
	if err != nil {
		//fmt.Println("something wrong in get xp_cmdshell")
	} else {
		info := GetSummary(Qresult, Columns)

		if info == "1" {
			res.XpCmdShell = "exsit"
		} else {
			res.XpCmdShell = "none"
		}
	}

	err, Qresult, Columns = MssqlQuery(SqlCon, "select count(*) from master.dbo.sysobjects where xtype='x' and name='SP_OACREATE'")
	if err != nil {
		//fmt.Println("something wrong in get SP_OACREATE")
	} else {
		info := GetSummary(Qresult, Columns)

		if info == "1" {
			res.SP_OACREATE = "exsit"
		} else {
			res.SP_OACREATE = "none"
		}
	}
	return res

}
