package ufile

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

// 测试本地文件
func TestLFile(t *testing.T) {
	os.Remove("ufile-test.txt")

	// 准备测试文件
	f, err := os.Create("ufile-test.txt")
	if _, err := f.Write([]byte("741852")); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// 初始化配置
	ufile, err := NewUFile(".", 1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err := ufile.Add("ufile-test.txt", 0, 123); err != nil {
		t.Fatal(err)
	}

	// 测试读取
	res := <-ufile.ResChan
	if res.RawPath != "ufile-test.txt" ||
		res.Err != nil ||
		strings.HasSuffix(res.Path, "ufile-test.txt") == false ||
		res.Rc == nil ||
		res.Userdata != 123 {
		t.Fatal("err")
	}
	data, err := ioutil.ReadAll(res.Rc)
	if err != nil {
		t.Fatal(err)
	}
	res.Rc.Close()
	if bytes.Equal(data, []byte("741852")) == false {
		t.Fail()
	}

	// 第二次写入
	f, err = os.Create("ufile-test.txt")
	if _, err := f.Write([]byte("1234567890")); err != nil {
		t.Fatal(err)
	}
	f.Close()

	// 第二次读取
	res = <-ufile.ResChan
	if res.Err != nil {
		t.Fatal("res.Err!=nil")
	}
	if res.RawPath != "ufile-test.txt" ||
		res.Err != nil ||
		strings.HasSuffix(res.Path, "ufile-test.txt") == false ||
		res.Rc == nil ||
		res.Userdata != 123 {
		t.Fatal("err")
	}
	data, err = ioutil.ReadAll(res.Rc)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(data, []byte("1234567890")) == false {
		t.Fatal("第二次匹配错误")
	}
	res.Rc.Close()

	if err := ufile.Remove("ufile-test.txt"); err != nil {
		t.Fatal("删除错误：", err)
	}

	// 读出可能存在的修改内容
	func() {
		timeout := time.After(2 * time.Second)
		for {
			select {
			case res = <-ufile.ResChan:
			case <-timeout:
				return
			}
		}
	}()

	// 第3次写入
	f, err = os.Create("ufile-test.txt")
	if _, err := f.Write([]byte("1234567890")); err != nil {
		t.Fatal(err)
	}
	f.Close()

	timeout2 := time.After(3 * time.Second)

	select {
	case res = <-ufile.ResChan:
		t.Fatal("第三次匹配错误：", res)
	case <-timeout2:
	}

	ufile.Close()

	os.Remove("ufile-test.txt")

}

// 测试http文件
func TestHFile(t *testing.T) {

	ufile, err := NewUFile("", 1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if err := ufile.Add("https://www.baidu.com/", 1 * time.Second, 123); err != nil {
		t.Fatal(err)
	}

	time.AfterFunc(10 * time.Second, func() {
		select {
		case ufile.ResChan <- nil:
			t.Fail()
		default:
		}
	})

	for i := 0; i < 2; i++ {

		res := <-ufile.ResChan

		if res.Err != nil {
			t.Fatal("res.Err!=nil")
		}

		if res.RawPath != "https://www.baidu.com/" ||
			res.Err != nil ||
			res.Path != "https://www.baidu.com/" ||
			res.Rc == nil ||
			res.Userdata != 123 {
			t.Fatal("err")
		}

		data, err := ioutil.ReadAll(res.Rc)
		if err != nil {
			t.Fatal(err)
		}
		res.Rc.Close()

		if bytes.Contains(data, []byte("baidu")) == false {
			t.Fail()
		}
	}

	ufile.Remove("https://www.baidu.com/")

	timeout := time.After(3 * time.Second)

	select {
	case <-ufile.ResChan:
		t.Fail()
	case <-timeout:
	}

	ufile.Close()
}
