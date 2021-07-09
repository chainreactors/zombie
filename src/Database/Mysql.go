package Database

import (
	"Zombie/src/Utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type MysqlService struct {
	Utils.IpInfo
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
}

func MysqlQuery(SqlCon *sql.DB, Query string) (err error, Qresult []map[string]string, Columns []string) {

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

func MysqlConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=utf8", User,
		Password, info.Ip, info.Port, Utils.Timeout, Utils.Timeout, Utils.Timeout)
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

func MysqlConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, db := MysqlConnect(User, Password, info)

	if err == nil {
		result.Result = res
		_ = db.Close()
	}

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

func (s *MysqlService) GetInfo() bool {
	defer s.SqlCon.Close()

	res := GetMysqlBaseInfo(s.SqlCon)

	if res == nil {
		return false
	}

	res.Count = GetMysqlSummary(s.SqlCon)

	res = GetMysqlVulnableInfo(s.SqlCon, res)
	res = GetMysqlGeneralLog(s.SqlCon, res)
	s.MysqlInf = *res

	//将结果放入管道
	Utils.TDatach <- *s
	return true
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

func GetMysqlSummary(SqlCon *sql.DB) string {
	err, Qresult, Columns := MysqlQuery(SqlCon, "select sum(table_rows) from  information_schema.tables where table_rows is not null")

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
