package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gamexg/goio/mycipher"
)

func Decode(filename string, conf interface{}) error {

	config_path, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	//config_dir := filepath.Dir(config_path)

	// 打开配置文件
	configFile, err := os.Open(config_path)
	if err != nil {
		return fmt.Errorf("%v %v", config_path, err)
	}
	defer configFile.Close()

	buf := make([]byte, 3)
	if _, err := io.ReadFull(configFile, buf); err != nil {
		return fmt.Errorf("%v %v", config_path, err)
	}
	if bytes.Equal(buf, []byte{0xEF, 0xBB, 0xBF}) == false {
		configFile.Seek(0, 0)
	}

	_, err = toml.DecodeReader(configFile, conf)
	if err != nil {
		return err
	}
	return nil
}

func DecodeCipherFile(filename, key string, conf interface{}) error {
	config_path, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	//config_dir := filepath.Dir(config_path)

	// 打开配置文件
	configFile, err := os.Open(config_path)
	if err != nil {
		return fmt.Errorf("%v %v", config_path, err)
	}
	defer configFile.Close()

	//解密
	cr, err := mycipher.NewCipherRead(key, configFile)
	if err != nil {
		return fmt.Errorf("解密 %v 错误，%v", config_path, err)
	}

	confBuf := bytes.Buffer{}
	// 存到内存
	if _, err := io.Copy(&confBuf, cr); err != nil {
		return fmt.Errorf("读取 %v 错误，%v", config_path, err)
	}

	if confBuf.Len() >= 3 && bytes.Equal(confBuf.Bytes()[:3], []byte{0xEF, 0xBB, 0xBF}) {
		buf := make([]byte, 3)
		if _, err := io.ReadFull(&confBuf, buf); err != nil {
			return fmt.Errorf("%v %v", config_path, err)
		}
	}

	_, err = toml.DecodeReader(&confBuf, conf)
	if err != nil {
		return err
	}
	return nil
}

//加密文件
func CipherFile(rawFilename, key string, cipherFilename string) error {
	rf, err := os.OpenFile(rawFilename, os.O_RDONLY, 0660)
	if err != nil {
		return fmt.Errorf("无法打开 %v 文件,%v\r\n", rawFilename, err)
	}
	defer rf.Close()

	wf, err := os.OpenFile(cipherFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return fmt.Errorf("无法打开 %v 文件,%v\r\n", cipherFilename, err)
	}
	defer wf.Close()

	c, err := mycipher.NewCipherWrite(key, wf)
	if err != nil {
		return fmt.Errorf("创建加密起失败，%v", err)
	}

	if _, err := io.Copy(c, rf); err != nil {
		return fmt.Errorf("写失败，%v", err)
	}

	return nil
}
