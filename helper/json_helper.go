package helper

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
)

// json 转换为字符串
func Json2String(js *simplejson.Json) (string, error) {
	result, err := json.Marshal(js.Interface())
	if err != nil {
		return "", err
	}
	return string(result), err
}

// string 转换为 json
func String2Json(str string) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(str))
}

// byte 转换为 json
func Byte2Json(body []byte) (*simplejson.Json, error) {
	return simplejson.NewJson(body)
}

// 对象转换为json string
func Object2JsonString(args interface{}) (string, error) {
	result, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	return string(result), err
}
