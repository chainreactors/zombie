package Core

import (
	"Zombie/src/ExecAble"
	"Zombie/src/Utils"
	"context"
	"fmt"
	"sync"
	"time"
)

var Summary int
var CountChan = make(chan int, 60)

type PoolPara struct {
	Ctx      context.Context
	Taskchan chan Utils.ScanTask
	Wgs      *sync.WaitGroup
}

var FlagUserName string

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
			if Utils.Proc != 0 {
				CountChan <- 1
			}

			res := Utils.BruteRes{}
			CurCon := ExecDispatch(task)

			if CurCon == nil {
				continue
			}

			alive := CurCon.Connect()

			res.Result = alive
			if !alive {
				switch CurCon.(type) {

				case *ExecAble.RedisService:
					res.Additional += CurCon.(*ExecAble.RedisService).Additional
				}
				continue
			}
			CurCon.DisConnect()
			if res.Result {
				output := Utils.OutputRes{
					TargetInfo: Utils.TargetInfo{
						IpServerInfo: Utils.IpServerInfo{
							Server: task.Server,
							IpInfo: Utils.IpInfo{
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

				if Utils.O2File {
					Utils.TDatach <- output
				}

				if !Utils.Simple {
					Utils.ChildCancel()
				}
			}
		// 加入连接超时，过长直接断开
		case <-time.After(2 * time.Duration(Utils.Timeout) * time.Second):
			continue

		}
	}

}

func Process(ct chan int) {

	pr := 0

	for i := range ct {
		pr += i
		if pr%Utils.Proc == 0 {
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
