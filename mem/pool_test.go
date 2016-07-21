package mem

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	if PoolByte4096.size != 4096 {
		t.Error("poolByte4096.size!=4096")
	}

	b1 := PoolByte4096.Get()
	b2 := PoolByte4096.Get()

	if len(b1) != 4096 || len(b2) != 4096 {
		t.Error("len(b1)!=4096 || len(b2)!=4096")
	}
	b1[1] = 123

	PoolByte4096.Put(b1)
	b3 := PoolByte4096.Get()
	if b3[1] != 123 {
		t.Error("b3[1]!=123")
	}

	PoolByte4096.Put(b1[:10])
	b4 := PoolByte4096.Get()
	if b4[1] != 123 {
		t.Error("b4[1]!=123")
	}

	PoolByte4096.Put(b1[10:])
	b5 := PoolByte4096.Get()
	if len(b5) != 4096 {
		t.Error("len(b5)!=4096")
	}

	//    t.Error("Test");
}
func TestPool2(t *testing.T) {

	b1 := Get(4096)
	b2 := Get(4096)

	if len(b1) != 4096 || len(b2) != 4096 {
		t.Error("len(b1)!=4096 || len(b2)!=4096")
	}
	b1[1] = 123

	Put(b1)
	b3 := Get(4096)
	if b3[1] != 123 {
		t.Error("b3[1]!=123")
	}

	Put(b1[:10])
	b4 := Get(4096)
	if b4[1] != 123 {
		t.Error("b4[1]!=123")
	}

	Put(b1[10:])
	b5 := Get(4096)
	if len(b5) != 4096 {
		t.Error("len(b5)!=4096")
	}

	//    t.Error("Test");
}

func TestPool3(t *testing.T) {
	sTime := time.Now()
	for i := 0; i < 1000000; i++ {
		_ = Get(2048)
	}
	// 只 Get 耗时： 2.1222792s
	eTime := time.Now()
	fmt.Println("只 Get 耗时：", eTime.Sub(sTime))
	fmt.Println("平均:", eTime.Sub(sTime) / 1000000)

	sTime = time.Now()
	for i := 0; i < 1000000; i++ {
		_ = make([]byte, 2048)
	}
	// make 耗时： 63.0036ms
	eTime = time.Now()
	fmt.Println("make 耗时：", eTime.Sub(sTime))
	fmt.Println("平均:", eTime.Sub(sTime) / 1000000)

	c := make(chan []byte, 10)
	go func() {
		for buf := range c {
			Put(buf)
		}
	}()

	sTime = time.Now()
	for i := 0; i < 1000000; i++ {
		c <- Get(2048)
	}
	// Get 另一线程 Put 耗时： 533.5768ms
	eTime = time.Now()
	fmt.Println("Get 另一线程 Put 耗时：", eTime.Sub(sTime))
	fmt.Println("平均:", eTime.Sub(sTime) / 1000000)

	close(c)

	c = make(chan []byte)

	sTime = time.Now()
	for i := 0; i < 1000000; i++ {
		Put(Get(2048))
	}
	// Get Put 耗时： 167.5262ms
	eTime = time.Now()
	fmt.Println("Get Put 耗时：", eTime.Sub(sTime))
	fmt.Println("平均:", eTime.Sub(sTime) / 1000000)

	close(c)
}

func TestPool4(t *testing.T) {
	c := make(chan []byte, 10)
	get := func() []byte {
		select {
		case b := <-c:
			return b
		default:
			return make([]byte, 2048)
		}
	}
	put := func(b []byte) {
		select {
		case c <- b:
		default:
		}
	}

	sTime := time.Now()
	for i := 0; i < 1000000; i++ {
		put(get())
	}
	// chan 耗时： 142.095ms
	eTime := time.Now()
	fmt.Println("chan 耗时：", eTime.Sub(sTime))
	fmt.Println("平均:", eTime.Sub(sTime) / 1000000)

}
