package ExecAble

import (
	"Zombie/src/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
)

type OracleService struct {
	Utils.IpInfo
	Username string `json:"username"`
	Password string `json:"password"`
	Input    string
	SqlCon   *sql.DB
}

func (s *OracleService) Query() bool {
	panic("implement me")
}

func OracleConnect(User string, Password string, info Utils.IpInfo) (err error, result bool, db *sql.DB) {

	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%d/%s?Connection TimeOut=%v&Connection Pool Timeout=%v", User, Password, info.Ip, info.Port, info.Instance, Utils.Timeout, Utils.Timeout)

	db, err = sql.Open("oracle", dataSourceName)

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

func (s *OracleService) Connect() bool {
	err, _, db := OracleConnect(s.Username, s.Password, s.IpInfo)
	if err == nil {
		s.SqlCon = db
		return true
	}
	return false
}

func (s *OracleService) DisConnect() bool {
	s.SqlCon.Close()
	return false
}

func (s *OracleService) SetQuery(query string) {
	s.Input = query
}

func (s *OracleService) Output(res interface{}) {
	finres := res.(OracleService)
	//MysqlCollectInfo := ""
	//MysqlCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\tOS: %v\nSummary: %v\n", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
	//MysqlCollectInfo += fmt.Sprintf("general_log: %v\tgeneral_log_file: %v\n", finres.GeneralLog, finres.GeneralLogFile)
	//MysqlCollectInfo += fmt.Sprintf("plugin_dir: %v\tsecure_file_priv: %v\n", finres.PluginPath, finres.SecureFilePriv)
	//MysqlCollectInfo += "\n"
	//fmt.Println(MysqlCollectInfo)
	switch Utils.FileFormat {
	case "raw":
		//Utils.TDatach <- MysqlCollectInfo
	case "json":
		jsons, errs := json.Marshal(finres)
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}
		Utils.TDatach <- jsons

	}

}

func (s *OracleService) GetInfo() bool {
	defer s.SqlCon.Close()

	return true
}
