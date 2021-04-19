package Utils

import "fmt"

type MysqlInf struct {
	Version string
	OS      string
	Count   string
}

func OutPutQuery(Qresult []map[string]string, Columns []string) {
	for _, cname := range Columns {
		fmt.Print(cname + "\t")
	}
	fmt.Print("\n")
	for _, items := range Qresult {
		for _, cname := range Columns {
			fmt.Print(items[cname] + "\t")
		}
		fmt.Print("\n")
	}
}

func GetBaseInfo(Qresult []map[string]string, Columns []string) MysqlInf {

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

func GetSummary(Qresult []map[string]string, Columns []string) string {
	if len(Qresult) == 1 && len(Columns) == 1 {
		return Qresult[0][Columns[0]]
	}
	return ""
}
