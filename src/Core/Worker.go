package Core

import (
	"Zombie/src/Utils"
	"context"
	"fmt"
	"sync"
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
			Bres := ""
			if !ok {
				return
			}
			CountChan <- 1
			err, res := DefaultScan2(task)
			if err != nil {
				if task.Server == "SMB" && task.Password != "" {
					if res.Additional != "" {
						Bres = fmt.Sprintf("%s:%d\t\tVersion:%s", task.Info.Ip, task.Info.Port, res.Additional)
						fmt.Println(Bres)
					}
				}
				continue
			}

			if res.Result {
				output := Utils.OutputRes{
					Type:       task.Server,
					IP:         task.Info.Ip,
					Port:       task.Info.Port,
					Username:   task.Username,
					Password:   task.Password,
					Additional: res.Additional,
				}

				FlagUserName = task.Username

				if Utils.O2File {
					Utils.TDatach <- output
				}

				if !Utils.Simple {
					Utils.ChildCancel()
				}
			}

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

func DefaultScan2(task Utils.ScanTask) (error, Utils.BruteRes) {
	err, result := BruteDispatch(task)

	return err, result
}
