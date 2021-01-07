package Core

import (
	"Zombie/src/Utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetUserList(username string)(UserList []string){
	if strings.Contains(username,","){
		userslice := strings.Split(username,",")
		for _, user := range userslice{
			UserList = append(UserList, user)
		}
	}else {
		UserList = append(UserList, username)
	}
	return UserList
}

func GetPassList(password string)(PassList []string){
	PassList = append(PassList, "")
	if strings.Contains(password,","){
		passslice := strings.Split(password,",")
		for _, pass := range passslice{
			PassList = append(PassList, pass)
		}
	}else {
		PassList = append(PassList, password)
	}
	return PassList
}

func GetIpList(ipstring string)(IpList []string){
	if strings.Contains(ipstring,","){
		ipslice := strings.Split(ipstring,",")
		for _, ip := range ipslice{
			IpList = append(IpList, ip)
		}
	}else {
		IpList = append(IpList, ipstring)
	}
	return IpList
}

func GetIpInfoList(iplist []string,Server string)(IpInfoList []Utils.IpInfo){
	for _, ip := range iplist{
		if strings.Contains(ip, ":"){
			SplitIp := strings.Split(ip,":")
			port, err := strconv.Atoi(SplitIp[1])
			if err != nil {
				fmt.Println("Please check your address")
				os.Exit(0)
			}
			target := Utils.IpInfo{
				Ip:   SplitIp[0],
				Port: port,
			}
			IpInfoList = append(IpInfoList, target)
		}else{
			target := Utils.IpInfo{
				Ip:   ip,
				Port: Utils.ServerPort[Server],
			}
			IpInfoList = append(IpInfoList, target)
		}
	}
	return IpInfoList
}