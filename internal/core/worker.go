package core

import (
	"context"
	"fmt"
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg/utils"
	"sync"
	"time"
)

var Summary int
var CountChan = make(chan int, 60)
var NotHoney sync.Map

type PoolPara struct {
	Ctx      context.Context
	Taskchan chan utils.ScanTask
	Wgs      *sync.WaitGroup
}

type HoneyPara struct {
	Task utils.ScanTask
}

var FlagUserName string

func HoneyTest(WorkerPara *HoneyPara) {
	CurCon := ExecDispatch(WorkerPara.Task)
	if CurCon == nil {
		return
	}

	alive := CurCon.Connect()

	if !alive {
		NotHoney.Store(WorkerPara.Task, true)
	} else {
		fmt.Printf("%s:%v\t%v maybe honey pot\n", WorkerPara.Task.Ip, WorkerPara.Task.Port, WorkerPara.Task.Server)
	}
	return
}

func BruteWork(WorkerPara *PoolPara) {
	defer WorkerPara.Wgs.Done()

	for {
		select {
		case <-WorkerPara.Ctx.Done():
			return
		case task, ok := <-WorkerPara.Taskchan:
			if !ok {
				return
			}
			if utils.Proc != 0 {
				CountChan <- 1
			}

			res := utils.BruteRes{}
			CurCon := ExecDispatch(task)

			if CurCon == nil {
				continue
			}

			alive := CurCon.Connect()

			res.Result = alive
			if !alive {
				switch CurCon.(type) {

				case *plugin.RedisService:
					res.Additional += CurCon.(*plugin.RedisService).Additional
				}
				continue
			}
			CurCon.DisConnect()
			if res.Result {
				output := utils.OutputRes{
					TargetInfo: utils.TargetInfo{
						IpServerInfo: utils.IpServerInfo{
							Server: task.Server,
							IpInfo: utils.IpInfo{
								Ip:   task.Ip,
								Port: task.Port,
							},
						},
						Username: task.Username,
						Password: task.Password,
					},
					Additional: res.Additional,
				}
				if task.Server == "ORACLE" {
					output.Additional = task.Instance
				}

				FlagUserName = task.Username

				if utils.O2File {
					utils.TDatach <- output
				}

				if !utils.Simple {
					utils.ChildCancel()
				}
			}
		// 加入连接超时，过长直接断开
		case <-time.After(2 * time.Duration(utils.Timeout) * time.Second):
			continue

		}
	}

}

func Process(ct chan int) {

	pr := 0

	for i := range ct {
		pr += i
		if pr%utils.Proc == 0 {
			fmt.Printf("(%d/%d)\n", pr, Summary)
		}

	}
	return
}

//
//func DefaultScan2(task Utils.ScanTask) (error, Utils.BruteRes) {
//	err, result := BruteDispatch(task)
//
//	return err, result
//}
