package pkg

import (
	"os"
)

//func UpdatePass(CurTask ScanTask) ScanTask {
//	if strings.Contains(CurTask.Password, "%user%") {
//		CurTask.Password = strings.Replace(CurTask.Password, "%user%", CurTask.Username, 1)
//	}
//
//	return CurTask
//}

//func GetUP(up string) (user, pass string, err error) {
//	ualist := strings.Split(up, " ")
//
//	if len(ualist) == 1 {
//		return ualist[0], "", nil
//	} else if len(ualist) == 2 {
//		return ualist[0], ualist[1], nil
//	}
//
//	return "", "", fmt.Errorf("Something error!")
//}
//
//func SliceContains(s []string, e string) bool {
//	for _, v := range s {
//		if v == e {
//			return true
//		}
//	}
//	return false
//}
//
//func SliceLike(s []string, e string) bool {
//	for _, v := range s {
//		e = strings.ToUpper(e)
//		if strings.Contains(e, v) {
//			return true
//		}
//	}
//	return false
//}

func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	isPipedFromChrDev := (stat.Mode() & os.ModeCharDevice) == 0
	isPipedFromFIFO := (stat.Mode() & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}
