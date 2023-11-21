package rsync

import (
	"github.com/chainreactors/zombie/pkg"
	"strconv"
	"strings"
)

type RsyncPlugin struct {
	*pkg.Task
}

func (s *RsyncPlugin) Unauth() (bool, error) {
	//TODO implement me
	panic("implement me")
}

//func (s *RsyncPlugin) Query() bool {
//	return false
//}
//
//func (s *RsyncPlugin) GetInfo() bool {
//	return false
//}

func (s *RsyncPlugin) Login() error {
	res, Libs := RsyncDetect(s.IP, s.Port)
	version := strings.Split(res, "\n")[0]
	SmallVersion, _ := strconv.ParseFloat(strings.Split(version, " ")[1], 64)
	err := RsyncLogin(s.IP, s.Port, s.Username, s.Password, Libs[0], SmallVersion)
	if err != nil {
		return err
	}

	return nil

}

func (s *RsyncPlugin) Name() string {
	return s.Service.String()
}

func (s *RsyncPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

func (s *RsyncPlugin) Close() error {
	return pkg.NilConnError{s.Service}
}

//func (s *RsyncPlugin) SetQuery(query string) {
//}
//
//func (s *RsyncPlugin) Output(res interface{}) {
//
//}
