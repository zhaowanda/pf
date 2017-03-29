package convert

import (
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	fast "github.com/valyala/fasthttp"
	"github.com/zhaowanda/pf/helper"
)

// 定义转换接口
type HttpMessageConvert interface {
	MessageConvert(ctx *fast.RequestCtx, args interface{}, errMsg error)
}

// 定义抽象方法
type HttpMessageConvertFunc func(ctx *fast.RequestCtx, args interface{}, errMsg error)

// 实现接口
func (hmc HttpMessageConvertFunc) MessageConvert(ctx *fast.RequestCtx, args interface{}, errMsg error) {
	hmc(ctx, args, errMsg)
}

// json 转换器
func JsonHttpMessageConvert(ctx *fast.RequestCtx, args interface{}, errMsg error) {
	result := simplejson.New()
	if errMsg != nil {
		result.Set("error", errMsg.Error())
	} else {
		result.Set("data", args)
	}
	resultString, err := helper.JO2Str(result)
	if err != nil {
		fmt.Fprint(ctx, err)
		return
	}
	fmt.Fprint(ctx, resultString)
}
