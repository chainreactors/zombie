package main

import (
	"fmt"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/neutron/common"
	"github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/utils/iutils"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

type Option struct {
	IP       string `short:"i" long:"ip" alias:"ipp" description:"String, input ip"`
	Template string `short:"t" long:"template" description:"File, input template"`
}

func init() {
	common.NeutronLog = logs.Log
	logs.Log.SetLevel(logs.Debug)
}

func main() {
	var opt Option
	parser := flags.NewParser(&opt, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Println(err.Error())
		}
		return
	}
	if opt.IP == "" {
		fmt.Println("please input ip")
		return
	}

	var template *templates.Template
	if opt.Template != "" {
		content, err := ioutil.ReadFile(opt.Template)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(content, &template)
		if err != nil {
			panic(err)
		}
		err = template.Compile(nil)
		if err != nil {
			panic(err)
		}
		logs.Log.Importantf("load template %s:%s success", opt.Template, template.Id)
	} else {
		panic("please choice template file")
	}
	logs.Log.Importantf("start execute %s", template.Id)
	res, err := template.Execute(opt.IP, nil)
	if err != nil {
		panic(err)
	}
	logs.Log.Importantf("execute %s finished", template.Id)
	if res == nil {
		logs.Log.Warn("no result")
		return
	}
	logs.Log.Infof(opt.IP, template.Id, iutils.ToString(res.PayloadValues["username"]), iutils.ToString(res.PayloadValues["password"]))
}
