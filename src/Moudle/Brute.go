package Moudle

import (
	"Zombie/src/Core"
	"Zombie/src/Utils"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var FileHandle *os.File
var O2File bool
var Datach = make(chan string, 1000)

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

	Utils.Timeout = ctx.Int("t")

	//ExpireTime := GetExpireTime(len(IpList), len(UserList), len(PassList))

	if ctx.IsSet("file") {
		initFile(ctx.String("file"))
	}

	TaskList := Core.GenerateTask(UserList, PassList, IpList, CurServer)

	wgs := &sync.WaitGroup{}

	scanPool, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		//defer cancel()
		tc := i.(Utils.ScanTask)
		defaultScan(tc)
		//for {
		//	select {
		//	case <-ctx.Done():
		//		fmt.Println("timeout")
		//		wgs.Done()
		//		return
		//	default:
		//	}
		//}
		wgs.Done()
	}, ants.WithExpiryDuration(2*time.Second))
	//,ants.WithExpiryDuration(2)

	for task := range TaskList {
		wgs.Add(1)
		_ = scanPool.Invoke(task)
	}

	//waitTimeout(wgs, time.Duration(ExpireTime)*Utils.Timeout)

	wgs.Wait()

	fmt.Printf("ScanSum is : %d\n", Core.ScanSum)

	return err
}

func defaultScan(task Utils.ScanTask) {

	err, result := Core.BruteDispatch(task)
	if err == nil && result {
		res := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", task.Info.Ip, task.Info.Port, task.Username, task.Password, task.Server)
		if O2File {
			Datach <- res
		}
		fmt.Println(res)
	}
}

func initFile(Filename string) {
	var err error

	if Filename != "" {
		O2File = true
		if checkFileIsExist(Filename) { //如果文件存在
			FileHandle, err = os.OpenFile(Filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
			//fmt.Println("文件存在")
			if err != nil {
				os.Exit(0)
			}
			//io.WriteString(FileHandle, "123")
		} else {
			FileHandle, err = os.Create(Filename) //创建文件
			//fmt.Println("文件不存在")
			if err != nil {
				os.Exit(0)
			}
			//io.WriteString(FileHandle, "123")
		}

	}
	go write2File(FileHandle, Datach)
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func write2File(FileHandle *os.File, Datach chan string) {
	for res := range Datach {
		FileHandle.WriteString(res)

	}
}
