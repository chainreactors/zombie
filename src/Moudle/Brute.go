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

func Brute(ctx *cli.Context) (err error) {
	var CurServer string
	var UserList, PassList []string
	var IpList []Utils.IpInfo

	if !ctx.IsSet("ip") {
		fmt.Println("Please check the address")
		os.Exit(0)
	}

	IpSlice := Core.GetIpList(ctx.String("ip"))

	Ip := IpSlice[0]
	if ctx.IsSet("server") {
		ServerName := strings.ToUpper(ctx.String("server"))
		if _, ok := Utils.ServerPort[ServerName]; ok {
			CurServer = ctx.String("server")
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

	Utils.Timeout = time.Duration(ctx.Int("t")) * time.Second

	if CurServer == "SSH" {
		Utils.Timeout += 10 * time.Second
	}

	ExpireTime := GetExpireTime(len(IpList), len(UserList), len(PassList))

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

	waitTimeout(wgs, time.Duration(ExpireTime)*Utils.Timeout)

	//wgs.Wait()

	fmt.Printf("ScanSum is : %d\n", Core.ScanSum)

	return err
}

func defaultScan(task Utils.ScanTask) {

	err, result := Core.BruteDispatch(task)
	if err == nil && result {
		fmt.Printf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", task.Info.Ip, task.Info.Port, task.Username, task.Password, task.Server)
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		fmt.Println("The program timeout ended,maybe the targets contains someone don't belongs to this server")
		return true // timed out
	}
}

func GetExpireTime(a int, b int, c int) (res int) {
	summary := a * b * c
	res = summary/100 + 10
	fmt.Printf("The summary is: %d\n", summary)
	return res
}
