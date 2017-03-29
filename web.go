package web

import (
	fast "github.com/valyala/fasthttp"
	"log"
	"os"
	"github.com/zhaowanda/pf/mvc/router"
	"github.com/zhaowanda/pf/config/yaml"
        "github.com/zhaowanda/pf/config"
)

//  提供restful接口,页面后期进行整合

func WebMain() {
	ipAddr, err := config.GetConfig("server.boot.ListenAndServe", yaml.YamlConfigure)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}
	log.Fatal(fast.ListenAndServe(ipAddr.(string), webrouter.InitRouter().Handler))
}
