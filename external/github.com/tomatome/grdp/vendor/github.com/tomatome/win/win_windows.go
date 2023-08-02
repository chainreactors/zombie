package win

import (
	"fmt"
	"math"
	"syscall"
)

func doLoadLibrary(name string) uintptr {
	lib, _ := syscall.LoadLibrary(name)
	return uintptr(lib)
}

func doGetProcAddress(lib uintptr, name string) uintptr {
	addr, _ := syscall.GetProcAddress(syscall.Handle(lib), name)
	return uintptr(addr)
}

func syscall3(trap, nargs, a1, a2, a3 uintptr) uintptr {
	ret, _, _ := syscall.Syscall(trap, nargs, a1, a2, a3)
	return ret
}

func syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(trap, nargs, a1, a2, a3, a4, a5, a6)
	return ret
}

func syscall9(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) uintptr {
	ret, _, _ := syscall.Syscall9(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	return ret
}

func syscall12(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12 uintptr) uintptr {
	ret, _, _ := syscall.Syscall12(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
	return ret
}

func syscall15(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) uintptr {
	ret, _, _ := syscall.Syscall15(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
	return ret
}

func syscallN(trap uintptr, args []uintptr) uintptr {
	n := len(args)
	num := int(math.Max(1.0, math.Ceil(float64(n)/3.0))) * 3
	for i := n; i < num; i++ {
		args = append(args, 0)
	}
	if num == 3 {
		ret, _, _ := syscall.Syscall(trap, uintptr(n),
			args[0], args[1], args[2])
		return ret
	} else if num == 6 {
		ret, _, _ := syscall.Syscall6(trap, uintptr(n),
			args[0], args[1], args[2],
			args[3], args[4], args[5])
		return ret
	} else if num == 9 {
		ret, _, _ := syscall.Syscall9(trap, uintptr(n),
			args[0], args[1], args[2],
			args[3], args[4], args[5],
			args[6], args[7], args[8])
		return ret
	} else if num == 12 {
		ret, _, _ := syscall.Syscall12(trap, uintptr(n),
			args[0], args[1], args[2],
			args[3], args[4], args[5],
			args[6], args[7], args[8],
			args[9], args[10], args[11])
		return ret
	} else if num == 15 {
		ret, _, _ := syscall.Syscall15(trap, uintptr(n),
			args[0], args[1], args[2],
			args[3], args[4], args[5],
			args[6], args[7], args[8],
			args[9], args[10], args[11],
			args[12], args[13], args[14])
		return ret
	} else {
		panic(fmt.Errorf("Too many syscall arguments: %d", n))
	}
}

func getUintptrFromBool(b bool) uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}
