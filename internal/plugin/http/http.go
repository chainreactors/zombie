package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/http"
)

type nilLog struct {
}

func (l nilLog) Print(v ...interface{}) {

}

type HttpPlugin struct {
	*pkg.Task
	HttpInf
}

func (s *HttpPlugin) Unauth() (bool, error) {
	//TODO implement me
	panic("implement me")
}

type HttpInf struct {
	Path string `json:"path"`
}

func (s *HttpPlugin) Name() string {
	return s.Service.String()
}

func (s *HttpPlugin) Login() error {
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

//func (s *HttpPlugin) Unauth() (bool, error) {
//	url := fmt.Sprintf("%s://%s:%s%s", s.Service, s.IP, s.Port, s.Path)
//	client := &http.Client{}
//	req, _ := http.NewRequest("GET", url, nil)
//
//	UsernameAndPassword := fmt.Sprintf("%s:%s", "", "")
//	base64_str := base64.StdEncoding.EncodeToString([]byte(UsernameAndPassword))
//	Authorization := fmt.Sprintf("Basic %s", base64_str)
//	req.Header.Add("Authorization", Authorization)
//	resp, err := client.Do(req)
//	if err != nil {
//		return false, errors.New("timeout")
//	}
//	if resp.StatusCode != 401 {
//		return true, nil
//	} else {
//		result := fmt.Sprintf("StatusCode -> %d", resp.StatusCode)
//		return false, errors.New(result)
//	}
//}

func (s *HttpPlugin) GetBasic() *pkg.Basic {
	// todo list dbs
	return &pkg.Basic{}
}

func (s *HttpPlugin) Close() error {
	return pkg.NilConnError{s.Service}
}
