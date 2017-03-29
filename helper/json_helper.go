package helper

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
)

// json 转换为字符串
func JO2Str(js *simplejson.Json) (string, error) {
	result, err := json.Marshal(js.Interface())
	if err != nil {
		return "", err
	}
	return string(result), err
}

// string 转换为 json
func Str2JO(str string) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(str))
}

// byte 转换为 json
func Byte2JO(body []byte) (*simplejson.Json, error) {
	return simplejson.NewJson(body)
}

// 对象转换为json string
func Obj2Json(args interface{}) (string, error) {
	result, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	return string(result), err
}
// 参数应该为偶数个
func Gen(args ...interface{}) *simplejson.Json {
	result := simplejson.New()
	for i := 0; i < len(args); i += 2 {
		if (args[i + 1] == nil) {
			result.Del(args[i].(string))
		} else {
			result.Set(args[i].(string), args[i + 1]);
		}
	}
	return result
}
