package Core

import (
	"Zombie/src/Utils"
)

var UPList []string

func GenerateTask(UserList []string, PassList []string, info Utils.IpServerInfo) chan Utils.ScanTask {
	TaskList := make(chan Utils.ScanTask)
	var UserIsDefault bool
	var PassIsDefault bool

	if len(UserList) == 0 {
		UserIsDefault = true
	}

	if len(PassList) == 0 {
		PassIsDefault = false
	}

	go func() {

		if len(UPList) == 0 {
			if UserIsDefault {
				if defaultuser, ok := Utils.DefaultUserDict[info.Server]; ok {
					UserList = defaultuser
				} else {
					UserList = []string{"admin"}
				}
			}

			if PassIsDefault {
				PassList = Utils.DefaultPasswords[info.Server]
			}

			for _, username := range UserList {
				for _, password := range PassList {
					NewTask := Utils.ScanTask{
						TargetInfo: Utils.TargetInfo{
							IpServerInfo: info,
							Username:     username,
							Password:     password,
						},
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
					TargetInfo: Utils.TargetInfo{
						IpServerInfo: info,
						Username:     username,
						Password:     password,
					},
				}
				TaskList <- NewTask

			}

		}

		close(TaskList)
	}()

	return TaskList
}

func GenerateRandTask(ipinfo []Utils.IpServerInfo) chan Utils.ScanTask {
	TaskList := make(chan Utils.ScanTask)
	go func() {
		for _, info := range ipinfo {
			NewTask := Utils.ScanTask{
				TargetInfo: Utils.TargetInfo{
					IpServerInfo: info,
					Username:     "admin",
					Password:     Utils.RandStringBytesMaskImprSrcUnsafe(12),
				},
			}
			if info.Server == "ORACLE" {
				NewTask.Instance = "orcl"
			}
			TaskList <- NewTask
		}
		close(TaskList)
	}()
	return TaskList
}

func GenerateTaskSimple(UserList []string, PassList []string, ipinfo []Utils.IpServerInfo) chan Utils.ScanTask {

	TaskList := make(chan Utils.ScanTask)
	var UserIsDefault bool
	var PassIsDefault bool

	if len(UserList) == 0 {
		UserIsDefault = true
	}

	if len(PassList) == 0 {
		PassIsDefault = true
	}

	go func() {
		for _, info := range ipinfo {
			if UserIsDefault {
				if defaultuser, ok := Utils.DefaultUserDict[info.Server]; ok {
					UserList = defaultuser
				} else {
					UserList = []string{"admin"}
				}
			}

			if PassIsDefault {
				PassList = Utils.DefaultPasswords[info.Server]
			}

			Summary = len(PassList) * len(UserList) * len(ipinfo)

			if len(UPList) == 0 {
				//PassList = append(PassList, Utils.RandStringBytesMaskImprSrcUnsafe(8))
				for _, username := range UserList {
					for _, password := range PassList {
						NewTask := Utils.ScanTask{
							TargetInfo: Utils.TargetInfo{
								IpServerInfo: info,
								Username:     username,
								Password:     password,
							},
						}
						if info.Server == "ORACLE" {
							for _, ins := range Utils.Instance {
								NewTask.Instance = ins
								TaskList <- NewTask
							}

						} else {
							TaskList <- NewTask
						}

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
						TargetInfo: Utils.TargetInfo{
							IpServerInfo: info,
							Username:     username,
							Password:     password,
						},
					}
					if info.Server == "ORACLE" {
						for _, ins := range Utils.Instance {
							NewTask.Instance = ins
							TaskList <- NewTask
						}

					} else {
						TaskList <- NewTask
					}

				}

			}
		}

		close(TaskList)
	}()

	return TaskList
}
