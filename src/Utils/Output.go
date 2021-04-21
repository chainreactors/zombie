package Utils

import "fmt"

type MysqlInf struct {
	Version string
	OS      string
	Count   string
}

type MssqlInf struct {
	Version string
	Count   int
	OS      string
}

type PostgreInf struct {
	Version string
	Count   int
	OS      string
}

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
