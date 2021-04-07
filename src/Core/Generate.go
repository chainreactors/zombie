package Core

import "Zombie/src/Utils"

func GenerateTask(UserList []string, PassList []string, info Utils.IpInfo, CurServer string) chan Utils.ScanTask {
	TaskList := make(chan Utils.ScanTask)
	go func() {

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

		close(TaskList)
	}()

	return TaskList
}
