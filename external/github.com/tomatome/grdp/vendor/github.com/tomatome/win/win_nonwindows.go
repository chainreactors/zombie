// +build !windows

package win

import (
	"syscall"
)

func syscall3(trap, nargs, a1, a2, a3 uintptr) uintptr {
	ret, _, _ := syscall.Syscall(trap, a1, a2, a3)
	return ret
}
