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
	var Curtask Utils.ScanTask

	if ctx.IsSet("file") {

	} else {
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
				fmt.Println("the Database isn't be supported")
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
				fmt.Println("Please input the type of Database")
				os.Exit(0)
			}
		} else {
			fmt.Println("Please input the type of Database")
			os.Exit(0)
		}

		CurServer = strings.ToUpper(CurServer)

		IpList := Core.GetIpInfoList(IpSlice, CurServer)

		Curtask = Utils.ScanTask{
			Info:     IpList[0],
			Username: ctx.String("username"),
			Password: ctx.String("password"),
			Server:   CurServer,
		}
	}

	CurCon := Core.ExecDispatch(Curtask)

	alive := CurCon.Connect()

	if !alive {
		fmt.Printf("can't connect to db")
		os.Exit(0)
	}

	IsAuto := ctx.Bool("auto")

	if IsAuto {
		CurCon.GetInfo()
	} else {
		CurCon.SetQuery(ctx.String("input"))
		CurCon.Query()
	}

	return err
}
