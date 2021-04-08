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

			err, res := DefaultScan2(task)
			if err != nil {
				continue
			}
			if res {
				Bres := fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\n", task.Info.Ip, task.Info.Port, task.Username, task.Password, task.Server)
				FlagUserName = task.Username
				if Utils.O2File {
					Utils.Datach <- Bres
				}
				fmt.Println(Bres)
				Utils.ChildCancel()
			}

		}
	}

}

func DefaultScan2(task Utils.ScanTask) (error, bool) {
	err, result := BruteDispatch(task)

	return err, result
}
