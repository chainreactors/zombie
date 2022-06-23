package Core

import (
	"Zombie/src/Utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
					TargetInfo: Utils.TargetInfo{
						IpServerInfo: Utils.IpServerInfo{
							Server: info.Server,
						},
						Username: info.Username,
						Password: info.Password,
					},
				}
				Curtask.Ip = info.Ip
				Curtask.Port = info.Port

				address := fmt.Sprintf("%v:%v", info.Ip, info.Port)

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
					TargetInfo: Utils.TargetInfo{
						IpServerInfo: Utils.IpServerInfo{
							Server: la[4],
						},
						Username: strings.Split(la[2], ":")[1],
						Password: strings.Split(la[3], ":")[1],
					},
				}
				if strings.Split(la[3], ":")[1] == "" && (la[4] == "SMB" || la[4] == "RDP") {
					address := fmt.Sprintf("%v:%v", Curtask.Ip, Curtask.Port)
					fmt.Printf("%s is %s, but it's passed by empty password, please check\n", address, la[4])
					IPStore[address] = 1
					continue
				}
				IpPo := strings.Split(la[0], ":")
				Curtask.Ip = IpPo[0]
				Curtask.Port, _ = strconv.Atoi(IpPo[1])

				address := fmt.Sprintf("%v:%v", Curtask.Ip, Curtask.Port)

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

func OutPutRes(reslist *[]Utils.OutputRes, cblist *[]Utils.Codebook, file string) error {

	dir, name := filepath.Split(file)

	clean_file := dir + ".clean_" + name
	codebook := dir + ".cb.log"
	hanlder := Utils.InitFile(clean_file)
	outputlist := *reslist
	for _, info := range outputlist {

		resstr := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t\n", info.Ip, info.Port, info.Username, info.Password, info.Server)
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

	var oldcb []Utils.Codebook
	if Utils.CheckFileIsExist(codebook) {

		cbytes, err := ioutil.ReadFile(codebook)
		if err != nil {
			return err
		}
		if len(cbytes) != 0 {
			err = json.Unmarshal(cbytes, &oldcb)
			if err != nil {
				return err
			}
		}
		oldcb = append(oldcb, *cblist...)
		oldcb = RemoveDuplicateCodeBook(oldcb)
		newcb, err := json.Marshal(oldcb)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(codebook, newcb, 0666)
		if err != nil {
			return err
		}
	} else {
		newcb, err := json.Marshal(*cblist)
		if err != nil {
			return err
		}

		cbhandle := Utils.InitFile(codebook)
		cbhandle.Write(newcb)
	}

	return nil
}

func RemoveDuplicateCodeBook(cb []Utils.Codebook) []Utils.Codebook {
	result := make([]Utils.Codebook, 0, len(cb))
	temp := map[string]struct{}{}
	for _, item := range cb {
		if _, ok := temp[item.Username+item.Password+item.Server]; !ok {
			temp[item.Username+item.Password+item.Server] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
