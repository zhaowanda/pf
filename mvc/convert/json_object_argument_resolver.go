package convert

import (
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	fast "github.com/valyala/fasthttp"
	"github.com/zhaowanda/pf/helper"
	"net/url"
	"strings"
)

// 参数转换成json
func ArgumentResolver(ctx *fast.RequestCtx) (*simplejson.Json, error) {
	param := simplejson.New()
	if ctx.IsGet() { // 判断是否为get请求
		requestUrl, err := url.Parse(string(ctx.Request.RequestURI()))
		if err != nil {
			err = errors.New("解析url出错")
			return param, err
		}
		if requestUrl.RawQuery == "" {
			return param, err
		}
		params := parseParam(requestUrl.RawQuery, "&")
		for _, value := range params {
			kv := parseParam(value, "=")
			param.Set(kv[0], kv[1])
		}
		return param, nil
	} else if ctx.IsPost() { // 判断是否为post请求
		paramBody, error := helper.Byte2JO(ctx.Request.Body())
		if error != nil {
			return param, error
		}
		requestUrl, err := url.Parse(string(ctx.Request.RequestURI()))
		if err != nil {
			err = errors.New("解析url出错")
			return param, err
		}
		if requestUrl.RawQuery == "" {
			return param, err
		}
		params := parseParam(requestUrl.RawQuery, "&")
		for _, value := range params {
			kv := parseParam(value, "=")
			paramBody.Set(kv[0], kv[1])
		}
		return paramBody, nil
	} else {
		return param, nil
	}

}

// 解析request param参数
func parseParam(param string, regexp string) []string {
	return strings.Split(param, regexp)
}
