package Utils

import (
	"fmt"
	"strings"
)

func UpdatePass(CurTask ScanTask) ScanTask {
	if strings.Contains(CurTask.Password, "%user%") {
		CurTask.Password = strings.Replace(CurTask.Password, "%user%", CurTask.Username, 1)
	}

	return CurTask
}

func GetUP(up string) (user, pass string, err error) {
	ualist := strings.Split(up, " ")

	if len(ualist) == 1 {
		return ualist[0], "", nil
	} else if len(ualist) == 2 {
		return ualist[0], ualist[1], nil
	}

	return "", "", fmt.Errorf("Something error!")
}