package ExecAble

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
		switch Utils.OutputType {
		case "Brute":
			finres := res.(Utils.OutputRes)
			resstr := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t%s\n", finres.IP, finres.Port, finres.Username, finres.Password, finres.Type, finres.Additional)
			fmt.Println(resstr)
			switch Utils.FileFormat {
			case "raw":
				FileHandle.WriteString(resstr)
			case "json":

				jsons, errs := json.Marshal(finres)
				if errs != nil {
					fmt.Println(errs.Error())
					continue
				}
				FileHandle.WriteString(string(jsons) + ",")
			}
		default:
			switch Utils.FileFormat {
			case "raw":
				FileHandle.WriteString(res.(string))
			case "json":
				FileHandle.WriteString(res.(string) + ",")
			}

		}
	}

}
