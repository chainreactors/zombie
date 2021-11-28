package Module

import (
	"Zombie/src/Core"
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Exec(ctx *cli.Context) (err error) {
	var CurServer string
	var CurtaskList []Utils.ScanTask
	Utils.Timeout = 2

	//ctx.String("InputFile")
	if ctx.IsSet("InputFile") {
		taskList, err := Core.CleanRes(ctx.String("InputFile"))
		CurtaskList = *taskList
		if err != nil {
			return err
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
			if _, ok := Utils.ServerPort[ServerName]; ok {
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

			if _, ok := Utils.PortServer[port]; ok {
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
			Input:    ctx.String("input"),
		}
		CurtaskList = append(CurtaskList, Curtask)

	}
	//初始化文件
	input := ctx.String("input")
	Utils.File = ctx.String("OutputFile")
	Utils.FileFormat = ctx.String("type")
	Utils.IsAuto = ctx.Bool("auto")
	Utils.Thread = ctx.Int("thread")

	//dir := "./res"
	//exist, _ := Utils.PathExists(dir)
	//
	//if !exist {
	//	err := os.Mkdir(dir, os.ModePerm)
	//	if err != nil {
	//		fmt.Printf("mkdir res dir failed![%v]\n", err)
	//	} else {
	//		fmt.Printf("mkdir res dir success!\n")
	//	}
	//}

	if Utils.File != "null" && Utils.IsAuto {
		Utils.FileHandle = Utils.InitFile(Utils.File)
		Utils.OutputType = CurtaskList[0].Server
		go ExecAble.QueryWrite3File(Utils.FileHandle, Utils.TDatach)

	}

	wgs := &sync.WaitGroup{}
	scanPool, _ := ants.NewPoolWithFunc(Utils.Thread, func(i interface{}) {
		par := i.(Utils.ScanTask)
		StartExec(par)
		wgs.Done()
	}, ants.WithExpiryDuration(2*time.Second))

	for _, Curtask := range CurtaskList {
		Curtask.Input = input
		//CurCon := Core.ExecDispatch(Curtask)
		//
		//alive := CurCon.Connect()
		//
		//if !alive {
		//	fmt.Printf("%v:%v can't connect to db\n", Curtask.Info.Ip, Curtask.Info.Port)
		//	continue
		//}
		//
		//if Utils.IsAuto {
		//	CurCon.GetInfo()
		//} else {
		//	CurCon.SetQuery(ctx.String("input"))
		//	CurCon.Query()
		//}
		wgs.Add(1)
		_ = scanPool.Invoke(Curtask)

	}

	wgs.Wait()

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

func StartExec(task Utils.ScanTask) {
	CurCon := Core.ExecDispatch(task)

	if CurCon == nil {
		return
	}

	alive := CurCon.Connect()

	if !alive {
		fmt.Printf("%v:%v can't connect to target\n", task.Info.Ip, task.Info.Port)
		return
	}

	if task.Input != "" {

		CurCon.SetQuery(task.Input)

	}

	if Utils.IsAuto {
		CurCon.GetInfo()
	} else {
		CurCon.Query()
	}
}
