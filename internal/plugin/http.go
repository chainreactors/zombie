package plugin

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
)

type Http struct {
	*pkg.Task
	HttpInf
}

type HttpInf struct {
	Path string `json:"path"`
}

func (s *Http) Query() bool {
	return false
}

func (s *Http) GetInfo() bool {
	return false
}

func (s *Http) Connect() error {

	url := fmt.Sprintf("%s://%s:%s%s", s.Service, s.IP, s.Port, s.Path)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	UsernameAndPassword := fmt.Sprintf("%s:%s", s.Username, s.Password)
	base64_str := base64.StdEncoding.EncodeToString([]byte(UsernameAndPassword))
	Authorization := fmt.Sprintf("Basic %s", base64_str)
	req.Header.Add("Authorization", Authorization)
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("timeout")
	}
	if resp.StatusCode != 401 {
		return nil
	} else {
		result := fmt.Sprintf("StatusCode -> %d", resp.StatusCode)
		return errors.New(result)
	}
}

func (s *Http) Close() error {
	return NilConnError{s.Service}
}

func (s *Http) SetQuery(query string) {
	//s.Input = query
}

func (s *Http) Output(res interface{}) {

}
