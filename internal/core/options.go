package core

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg/utils"
	"strings"
)

type Option struct {
	IP           string `short:"i"`
	IPFile       string `short:"I"`
	Username     string `short:"u" long:"user"`
	UsernameFile string `short:"U"`
	Password     string `short:"p" long:"pass"`
	PasswordFile string `short:"P"`
	GogoFile     string `long:"go"`
	ServiceName  string `short:"s"`
	OutputFile   string `short:"f"`
	OutputType   string `short:"o"`
	Threads      int    `short:"t" default:"100"`
	Timeout      int    `short:"d" long:"timeout" default:"2"`
	Mod          string `short:"m" default:"clusterbomb"`
	Debug        bool   `long:"debug"`
}

func (o Option) Validate() error {
	if _, ok := utils.ServicePortMap[strings.ToUpper(o.ServiceName)]; !ok {
		return fmt.Errorf("not support %s plugin", o.ServiceName)
	}
	return nil
}
