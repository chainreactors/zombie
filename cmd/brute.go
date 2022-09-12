package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	core2 "github.com/chainreactors/zombie/internal/core"
	exec2 "github.com/chainreactors/zombie/internal/exec"
	utils2 "github.com/chainreactors/zombie/pkg/utils"
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
	var IpList []utils2.IpInfo
	var fromgt bool
	var ipserverinfo []utils2.IpServerInfo

	fromgt = ctx.IsSet("gt")

	//if Utils.HasStdin() {
	//	stdinip, _ := Core.ReadStdin(os.Stdin)
	//	IpSlice = append(IpSlice, stdinip...)
	//}

	if ctx.IsSet("ip") {
		IpSlice = append(IpSlice, core2.GetIPList(ctx.String("ip"))...)
	}

	if ctx.IsSet("IP") {
		IpdictSlice, err := core2.ReadIPDict(ctx.String("IP"))
		if err != nil {
			return err
		}
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
			if _, ok := utils2.ServerPort[ServerName]; ok {
				CurServer = ServerName
			} else {
				return fmt.Errorf("the ExecAble isn't be supported")
			}

		} else if strings.Contains(Ip, ":") {
			Temp := strings.Split(Ip, ":")
			Sport := Temp[1]
			port, err := strconv.Atoi(Sport)
			if err != nil {
				return err
			}

			if _, ok := utils2.PortServer[port]; ok {
				CurServer = utils2.PortServer[port]
				fmt.Println("Use default server")
			} else {
				return fmt.Errorf("Please input the type of ExecAble")

			}
		} else {
			return fmt.Errorf("Please input the type of ExecAble")
		}

		CurServer = strings.ToUpper(CurServer)

		IpList = core2.GetIpInfoList(IpSlice, CurServer)
		ipserverinfo = GenerIPServerInfo(IpList, CurServer)
	} else {
		fmt.Println("[+] Read from gt result")

		ipserverinfo = GenFromGT(ctx.String("gt"), ctx.String("ss"))
	}

	if ctx.IsSet("uppair") {
		uppair := ctx.String("uppair")
		core2.UPList, _ = core2.GetUAList(uppair)
	} else {

		if ctx.IsSet("username") && !ctx.IsSet("userdict") {
			username := ctx.String("username")
			UserList = core2.GetUserList(username)
		} else if ctx.IsSet("userdict") {
			UserList, _ = core2.ReadUserDict(ctx.String("userdict"))
		} else {
			fmt.Println("[+] Use default user dict")
		}

		if ctx.IsSet("password") && !ctx.IsSet("passdict") {
			password := ctx.String("password")
			PassList = core2.GetPassList(password)
		} else if ctx.IsSet("passdict") {
			PassList, _ = core2.ReadPassDict(ctx.String("passdict"))
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
		utils2.Instance = core2.GetUserList(ctx.String("instance"))
	}

	utils2.Timeout = ctx.Int("timeout")
	utils2.Thread = ctx.Int("thread")
	utils2.Simple = ctx.Bool("simple")
	utils2.Proc = ctx.Int("proc")
	utils2.FileFormat = ctx.String("type")
	utils2.File = ctx.String("file")
	utils2.OutputType = "Brute"

	if utils2.File == "./.res.log" {
		utils2.File = getExcPath() + "/.res.log"
	}

	if utils2.File != "null" {
		utils2.FileHandle = utils2.InitFile(utils2.File)
		go exec2.QueryWrite3File(utils2.FileHandle, utils2.TDatach)
	}

	core2.Summary = len(UserList) * len(PassList) * len(IpList)

	if utils2.Proc != 0 {
		go core2.Process(core2.CountChan)
	}

	//ipserverinfo = HoneyPotTest(ipserverinfo)

	if utils2.Simple {
		err = StartTaskSimple(UserList, PassList, ipserverinfo)
	} else {
		err = StartTask(UserList, PassList, ipserverinfo)
	}
	if err != nil {
		return err
	}
	utils2.FileHandle.Close()

	close(core2.CountChan)

	fmt.Println("start analysis brute res")

	cblist, reslist := exec2.CleanBruteRes(&utils2.BrutedList)

	core2.OutPutRes(&reslist, &cblist, utils2.File)

	return nil
}

