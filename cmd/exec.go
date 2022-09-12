package cmd

import (
	"encoding/json"
	"fmt"
	core2 "github.com/chainreactors/zombie/internal/core"
	"github.com/chainreactors/zombie/internal/plugin"
	utils2 "github.com/chainreactors/zombie/pkg/utils"
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
	var CurtaskList []utils2.ScanTask
	utils2.Timeout = 2

	//ctx.String("InputFile")
	if ctx.IsSet("InputFile") {
		taskList, err := core2.CleanRes(ctx.String("InputFile"))
		CurtaskList = *taskList
		if err != nil {
			return err
		}
	} else {
		if strings.Contains(ctx.String("ip"), ",") {
			fmt.Println("Exec Module only support single ip")
			os.Exit(0)
		}

		IpSlice := core2.GetIpList(ctx.String("ip"))

		Ip := IpSlice[0]
		if ctx.IsSet("server") {
			ServerName := strings.ToUpper(ctx.String("server"))
			if _, ok := utils2.ServerPort[ServerName]; ok {
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

			if _, ok := utils2.PortServer[port]; ok {
				CurServer = utils2.PortServer[port]
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

		IpList := core2.GetIpInfoList(IpSlice, CurServer)

		Curtask := utils2.ScanTask{
			TargetInfo: utils2.TargetInfo{
				IpServerInfo: utils2.IpServerInfo{
					IpInfo: IpList[0],
					Server: CurServer,
				},
				Username: ctx.String("username"),
				Password: ctx.String("password"),
			},
			Input: ctx.String("input"),
		}
		CurtaskList = append(CurtaskList, Curtask)

	}
	//初始化文件
	input := ctx.String("input")
	utils2.File = ctx.String("OutputFile")
	utils2.FileFormat = ctx.String("type")
	utils2.IsAuto = ctx.Bool("auto")
	utils2.Thread = ctx.Int("thread")

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

	if utils2.File != "null" {
		utils2.FileHandle = utils2.InitFile(utils2.File)
		utils2.OutputType = CurtaskList[0].Server
		go plugin.QueryWrite3File(utils2.FileHandle, utils2.TDatach)

	}

	wgs := &sync.WaitGroup{}
	scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
		par := i.(utils2.ScanTask)
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
	if utils2.FileFormat == "json" {
		final := utils2.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		utils2.FileHandle.WriteString(string(jsons) + "}")
	}
	fmt.Println("All Task Done!!!!")

	return err
}

func StartExec(task utils2.ScanTask) {
	CurCon := core2.ExecDispatch(task)

	if CurCon == nil {
		return
	}

	alive := CurCon.Connect()

	if !alive {
		fmt.Printf("%v:%v can't connect to target\n", task.Ip, task.Port)
		return
	}

	if task.Input != "" {

		CurCon.SetQuery(task.Input)

	}

	if utils2.IsAuto {
		CurCon.GetInfo()
	} else {
		CurCon.Query()
	}
}
