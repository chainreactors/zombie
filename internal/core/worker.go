package core

import (
	"errors"
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg"
)

var ErrNoUnauth = errors.New("cannot unauth login")

func Unauth(task *pkg.Task) *pkg.Result {
	conn, err := plugin.Dispatch(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	ok, err := conn.Unauth()
	if err != nil {
		return pkg.NewResult(task, err)
	}
	if !ok {
		return pkg.NewResult(task, ErrNoUnauth)
	}
	return pkg.NewResult(task, nil)
}

func Brute(task *pkg.Task) *pkg.Result {
	conn, err := plugin.Dispatch(task)
	if err != nil {
		return pkg.NewResult(task, err)
	}
	err = conn.Login()
	if err != nil {
		return pkg.NewResult(task, err)
	}
	defer conn.Close()

	return pkg.NewResult(task, nil)
}
