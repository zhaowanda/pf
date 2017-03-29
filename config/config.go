package config

import (
	"strings"
	"reflect"
	"errors"
)

var (
	ConfigCache = make(map[string]interface{})
)

// 定义配置文件帮助类
type ConfigEntity struct {
	Config map[string]interface{}
}

// 获取 yaml 配置
func GetConfig(conf string, config *ConfigEntity) (interface{}, error) {
	if len(conf) == 0 {
		return nil, errors.New("conf.name is null")
	}
	result, ok := ConfigCache[conf]
	if ok {
		return result, nil
	}
	if config.Config == nil {
		return nil, errors.New("conf is null")
	}
	obj := config.Config
	arr := strings.Split(conf, ".")
	for _, value := range arr {
		if strings.EqualFold(reflect.Map.String(), reflect.TypeOf(obj).Kind().String()) {
			object := obj[value]
			if object == nil {
				result = ""
				continue
			}
			if strings.EqualFold(reflect.Map.String(), reflect.TypeOf(object).Kind().String()) {
				obj = object.(map[string]interface{})
				result = obj
			} else {
				result = object
			}
		} else {
			result = obj
		}
	}
	ConfigCache[conf] = result
	return result, nil
}
