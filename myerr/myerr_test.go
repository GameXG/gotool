package myerr

import (
	"strings"
	"testing"
)

func TestErr(t *testing.T) {

	func() {
		defer func() {
			r := recover()
			PrintRStack(r)
		}()
		panic("")
	}()

	func() {
		defer func() {
			r := recover()
			e := ReturnRStack(r)

			if strings.Contains(e.Error(), "Stack") == false {
				t.Error("Stack")
			}
		}()
		panic("")
	}()

	func() {
		defer func() {
			r := recover()
			e := ReturnRStack(r)

			if strings.Contains(e.Error(), "Stack") || strings.Contains(e.Error(), "112233445566778899") == false {
				t.Error("000")
			}
		}()
		panic(NewPErr("112233445566%v", "778899"))
	}()

	PrintRStack(nil)

	if e := ReturnRStack(nil); e != nil {
		t.Error("e!=nil")
	}
}
