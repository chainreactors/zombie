package exec

import (
	utils2 "Zombie/v1/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type MysqlService struct {
	utils2.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	MysqlInf
	Input  string
	SqlCon *sql.DB
}

type MysqlInf struct {
	Version        string `json:"version"`
	OS             string `json:"os"`
	Count          string `json:"count"`
	GeneralLog     string `json:"general_log"`
	GeneralLogFile string `json:"general_log_file"`
	SecureFilePriv string `json:"secure_file_priv"`
	PluginPath     string `json:"plugin_path"`
	vb             []MysqlValuable
}

type MysqlValuable struct {
	STName     string
	ColumnName string
}

func MysqlQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {

	err = SqlCon.Ping()
	if err == nil {
		rows, err := SqlCon.Query(Query)
		if err == nil {
			Qresult, Columns = DoRowsMapper(rows)

		} else {
			if !utils2.IsAuto {
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

func MysqlConnect(User string, Password string, info utils2.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8", User,
		Password, info.Ip, info.Port, utils2.Timeout, utils2.Timeout, utils2.Timeout)
	db, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		result = false
		return err, result, nil
	}

	db.SetMaxOpenConns(60)
	db.SetMaxIdleConns(60)

	err = db.Ping()

	if err == nil {
		result = true
	}
	return err, result, db
}

func (s *MysqlService) Connect() bool {
	err, _, db := MysqlConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SqlCon = db
		return true
	}
	return false
}

func (s *MysqlService) DisConnect() bool {
	s.SqlCon.Close()
	return false
}

func (s *MysqlService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetMysqlBaseInfo(s.SqlCon)

	if res == nil {
		return false
	}

	res.Count = GetMysqlSummary(s.SqlCon)

	res = GetMysqlVulnableInfo(s.SqlCon, res)
	res = GetMysqlGeneralLog(s.SqlCon, res)
	res.vb = *FindValuableTable(s.SqlCon)
	s.MysqlInf = *res

	//将结果放入管道
	s.Output(*s)
	return true
}

func (s *MysqlService) Output(res interface{}) {
	finres := res.(MysqlService)
	MysqlCollectInfo := ""
	MysqlCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\tOS: %v\nSummary: %v\n", finres.Ip, utils2.OutputType, finres.Version, finres.OS, finres.Count)
	MysqlCollectInfo += fmt.Sprintf("general_log: %v\tgeneral_log_file: %v\n", finres.GeneralLog, finres.GeneralLogFile)
	MysqlCollectInfo += fmt.Sprintf("plugin_dir: %v\tsecure_file_priv: %v\n", finres.PluginPath, finres.SecureFilePriv)
	for _, info := range finres.vb {
		MysqlCollectInfo += fmt.Sprintf("%v:%v\t", info.STName, info.ColumnName)
	}
	MysqlCollectInfo += "\n"
	fmt.Println(MysqlCollectInfo)
	switch utils2.FileFormat {
	case "raw":
		utils2.TDatach <- MysqlCollectInfo
	case "json":
		jsons, errs := json.Marshal(finres)
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}
		utils2.TDatach <- jsons

	}

}

func (s *MysqlService) Query() bool {
	err, Qresult, Columns := MysqlQuery(s.SqlCon, s.Input)

	if err != nil {
		fmt.Println("something wrong")
		return false
	} else {
		OutPutQuery(Qresult, Columns, true)
	}
	return true
}

func (s *MysqlService) SetQuery(query string) {
	s.Input = query
}

func GetMysqlBaseInfo(SqlCon *sql.DB) *MysqlInf {
	err, Qresult, Columns := MysqlQuery(SqlCon, "select VERSION(),@@version_compile_os")
	if err != nil {
		fmt.Println("something wrong at get version")
		return nil
	}
	res := HandleBaseInfo(Qresult, Columns)

	return &res
}

func FindValuableTable(SqlCon *sql.DB) *[]MysqlValuable {
	err, Qresult, _ := MysqlQuery(SqlCon, "select concat(table_schema,\"->\",table_name),column_name from information_schema.columns where "+
		"table_schema != \"information_schema\" &&  table_schema != \"mysql\" && table_schema != \"performance_schema\"")

	if err != nil {
		return nil
	}

	vb := HandleMysqlValuable(Qresult)
	return &vb
}

func GetMysqlSummary(SqlCon *sql.DB) string {
	err, Qresult, Columns := MysqlQuery(SqlCon, "select sum(table_rows) from  information_schema.tables where table_rows is not null")

	if err == nil {
		Count := GetSummary(Qresult, Columns)
		return Count
	}

	//某些场景下使用is not null 会出现bug，所以增加一条，失败后是
	err, Qresult, Columns = MysqlQuery(SqlCon, "select sum(table_rows) from  information_schema.tables")

	if err == nil {
		Count := GetSummary(Qresult, Columns)
		return Count
	}
	return ""
}

func GetMysqlGeneralLog(SqlCon *sql.DB, res *MysqlInf) *MysqlInf {
	err, Qresult, Columns := MysqlQuery(SqlCon, "show VARIABLES like 'general%'")
	if err != nil {
		//fmt.Println("something wrong in get general log")
	} else {
		flag := 0
		//Utils.OutPutQuery(Qresult, Columns, false)
		for _, items := range Qresult {
			for _, cname := range Columns {
				if flag == 1 {
					res.GeneralLog = items[cname]
					flag = 0
				} else {
					res.GeneralLogFile = items[cname]
					flag = 0
				}
				if strings.Contains(items[cname], "general_log_file") {
					flag = 2
				} else if strings.Contains(items[cname], "general_log") {
					flag = 1
				}
			}
		}
	}
	return res

}

func GetMysqlVulnableInfo(SqlCon *sql.DB, res *MysqlInf) *MysqlInf {
	err, Qresult, Columns := MysqlQuery(SqlCon, "SHOW VARIABLES LIKE \"secure_file_priv\"")
	if err != nil {
		//获取失败
		//fmt.Println("\nsomething wrong in get secure_file_priv")
	} else {
		if len(Qresult) == 1 && len(Columns) == 2 {
			//MysqlCollectInfo += fmt.Sprint("\n" + Qresult[0][Columns[0]] + ":\t" + Qresult[0][Columns[1]])
			res.SecureFilePriv = Qresult[0][Columns[1]]
		}
	}

	err, Qresult, Columns = MysqlQuery(SqlCon, "show variables like '%plugin%'")
	if err != nil {
		//获取失败
		//fmt.Println("\nsomething wrong in get plugin dir")
	} else {
		if len(Qresult) == 1 && len(Columns) == 2 {
			//MysqlCollectInfo += fmt.Sprint("\n" + Qresult[0][Columns[0]] + ":\t" + Qresult[0][Columns[1]])
			res.PluginPath = Qresult[0][Columns[1]]
		}
	}
	return res

}

func HandleMysqlValuable(Qresult []map[string]string) []MysqlValuable {
	var fin []MysqlValuable
	for _, items := range Qresult {
		if utils2.SliceLike(utils2.ValueableSlice, items["column_name"]) {
			temp := MysqlValuable{
				STName:     items["concat(table_schema,\"->\",table_name)"],
				ColumnName: items["column_name"],
			}
			fin = append(fin, temp)
		}
	}
	return fin

}

func HandleBaseInfo(Qresult []map[string]string, Columns []string) MysqlInf {

	res := MysqlInf{}

	for _, items := range Qresult {
		for _, cname := range Columns {
			if cname == "VERSION()" {
				res.Version = items[cname]
			} else {
				res.OS = items[cname]
			}
		}

	}
	return res
}
