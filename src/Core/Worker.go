package Core

import (
	"Zombie/src/Utils"
	"context"
	"fmt"
	"sync"
)

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

			//if Utils.Simple == 1 {
			//	fmt.Printf("Now Processing %s:%d, Server: %s\n", task.Info.Ip, task.Info.Port, task.Server)
			//}
			err, res := DefaultScan2(task)
			if err != nil {
				continue
			}
			if res.Result {
				Bres := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t%s", task.Info.Ip, task.Info.Port, task.Username, task.Password, task.Server, res.Additional)
				FlagUserName = task.Username
				if Utils.O2File {
					Utils.Datach <- Bres
				}
				fmt.Println(Bres)

				if !Utils.Simple {
					Utils.ChildCancel()
				}
			}

		}
	}

}

func DefaultScan2(task Utils.ScanTask) (error, Utils.BruteRes) {
	err, result := BruteDispatch(task)

	return err, result
}
