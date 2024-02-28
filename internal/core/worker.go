package core

import (
	"errors"
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg"
)

var ErrNoUnauth = errors.New("cannot unauth login")

func Unauth(task *pkg.Task) *pkg.Result {
	conn := plugin.Dispatch(task)
	ok, err := conn.Unauth()
	if err != nil {
		return pkg.NewResult(task, err)
	}
	if !ok {
		return pkg.NewResult(task, ErrNoUnauth)
	}
	return conn.GetResult()
}

func Brute(task *pkg.Task) *pkg.Result {
	conn := plugin.Dispatch(task)
	err := conn.Login()
	if err != nil {
		return pkg.NewResult(task, err)
	}
	defer conn.Close()
	return conn.GetResult()
}
