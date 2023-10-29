package lib

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/md4"
	"log"
	"net"
	"strings"
	"time"
)

func VersionAndLib(ip string, port string) (string, []string) {
	s := "@RSYNCD: 31.0"
	conn, err := net.DialTimeout("tcp", ip+":"+port, 8*time.Second)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Write(strhex(s))
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

		//fmt.Printf("%v %v",len1,string(rev1))
		modulename := strings.Split(strings.Replace(Libs, " ", "", len(Libs)), "\n")
		//fmt.Println(modulename)
		for _, v := range modulename {
			realname := strings.Split(v, "\t")
			if realname[0] != "" && strings.Contains(realname[0], "@RSYNCD:EXIT") == false {
				Lib[i] = realname[0]
				i++
			} else if strings.Contains(realname[0], "@RSYNCD:EXIT") {
				break
			}
		}

		break

	}

	return version, Lib
}

func HighVersion(ip, port, user, passwd string, mod string) error {

	s := strhex("@RSYNCD: 31.0")

	conn, err := net.DialTimeout("tcp", ip+":"+port, 8*time.Second)
	defer conn.Close()

	if err != nil {
		return err
	}
	//1.发送版本信息
	_, err = conn.Write(s)
	if err != nil {
		return err
	}
	var rev = make([]byte, 1024)
	_, err = conn.Read(rev)
	if err != nil {
		return err
	}

	//	fmt.Printf("rev: %s\nres: %s\n",rev,mod)
	module := mod + "\n"

	_, err = conn.Write([]byte(module))

	var rev2 = make([]byte, 1024)
	_, err = conn.Read(rev2)
	if err != nil {
		return err
	}

	//get challenge code
	challenge := strings.Split(string(rev2), " ")
	c := challenge[len(challenge)-1]
	//	fmt.Printf("challenge: %s\n",c)

	//截取字符串
	//c1 := strings.Replace(passwd+c, "\n", "", -1)[:26]
	c1 := strings.Split(passwd+c, "\n")
	c2 := c1[0]

	//md5加密
	md := md5.New()
	md.Write([]byte(c2))
	//md5校验和获取
	str := md.Sum(nil)

	auth_send_data := base64.StdEncoding.EncodeToString(str)
	a := strings.Replace(auth_send_data, "==", "", len(auth_send_data))
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
	//判断是否爆破成功
	if strings.Contains(string(rev3), "OK") {
		return nil
	}
	return errors.New("connect error")
}

func LowVersion(ip, port, user, passwd string, mod string) error {

	s := strhex("@RSYNCD: 31.0")

	conn, err := net.DialTimeout("tcp", ip+":"+port, 8*time.Second)
	defer conn.Close()

	if err != nil {
		return err
	}
	//1.发送版本信息
	_, err = conn.Write(s)
	if err != nil {
		return err
	}
	var rev = make([]byte, 1024)
	_, err = conn.Read(rev)
	if err != nil {
		return err
	}

	//	fmt.Printf("rev: %s\nres: %s\n",rev,mod)
	module := mod + "\n"

	_, err = conn.Write([]byte(module))

	var rev2 = make([]byte, 1024)
	_, err = conn.Read(rev2)
	if err != nil {
		return err
	}

	//get challenge code
	challenge := strings.Split(string(rev2), " ")
	c := challenge[len(challenge)-1]
	//	fmt.Printf("challenge: %s\n",c)

	//截取字符串
	//c1 := strings.Replace(passwd+c, "\n", "", -1)[:26]
	c1 := strings.Split(passwd+c, "\n")
	c2 := c1[0]

	md := md4.New()
	md.Write([]byte(c2))
	str := md.Sum(nil)

	auth_send_data := base64.StdEncoding.EncodeToString(str)
	a := strings.Replace(auth_send_data, "==", "", len(auth_send_data))
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
	//判断是否爆破成功
	if strings.Contains(string(rev3), "OK") {
		return nil
	}
	return errors.New("connect error")
}

func strhex(str string) []byte {
	hex_s := hex.EncodeToString([]byte(str + "\n"))

	dst := make([]byte, hex.DecodedLen(len(hex_s)))
	n, err := hex.Decode(dst, []byte(hex_s))

	if err != nil {
		log.Fatal(err)
		return nil
	}
	return dst[:n]

}
