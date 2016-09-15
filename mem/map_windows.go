// +build windows

package mem

import (
	"fmt"
	"syscall"
	"unsafe"
)

const NAME = "CAC5ADA322594E84B6C1120F0B52FD7B"

func ReadMem() ([]byte, error) {
	return ReadNameMem("")
}
func ReadNameMem(name string) ([]byte, error) {
	fh := int32(-1)

	nname := NAME
	if name != "" {
		nname = NAME + "_" + name
	}

	h, err := syscall.CreateFileMapping(syscall.Handle(fh),
		nil, uint32(syscall.PAGE_READONLY), 0, uint32(1024), syscall.StringToUTF16Ptr(nname))
	if h == 0 {
		return nil, err
	}

	addr, err := syscall.MapViewOfFile(h, uint32(syscall.FILE_MAP_READ), 0,
		0, uintptr(1024))
	if addr == 0 {
		return nil, err
	}

	l := *((*byte)(unsafe.Pointer(addr)))

	if l == 0 {
		return nil, fmt.Errorf("配置长度为0错误。")
	}

	b := *((*[1024]byte)(unsafe.Pointer(addr)))

	res := make([]byte, uint8(l))
	copy(res, b[1:l+1])

	return res, nil

}

func EncipherMem(b []byte) {
	l := len(b)
	b[0] = b[0] ^ 155
	for i := 1; i < l; i++ {
		b[i] = b[i] ^ b[0]
	}
}

func DecryptMem(b []byte) {
	l := len(b)
	for i := 1; i < l; i++ {
		b[i] = b[i] ^ b[0]
	}
	b[0] = b[0] ^ 155
}
