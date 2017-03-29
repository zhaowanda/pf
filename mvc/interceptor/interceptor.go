package interceptor

import (
	"container/list"
	"errors"
	"fmt"
	fast "github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strings"
)
// 初始化拦截器存储地方
var (
	interceptorFactory = make(map[string]InterceptorBean)
)

// http 请求拦截器
type Interceptor interface {
	// 拦截 http 请求, 根据需要做一些判断, 返回是否允许后续逻辑继续处理请求, 如返回 false 则表示请求到此为止.
	// 请注意, 后续逻辑需要读取 ctx.Body 里的内容, 请谨慎读取!
	Intercept(ctx *fast.RequestCtx) (bool, error)
}

type InterceptorFunc func(ctx *fast.RequestCtx) (bool, error)

//  拦截url接口
type Patterns interface {
	AddIncludePattern(name string, pattern string) InterceptorBean
	AddExcludePattern(name string, pattern string) InterceptorBean
	ExecutedInterceptor(ctx *fast.RequestCtx, interceptor InterceptorBean) (bool, error)
}

//  interceptor entity
type InterceptorBean struct {
	Interceptor     InterceptorFunc
	includePatterns *list.List
	excludePatterns *list.List
}

// 拦截器的实现
func (fn InterceptorFunc) Intercept(ctx *fast.RequestCtx) (bool, error) {
	return fn(ctx)
}

// 注册单个interceptor
func Register(name string, interceptor InterceptorBean) {
	_, interceptors := interceptorFactory[name]
	if interceptors {
		log.Println(fmt.Sprintf("InterceptorBean named %s already registered", name))
	} else {
		interceptorFactory[name] = interceptor
	}
}

// 批量注册interceptor
func Registers(interceptors map[string]InterceptorBean) {
	for key, _ := range interceptors {
		Register(key, interceptors[key])
	}
}

// 获取interceptor
func GetInterceptor(name string) (InterceptorBean, bool) {
	interceptor, ok := interceptorFactory[name]
	return interceptor, ok

}

// 添加需要拦截的url
func (ib InterceptorBean) AddIncludePattern(name string, pattern string) InterceptorBean {
	interceptor, ok := GetInterceptor(name)
	if !ok {
		Register(name, InterceptorBean{})
	}
	include := interceptor.includePatterns
	if include == nil {
		include = list.New()
	}

	include.PushBack(pattern)
	interceptor.includePatterns = include
	interceptorFactory[name] = interceptor
	return interceptor
}

// 添加不需要拦截的url
func (ib InterceptorBean) AddExcludePattern(name string, pattern string) InterceptorBean {
	interceptor, ok := GetInterceptor(name)
	if !ok {
		Register(name, InterceptorBean{})
	}
	exclude := interceptor.excludePatterns
	if exclude == nil {
		exclude = list.New()
	}
	exclude.PushBack(pattern)
	//log.Println("len:", exclude.Len())

	interceptor.excludePatterns = exclude
	interceptorFactory[name] = interceptor
	return interceptor
}

// 拦截执行器
func (ib InterceptorBean) ExecutedInterceptor(ctx *fast.RequestCtx, interceptor InterceptorBean) (bool, error) {
	ctx.Request.URI()

	fullRequestUrl, err := url.Parse(string(ctx.Request.RequestURI()))
	if err != nil {
		return false, errors.New("Url is incorrect")
	}

	requestUrl := fullRequestUrl.Path
	// todo 解析request url 层级关系 支持通配符方式来拦截相关url
	inc := interceptor.Interceptor
	if inc == nil {
		return false, errors.New("Interceptor not allowed nil")
	}
	exclude := interceptor.excludePatterns
	isInclude := false

	if exclude == nil {

	} else {
		for it := exclude.Front(); it != nil; it = it.Next() {
			pattern := it.Value.(string)
			pattern = strings.Replace(pattern, "*", "", len(pattern))
			if strings.HasPrefix(requestUrl, pattern) {
				isInclude = true
				break
			}
		}
	}
	if isInclude {
		return true, nil
	}
	include := interceptor.includePatterns
	if include == nil {
		return true, nil
	} else {
		for it := include.Front(); it != nil; it = it.Next() {
			pattern := it.Value.(string)
			pattern = strings.Replace(pattern, "*", "", len(pattern))
			if strings.HasPrefix(requestUrl, pattern) {
				isInclude = true
				break
			}
		}
	}
	if !isInclude {
		return true, nil
	}
	interceptFunc := InterceptorFunc(interceptor.Interceptor)
	flag, err := interceptFunc.Intercept(ctx)
	return flag, err
}
