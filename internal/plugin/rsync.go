package plugin

import (
	"crypto/md5"
	"errors"
	"fmt"
	parse "github.com/chainreactors/parsers"
	"github.com/chainreactors/zombie/pkg"
	"golang.org/x/crypto/md4"
	"net"
	"strconv"
	"strings"
	"time"
)

type RsyncService struct {
	*pkg.Task
}

func (s *RsyncService) Query() bool {
	return false
}

func (s *RsyncService) GetInfo() bool {
	return false
}

func (s *RsyncService) Connect() error {
	res, Libs := RsyncDetect(s.IP, s.Port)
	version := strings.Split(res, "\n")[0]
	SmallVersion, _ := strconv.ParseFloat(strings.Split(version, " ")[1], 64)
	err := RsyncLogin(s.IP, s.Port, s.Username, s.Password, Libs[0], SmallVersion)
	if err != nil {
		return err
	}

	return nil

}

func (s *RsyncService) Close() error {
	return pkg.NilConnError{s.Service}
}

func (s *RsyncService) SetQuery(query string) {
}

func (s *RsyncService) Output(res interface{}) {

}

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

	module := mod + "\n"

	_, err = conn.Write([]byte(module))

	var rev2 = make([]byte, 1024)
	_, err = conn.Read(rev2)
	if err != nil {
		return err
	}

	challenge := strings.Split(string(rev2), " ")
	c := challenge[len(challenge)-1]
	c1 := strings.Split(passwd+c, "\n")
	c2 := c1[0]

	var str []byte
	if SmallVersion >= 30 {
		md := md5.New()
		md.Write([]byte(c2))
		str = md.Sum(nil)
	} else {
		md := md4.New()
		md.Write([]byte(c2))
		str = md.Sum(nil)
	}

	AutoData := parse.Base64Encode(str)
	a := strings.Replace(AutoData, "==", "", len(AutoData))
	payload := user + " " + a + "\n"

	_, err = conn.Write([]byte(payload))
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
	return errors.New("connect error")
}
