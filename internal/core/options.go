package core

import (
	"github.com/chainreactors/ipcs"
)

type Option struct {
	IP            string `short:"i"`
	IPFile        string `short:"I"`
	Username      string `short:"u" long:"user"`
	UsernameFile  string `short:"U"`
	Password      string `short:"p" long:"pass"`
	PasswordFile  string `short:"P"`
	GogoFile      string `long:"go"`
	ServiceName   string `short:"s" long:"service"`
	FilterService string `short:"S" long:"filter-service"`
	OutputFile    string `short:"f" long:"file"`
	OutputType    string `short:"o" long:"output"`
	Threads       int    `short:"t" default:"100"`
	Timeout       int    `short:"d" long:"timeout" default:"2"`
	Mod           string `short:"m" default:"clusterbomb"`
	Debug         bool   `long:"debug"`
}

func (o Option) Validate() error {
	//if _, ok := utils.ServicePortMap[strings.ToUpper(o.ServiceName)]; !ok {
	//	return fmt.Errorf("not support %s plugin", o.ServiceName)
	//}
	return nil
}

type Target struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Service string `json:"service"`
}

func (t Target) Addr() *ipcs.Addr {
	return &ipcs.Addr{IP: ipcs.NewIP(t.IP), Port: t.Port}
}
