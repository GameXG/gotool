package mem

import (
	"bytes"
	"fmt"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestMap(t *testing.T) {
	fh := int(-1)
	h, err := syscall.CreateFileMapping(syscall.Handle(fh),
		nil, uint32(windows.PAGE_READWRITE), 0, uint32(1024), syscall.StringToUTF16Ptr(NAME))
	if h == 0 {
		t.Fatal(err)
	}
	defer syscall.CloseHandle(h)

	addr, err := syscall.MapViewOfFile(h, uint32(syscall.FILE_MAP_WRITE | syscall.FILE_MAP_READ), 0,
		0, uintptr(1024))
	if addr == 0 {
		t.Fatal(err)
	}
	defer syscall.UnmapViewOfFile(addr)

	b1 := ((*[1024]byte)(unsafe.Pointer(addr)))

	s := "001234567"
	b2 := []byte(s)

	for i := 1; i <= len(b2); i++ {
		(*b1)[i] = b2[i - 1]
	}

	(*b1)[0] = uint8(len(b2))

	s1, err := ReadMem()
	if err != nil {
		t.Error(err)
	}

	if bytes.Equal(b2, s1) != true {
		t.Error(b2, "!=", s1)
	}

}

func TestC(t *testing.T) {
	b1 := []byte("127.0.0.1|6789|98.126.244.30|2200")
	b2 := make([]byte, len(b1))
	copy(b2, b1)

	fmt.Printf("%#v\r\n", b1)
	EncipherMem(b1)
	fmt.Printf("%#v\r\n", b1)
	DecryptMem(b1)

	if bytes.Equal(b1, b2) != true {
		t.Error(b1, "!=", b2)
	}
}
func TestC1(t *testing.T) {
	b1 := []byte{0xaa, 0x98, 0x9d, 0x84, 0x9a, 0x84, 0x9a, 0x84, 0x9b, 0xd6, 0x9c, 0x9d, 0x92}

	DecryptMem(b1)

	if string(b1) != "127.0.0.1|678" {
		t.Fatal(string(b1), "!=", "127.0.0.1|678")
	}
}
