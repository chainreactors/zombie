package core

import (
	"fmt"
	"github.com/chainreactors/zombie/internal/plugin"
	"github.com/chainreactors/zombie/pkg"
)

func Brute(task *pkg.Task) error {
	conn := plugin.Dispatch(task)
	if conn == nil {
		return fmt.Errorf("not support service " + task.Service)
	}
	err := conn.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
