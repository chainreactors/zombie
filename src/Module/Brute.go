package Module

import (
	"Zombie/src/Core"
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Brute(ctx *cli.Context) (err error) {
	var CurServer string
	var UserList, PassList, IpSlice []string
	var IpList []Utils.IpInfo
	var fromgt bool
	var ipserverinfo []Utils.IpServerInfo

	fromgt = ctx.IsSet("gt")

	if ctx.IsSet("ip") && (!ctx.IsSet("IP") || !ctx.IsSet("gt")) {
		IpSlice = Core.GetIPList(ctx.String("ip"))
	} else if ctx.IsSet("IP") && !ctx.IsSet("gt") {
		IpSlice, _ = Core.ReadIPDict(ctx.String("IP"))
	} else if ctx.IsSet("gt") {
		fmt.Println("Read from gt result")
	} else {
		fmt.Println("please check the ip")
		os.Exit(0)
	}

	if !fromgt {

		Ip := IpSlice[0]
		if ctx.IsSet("server") {
			ServerName := strings.ToUpper(ctx.String("server"))
			if _, ok := Utils.ServerPort[ServerName]; ok {
				CurServer = ServerName
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

		IpList = Core.GetIpInfoList(IpSlice, CurServer)
		ipserverinfo = GenerIPServerInfo(IpList, CurServer)
	} else {
		ipserverinfo = GenFromGT(ctx.String("gt"))
	}

	if ctx.IsSet("uppair") {
		uppair := ctx.String("uppair")
		Core.UPList, _ = Core.GetUAList(uppair)
	} else {

		if ctx.IsSet("username") && !ctx.IsSet("userdict") {
			username := ctx.String("username")
			UserList = Core.GetUserList(username)
		} else if ctx.IsSet("userdict") {
			UserList, _ = Core.ReadUserDict(ctx.String("userdict"))
		} else {
			fmt.Println("[+] Use default user dict")
		}

		if ctx.IsSet("password") && !ctx.IsSet("passdict") {
			password := ctx.String("password")
			PassList = Core.GetPassList(password)
		} else if ctx.IsSet("passdict") {
			PassList, _ = Core.ReadPassDict(ctx.String("passdict"))
		} else {
			fmt.Println("[+] Use default password dict")
		}
	}

	Utils.Timeout = ctx.Int("timeout")
	Utils.SSL = ctx.Bool("ssl")
	Utils.Thread = ctx.Int("thread")
	Utils.Simple = ctx.Bool("simple")
	Utils.Proc = ctx.Int("proc")
	Utils.FileFormat = ctx.String("type")
	Utils.File = ctx.String("file")
	Utils.OutputType = "Brute"

	if Utils.File != "null" {
		Utils.FileHandle = Utils.InitFile(Utils.File)
		go ExecAble.QueryWrite3File(Utils.FileHandle, Utils.TDatach)
	}

	Core.Summary = len(UserList) * len(PassList) * len(IpList)

	if Utils.Proc != 0 {
		go Core.Process(Core.CountChan)
	}

	if Utils.Simple {
		err = StartTaskSimple(UserList, PassList, ipserverinfo)
	} else {
		err = StartTask(UserList, PassList, ipserverinfo)
	}
	Utils.FileHandle.Close()

	close(Core.CountChan)

	reslist, err := Core.CleanRes(Utils.File)
	if err != nil {
		return err
	}
	Core.OutPutRes(reslist, Utils.File)

	return err
}

func StartTask(UserList []string, PassList []string, IpServerList []Utils.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())
	for _, ipinfo := range IpServerList {

		fmt.Printf("Now Processing %s:%d, ExecAble: %s\n", ipinfo.Ip, ipinfo.Port, ipinfo.Server)

		Utils.ChildContext, Utils.ChildCancel = context.WithCancel(rootContext)

		TaskList := Core.GenerateTask(UserList, PassList, ipinfo)

		wgs := &sync.WaitGroup{}
		PrePara := Core.PoolPara{
			Ctx:      Utils.ChildContext,
			Taskchan: TaskList,
			Wgs:      wgs,
		}

		scanPool, _ := ants.NewPoolWithFunc(Utils.Thread, func(i interface{}) {
			par := i.(Core.PoolPara)
			Core.BruteWork(&par)
		}, ants.WithExpiryDuration(2*time.Second))

		for i := 0; i < Utils.Thread; i++ {
			wgs.Add(1)
			_ = scanPool.Invoke(PrePara)
		}
		wgs.Wait()

		RandomTask := Utils.ScanTask{
			Info:     ipinfo.IpInfo,
			Username: Core.FlagUserName,
			Password: Utils.RandStringBytesMaskImprSrcUnsafe(12),
			Server:   ipinfo.Server,
		}

		CurCon := Core.ExecDispatch(RandomTask)

		alive := CurCon.Connect()

		if alive {
			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Info.Ip, RandomTask.Info.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Info.Ip, RandomTask.Info.Port)
		}
	}

	fmt.Println("All Task done")

	time.Sleep(1000 * time.Millisecond)
	if Utils.FileFormat == "json" {
		final := Utils.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		Utils.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func StartTaskSimple(UserList []string, PassList []string, IpServerList []Utils.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())

	TaskList := Core.GenerateTaskSimple(UserList, PassList, IpServerList)

	wgs := &sync.WaitGroup{}
	PrePara := Core.PoolPara{
		Ctx:      rootContext,
		Taskchan: TaskList,
		Wgs:      wgs,
	}

	scanPool, _ := ants.NewPoolWithFunc(Utils.Thread, func(i interface{}) {
		par := i.(Core.PoolPara)
		Core.BruteWork(&par)
	}, ants.WithExpiryDuration(2*time.Second))

	for i := 0; i < Utils.Thread; i++ {
		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}
	wgs.Wait()

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("All Task done")

	if Utils.FileFormat == "json" {
		final := Utils.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		Utils.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func GenerIPServerInfo(ipinfo []Utils.IpInfo, server string) (ipserverinfo []Utils.IpServerInfo) {
	for _, info := range ipinfo {
		isinfo := Utils.IpServerInfo{}
		isinfo.IpInfo = info
		isinfo.Server = server
		ipserverinfo = append(ipserverinfo, isinfo)
	}

	return ipserverinfo
}

func GenFromGT(gtfile string) (ipserverinfo []Utils.IpServerInfo) {

	bytes, err := ioutil.ReadFile(gtfile)
	if err != nil {
		println(gtfile + " open failed")
		//panic(dictPath + " open failed")
		os.Exit(0)
	}

	if err := json.Unmarshal(bytes, &ipserverinfo); err != nil {
		println(" Unmarshal failed")
		os.Exit(0)
	}

	return ipserverinfo
}
