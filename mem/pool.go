package mem

import (
	"runtime/debug"
	"sync"

	"fmt"
)

/*
内存池



*/

func init() {
	PoolByte4096 = NewPool(4096)
	PoolByte2048 = NewPool(2048)
	PoolByte1024 = NewPool(1024)
	PoolByteforward = PoolByte4096
}

// 全局内存池
var PoolByteforward *Pool
var PoolByte4096 *Pool
var PoolByte2048 *Pool
var PoolByte1024 *Pool

// 这个封装的目的是防止返回的切片尺寸不对
type Pool struct {
	pool *sync.Pool
	size int
}

// 新建一个池
func NewPool(size int) *Pool {
	p := Pool{}
	p.size = size
	p.pool = &sync.Pool{New: func() interface{} {
		return make([]byte, size)
	}}

	return &p
}

func (p *Pool) Get() []byte {
	// sync.pool 是线程安全的
	return p.pool.Get().([]byte)
}

// 回填池
//可以放心尺寸的问题，切片尺寸有问题会被放弃
func (p *Pool) Put(buf []byte) {
	if cap(buf) == p.size {
		p.pool.Put(buf[:p.size])
	}
}

// 自动选择合适的池
func Get(l int) []byte {
	if l > 4096 {
		return make([]byte, l)
	} else if l > 2048 {
		return PoolByte4096.Get()[:l]
	} else if l > 1024 {
		return PoolByte2048.Get()[:l]
	} else {
		return PoolByte1024.Get()[:l]
	}
}

// 回填池
//可以放心尺寸的问题，切片尺寸有问题会被放弃
func Put(buf []byte) error {
	l := cap(buf)
	switch l {
	case 4096:
		PoolByte4096.Put(buf)
	case 2048:
		PoolByte2048.Put(buf)
	case 1024:
		PoolByte1024.Put(buf)
	default:
		return fmt.Errorf("内存池释放错误，尺寸 %v 不匹配。调用上下文： %v \r\n", l, string(debug.Stack()))
	}
	return nil
}
