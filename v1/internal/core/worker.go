package core

import (
	"Zombie/v1/internal/exec"
	utils2 "Zombie/v1/pkg/utils"
	"context"
	"fmt"
	"sync"
	"time"
)

var Summary int
var CountChan = make(chan int, 60)
var NotHoney sync.Map

type PoolPara struct {
	Ctx      context.Context
	Taskchan chan utils2.ScanTask
	Wgs      *sync.WaitGroup
}

type HoneyPara struct {
	Task utils2.ScanTask
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
			if utils2.Proc != 0 {
				CountChan <- 1
			}

			res := utils2.BruteRes{}
			CurCon := ExecDispatch(task)

			if CurCon == nil {
				continue
			}

			alive := CurCon.Connect()

			res.Result = alive
			if !alive {
				switch CurCon.(type) {

				case *exec.RedisService:
					res.Additional += CurCon.(*exec.RedisService).Additional
				}
				continue
			}
			CurCon.DisConnect()
			if res.Result {
				output := utils2.OutputRes{
					TargetInfo: utils2.TargetInfo{
						IpServerInfo: utils2.IpServerInfo{
							Server: task.Server,
							IpInfo: utils2.IpInfo{
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

				if utils2.O2File {
					utils2.TDatach <- output
				}

				if !utils2.Simple {
					utils2.ChildCancel()
				}
			}
		// 加入连接超时，过长直接断开
		case <-time.After(2 * time.Duration(utils2.Timeout) * time.Second):
			continue

		}
	}

}

func Process(ct chan int) {

	pr := 0

	for i := range ct {
		pr += i
		if pr%utils2.Proc == 0 {
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
