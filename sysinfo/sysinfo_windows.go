package sysinfo

import (
	"fmt"
	"syscall"

	"runtime"
	"unsafe"

	"github.com/blang/semver"
)

func GetSemanticVersion() (semver.Version, error) {
	str, err := GetString()
	if err != nil {
		return semver.Version{}, err
	}

	return semver.Make(str)
}

//TODO: windows 10 之后微软修改了函数，返回值错误。
// http://chenjava.blog.51cto.com/374566/1628084
func GetString() (string, error) {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return "", fmt.Errorf("Error loading kernel32.dll: %v", err)
	}
	p, err := dll.FindProc("GetVersion")
	if err != nil {
		return "", fmt.Errorf("Error finding GetVersion procedure: %v", err)
	}
	// The error is always non-nil
	v, _, _ := p.Call()
	return fmt.Sprintf("%d.%d.%d", byte(v), byte(v >> 8), uint16(v >> 16)), nil
}

func Is64Sys() (bool, error) {
	if runtime.GOARCH == `amd64` {
		return true, nil
	} else {
		var mod = syscall.NewLazyDLL("kernel32.dll")
		if mod == nil {
			return false, fmt.Errorf("载入 kernel32.dll 失败。")
		}

		var proc = mod.NewProc("IsWow64Process")
		if proc == nil {
			return false, nil
		}

		is64 := 0
		h, err := syscall.GetCurrentProcess()
		if err != nil {
			panic(err)
		}
		rl, _, _ := proc.Call(uintptr(h), uintptr(unsafe.Pointer(&is64)))
		if rl != 1 {
			return false, fmt.Errorf("IsWow64Process 调用失败。")
		}

		if is64 == 0 {
			return false, nil
		} else {
			return true, nil
		}

	}
}
