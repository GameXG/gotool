package zip

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type zipWrite struct {
	zipName  string
	zipType  string // 压缩类型
	zipLevel int    // 压缩级别
	rw       io.Writer
	ow       io.Writer
	of       Flusher
	oc       io.Closer
}

type Flusher interface {
	Flush() error
}

// 注意，生成的对象非线程安全
// 引用库的话：不再使用时应该调用 Close 方法，Close 只是关闭 压缩流，不会关闭低层流。
func NewZipWrite(w io.Writer, name string) (io.WriteCloser, error) {
	ss := strings.SplitN(name, ":", 2)

	zipType := strings.TrimSpace(ss[0])
	zipArgs := "-1" //默认压缩级别

	switch zipType {
	case "zlib", "gzip", "deflate":
	default:
		return nil, fmt.Errorf("未知的压缩类型 %v", zipType)
	}

	if len(ss) > 1 {
		zipArgs = ss[1]
	}

	ZipLevel, err := strconv.Atoi(zipArgs)
	if err != nil {
		return nil, fmt.Errorf("压缩参数错误，%v 无法转换成为数字，%v", zipArgs, err)
	}

	return &zipWrite{
		zipName:  name,
		zipType:  zipType,
		zipLevel: ZipLevel,
		rw:       w,
	}, nil
}
func (z *zipWrite) init() error {
	if z.ow == nil {
		// zip库实现问题，需要延迟求值
		switch z.zipType {
		case "zlib":
			zw, err := zlib.NewWriterLevel(z.rw, z.zipLevel)
			if err != nil {
				return fmt.Errorf("创建压缩写失败，%v", err)
			}
			z.ow = zw
			z.of = zw
			z.oc = zw
		case "gzip":
			zw, err := gzip.NewWriterLevel(z.rw, z.zipLevel)
			if err != nil {
				return fmt.Errorf("创建压缩写失败，%v", err)
			}
			z.ow = zw
			z.of = zw
			z.oc = zw

		case "deflate":
			zw, err := flate.NewWriter(z.rw, z.zipLevel)
			if err != nil {
				return fmt.Errorf("创建压缩写失败，%v", err)
			}
			z.ow = zw
			z.of = zw
			z.oc = zw
		default:
			return fmt.Errorf("未知的压缩类型:%v。", z.zipLevel)
		}
	}
	return nil
}

func (z *zipWrite) Write(b []byte) (int, error) {
	if z.ow == nil {
		if err := z.init(); err != nil {
			return 0, err
		}
	}

	n, err := z.ow.Write(b)
	if err != nil {
		return n, err
	}

	if z.of != nil {
		return n, z.of.Flush()
	}
	return n, nil
}

func (z *zipWrite) Close() error {
	if z.ow == nil {
		if err := z.init(); err != nil {
			return err
		}
	}

	if z.oc != nil {
		return z.oc.Close()
	}
	return nil
}

type zipRead struct {
	zipName string
	zipType string // 压缩类型
	zipArgs string
	rr      io.Reader
	or      io.Reader
	oc      io.Closer
}

func NewZipRead(r io.Reader, name string) (io.ReadCloser, error) {
	ss := strings.SplitN(name, ":", 2)

	zipType := strings.TrimSpace(ss[0])
	zipArgs := ""

	switch zipType {
	case "zlib", "gzip", "deflate":
	default:
		return nil, fmt.Errorf("未知的压缩类型 %v", zipType)
	}

	if len(ss) > 1 {
		zipArgs = ss[1]
	}

	return &zipRead{
		zipName: name,
		zipType: zipType,
		zipArgs: zipArgs,
		rr:      r,
	}, nil
}

func (z *zipRead) init() error {
	if z.or == nil {
		switch z.zipType {
		case "zlib":
			zr, err := zlib.NewReader(z.rr)
			if err != nil {
				return err
			}
			z.or = zr
			z.oc = zr
		case "gzip":
			zr, err := gzip.NewReader(z.rr)
			if err != nil {
				return err
			}
			z.or = zr
			z.oc = zr
		case "deflate":
			zr := flate.NewReader(z.rr)
			z.or = zr
			z.oc = zr
		default:
			return fmt.Errorf("未知的压缩类型:%v。", z.zipType)
		}
	}
	return nil
}

func (z *zipRead) Read(b []byte) (int, error) {
	if z.or == nil {
		if err := z.init(); err != nil {
			return 0, err
		}
	}
	return z.or.Read(b)
}

func (z *zipRead) Close() error {
	if z.or == nil {
		if err := z.init(); err != nil {
			return err
		}
	}
	if z.oc != nil {
		return z.oc.Close()
	}
	return nil

}
