package core

import (
	"github.com/chainreactors/zombie/pkg/utils"
)

//func HoneyTest(WorkerPara *HoneyPara) {
//	CurCon := PluginDispatch(WorkerPara.Task)
//	if CurCon == nil {
//		return
//	}
//
//	alive := CurCon.Connect()
//
//	if !alive {
//		NotHoney.Store(WorkerPara.Task, true)
//	} else {
//		fmt.Printf("%s:%v\t%v maybe honey pot\n", WorkerPara.Task.Ip, WorkerPara.Task.Port, WorkerPara.Task.Server)
//	}
//	return
//}

func Brute(task *utils.Task) *utils.Result {
	conn := PluginDispatch(task)

	err := conn.Connect()
	if err != nil {
		return &utils.Result{
			Task: task,
			Err:  err,
			OK:   false,
		}
	}
	defer conn.Close()
	return &utils.Result{
		Task: task,
		Err:  err,
		OK:   true,
	}
}

//
//func DefaultScan2(task Utils.ScanTask) (error, Utils.BruteRes) {
//	err, result := BruteDispatch(task)
//
//	return err, result
//}
