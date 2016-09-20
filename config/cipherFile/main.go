package main

import (
	"fmt"

	"github.com/gamexg/gotool/config"
)

func main() {
	if err := config.CipherFile("newsocket.toml", "key", "newsocket.data"); err != nil {
		fmt.Printf("加密失败，%v\r\n", err)
	} else {
		fmt.Println("成功")
	}

}
