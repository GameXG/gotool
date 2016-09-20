package config

import (
	"os"
	"testing"
)

type ServerConfig struct {
	LAddr              string
	RAddr              string
	rport              int
	Socks5Timeout      int
	ForwardTimeout     int
	ForwardDelay       bool
	PansiChar          string
	SoftMaxConnCount   int32
	HardMaxConnCount   int32
	HardMaxConnMessage string
	User               string
	UpConnCountUrl     string
	Line               string
}

func TestDecodeCipherFile(t *testing.T) {
	rf, err := os.OpenFile("tconfig.toml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := rf.Write([]byte(`
# 配置文件

#LAddr=":6789"
# 本地 socks5 监听地址 废弃

#RAddr="127.0.0.1:2001"
# 中转服务器地址 废弃

Socks5Timeout=123
# socks5 命令超时时间，单位秒
# 注意，目前还包括连接目标网站的时间，需要注意不能太小。

ForwardTimeout=456
# socks5 转发超时，单位秒

ForwardDelay=true
#socks5协议协商完成后，纯转发请求时是否开启 Delay ，true 开启，增加吞吐量，增加响应延迟。

HardMaxConnMessage = "159357"
# 硬性连接数限制消息框内容。

SoftMaxConnCount = 789
HardMaxConnCount = 147

UpConnCountUrl = "http://1111111"
`)); err != nil {
		t.Fatal(err)
	}

	rf.Close()

	if err := CipherFile("tconfig.toml", "123456", "tconfig.data"); err != nil {
		t.Fatal(err)
	}

	conf := ServerConfig{}

	if err := DecodeCipherFile("tconfig.data", "123456", &conf); err != nil {
		t.Fatal(err)
	}

	if conf.Socks5Timeout != 123 || conf.ForwardTimeout != 456 || conf.ForwardDelay != true || conf.HardMaxConnMessage != "159357" ||
		conf.SoftMaxConnCount != 789 || conf.HardMaxConnCount != 147 || conf.UpConnCountUrl != "http://1111111" {
		t.Fatal(conf)
	}

}
