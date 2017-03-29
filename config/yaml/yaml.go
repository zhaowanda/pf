package yaml

import (
	"os"
	"log"
	"io"
	yaml "github.com/wendal/goyaml2"
	"config"
)

// 解析 yaml 配置文件
func InitYamlConfig(path string) (*config.ConfigEntity, error) {
	configure := &config.ConfigEntity{}
	configure.Config = make(map[string]interface{})
	fi, err := os.Open(path)
	defer fi.Close()
	if err != nil {
		log.Println("读取文件出错, error:", err)
		return nil, err
	}
	reader := io.Reader(fi)
	obj, err := yaml.Read(reader)
	if err != nil {
		log.Println("解析yaml文件出错,error:", err)
		return nil, err
	}
	configure.Config = obj.(map[string]interface{})
	return configure, err
}
