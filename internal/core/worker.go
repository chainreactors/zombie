package core

import (
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg/utils"
)

func Brute(task *utils.Task) *utils.Result {
	conn := plugin.Dispatch(task)

	result := &utils.Result{
		Task: task,
	}
	err := conn.Connect()
	if err != nil {
		result.Err = err
		return result
	}
	defer conn.Close()
	result.OK = true
	return result
}
