package cmd

//
//func Brute(opt Option) (err error) {
//	//var targets []utils2.Task
//	var users, pwds []string
//	var addrs ipcs.Addrs
//	//fromgt = ctx.IsSet("gt")
//
//	//if Utils.HasStdin() {
//	//	stdinip, _ := Core.ReadStdin(os.Stdin)
//	//	IpSlice = append(IpSlice, stdinip...)
//	//}
//
//	if opt.GogoFile != "" {
//
//	} else {
//		var ipslice []string
//		if opt.IP != "" {
//			ipslice = strings.Split(opt.IP, ",")
//		} else if opt.IPFile != "" {
//			ipf, err := os.OK(opt.IPFile)
//			if err != nil {
//				return err
//			}
//			ipslice = words.NewWorderWithFile(ipf).All()
//		}
//
//		if len(ipslice) == 0 {
//			return fmt.Errorf("not any ip input")
//		}
//
//		if strings.Contains(ipslice[0], ":") {
//			addrs = ipcs.NewAddrs(ipslice)
//		} else {
//			addrs = ipcs.NewAddrsWithDefaultPort(ipslice, strconv.Itoa(utils2.ServerToPortMap[opt.ServiceName]))
//		}
//	}
//
//	if opt.Username != "" {
//		users = strings.Split(opt.Username, ",")
//	} else if opt.UsernameFile != "" {
//		userf, err := os.OK(opt.UsernameFile)
//		if err != nil {
//			return err
//		}
//		users = words.NewWorderWithFile(userf).All()
//	} else {
//		users = utils2.DefaultUsernames[opt.ServiceName]
//	}
//
//	if opt.Password != "" {
//		pwds = strings.Split(opt.Password, ",")
//	} else if opt.PasswordFile != "" {
//		pwdf, err := os.OK(opt.PasswordFile)
//		if err != nil {
//			return err
//		}
//		pwds = words.NewWorderWithFile(pwdf).All()
//	} else {
//		pwds = utils2.DefaultPasswords[opt.ServiceName]
//	}
//	core2.Total = len(addrs) * len(pwds) * len(users)
//
//	//if ctx.IsSet("cb") {
//	//	u, p := GenFromCb(ctx.String("cb"), ctx.String("ss"))
//	//	UserList = append(UserList, u...)
//	//	PassList = append(PassList, p...)
//	//}
//
//	//if ctx.IsSet("instance") {
//	//	utils2.Instance = core2.GetUserList(ctx.String("instance"))
//	//}
//
//	//utils2.Timeout = ctx.Int("timeout")
//	//utils2.Thread = ctx.Int("thread")
//	//utils2.Simple = ctx.Bool("simple")
//	//utils2.Proc = ctx.Int("proc")
//	//utils2.FileFormat = ctx.String("type")
//	//utils2.File = ctx.String("file")
//	//utils2.OutputType = "Brute"
//
//	//if utils2.File == "./.res.log" {
//	//	utils2.File = getExcPath() + "/.res.log"
//	//}
//
//	//if utils2.File != "null" {
//	//	utils2.FileHandle = utils2.InitFile(utils2.File)
//	//	go exec2.QueryWrite3File(utils2.FileHandle, utils2.TDatach)
//	//}
//
//	//ipserverinfo = HoneyPotTest(ipserverinfo)
//
//	if utils2.Simple {
//		err = StartTaskSimple(UserList, PassList, ipserverinfo)
//	} else {
//		err = StartTask(UserList, PassList, ipserverinfo)
//	}
//	if err != nil {
//		return err
//	}
//
//	//fmt.Println("start analysis brute res")
//	//
//	//cblist, reslist := exec2.CleanBruteRes(&utils2.BrutedList)
//	//
//	//core2.OutPutRes(&reslist, &cblist, utils2.File)
//
//	return nil
//}
//
////func HoneyPotTest(IpServerList []utils2.IpServerInfo) []utils2.IpServerInfo {
////
////	tasklist := core2.GenerateRandTask(IpServerList)
////
////	wgs := &sync.WaitGroup{}
////
////	scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
////		par := i.(core2.HoneyPara)
////		core2.HoneyTest(&par)
////		wgs.Done()
////	}, ants.WithExpiryDuration(2*time.Second))
////
////	for target := range tasklist {
////		PrePara := core2.HoneyPara{
////			Task: target,
////		}
////
////		wgs.Add(1)
////		_ = scanPool.Invoke(PrePara)
////	}
////
////	wgs.Wait()
////	scanPool.Release()
////	fmt.Println("Honey Pot Check  done")
////
////	var aliveinfo []utils2.IpServerInfo
////	core2.NotHoney.Range(func(key, value interface{}) bool {
////		if value.(bool) == true {
////			aliveinfo = append(aliveinfo, key.(utils2.ScanTask).IpServerInfo)
////		}
////		return true
////	})
////	return aliveinfo
////}
//
//func StartTask(UserList []string, PassList []string, IpServerList []utils2.IpServerInfo) error {
//	rootContext, rootCancel := context.WithCancel(context.Background())
//	for _, ipinfo := range IpServerList {
//
//		fmt.Printf("Now Processing %s:%d, ExecAble: %s\n", ipinfo.Ip, ipinfo.Port, ipinfo.Server)
//
//		utils2.ChildContext, utils2.ChildCancel = context.WithCancel(rootContext)
//
//		TaskList := core2.GenerateTask(UserList, PassList, ipinfo)
//
//		wgs := &sync.WaitGroup{}
//		PrePara := core2.PoolPara{
//			Ctx:      utils2.ChildContext,
//			Taskchan: TaskList,
//			Wgs:      wgs,
//		}
//
//		scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
//			par := i.(core2.PoolPara)
//			core2.BruteWork(&par)
//		}, ants.WithExpiryDuration(2*time.Second))
//
//		for i := 0; i < utils2.Thread; i++ {
//			wgs.Add(1)
//			_ = scanPool.Invoke(PrePara)
//		}
//		wgs.Wait()
//
//		RandomTask := utils2.ScanTask{
//			TargetInfo: utils2.TargetInfo{
//				IpServerInfo: ipinfo,
//				Username:     core2.FlagUserName,
//				Password:     utils2.RandStringBytesMaskImprSrcUnsafe(12),
//			},
//		}
//
//		CurCon := core2.ExecDispatch(RandomTask)
//
//		alive := CurCon.Connect()
//
//		if alive {
//			fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", RandomTask.Ip, RandomTask.Port, RandomTask.Username, RandomTask.Password, RandomTask.Server)
//			fmt.Sprintf("%s:%d\t is it a honeypot?", RandomTask.Ip, RandomTask.Port)
//		}
//	}
//
//	fmt.Println("All Task done")
//
//	time.Sleep(1000 * time.Millisecond)
//	if utils2.FileFormat == "json" {
//		final := utils2.Result{}
//		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
//		if errs != nil {
//			fmt.Println(errs.Error())
//		}
//		utils2.FileHandle.WriteString(string(jsons) + "]")
//	}
//
//	rootCancel()
//
//	return nil
//}
//
//func StartTaskSimple(UserList []string, PassList []string, IpServerList []utils2.IpServerInfo) error {
//	rootContext, rootCancel := context.WithCancel(context.Background())
//
//	TaskList := core2.GenerateTaskSimple(UserList, PassList, IpServerList)
//
//	var wgs sync.WaitGroup
//	scanPool, _ := ants.NewPoolWithFunc(utils2.Thread, func(i interface{}) {
//		core2.BruteWork(rootContext)
//		wgs.Done()
//	}, ants.WithExpiryDuration(2*time.Second))
//
//	for i := 0; i < utils2.Thread; i++ {
//		wgs.Add(1)
//		_ = scanPool.Invoke(PrePara)
//	}
//	wgs.Wait()
//
//	time.Sleep(1000 * time.Millisecond)
//
//	fmt.Println("All Task done")
//
//	if utils2.FileFormat == "json" {
//		final := utils2.Result{}
//		jsons, errs := json.Marshal(final) //转换成JSON返回的是byte[]
//		if errs != nil {
//			fmt.Println(errs.Error())
//		}
//		utils2.FileHandle.WriteString(string(jsons) + "]")
//	}
//
//	rootCancel()
//
//	return nil
//}
//
//func GenerIPServerInfo(ipinfo []utils2.Task, server string) (ipserverinfo []utils2.IpServerInfo) {
//	for _, info := range ipinfo {
//		isinfo := utils2.IpServerInfo{}
//		isinfo.Task = info
//		isinfo.Server = server
//		ipserverinfo = append(ipserverinfo, isinfo)
//	}
//
//	return ipserverinfo
//}
//
//func GenFromCb(cbfile string, server string) (userlist, passlist []string) {
//	var cblist []utils2.Codebook
//	cbbytes, err := ioutil.ReadFile(cbfile)
//	if err != nil {
//		println(cbfile + " open failed")
//		//panic(dictPath + " open failed")
//		os.Exit(0)
//	}
//
//	if err := json.Unmarshal(cbbytes, &cblist); err != nil {
//		println(" Unmarshal failed")
//		os.Exit(0)
//	}
//
//	if server != "all" {
//		var temp []utils2.Codebook
//		for _, info := range cblist {
//			if strings.HasPrefix(server, "~") {
//				if info.Server == strings.ToUpper(server[1:]) {
//					continue
//				}
//				temp = append(temp, info)
//			} else {
//				if info.Server != strings.ToUpper(server) {
//					continue
//				}
//				temp = append(temp, info)
//			}
//
//		}
//		cblist = temp
//	}
//
//	for _, info := range cblist {
//		userlist = append(userlist, info.Username)
//		passlist = append(passlist, info.Password)
//	}
//
//	userlist = utils2.RemoveDuplicateElement(userlist)
//	passlist = utils2.RemoveDuplicateElement(passlist)
//
//	return
//}
//
//func GenFromGT(gtfile string, server string) (ipserverinfo []utils2.IpServerInfo) {
//
//	bytes, err := ioutil.ReadFile(gtfile)
//	if err != nil {
//		println(gtfile + " open failed")
//		//panic(dictPath + " open failed")
//		os.Exit(0)
//	}
//
//	if err := json.Unmarshal(bytes, &ipserverinfo); err != nil {
//		println(" Unmarshal failed")
//		os.Exit(0)
//	}
//
//	if server != "all" {
//		var temp []utils2.IpServerInfo
//
//		for _, info := range ipserverinfo {
//			if strings.HasPrefix(server, "~") {
//				if info.Server == strings.ToUpper(server[1:]) {
//					continue
//				}
//				temp = append(temp, info)
//			} else {
//
//				if info.Server != strings.ToUpper(server) {
//					continue
//				}
//				temp = append(temp, info)
//			}
//
//		}
//		return temp
//	}
//
//	return ipserverinfo
//}
//
//func getExcPath() string {
//	file, _ := exec.LookPath(os.Args[0])
//	// 获取包含可执行文件名称的路径
//	path, _ := filepath.Abs(file)
//	// 获取可执行文件所在目录
//	index := strings.LastIndex(path, string(os.PathSeparator))
//	ret := path[:index]
//	return strings.Replace(ret, "\\", "/", -1)
//}
