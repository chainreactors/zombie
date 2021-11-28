package Core

import (
	"Zombie/src/Utils"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func CleanRes(Filename string) (*[]Utils.ScanTask, error) {
	var eachJson []Utils.OutputRes
	IPStore := make(map[string]int)
	TestList, _ := GetUAList(Filename)
	var CurtaskList []Utils.ScanTask
	if len(TestList) == 1 {
		if strings.HasPrefix(TestList[0], "[") {
			fmt.Println("start analysis json result")
			plain := TestList[0][1 : len(TestList[0])-1]
			plain = "[" + plain + "]"

			if err := json.Unmarshal([]byte(plain), &eachJson); err != nil {
				return nil, err
			}
			for _, info := range eachJson {
				Curtask := Utils.ScanTask{
					Username: info.Username,
					Password: info.Password,
					Server:   info.ScanTask.Input,
				}
				Curtask.Info.Ip = info.ScanTask.Info.Ip
				Curtask.Info.Port = info.ScanTask.Info.Port

				address := fmt.Sprintf("%v:%v", info.ScanTask.Info.Ip, info.ScanTask.Info.Port)

				if IPStore[address] == 1 {
					continue
				}

				CurtaskList = append(CurtaskList, Curtask)
				IPStore[address] = 1
			}
			CurtaskList = CurtaskList[:len(CurtaskList)-1]

		}
	} else {
		fmt.Println("start analysis raw result(clean)")
		for _, test := range TestList {
			la := strings.Split(test, "\t")

			if len(la) == 6 {
				Curtask := Utils.ScanTask{
					Username: strings.Split(la[2], ":")[1],
					Password: strings.Split(la[3], ":")[1],
					Server:   la[4],
				}
				if strings.Split(la[3], ":")[1] == "" && (la[4] == "SMB" || la[4] == "RDP") {
					address := fmt.Sprintf("%v:%v", Curtask.Info.Ip, Curtask.Info.Port)
					fmt.Printf("%s is %s, but it's passed by empty password, please check\n", address, la[4])
					IPStore[address] = 1
					continue
				}
				IpPo := strings.Split(la[0], ":")
				Curtask.Info.Ip = IpPo[0]
				Curtask.Info.Port, _ = strconv.Atoi(IpPo[1])

				address := fmt.Sprintf("%v:%v", Curtask.Info.Ip, Curtask.Info.Port)

				if _, ok := IPStore[address]; ok {
					continue
				}

				CurtaskList = append(CurtaskList, Curtask)
				IPStore[address] = 1
			}
			continue
		}
	}
	return &CurtaskList, nil
}

func OutPutRes(reslist *[]Utils.OutputRes, file string) error {

	dir, name := filepath.Split(file)

	clean_file := dir + ".clean_" + name
	hanlder := Utils.InitFile(clean_file)
	outputlist := *reslist
	for _, info := range outputlist {

		resstr := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t\n", info.ScanTask.Info.Ip, info.ScanTask.Info.Port, info.Username, info.Password, info.ScanTask.Server)
		//fmt.Println(resstr)
		switch Utils.FileFormat {
		case "raw":
			hanlder.WriteString(resstr)
		case "json":

			jsons, errs := json.Marshal(resstr)
			if errs != nil {
				fmt.Println(errs.Error())
				continue
			}
			hanlder.WriteString(string(jsons) + ",")
		default:
			hanlder.WriteString(resstr)
		}

	}
	return nil
}
