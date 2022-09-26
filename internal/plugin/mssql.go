package plugin

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlService struct {
	*utils.Task
	MssqlInf
	Input string
	conn  *sql.DB
}

type MssqlValuable struct {
	STName     string
	ColumnName string
}

type MssqlInf struct {
	Version     string `json:"version"`
	Count       int    `json:"count"`
	OS          string `json:"os"`
	XpCmdShell  string `json:"xp_cmdshell"`
	SP_OACREATE string `json:"sp_oacreate"`
	vb          []MssqlValuable
}

func MssqlConnect(task *utils.Task) (conn *sql.DB, err error) {
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;connection timeout=%v;encrypt=disable", task.IP,
		task.Port, task.Username, task.Password, "master", task.Timeout)

	//time.Duration(Utils.Timeout)*time.Second
	conn, err = sql.Open("mssql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//func MssqlQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {
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
//			return err, Qresult, Columns
//		}
//	} else {
//		fmt.Println("connect failed,please check your input.")
//		return err, Qresult, Columns
//	}
//
//	return err, Qresult, Columns
//}

func (s *MssqlService) Query() bool {
	//err, Qresult, Columns := MssqlQuery(s.conn, s.Input)
	//
	//if err != nil {
	//	fmt.Println("something wrong")
	//	return false
	//} else {
	//	OutPutQuery(Qresult, Columns, true)
	//}

	return true
}

func (s *MssqlService) SetQuery(query string) {
	s.Input = query
}

func (s *MssqlService) Connect() error {
	conn, err := MssqlConnect(s.Task)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *MssqlService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return NilConnError{s.Service}
}

func (s *MssqlService) GetInfo() bool {
	//defer s.conn.Close()
	//
	//res := GetMssqlBaseInfo(s.conn)
	//
	//if res == nil {
	//	return false
	//}
	//
	//res.Count = GetMssqlSummary(s.conn)
	//
	//res = GetMssqlVulnableInfo(s.conn, res)
	//res.vb = *FindMssqlValuableTable(s.conn)
	//s.MssqlInf = *res
	////将结果放入管道
	//s.Output(*s)
	//
	return true
}

func (s *MssqlService) Output(res interface{}) {
	//finres := res.(MssqlService)
	//MsCollectInfo := ""
	//MsCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", finres.IP, utils.OutputType, finres.Version, finres.OS, finres.Count)
	//MsCollectInfo += fmt.Sprintf("\nSP_OACREATE: %v", finres.SP_OACREATE)
	//MsCollectInfo += fmt.Sprintf("\nxp_cmdshell: %v\n", finres.XpCmdShell)
	//for _, info := range finres.vb {
	//	MsCollectInfo += fmt.Sprintf("%v:%v\t", info.STName, info.ColumnName)
	//}
	//MsCollectInfo += "\n"
	//fmt.Println(MsCollectInfo)
	//switch utils.FileFormat {
	//case "raw":
	//	utils.TDatach <- MsCollectInfo
	//case "json":
	//
	//	jsons, errs := json.Marshal(res)
	//	if errs != nil {
	//		fmt.Println(errs.Error())
	//		return
	//	}
	//	utils.TDatach <- jsons
	//}
}

//func GetMssqlSummary(SqlCon *sql.DB) int {
//
//	var db []string
//	var sum int
//
//	err, Qresult, Columns := MssqlQuery(SqlCon, "select name from sysdatabases where dbid>4")
//
//	if err != nil {
//		fmt.Println("something wrong")
//		return 0
//	}
//
//	for _, items := range Qresult {
//		for _, cname := range Columns {
//			db = append(db, items[cname])
//		}
//	}
//
//	for _, dbname := range db {
//		curinput := fmt.Sprintf("use %s;select SUM(i.rows )as [RowCount] from sys.tables as t, sysindexes as i where t.object_id = i.id and i.indid <=1;", dbname)
//
//		_, Qresult, Columns = MssqlQuery(SqlCon, curinput)
//
//		CurIntSum := GetSummary(Qresult, Columns)
//
//		CurSum, err := strconv.Atoi(CurIntSum)
//
//		if err != nil {
//			continue
//		}
//
//		sum += CurSum
//
//	}
//
//	return sum
//
//}

//func GetMssqlBaseInfo(SqlCon *sql.DB) *MssqlInf {
//
//	res := MssqlInf{}
//
//	err, Qresult, Columns := MssqlQuery(SqlCon, "select @@version")
//	if err != nil {
//		fmt.Println("something wrong at get version")
//		return nil
//	}
//
//	info := GetSummary(Qresult, Columns)
//
//	infolist := strings.Split(info, "\n")
//
//	for _, in := range infolist {
//		if strings.Contains(in, "SQL") {
//			res.Version = in
//		} else if strings.Contains(in, "bit") {
//			res.OS = in
//		}
//	}
//
//	return &res
//}
//
//func GetMssqlVulnableInfo(SqlCon *sql.DB, res *MssqlInf) *MssqlInf {
//	err, Qresult, Columns := MssqlQuery(SqlCon, "select count(*) from master.dbo.sysobjects where xtype='x' and name='xp_cmdshell'")
//	if err != nil {
//		//fmt.Println("something wrong in get xp_cmdshell")
//	} else {
//		info := GetSummary(Qresult, Columns)
//
//		if info == "1" {
//			res.XpCmdShell = "exsit"
//		} else {
//			res.XpCmdShell = "none"
//		}
//	}
//
//	err, Qresult, Columns = MssqlQuery(SqlCon, "select count(*) from master.dbo.sysobjects where xtype='x' and name='SP_OACREATE'")
//	if err != nil {
//		//fmt.Println("something wrong in get SP_OACREATE")
//	} else {
//		info := GetSummary(Qresult, Columns)
//
//		if info == "1" {
//			res.SP_OACREATE = "exsit"
//		} else {
//			res.SP_OACREATE = "none"
//		}
//	}
//	return res
//
//}
//
//func FindMssqlValuableTable(SqlCon *sql.DB) *[]MssqlValuable {
//	err, Qresult, Columns := MssqlQuery(SqlCon, "SELECT concat(TABLE_SCHEMA,'->',TABLE_NAME),COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS")
//
//	if err != nil {
//		return nil
//	}
//
//	vb := HandleMssqlValuable(Qresult, Columns)
//	return &vb
//}
//
//func HandleMssqlValuable(Qresult []map[string]string, Columns []string) []MssqlValuable {
//	var fin []MssqlValuable
//	for _, items := range Qresult {
//		if utils.SliceLike(utils.ValueableSlice, items["COLUMN_NAME"]) {
//			temp := MssqlValuable{
//				STName:     items[Columns[0]],
//				ColumnName: items[Columns[1]],
//			}
//			fin = append(fin, temp)
//		}
//	}
//	return fin
//
//}
