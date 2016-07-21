package zip

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

func TestNonBlocking(t *testing.T) {
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:1", "gzip:1", "deflate:1"}

	// 产生随机测试数据
	testDatas := make([][]byte, 1000)
	for i := range testDatas {
		t := make([]byte, rand.Intn(1024) + 1)
		rand.Read(t)
		testDatas[i] = t
	}
	l, err := net.Listen("tcp", "127.0.0.1:15634")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for _, zName := range zipTypes {
		fmt.Println("type:", zName)

		cd := make(chan []byte)
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			c, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()

			r, err := NewZipRead(c, zName)
			if err != nil {
				t.Fatal(err)
			}

			for data := range cd {
				buf := make([]byte, len(data))

				if _, err := io.ReadFull(r, buf); err != nil {
					t.Fatal(err)
				} else if bytes.Equal(data, buf) == false {
					t.Fatal(data, "!=", buf)
				}
			}

		}()

		c, err := net.DialTimeout("tcp", l.Addr().String(), 3 * time.Second)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		w, err := NewZipWrite(c, zName)
		if err != nil {
			t.Fatal(err)
		}
		for _, data := range testDatas {
			//fmt.Println("write:", len(data))
			if _, err := w.Write(data); err != nil {
				t.Fatal(err)
			}
			//fmt.Println("write ok ")

			cd <- data
		}
		close(cd)

		wg.Wait()
	}
}

func TestBlocking(t *testing.T) {
	zipTypes := []string{"zlib", "gzip", "deflate", "zlib:1", "gzip:1", "deflate:1"}

	// 产生随机测试数据
	testDatas := make([][]byte, 1000)
	for i := range testDatas {
		t := make([]byte, rand.Intn(1024) + 1)
		rand.Read(t)
		testDatas[i] = t
	}

	for _, zName := range zipTypes {
		fmt.Println("type:", zName)

		r, w := io.Pipe()

		zr, err := NewZipRead(r, zName)
		if err != nil {
			t.Fatal(err)
		}
		zw, err := NewZipWrite(w, zName)
		if err != nil {
			t.Fatal(err)
		}

		c := make(chan []byte)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			for data := range c {
				buf := make([]byte, len(data))
				if _, err := io.ReadFull(zr, buf); err != nil {
					t.Fatal(err)
				} else if bytes.Equal(buf, data) != true {
					t.Fatal(buf, "!=", data)
				}
			}
		}()

		for _, data := range testDatas {
			c <- data
			if _, err := zw.Write(data); err != nil {
				t.Fatal(err)
			}
		}
		close(c)
		wg.Wait()
	}
}
