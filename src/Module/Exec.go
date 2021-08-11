package Module

import (
	"Zombie/src/Core"
	"Zombie/src/ExecAble"
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
				CurtaskList = CurtaskList[:len(CurtaskList)-1]

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
			fmt.Println("Exec Module only support single ip")
			os.Exit(0)
		}

		IpSlice := Core.GetIpList(ctx.String("ip"))

		Ip := IpSlice[0]
		if ctx.IsSet("server") {
			ServerName := strings.ToUpper(ctx.String("server"))
			if _, ok := Utils.ExecPort[ServerName]; ok {
				CurServer = ctx.String("server")
			} else {
				fmt.Println("the ExecAble isn't be supported")
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
				fmt.Println("Please input the type of ExecAble")
				os.Exit(0)
			}
		} else {
			fmt.Println("Please input the type of ExecAble")
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
	Utils.FileFormat = ctx.String("type")
	Utils.IsAuto = ctx.Bool("auto")

	dir := "./res"
	exist, _ := Utils.PathExists(dir)

	if !exist {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir res dir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir res dir success!\n")
		}
	}

	if Utils.File != "null" && Utils.IsAuto {
		initFile(Utils.File)
		Utils.OutputType = CurtaskList[0].Server
		go ExecAble.QueryWrite3File(Utils.FileHandle, Utils.TDatach)

	}

	for _, Curtask := range CurtaskList {

		CurCon := Core.ExecDispatch(Curtask)

		alive := CurCon.Connect()

		if !alive {
			fmt.Printf("%v:%v can't connect to db\n", Curtask.Info.Ip, Curtask.Info.Port)
			continue
		}

		if Utils.IsAuto {
			CurCon.GetInfo()
		} else {
			CurCon.SetQuery(ctx.String("input"))
			CurCon.Query()
		}
	}

	time.Sleep(1000 * time.Millisecond)
	if Utils.FileFormat == "json" {
		final := Utils.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		Utils.FileHandle.WriteString(string(jsons) + "}")
	}
	fmt.Println("All Task Done!!!!")

	return err
}
