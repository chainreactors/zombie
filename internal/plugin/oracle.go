package plugin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	_ "github.com/sijms/go-ora/v2"
)

type OracleService struct {
	*utils.Task
	Input string
	conn  *sql.DB
}

func (s *OracleService) Query() bool {
	panic("implement me")
}

func OracleConnect(info *utils.Task) (conn *sql.DB, err error) {
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%d/%s?Connection TimeOut=%v&Connection Pool Timeout=%v", info.Username, info.Password, info.IP.String(), info.Port, info.Instance, info.Timeout, info.Timeout)

	conn, err = sql.Open("oracle", dataSourceName)
	if err != nil {
		return nil, err
	}

	//conn.SetMaxOpenConns(60)
	//conn.SetMaxIdleConns(60)

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *OracleService) Connect() error {
	conn, err := OracleConnect(s.Task)
	if err != nil {
		return err
	}
	s.conn = conn
	return err
}

func (s *OracleService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return NilConnError{s.Service}
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
	switch utils.FileFormat {
	case "raw":
		//Utils.TDatach <- MysqlCollectInfo
	case "json":
		jsons, errs := json.Marshal(finres)
		if errs != nil {
			fmt.Println(errs.Error())
			return
		}
		utils.TDatach <- jsons

	}

}

func (s *OracleService) GetInfo() bool {
	defer s.conn.Close()

	return true
}
