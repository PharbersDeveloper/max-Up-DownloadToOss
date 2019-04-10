package main

import (
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmMax/BmFactory"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmMax/BmMaxDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	//"os"
)

func main() {
	version := "v2"
	fmt.Println("pod archi begins")

	fac := BmFactory.BmTable{}
	var pod = BmMaxDefine.Pod{ Name: "swl test", Factory:fac }
	pod.RegisterSerFromYAML("resource/def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig("BM_HOME")
	//bmRouter.Port = "2019"
	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)
	pod.RegisterAllFunctions(version, api)

	handler := api.Handler().(*httprouter.Router)
	//pod.RegisterPanicHandler(handler)
	http.ListenAndServe(":"+bmRouter.Port, handler)

	fmt.Println("pod archi ends")
}
