package Utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

var FileHandle *os.File
var O2File bool
var BDatach = make(chan OutputRes, 1000)
var QDatach = make(chan string, 1000)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func BruteWrite2File(FileHandle *os.File, Datach chan OutputRes) {

	switch FileFormat {
	case "raw":
		for res := range Datach {
			FileHandle.WriteString(fmt.Sprintf("%s:%d\t\tusername:%s\tpassword:%s\t%s\tsuccess\t%s\n", res.IP, res.Port, res.Username, res.Password, res.Type, res.Additional))
		}
	case "json":
		FileHandle.WriteString("{")
		for res := range Datach {
			jsons, errs := json.Marshal(res)
			if errs != nil {
				fmt.Println(errs.Error())
			}
			FileHandle.WriteString(string(jsons) + ",")
		}
	}

}

func QueryWrite2File(FileHandle *os.File, QDatach chan string) {

	for res := range QDatach {
		FileHandle.WriteString(res + "\n")
		fmt.Println(res)
	}
}

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
