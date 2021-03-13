package Moudle

import (
	"Zombie/src/Core"
	"Zombie/src/Utils"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
)

func Exec(ctx *cli.Context) (err error) {
	var CurServer string

	if !ctx.IsSet("ip") {
		fmt.Println("Please check the address")
		os.Exit(0)
	}
	if strings.Contains(ctx.String("ip"), ",") {
		fmt.Println("Exec Moudle only support single ip")
		os.Exit(0)
	}

	IpSlice := Core.GetIpList(ctx.String("ip"))

	Ip := IpSlice[0]
	if ctx.IsSet("server") {
		ServerName := strings.ToUpper(ctx.String("server"))
		if _, ok := Utils.ExecPort[ServerName]; ok {
			CurServer = ctx.String("server")
		} else {
			fmt.Println("the Server isn't be supported")
			os.Exit(0)
		}

	} else if strings.Contains(Ip, ":") {
		Temp := strings.Split(Ip, ":")
		Sport := Temp[1]
		port, err := strconv.Atoi(Sport)
		if err != nil {
			fmt.Println("Please check your address")
			os.Exit(0)
		}

		if _, ok := Utils.ExecServer[port]; ok {
			CurServer = Utils.PortServer[port]
			fmt.Println("Use default server")
		} else {
			fmt.Println("Please input the type of Server")
			os.Exit(0)
		}
	} else {
		fmt.Println("Please input the type of Server")
		os.Exit(0)
	}

	if !ctx.IsSet("username") && !ctx.IsSet("password") {
		fmt.Println("please input username and password, if the server don't need username,input anything to place")
		os.Exit(0)
	}

	if !ctx.IsSet("query") {
		fmt.Println("please input your query")
		os.Exit(0)
	}

	CurServer = strings.ToUpper(CurServer)

	IpList := Core.GetIpInfoList(IpSlice, CurServer)

	Curtask := Utils.ScanTask{
		Info:     IpList[0],
		Username: ctx.String("username"),
		Password: ctx.String("password"),
		Server:   CurServer,
	}

	err, Qresult, Columns := Core.ExecDispatch(Curtask, ctx.String("query"))

	if err != nil {
		fmt.Println("something wrong")
		os.Exit(0)
	} else {
		for _, cname := range Columns {
			fmt.Print(cname + "\t")
		}
		fmt.Print("\n")
		for _, items := range Qresult {
			for _, cname := range Columns {
				fmt.Print(items[cname] + "\t")
			}
			fmt.Print("\n")
		}
	}

	return err
}
