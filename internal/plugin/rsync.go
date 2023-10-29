package plugin

import (
	"github.com/chainreactors/zombie/internal/plugin/lib"
	"github.com/chainreactors/zombie/pkg"
	"strconv"
	"strings"
)

type RsyncService struct {
	*pkg.Task
}

func (s *RsyncService) Query() bool {
	return false
}

func (s *RsyncService) GetInfo() bool {
	return false
}

func (s *RsyncService) Connect() error {
	res, Libs := lib.VersionAndLib(s.IP, s.Port)
	version := strings.Split(res, "\n")[0]
	small_version, _ := strconv.ParseFloat(strings.Split(version, " ")[1], 64)
	if small_version > 30 {
		err := lib.HighVersion(s.IP, s.Port, s.Username, s.Password, Libs[0])
		if err != nil {
			return err
		}
	} else {
		err := lib.LowVersion(s.IP, s.Port, s.Username, s.Password, Libs[0])
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *RsyncService) Close() error {
	return NilConnError{s.Service}
}

func (s *RsyncService) SetQuery(query string) {
}

func (s *RsyncService) Output(res interface{}) {

}
