package BmMaxDefine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/alfredyang1986/BmServiceDef/BmMiddleware"
	"github.com/alfredyang1986/BmServiceDef/BmDataStorage"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmFactory"
	"github.com/alfredyang1986/BmServiceDef/BmHandler"
	//"github.com/alfredyang1986/BmServiceDef/BmPanic"
	"github.com/alfredyang1986/BmServiceDef/BmResource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmSingleton"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/yaml.v2"
)

type Pod struct {
	Name string
	Res  map[string]interface{}
	conf Conf

	Factory BmFactory.BmTable

	Storages     map[string]BmDataStorage.BmStorage
	Resources    map[string]BmResource.BmRes
	Daemons      map[string]BmDaemons.BmDaemon
	Handler      map[string]BmHandler.BmHandler
	Middleware   map[string]BmMiddleware.BmMiddleware
	PanicHandler BmHandler.BmPanicHandler
}

func (p *Pod) RegisterSerFromYAML(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error")
	}
	//check(err)

	p.conf = Conf{}
	err = yaml.Unmarshal(data, &p.conf)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		//panic(BmPanic.ALFRED_TEST_ERROR)
	}

	p.CreateDaemonInstances()
	p.CreateStorageInstances()
	p.CreateResourceInstances()
	p.CreateFunctionInstances()
	p.CreateMiddleInstances()
	//p.CreatePanicHandleInstances()
}

func (p *Pod) CreateDaemonInstances() {
	if p.Daemons == nil {
		p.Daemons = make(map[string]BmDaemons.BmDaemon)
	}

	for _, d := range p.conf.Daemons {
		any := p.Factory.GetDaemonByName(d.Name)
		name := d.Method
		args := d.Args
		inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, name, args)
		p.Daemons[d.Name] = inc.Interface()
	}
}

func (p *Pod) CreateStorageInstances() {

	if p.Storages == nil {
		p.Storages = make(map[string]BmDataStorage.BmStorage)
	}

	for _, s := range p.conf.Storages {
		any := p.Factory.GetStorageByName(s.Name)
		name := s.Method
		var args []BmDaemons.BmDaemon
		for _, d := range s.Daemons {
			tmp := p.Daemons[d]
			args = append(args, tmp)
		}

		inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, name, args)
		p.Storages[s.Name] = inc.Interface()
	}
}

func (p *Pod) CreateResourceInstances() {
	if p.Resources == nil {
		p.Resources = make(map[string]BmResource.BmRes)
	}

	for _, r := range p.conf.Resources {
		any := p.Factory.GetResourceByName(r.Name)
		name := r.Method
		var args []BmDataStorage.BmStorage
		for _, s := range r.Storages {
			tmp := p.Storages[s] //BmFactory.GetStorageByName(s)
			args = append(args, tmp)
		}

		for _, s := range r.Friendly {
			tmp := p.Resources[s] //BmFactory.GetStorageByName(s)
			args = append(args, tmp)
		}

		inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, name, args)
		p.Resources[r.Name] = inc.Interface()
	}
}

func (p *Pod) CreateFunctionInstances() {
	if p.Handler == nil {
		p.Handler = make(map[string]BmHandler.BmHandler)
	}

	for _, r := range p.conf.Functions {
		any := p.Factory.GetFunctionByName(r.Name)
		constuctor := r.Create
		var args []BmDaemons.BmDaemon
		for _, d := range r.Daemons {
			tmp := p.Daemons[d]
			args = append(args, tmp)
		}

		inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, constuctor, args, r.Method, r.Http, r.Args)
		p.Handler[r.Name] = inc.Interface().(BmHandler.BmHandler)
	}
}

func (p *Pod) CreateMiddleInstances() {
	if p.Middleware == nil {
		p.Middleware = make(map[string]BmMiddleware.BmMiddleware)
	}

	for _, r := range p.conf.Middlewares {
		any := p.Factory.GetMiddlewareByName(r.Name)
		constuctor := r.Create
		var args []BmDaemons.BmDaemon
		for _, d := range r.Daemons {
			tmp := p.Daemons[d]
			args = append(args, tmp)
		}

		inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, constuctor, args, r.Args)
		v := inc.Interface()
		p.Middleware[r.Name] = v.(BmMiddleware.BmMiddleware)
	}
}
/*
func (p *Pod) CreatePanicHandleInstances() {
	if p.PanicHandler == nil {
		p.PanicHandler = *new(BmHandler.BmPanicHandler)
	}
	r := p.conf.Panic
	any := p.Factory.GetFunctionByName(r.Name)
	constuctor := r.Create
	inc, _ := BmSingleton.GetFactoryInstance().ReflectFunctionCall(any, constuctor)
	v := inc.Interface()
	p.PanicHandler = v.(BmHandler.BmPanicHandler)
}
*/
func (p Pod) RegisterAllResource(api *api2go.API) {
	for _, ser := range p.conf.Services {
		md := p.Factory.GetModelByName(ser.Model)
		res := p.Resources[ser.Resource]
		api.AddResource(md.(jsonapi.MarshalIdentifier), res)
	}
}

func (p Pod) RegisterAllFunctions(prefix string, api *api2go.API) {
	handler := api.Handler().(*httprouter.Router)

	// Add initial and trailing slash to prefix
	prefixSlashes := strings.Trim(prefix, "/")
	if len(prefixSlashes) > 0 {
		prefixSlashes = "/" + prefixSlashes + "/"
	} else {
		prefixSlashes = "/"
	}

	for i, _ := range p.Handler {
		ifunc := p.Handler[i]
		if ifunc.GetHttpMethod() == "POST" {
			handler.POST(prefixSlashes+ifunc.GetHandlerMethod(), func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
				BmSingleton.GetFactoryInstance().ReflectFunctionCall(ifunc, ifunc.GetHandlerMethod(), writer, request, params)
			})
		} else {
			handler.GET(prefixSlashes+ifunc.GetHandlerMethod(), func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
				BmSingleton.GetFactoryInstance().ReflectFunctionCall(ifunc, ifunc.GetHandlerMethod(), writer, request, params)
			})
		}
	}

}

func (p Pod) RegisterAllMiddleware(api *api2go.API) {
	for _, mw := range p.Middleware {
		api.UseMiddleware(mw.DoMiddleware)
	}
}
/*
func (p Pod) RegisterPanicHandler(router *httprouter.Router) {
	router.PanicHandler = p.PanicHandler.HandlePanic
}
*/