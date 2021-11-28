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

func CleanBruteRes(BruList *[]Utils.OutputRes) (CleanedList []Utils.Codebook, CleanedRes []Utils.OutputRes) {
	IPStore := make(map[string]int)
	for _, res := range *BruList {
		address := fmt.Sprintf("%v:%v:%v", res.Ip, res.Port, res.Username)
		if _, ok := IPStore[address]; ok {
			continue
		}

		if res.Password == "" && (res.Server == "SMB" || res.Server == "RDP") {
			IPStore[address] = 1
			continue
		}
		cb := Utils.Codebook{
			Username: res.Username,
			Password: res.Password,
			Server:   res.Server,
		}
		CleanedList = append(CleanedList, cb)
		CleanedRes = append(CleanedRes, res)
		IPStore[address] = 1
	}
	CleanedList = RemoveCodeBookDu(CleanedList)
	return CleanedList, CleanedRes
}

func RemoveCodeBookDu(CBList []Utils.Codebook) []Utils.Codebook {
	result := make([]Utils.Codebook, 0, len(CBList))
	temp := map[string]struct{}{}
	for _, item := range CBList {
		if _, ok := temp[item.Username+item.Password]; !ok {
			temp[item.Username+item.Password] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func QueryWrite3File(FileHandle *os.File, QDatach chan interface{}) {

	for res := range QDatach {
		switch Utils.OutputType {
		case "Brute":
			finres := res.(Utils.OutputRes)
			*Utils.BrutedList = append(*Utils.BrutedList, finres)
			var resstr string
			if finres.Server == "ORACLE" {
				resstr = fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\tinstance:%s\t%s\tsuccess\n", finres.Ip, finres.Port, finres.Username, finres.Password, finres.Additional, finres.Server)
			} else {
				resstr = fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t%s\n", finres.Ip, finres.Port, finres.Username, finres.Password, finres.Server, finres.Additional)
			}
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
