package Core

import (
	"Zombie/src/Utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetUserList(username string) (UserList []string) {
	if strings.Contains(username, ",") {
		userslice := strings.Split(username, ",")
		for _, user := range userslice {
			UserList = append(UserList, user)
		}
	} else {
		UserList = append(UserList, username)
	}
	return UserList
}

func GetIPList(ipname string) (IPList []string) {
	suffix := ""

	if strings.Contains(ipname, ",") {
		ipslice := strings.Split(ipname, ",")
		for _, ip := range ipslice {

			if strings.Contains(ip, ":") {
				ipsuffix := strings.Split(ip, ":")
				suffix = ":" + ipsuffix[1]
				ip = ipsuffix[0]
			}

			if strings.Contains(ip, "/") {
				start, fin := Utils.GetIpRange(ip)
				for i := start; i <= fin; i++ {
					// 如果是广播地址或网络地址,则跳过
					if (i)%256 != 255 && (i)%256 != 0 {
						IPList = append(IPList, Utils.Int2ip(i)+suffix)
					}
				}
			} else {
				IPList = append(IPList, ip+suffix)
			}
		}
	} else {

		if strings.Contains(ipname, ":") {
			ipsuffix := strings.Split(ipname, ":")
			suffix = ":" + ipsuffix[1]
			ipname = ipsuffix[0]
		}

		if strings.Contains(ipname, "/") {
			start, fin := Utils.GetIpRange(ipname)
			for i := start; i <= fin; i++ {
				// 如果是广播地址或网络地址,则跳过
				if (i)%256 != 255 && (i)%256 != 0 {
					IPList = append(IPList, Utils.Int2ip(i)+suffix)
				}
			}
		} else {
			IPList = append(IPList, ipname+suffix)
		}
	}
	return IPList
}

func GetUAList(UAFile string) (UserPass []string, err error) {
	file, err := os.Open(UAFile)
	if err != nil {
		panic("please check your file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		up := strings.TrimSpace(scanner.Text())
		if up != "" {
			UserPass = append(UserPass, up)
		}
	}
	return UserPass, err
}

func GetPassList(password string) (PassList []string) {
	PassList = append(PassList, "")
	if strings.Contains(password, ",") {
		passslice := strings.Split(password, ",")
		for _, pass := range passslice {
			PassList = append(PassList, pass)
		}
	} else {
		PassList = append(PassList, password)
	}
	return PassList
}

func GetIpList(ipstring string) (IpList []string) {
	if strings.Contains(ipstring, ",") {
		ipslice := strings.Split(ipstring, ",")
		for _, ip := range ipslice {
			IpList = append(IpList, ip)
		}
	} else {
		IpList = append(IpList, ipstring)
	}
	return IpList
}

func GetIpInfoList(iplist []string, Server string) (IpInfoList []Utils.IpInfo) {
	for _, ip := range iplist {
		target := Utils.IpInfo{}

		if strings.Contains(ip, ":") {
			SplitIp := strings.Split(ip, ":")
			port, err := strconv.Atoi(SplitIp[1])
			if err != nil {
				fmt.Println("Please check your address")
				os.Exit(0)
			}
			target.Port = port
			target.Ip = SplitIp[0]
			IpInfoList = append(IpInfoList, target)
		} else {
			target := Utils.IpInfo{
				Ip:   ip,
				Port: Utils.ServerPort[Server],
			}
			IpInfoList = append(IpInfoList, target)
		}
	}
	return IpInfoList
}

func ReadUserDict(userDict string) (UserList []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		panic("please check your file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			UserList = append(UserList, user)
		}
	}
	return UserList, err
}

func ReadIPDict(IpDict string) (IPList []string, err error) {
	file, err := os.Open(IpDict)
	if err != nil {
		panic("please check your file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			IPList = append(IPList, user)
		}
	}
	return IPList, err
}

func ReadPassDict(passDict string) (password []string, err error) {
	file, err := os.Open(passDict)
	if err != nil {
		panic("please check your file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			password = append(password, passwd)
		}
	}
	password = append(password, "")
	return password, err
}
