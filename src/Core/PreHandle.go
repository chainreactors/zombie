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
		target := Utils.IpInfo{
			SSL: false,
		}

		if strings.HasPrefix(ip, "http") {
			ips := strings.Split(ip, "http://")
			ip = ips[1]
		} else if strings.HasPrefix(ip, "https") {
			ips := strings.Split(ip, "http://")
			ip = ips[1]
			target.SSL = true
		}

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
