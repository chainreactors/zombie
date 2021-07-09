package Database

import (
	"Zombie/src/Utils"
	"encoding/json"
	"fmt"
	"os"
)

func OutPutQuery(Qresult []map[string]string, Columns []string, title bool) {

	if title {
		for _, cname := range Columns {
			fmt.Print(cname + "\t")
		}
	}
	fmt.Print("\n")
	for _, items := range Qresult {
		for _, cname := range Columns {
			fmt.Print(items[cname] + "\t")
		}
		fmt.Print("\n")
	}
}

func GetSummary(Qresult []map[string]string, Columns []string) string {
	if len(Qresult) == 1 && len(Columns) == 1 {
		return Qresult[0][Columns[0]]
	}
	return ""
}

func QueryWrite3File(FileHandle *os.File, QDatach chan interface{}) {

	for res := range QDatach {
		flag := 0
		switch Utils.OutputType {
		case "MYSQL":
			flag += 1
			finres := res.(MysqlService)
			MysqlCollectInfo := ""
			MysqlCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\tOS: %v\nSummary: %v\n", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
			MysqlCollectInfo += fmt.Sprintf("general_log: %v\tgeneral_log_file: %v\n", finres.GeneralLog, finres.GeneralLogFile)
			MysqlCollectInfo += fmt.Sprintf("plugin_dir: %v\tsecure_file_priv: %v\n", finres.PluginPath, finres.SecureFilePriv)
			MysqlCollectInfo += "\n"
			fmt.Println(MysqlCollectInfo)
			switch Utils.FileFormat {
			case "raw":
				FileHandle.WriteString(MysqlCollectInfo)
			case "json":
				if flag == 1 {
					FileHandle.WriteString("{")
				}
				jsons, errs := json.Marshal(res)
				if errs != nil {
					fmt.Println(errs.Error())
					continue
				}
				FileHandle.WriteString(string(jsons) + ",")

			}

		case "POSTGRESQL":
			flag += 1
			finres := res.(PostgresService)
			PostCollectInfo := ""
			PostCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
			PostCollectInfo += "\n"
			fmt.Println(PostCollectInfo)
			switch Utils.FileFormat {
			case "raw":
				FileHandle.WriteString(PostCollectInfo)
			case "json":
				if flag == 1 {
					FileHandle.WriteString("{")
				}
				jsons, errs := json.Marshal(res)
				if errs != nil {
					fmt.Println(errs.Error())
					continue
				}
				FileHandle.WriteString(string(jsons) + ",")
			}
		case "MSSQL":
			flag += 1
			finres := res.(MssqlService)
			MsCollectInfo := ""
			MsCollectInfo += fmt.Sprintf("IP: %v\tServer: %v\nVersion: %v\nOS: %v\nSummary: %v", finres.Ip, Utils.OutputType, finres.Version, finres.OS, finres.Count)
			MsCollectInfo += fmt.Sprintf("\nSP_OACREATE: %v", finres.SP_OACREATE)
			MsCollectInfo += fmt.Sprintf("\nxp_cmdshell: %v", finres.XpCmdShell)
			MsCollectInfo += "\n"
			fmt.Println(MsCollectInfo)
			switch Utils.FileFormat {
			case "raw":
				FileHandle.WriteString(MsCollectInfo)
			case "json":
				if flag == 1 {
					FileHandle.WriteString("{")
				}
				jsons, errs := json.Marshal(res)
				if errs != nil {
					fmt.Println(errs.Error())
					continue
				}
				FileHandle.WriteString(string(jsons) + ",")
			}
		case "Brute":
			finres := res.(Utils.OutputRes)
			resstr := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t%s\n", finres.IP, finres.Port, finres.Username, finres.Password, finres.Type, finres.Additional)
			fmt.Println(resstr)
			switch Utils.FileFormat {
			case "raw":

				FileHandle.WriteString(resstr)
			case "json":
				flag += 1
				if flag == 1 {
					FileHandle.WriteString("{")
				}
				jsons, errs := json.Marshal(finres)
				if errs != nil {
					fmt.Println(errs.Error())
					continue
				}
				FileHandle.WriteString(string(jsons) + ",")
			}

		}
	}

}
