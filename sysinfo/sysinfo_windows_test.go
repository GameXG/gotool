package sysinfo

import (
	"testing"

	"runtime"

	"github.com/blang/semver"
)

func TestGetSemanticVersion(t *testing.T) {
	v, err := GetSemanticVersion()
	if err != nil {
		t.Fatal(err)
	}

	b, err := semver.Make("1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if v.GE(b) == false {
		t.Fatal(v, " < 1.0.0")
	}
}

func TestIs64(t *testing.T) {
	is64, err := Is64Sys()
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOARCH == "amd64" && is64 == false {
		t.Fatal("amd64!=64")
	}
}
