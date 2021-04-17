package Core

import "Zombie/src/Utils"

var UPList []string

func GenerateTask(UserList []string, PassList []string, info Utils.IpInfo, CurServer string) chan Utils.ScanTask {
	TaskList := make(chan Utils.ScanTask)
	go func() {

		if len(UPList) == 0 {

			for _, username := range UserList {
				for _, password := range PassList {
					NewTask := Utils.ScanTask{
						Info:     info,
						Username: username,
						Password: password,
						Server:   CurServer,
					}
					TaskList <- NewTask
				}
			}
		} else {
			//UPList = append(UPList, "root "+ Utils.RandStringBytesMaskImprSrcUnsafe(8))

			for _, up := range UPList {
				username, password, err := Utils.GetUP(up)
				if err != nil {
					continue
				}
				NewTask := Utils.ScanTask{
					Info:     info,
					Username: username,
					Password: password,
					Server:   CurServer,
				}
				TaskList <- NewTask

			}

		}

		close(TaskList)
	}()

	return TaskList
}

func GenerateTaskSimple(UserList []string, PassList []string, ipinfo []Utils.IpInfo, CurServer string) chan Utils.ScanTask {
	TaskList := make(chan Utils.ScanTask)
	go func() {
		for _, info := range ipinfo {
			if len(UPList) == 0 {
				//PassList = append(PassList, Utils.RandStringBytesMaskImprSrcUnsafe(8))
				for _, username := range UserList {
					for _, password := range PassList {
						NewTask := Utils.ScanTask{
							Info:     info,
							Username: username,
							Password: password,
							Server:   CurServer,
						}
						TaskList <- NewTask
					}
				}

			} else {
				//UPList = append(UPList, "root "+ Utils.RandStringBytesMaskImprSrcUnsafe(8))
				for _, up := range UPList {
					username, password, err := Utils.GetUP(up)
					if err != nil {
						continue
					}
					NewTask := Utils.ScanTask{
						Info:     info,
						Username: username,
						Password: password,
						Server:   CurServer,
					}
					TaskList <- NewTask

				}

			}
		}

		close(TaskList)
	}()

	return TaskList
}
