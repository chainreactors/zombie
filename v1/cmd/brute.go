package cmd

import (
	"Zombie/v1/internal/core"
	exec2 "Zombie/v1/internal/exec"
	"Zombie/v1/pkg/utils"
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
	var IpList []utils.IpInfo
	var fromgt bool
	var ipserverinfo []utils.IpServerInfo

	fromgt = ctx.IsSet("gt")

	//if Utils.HasStdin() {
	//	stdinip, _ := Core.ReadStdin(os.Stdin)
	//	IpSlice = append(IpSlice, stdinip...)
	//}

	if ctx.IsSet("ip") {
		IpSlice = append(IpSlice, core.GetIPList(ctx.String("ip"))...)
	}

	if ctx.IsSet("IP") {
		IpdictSlice, _ := core.ReadIPDict(ctx.String("IP"))
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
			if _, ok := utils.ServerPort[ServerName]; ok {
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

			if _, ok := utils.PortServer[port]; ok {
				CurServer = utils.PortServer[port]
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

		IpList = core.GetIpInfoList(IpSlice, CurServer)
		ipserverinfo = GenerIPServerInfo(IpList, CurServer)
	} else {
		fmt.Println("[+] Read from gt result")

		ipserverinfo = GenFromGT(ctx.String("gt"), ctx.String("ss"))
	}

	if ctx.IsSet("uppair") {
		uppair := ctx.String("uppair")
		core.UPList, _ = core.GetUAList(uppair)
	} else {

		if ctx.IsSet("username") && !ctx.IsSet("userdict") {
			username := ctx.String("username")
			UserList = core.GetUserList(username)
		} else if ctx.IsSet("userdict") {
			UserList, _ = core.ReadUserDict(ctx.String("userdict"))
		} else {
			fmt.Println("[+] Use default user dict")
		}

		if ctx.IsSet("password") && !ctx.IsSet("passdict") {
			password := ctx.String("password")
			PassList = core.GetPassList(password)
		} else if ctx.IsSet("passdict") {
			PassList, _ = core.ReadPassDict(ctx.String("passdict"))
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
		utils.Instance = core.GetUserList(ctx.String("instance"))
	}

	utils.Timeout = ctx.Int("timeout")
	utils.Thread = ctx.Int("thread")
	utils.Simple = ctx.Bool("simple")
	utils.Proc = ctx.Int("proc")
	utils.FileFormat = ctx.String("type")
	utils.File = ctx.String("file")
	utils.OutputType = "Brute"

	if utils.File == "./.res.log" {
		utils.File = getExcPath() + "/.res.log"
	}

	if utils.File != "null" {
		utils.FileHandle = utils.InitFile(utils.File)
		go exec2.QueryWrite3File(utils.FileHandle, utils.TDatach)
	}

	core.Summary = len(UserList) * len(PassList) * len(IpList)

	if utils.Proc != 0 {
		go core.Process(core.CountChan)
	}

	ipserverinfo = HoneyPotTest(ipserverinfo)

	if utils.Simple {
		err = StartTaskSimple(UserList, PassList, ipserverinfo)
	} else {
		err = StartTask(UserList, PassList, ipserverinfo)
	}
	utils.FileHandle.Close()

	close(core.CountChan)

	fmt.Println("start analysis brute res")

	cblist, reslist := exec2.CleanBruteRes(&utils.BrutedList)

	core.OutPutRes(&reslist, &cblist, utils.File)

	return err
}

func HoneyPotTest(IpServerList []utils.IpServerInfo) []utils.IpServerInfo {

	tasklist := core.GenerateRandTask(IpServerList)

	wgs := &sync.WaitGroup{}

	scanPool, _ := ants.NewPoolWithFunc(utils.Thread, func(i interface{}) {
		par := i.(core.HoneyPara)
		core.HoneyTest(&par)
		wgs.Done()
	}, ants.WithExpiryDuration(2*time.Second))

	for target := range tasklist {
		PrePara := core.HoneyPara{
			Task: target,
		}

		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}

	wgs.Wait()
	scanPool.Release()
	fmt.Println("Honey Pot Check  done")

	var aliveinfo []utils.IpServerInfo
	core.NotHoney.Range(func(key, value interface{}) bool {
		if value.(bool) == true {
			aliveinfo = append(aliveinfo, key.(utils.ScanTask).IpServerInfo)
		}
		return true
	})
	return aliveinfo
}

func StartTask(UserList []string, PassList []string, IpServerList []utils.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())
	for _, ipinfo := range IpServerList {

		fmt.Printf("Now Processing %s:%d, ExecAble: %s\n", ipinfo.Ip, ipinfo.Port, ipinfo.Server)

		utils.ChildContext, utils.ChildCancel = context.WithCancel(rootContext)

		TaskList := core.GenerateTask(UserList, PassList, ipinfo)

		wgs := &sync.WaitGroup{}
		PrePara := core.PoolPara{
			Ctx:      utils.ChildContext,
			Taskchan: TaskList,
			Wgs:      wgs,
		}

		scanPool, _ := ants.NewPoolWithFunc(utils.Thread, func(i interface{}) {
			par := i.(core.PoolPara)
			core.BruteWork(&par)
		}, ants.WithExpiryDuration(2*time.Second))

		for i := 0; i < utils.Thread; i++ {
			wgs.Add(1)
			_ = scanPool.Invoke(PrePara)
		}
		wgs.Wait()

		RandomTask := utils.ScanTask{
			TargetInfo: utils.TargetInfo{
				IpServerInfo: ipinfo,
				Username:     core.FlagUserName,
				Password:     utils.RandStringBytesMaskImprSrcUnsafe(12),
			},
		}

		CurCon := core.ExecDispatch(RandomTask)

		alive := CurCon.Connect()

		if alive {
			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Ip, RandomTask.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Ip, RandomTask.Port)
		}
	}

	fmt.Println("All Task done")

	time.Sleep(1000 * time.Millisecond)
	if utils.FileFormat == "json" {
		final := utils.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		utils.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func StartTaskSimple(UserList []string, PassList []string, IpServerList []utils.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())

	TaskList := core.GenerateTaskSimple(UserList, PassList, IpServerList)

	wgs := &sync.WaitGroup{}
	PrePara := core.PoolPara{
		Ctx:      rootContext,
		Taskchan: TaskList,
		Wgs:      wgs,
	}

	scanPool, _ := ants.NewPoolWithFunc(utils.Thread, func(i interface{}) {
		par := i.(core.PoolPara)
		core.BruteWork(&par)
	}, ants.WithExpiryDuration(2*time.Second))

	for i := 0; i < utils.Thread; i++ {
		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}
	wgs.Wait()

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("All Task done")

	if utils.FileFormat == "json" {
		final := utils.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		utils.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func GenerIPServerInfo(ipinfo []utils.IpInfo, server string) (ipserverinfo []utils.IpServerInfo) {
	for _, info := range ipinfo {
		isinfo := utils.IpServerInfo{}
		isinfo.IpInfo = info
		isinfo.Server = server
		ipserverinfo = append(ipserverinfo, isinfo)
	}

	return ipserverinfo
}

func GenFromCb(cbfile string, server string) (userlist, passlist []string) {
	var cblist []utils.Codebook
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
		var temp []utils.Codebook
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

	userlist = utils.RemoveDuplicateElement(userlist)
	passlist = utils.RemoveDuplicateElement(passlist)

	return
}

func GenFromGT(gtfile string, server string) (ipserverinfo []utils.IpServerInfo) {

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
		var temp []utils.IpServerInfo

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
