package oracle

import (
	"database/sql"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	_ "github.com/sijms/go-ora/v2"
)

type OraclePlugin struct {
	*pkg.Task
	Input string
	conn  *sql.DB
}

func (s *OraclePlugin) Unauth() (bool, error) {
	return false, nil
}

//func (s *OraclePlugin) Query() bool {
//	return true
//}

func (s *OraclePlugin) Login() error {
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%s/%s?Connection TimeOut=%v&Connection Pool Timeout=%v", s.Username, s.Password, s.IP, s.Port, s.Param["instance"], s.Timeout, s.Timeout)

	conn, err := sql.Open("oracle", dataSourceName)
	if err != nil {
		return err
	}

	//conn.SetMaxOpenConns(60)
	//conn.SetMaxIdleConns(60)

	err = conn.Ping()
	if err != nil {
		return err
	}

	s.conn = conn
	return err
}

func (s *OraclePlugin) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return pkg.NilConnError{s.Service}
}

func (s *OraclePlugin) Name() string {
	return s.Service.String()
}

func (s *OraclePlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

//func (s *OraclePlugin) SetQuery(query string) {
//	s.Input = query
//}

//func (s *OraclePlugin) Output(res interface{}) {
//	//finres := res.(OracleService)
//	////MysqlCollectInfo := ""
//	////MysqlCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\tOS: %v\nSummary: %v\n", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
//	////MysqlCollectInfo += fmt.Sprintf("general_log: %v\tgeneral_log_file: %v\n", finres.GeneralLog, finres.GeneralLogFile)
//	////MysqlCollectInfo += fmt.Sprintf("plugin_dir: %v\tsecure_file_priv: %v\n", finres.PluginPath, finres.SecureFilePriv)
//	////MysqlCollectInfo += "\n"
//	////fmt.Println(MysqlCollectInfo)
//	//switch utils.FileFormat {
//	//case "raw":
//	//	//Utils.TDatach <- MysqlCollectInfo
//	//case "json":
//	//	jsons, errs := json.Marshal(finres)
//	//	if errs != nil {
//	//		fmt.Println(errs.Error())
//	//		return
//	//	}
//	//	utils.TDatach <- jsons
//	//
//	//}
//
//}

//func (s *OraclePlugin) GetInfo() bool {
//	defer s.conn.Close()
//
//	return true
//}
