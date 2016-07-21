package myerr

import (
	"fmt"
	"runtime/debug"

	log "bitbucket.org/jack/log4go"
)

// 故意 panic 方式引发错误时和系统 panic 区分开。
type MyPErr struct {
	s string
}

func (e *MyPErr) Error() string {
	return e.s
}

func NewPErr(s string, a ...interface{}) *MyPErr {
	return &MyPErr{
		s: fmt.Sprintf(s, a...),
	}
}

// 排除 MyPErr ，之外的都打印 Stack
func PrintRStack(r interface{}) {
	if r != nil {
		if merr, ok := r.(*MyPErr); ok == true {
			log.Info("MyPErr:%v", merr.Error())
		} else if err, ok := r.(error); ok == true {
			log.Warn("panic:%v\r\nStack:%v", err, string(debug.Stack()))
		} else {
			log.Warn("panic:%v\r\nStack:%v", r, string(debug.Stack()))
		}
	}
}

// 除 MyPErr 之外的都带 Stack
func ReturnRStack(r interface{}) error {
	if r != nil {
		if merr, ok := r.(*MyPErr); ok == true {
			return merr
		} else if err, ok := r.(error); ok == true {
			return fmt.Errorf("panic:%v\r\nStack:%v", err, string(debug.Stack()))
		} else {
			return fmt.Errorf("panic:%v\r\nStack:%v", r, string(debug.Stack()))
		}
	}
	return nil
}
