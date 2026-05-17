package core

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"testing"
)

func TestTarget_CheckFinger(t1 *testing.T) {
	err := pkg.LoadFingers()
	if err != nil {
		t1.Log(err)
		return
	}
	target := &Target{
		IP:      "127.0.0.1",
		Port:    "6379",
		Service: "redis",
	}
	if open := target.CheckOpen(); open {
		fmt.Println("Open")
	}

	if matched := target.CheckFinger(); matched {
		fmt.Println("Matched")
	}
}
