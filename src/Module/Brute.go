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
	"os/exec"
	"path/filepath"
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

	//if Utils.HasStdin() {
	//	stdinip, _ := Core.ReadStdin(os.Stdin)
	//	IpSlice = append(IpSlice, stdinip...)
	//}

	if ctx.IsSet("ip") {
		IpSlice = append(IpSlice, Core.GetIPList(ctx.String("ip"))...)
	}

	if ctx.IsSet("IP") {
		IpdictSlice, _ := Core.ReadIPDict(ctx.String("IP"))
		IpSlice = append(IpSlice, IpdictSlice...)
	}
	if !fromgt {
		if len(IpSlice) == 0 {
			fmt.Println("Please check your input")
			return
		}

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
		fmt.Println("[+] Read from gt result")

		ipserverinfo = GenFromGT(ctx.String("gt"), ctx.String("ss"))
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

	if ctx.IsSet("cb") {
		u, p := GenFromCb(ctx.String("cb"), ctx.String("ss"))
		UserList = append(UserList, u...)
		PassList = append(PassList, p...)
	}

	if ctx.IsSet("instance") {
		Utils.Instance = Core.GetUserList(ctx.String("instance"))
	}

	Utils.Timeout = ctx.Int("timeout")
	Utils.Thread = ctx.Int("thread")
	Utils.Simple = ctx.Bool("simple")
	Utils.Proc = ctx.Int("proc")
	Utils.FileFormat = ctx.String("type")
	Utils.File = ctx.String("file")
	Utils.OutputType = "Brute"

	if Utils.File == "./.res.log" {
		Utils.File = getExcPath() + "/.res.log"
	}

	if Utils.File != "null" {
		Utils.FileHandle = Utils.InitFile(Utils.File)
		go ExecAble.QueryWrite3File(Utils.FileHandle, Utils.TDatach)
	}

	Core.Summary = len(UserList) * len(PassList) * len(IpList)

	if Utils.Proc != 0 {
		go Core.Process(Core.CountChan)
	}

	ipserverinfo = HoneyPotTest(ipserverinfo)

	if Utils.Simple {
		err = StartTaskSimple(UserList, PassList, ipserverinfo)
	} else {
		err = StartTask(UserList, PassList, ipserverinfo)
	}
	Utils.FileHandle.Close()

	close(Core.CountChan)

	fmt.Println("start analysis brute res")

	cblist, reslist := ExecAble.CleanBruteRes(&Utils.BrutedList)

	Core.OutPutRes(&reslist, &cblist, Utils.File)

	return err
}

func HoneyPotTest(IpServerList []Utils.IpServerInfo) []Utils.IpServerInfo {

	tasklist := Core.GenerateRandTask(IpServerList)

	wgs := &sync.WaitGroup{}

	scanPool, _ := ants.NewPoolWithFunc(Utils.Thread, func(i interface{}) {
		par := i.(Core.HoneyPara)
		Core.HoneyTest(&par)
		wgs.Done()
	}, ants.WithExpiryDuration(2*time.Second))

	for target := range tasklist {
		PrePara := Core.HoneyPara{
			Task: target,
		}

		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}

	wgs.Wait()
	scanPool.Release()
	fmt.Println("Honey Pot Check  done")

	var aliveinfo []Utils.IpServerInfo
	Core.NotHoney.Range(func(key, value interface{}) bool {
		if value.(bool) == true {
			aliveinfo = append(aliveinfo, key.(Utils.ScanTask).IpServerInfo)
		}
		return true
	})
	return aliveinfo
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
			TargetInfo: Utils.TargetInfo{
				IpServerInfo: ipinfo,
				Username:     Core.FlagUserName,
				Password:     Utils.RandStringBytesMaskImprSrcUnsafe(12),
			},
		}

		CurCon := Core.ExecDispatch(RandomTask)

		alive := CurCon.Connect()

		if alive {
			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Ip, RandomTask.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Ip, RandomTask.Port)
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

func GenFromCb(cbfile string, server string) (userlist, passlist []string) {
	var cblist []Utils.Codebook
	cbbytes, err := ioutil.ReadFile(cbfile)
	if err != nil {
		println(cbfile + " open failed")
		//panic(dictPath + " open failed")
		os.Exit(0)
	}

	if err := json.Unmarshal(cbbytes, &cblist); err != nil {
		println(" Unmarshal failed")
		os.Exit(0)
	}

	if server != "all" {
		var temp []Utils.Codebook
		for _, info := range cblist {
			if strings.HasPrefix(server, "~") {
				if info.Server == strings.ToUpper(server[1:]) {
					continue
				}
				temp = append(temp, info)
			} else {
				if info.Server != strings.ToUpper(server) {
					continue
				}
				temp = append(temp, info)
			}

		}
		cblist = temp
	}

	for _, info := range cblist {
		userlist = append(userlist, info.Username)
		passlist = append(passlist, info.Password)
	}

	userlist = Utils.RemoveDuplicateElement(userlist)
	passlist = Utils.RemoveDuplicateElement(passlist)

	return
}

func GenFromGT(gtfile string, server string) (ipserverinfo []Utils.IpServerInfo) {

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

	if server != "all" {
		var temp []Utils.IpServerInfo

		for _, info := range ipserverinfo {
			if strings.HasPrefix(server, "~") {
				if info.Server == strings.ToUpper(server[1:]) {
					continue
				}
				temp = append(temp, info)
			} else {

				if info.Server != strings.ToUpper(server) {
					continue
				}
				temp = append(temp, info)
			}

		}
		return temp
	}

	return ipserverinfo
}

func getExcPath() string {
	file, _ := exec.LookPath(os.Args[0])
	// 获取包含可执行文件名称的路径
	path, _ := filepath.Abs(file)
	// 获取可执行文件所在目录
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return strings.Replace(ret, "\\", "/", -1)
}
