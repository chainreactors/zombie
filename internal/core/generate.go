package core

import (
	utils2 "github.com/chainreactors/zombie/pkg/utils"
)

var UPList []string

func GenerateTask(UserList []string, PassList []string, info utils2.IpServerInfo) chan utils2.ScanTask {
	TaskList := make(chan utils2.ScanTask)
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
				if defaultuser, ok := utils2.DefaultUserDict[info.Server]; ok {
					UserList = defaultuser
				} else {
					UserList = []string{"admin"}
				}
			}

			if PassIsDefault {
				PassList = utils2.DefaultPasswords[info.Server]
			}

			for _, username := range UserList {
				for _, password := range PassList {
					NewTask := utils2.ScanTask{
						TargetInfo: utils2.TargetInfo{
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
				username, password, err := utils2.GetUP(up)
				if err != nil {
					continue
				}
				NewTask := utils2.ScanTask{
					TargetInfo: utils2.TargetInfo{
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

func GenerateRandTask(ipinfo []utils2.IpServerInfo) chan utils2.ScanTask {
	TaskList := make(chan utils2.ScanTask)
	go func() {
		for _, info := range ipinfo {
			NewTask := utils2.ScanTask{
				TargetInfo: utils2.TargetInfo{
					IpServerInfo: info,
					Username:     "admin",
					Password:     utils2.RandStringBytesMaskImprSrcUnsafe(12),
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

func GenerateTaskSimple(UserList []string, PassList []string, ipinfo []utils2.IpServerInfo) chan utils2.ScanTask {

	TaskList := make(chan utils2.ScanTask)
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
				if defaultuser, ok := utils2.DefaultUserDict[info.Server]; ok {
					UserList = defaultuser
				} else {
					UserList = []string{"admin"}
				}
			}

			if PassIsDefault {
				PassList = utils2.DefaultPasswords[info.Server]
			}

			Summary = len(PassList) * len(UserList) * len(ipinfo)

			if len(UPList) == 0 {
				//PassList = append(PassList, Utils.RandStringBytesMaskImprSrcUnsafe(8))
				for _, username := range UserList {
					for _, password := range PassList {
						NewTask := utils2.ScanTask{
							TargetInfo: utils2.TargetInfo{
								IpServerInfo: info,
								Username:     username,
								Password:     password,
							},
						}
						if info.Server == "ORACLE" {
							for _, ins := range utils2.Instance {
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
					username, password, err := utils2.GetUP(up)
					if err != nil {
						continue
					}
					NewTask := utils2.ScanTask{
						TargetInfo: utils2.TargetInfo{
							IpServerInfo: info,
							Username:     username,
							Password:     password,
						},
					}
					if info.Server == "ORACLE" {
						for _, ins := range utils2.Instance {
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
