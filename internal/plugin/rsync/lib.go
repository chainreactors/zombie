package rsync

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/chainreactors/utils/encode"
	"golang.org/x/crypto/md4"
	"net"
	"strings"
	"time"
)

func RsyncDetect(ip string, port string) (string, []string) {
	s := "@RSYNCD: 31.0"
	conn, err := net.DialTimeout("tcp", ip+":"+port, 8*time.Second)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Write([]byte(s + "\n"))
	if err != nil {
		fmt.Println(err)
	}

	var rev = make([]byte, 1024)
	_, err = conn.Read(rev)
	if err != nil {
		fmt.Println(err)
	}

	version := strings.TrimSpace(string(rev))

	s = "\n"
	_, err = conn.Write([]byte(s))

	var Lib = make([]string, 10)
	i := 0

	for true {
		var rev1 = make([]byte, 1024)
		_, err = conn.Read(rev1)
		if err != nil {
			fmt.Println(err)
		}

		Libs := strings.TrimSpace(string(rev1))

		ModuleName := strings.Split(strings.Replace(Libs, " ", "", len(Libs)), "\n")
		for _, v := range ModuleName {
			RealName := strings.Split(v, "\t")
			if RealName[0] != "" && strings.Contains(RealName[0], "@RSYNCD:EXIT") == false {
				Lib[i] = RealName[0]
				i++
			} else if strings.Contains(RealName[0], "@RSYNCD:EXIT") {
				break
			}
		}

		break

	}

	return version, Lib
}

func RsyncLogin(ip, port, user, passwd string, mod string, SmallVersion float64) error {
	s := []byte("@RSYNCD: 31." + "\n")

	conn, err := net.DialTimeout("tcp", ip+":"+port, 8*time.Second)
	defer conn.Close()

	if err != nil {
		return err
	}
	_, err = conn.Write(s)

	if err != nil {
		return err
	}
	var rev = make([]byte, 1024)
	_, err = conn.Read(rev)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(mod + "\n"))

	var rev2 = make([]byte, 1024)
	_, err = conn.Read(rev2)
	if err != nil {
		return err
	}

	challenge := strings.Split(string(rev2), " ")
	c := challenge[len(challenge)-1]
	c1 := strings.Split(passwd+c, "\n")

	var hash []byte
	if SmallVersion >= 30 {
		md := md5.New()
		md.Write([]byte(c1[0]))
		hash = md.Sum(nil)
	} else {
		md := md4.New()
		md.Write([]byte(c1[0]))
		hash = md.Sum(nil)
	}

	AutoData := encode.Base64Encode(hash)
	a := strings.Replace(AutoData, "==", "", len(AutoData))

	_, err = conn.Write([]byte(user + " " + a + "\n"))
	if err != nil {
		return err
	}
	var rev3 = make([]byte, 1024)
	_, err = conn.Read(rev3)
	if err != nil {
		return err
	}
	if strings.Contains(string(rev3), "OK") {
		return nil
	}
	return errors.New("rsync connect error")
}
