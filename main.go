package main

import (
	"fmt"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmFactory"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmMaxDefine"
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

	const (
		version = "v2"
		envHome = "PH_FUAD_HOME"
		projectName string = "FileUpAndDownLoad"
	)
	// 本机测试，添加上
	//os.Setenv(envHome, ".")

	phHome := os.Getenv(envHome)
	fac := BmFactory.BmTable{}
	var pod = BmMaxDefine.Pod{ Name: projectName, Factory:fac }
	pod.RegisterSerFromYAML(phHome + "/resource/def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(envHome)
	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)
	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)
	handler := api.Handler().(*httprouter.Router)
	http.ListenAndServe(":"+bmRouter.Port, handler)

	fmt.Println("pod archi ends")
}
