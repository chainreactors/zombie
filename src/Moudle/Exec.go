package Moudle

import (
	"Zombie/src/Core"
	"Zombie/src/Utils"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func Exec(ctx *cli.Context) (err error) {
	var CurServer string
	var CurtaskList []Utils.ScanTask

	if ctx.IsSet("InputFile") {
		IPStore := make(map[string]int)
		var eachJson []Utils.OutputRes
		TestList, _ := Core.GetUAList(ctx.String("InputFile"))
		if len(TestList) == 1 {
			if strings.HasPrefix(TestList[0], "{") {
				fmt.Println("start analysis json result")
				plain := TestList[0][1 : len(TestList[0])-1]
				plain = "[" + plain + "]"

				if err := json.Unmarshal([]byte(plain), &eachJson); err != nil {
					return err

				}
				for _, info := range eachJson {
					Curtask := Utils.ScanTask{
						Username: info.Username,
						Password: info.Password,
						Server:   info.Type,
					}
					Curtask.Info.Ip = info.IP
					Curtask.Info.Port = info.Port

					address := fmt.Sprintf("%v:%v", info.IP, info.Port)

					if IPStore[address] == 1 {
						continue
					}

					CurtaskList = append(CurtaskList, Curtask)
					IPStore[address] = 1
				}

			}
		} else {
			fmt.Println("start analysis raw result")
			for _, test := range TestList {
				la := strings.Split(test, "\t")

				if len(la) == 6 {
					Curtask := Utils.ScanTask{
						Username: strings.Split(la[2], ":")[1],
						Password: strings.Split(la[3], ":")[1],
						Server:   la[4],
					}
					IpPo := strings.Split(la[0], ":")
					Curtask.Info.Ip = IpPo[0]
					Curtask.Info.Port, _ = strconv.Atoi(IpPo[1])

					address := fmt.Sprintf("%v:%v", Curtask.Info.Ip, Curtask.Info.Port)

					if IPStore[address] == 1 {
						continue
					}

					CurtaskList = append(CurtaskList, Curtask)
					IPStore[address] = 1
				}
				continue
			}
		}

	} else {
		if strings.Contains(ctx.String("ip"), ",") {
			fmt.Println("Exec Moudle only support single ip")
			os.Exit(0)
		}

		IpSlice := Core.GetIpList(ctx.String("ip"))

		Ip := IpSlice[0]
		if ctx.IsSet("server") {
			ServerName := strings.ToUpper(ctx.String("server"))
			if _, ok := Utils.ExecPort[ServerName]; ok {
				CurServer = ctx.String("server")
			} else {
				fmt.Println("the Database isn't be supported")
				os.Exit(0)
			}

		} else if strings.Contains(Ip, ":") {
			Temp := strings.Split(Ip, ":")
			Sport := Temp[1]
			port, err := strconv.Atoi(Sport)
			if err != nil {
				fmt.Println("Please check your address")
				os.Exit(0)
			}

			if _, ok := Utils.ExecServer[port]; ok {
				CurServer = Utils.PortServer[port]
				fmt.Println("Use default server")
			} else {
				fmt.Println("Please input the type of Database")
				os.Exit(0)
			}
		} else {
			fmt.Println("Please input the type of Database")
			os.Exit(0)
		}

		CurServer = strings.ToUpper(CurServer)

		IpList := Core.GetIpInfoList(IpSlice, CurServer)

		Curtask := Utils.ScanTask{
			Info:     IpList[0],
			Username: ctx.String("username"),
			Password: ctx.String("password"),
			Server:   CurServer,
		}
		CurtaskList = append(CurtaskList, Curtask)

	}
	//初始化文件
	Utils.File = ctx.String("OutputFile")

	if Utils.File != "null" {
		ExecInitFile(Utils.File)
	}

	for _, Curtask := range CurtaskList[:len(CurtaskList)-1] {

		CurCon := Core.ExecDispatch(Curtask)

		alive := CurCon.Connect()

		if !alive {
			fmt.Printf("can't connect to db\n")
			continue
		}

		Utils.IsAuto = ctx.Bool("auto")

		if Utils.IsAuto {
			CurCon.GetInfo()
		} else {
			CurCon.SetQuery(ctx.String("input"))
			CurCon.Query()
		}
	}

	time.Sleep(1000 * time.Millisecond)
	fmt.Println("All Task Done!!!!")

	return err
}

func ExecInitFile(Filename string) {
	var err error

	if Filename != "" {
		Utils.O2File = true
		if Utils.CheckFileIsExist(Filename) { //如果文件存在
			Utils.FileHandle, err = os.OpenFile(Filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
			//fmt.Println("文件存在")
			if err != nil {
				os.Exit(0)
			}
			//io.WriteString(FileHandle, "123")
		} else {
			Utils.FileHandle, err = os.Create(Filename) //创建文件
			//fmt.Println("文件不存在")
			if err != nil {
				os.Exit(0)
			}
		}

	}
	go Utils.QueryWrite2File(Utils.FileHandle, Utils.QDatach)
}
