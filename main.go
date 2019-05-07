package main

import (
	"fmt"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmFactory"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmMaxDefine"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmSingleton"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"net/http"
	"os"

	//"os"
)

func main() {
	fmt.Println("pod archi begins")
	// 本机测试，添加上
	os.Setenv(BmSingleton.EnvHome, ".")

	phHome := os.Getenv(BmSingleton.EnvHome)
	fac := BmFactory.BmTable{}
	var pod = BmMaxDefine.Pod{ Name: BmSingleton.ProjectName, Factory:fac }
	pod.RegisterSerFromYAML(phHome + "/resource/def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(BmSingleton.EnvHome)
	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(BmSingleton.Version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)
	pod.RegisterAllFunctions(BmSingleton.Version, api)
	pod.RegisterAllMiddleware(api)
	handler := api.Handler().(*httprouter.Router)
	http.ListenAndServe(":"+bmRouter.Port, handler)

	fmt.Println("pod archi ends")
}