func HoneyPotTest(IpServerList []utils2.IpServerInfo) []utils2.IpServerInfo {

	tasklist := core2.GenerateRandTask(IpServerList)

	wgs := &sync.WaitGroup{}

	scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
		par := i.(core2.HoneyPara)
		core2.HoneyTest(&par)
		wgs.Done()
	}, ants.WithExpiryDuration(2*time.Second))

	for target := range tasklist {
		PrePara := core2.HoneyPara{
			Task: target,
		}

		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}

	wgs.Wait()
	scanPool.Release()
	fmt.Println("Honey Pot Check  done")

	var aliveinfo []utils2.IpServerInfo
	core2.NotHoney.Range(func(key, value interface{}) bool {
		if value.(bool) == true {
			aliveinfo = append(aliveinfo, key.(utils2.ScanTask).IpServerInfo)
		}
		return true
	})
	return aliveinfo
}

func StartTask(UserList []string, PassList []string, IpServerList []utils2.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())
	for _, ipinfo := range IpServerList {

		fmt.Printf("Now Processing %s:%d, ExecAble: %s\n", ipinfo.Ip, ipinfo.Port, ipinfo.Server)

		utils2.ChildContext, utils2.ChildCancel = context.WithCancel(rootContext)

		TaskList := core2.GenerateTask(UserList, PassList, ipinfo)

		wgs := &sync.WaitGroup{}
		PrePara := core2.PoolPara{
			Ctx:      utils2.ChildContext,
			Taskchan: TaskList,
			Wgs:      wgs,
		}

		scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
			par := i.(core2.PoolPara)
			core2.BruteWork(&par)
		}, ants.WithExpiryDuration(2*time.Second))

		for i := 0; i < utils2.Thread; i++ {
			wgs.Add(1)
			_ = scanPool.Invoke(PrePara)
		}
		wgs.Wait()

		RandomTask := utils2.ScanTask{
			TargetInfo: utils2.TargetInfo{
				IpServerInfo: ipinfo,
				Username:     core2.FlagUserName,
				Password:     utils2.RandStringBytesMaskImprSrcUnsafe(12),
			},
		}

		CurCon := core2.ExecDispatch(RandomTask)

		alive := CurCon.Connect()

		if alive {
			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Ip, RandomTask.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Ip, RandomTask.Port)
		}
	}

	fmt.Println("All Task done")

	time.Sleep(1000 * time.Millisecond)
	if utils2.FileFormat == "json" {
		final := utils2.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		utils2.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func StartTaskSimple(UserList []string, PassList []string, IpServerList []utils2.IpServerInfo) error {
	rootContext, rootCancel := context.WithCancel(context.Background())

	TaskList := core2.GenerateTaskSimple(UserList, PassList, IpServerList)

	wgs := &sync.WaitGroup{}
	PrePara := core2.PoolPara{
		Ctx:      rootContext,
		Taskchan: TaskList,
		Wgs:      wgs,
	}

	scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
		par := i.(core2.PoolPara)
		core2.BruteWork(&par)
	}, ants.WithExpiryDuration(2*time.Second))

	for i := 0; i < utils2.Thread; i++ {
		wgs.Add(1)
		_ = scanPool.Invoke(PrePara)
	}
	wgs.Wait()

	time.Sleep(1000 * time.Millisecond)

	fmt.Println("All Task done")

	if utils2.FileFormat == "json" {
		final := utils2.OutputRes{}
		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
		if errs != nil {
			fmt.Println(errs.Error())
		}
		utils2.FileHandle.WriteString(string(jsons) + "]")
	}

	rootCancel()

	return nil
}

func GenerIPServerInfo(ipinfo []utils2.IpInfo, server string) (ipserverinfo []utils2.IpServerInfo) {
	for _, info := range ipinfo {
		isinfo := utils2.IpServerInfo{}
		isinfo.IpInfo = info
		isinfo.Server = server
		ipserverinfo = append(ipserverinfo, isinfo)
	}

	return ipserverinfo
}

func GenFromCb(cbfile string, server string) (userlist, passlist []string) {
	var cblist []utils2.Codebook
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
		var temp []utils2.Codebook
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

	userlist = utils2.RemoveDuplicateElement(userlist)
	passlist = utils2.RemoveDuplicateElement(passlist)

	return
}

func GenFromGT(gtfile string, server string) (ipserverinfo []utils2.IpServerInfo) {

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
		var temp []utils2.IpServerInfo

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
