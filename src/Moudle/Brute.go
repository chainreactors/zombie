package Moudle

import (
	"Zombie/src/Core"
	"Zombie/src/Utils"
	"context"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
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

	if ctx.IsSet("ip") || !ctx.IsSet("IP") {
		IpSlice = Core.GetUserList(ctx.String("ip"))
	} else if ctx.IsSet("IP") {
		IpSlice, _ = Core.ReadUserDict(ctx.String("IP"))
	} else {
		fmt.Println("please check the ip")
		os.Exit(0)
	}

	Ip := IpSlice[0]
	if ctx.IsSet("server") {
		ServerName := strings.ToUpper(ctx.String("server"))
		if _, ok := Utils.ServerPort[ServerName]; ok {
			CurServer = ServerName
		} else {
			fmt.Println("the Server isn't be supported")
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
			fmt.Println("Please input the type of Server")
			os.Exit(0)
		}
	} else {
		fmt.Println("Please input the type of Server")
		os.Exit(0)
	}

	IpList = Core.GetIpInfoList(IpSlice, CurServer)

	//if !ctx.IsSet("username") && !ctx.IsSet("password") {
	//	fmt.Println("please input username and password, if the server don't need username,input anything to place")
	//	os.Exit(0)
	//} else {
	//	username := ctx.String("username")
	//	password := ctx.String("password")
	//	UserList = Core.GetUserList(username)
	//	PassList = Core.GetPassList(password)
	//}

	if ctx.IsSet("username") || !ctx.IsSet("userdict") {
		username := ctx.String("username")
		UserList = Core.GetUserList(username)
	} else if ctx.IsSet("userdict") {
		UserList, _ = Core.ReadUserDict(ctx.String("userdict"))
	} else {
		fmt.Println("please input username")
		os.Exit(0)
	}

	if ctx.IsSet("password") || !ctx.IsSet("passdict") {
		password := ctx.String("password")
		PassList = Core.GetPassList(password)
	} else if ctx.IsSet("passdict") {
		PassList, _ = Core.ReadPassDict(ctx.String("passdict"))
	} else {
		fmt.Println("please input user")
		os.Exit(0)
	}

	CurServer = strings.ToUpper(CurServer)

	Utils.Timeout = ctx.Int("timeout")
	Utils.Thread = ctx.Int("thread")

	//ExpireTime := GetExpireTime(len(IpList), len(UserList), len(PassList))

	if ctx.IsSet("file") {
		initFile(ctx.String("file"))
	}

	err = StartTask(UserList, PassList, IpList, CurServer)

	return err
}

func initFile(Filename string) {
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
			//io.WriteString(FileHandle, "123")
		}

	}
	go Utils.Write2File(Utils.FileHandle, Utils.Datach)
}

func StartTask(UserList []string, PassList []string, IpList []Utils.IpInfo, CurServer string) error {
	rootContext, rootCancel := context.WithCancel(context.Background())
	for _, ipinfo := range IpList {

		fmt.Printf("Now Processing %s:%d, Server: %s\n", ipinfo.Ip, ipinfo.Port, CurServer)

		Utils.ChildContext, Utils.ChildCancel = context.WithCancel(rootContext)

		TaskList := Core.GenerateTask(UserList, PassList, ipinfo, CurServer)

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
			Info:     ipinfo,
			Username: Core.FlagUserName,
			Password: Utils.RandStringBytesMaskImprSrcUnsafe(12),
			Server:   CurServer,
		}

		err, RanRes := Core.DefaultScan2(RandomTask)

		if err == nil && RanRes.Result {
			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Info.Ip, RandomTask.Info.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Info.Ip, RandomTask.Info.Port)
		}
	}

	rootCancel()

	return nil
}
