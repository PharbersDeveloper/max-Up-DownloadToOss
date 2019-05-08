package BmHandler

import (
	"fmt"
	"net/http"
	"reflect"
	//"time"
	//"strings"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"github.com/alfredyang1986/blackmirror/jsonapi/jsonapiobj"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"encoding/json"
	"io/ioutil"
	"github.com/manyminds/api2go"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmResource"
)

type AccountHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h AccountHandler) NewAccountHandler(args ...interface{}) AccountHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	return AccountHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h AccountHandler) AccountValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (api2go.Responder) {
	w.Header().Add("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	res := BmModel.Account{}
	json.Unmarshal(body, &res)

	var emailout BmModel.Account
	var phoneout BmModel.Account
	var resAccount BmModel.Account
	emailcond := bson.M{"email": res.Email, "password": res.Password}
	phonecond := bson.M{"phone": res.Phone, "password": res.Password}
	emailerr := h.db.FindOneByCondition(&res, &emailout,emailcond)
	phoneerr := h.db.FindOneByCondition(&res, &phoneout,phonecond)
	if (emailerr != nil && emailout.ID == "")&&(phoneerr != nil && phoneout.ID == ""){
		response := map[string]interface{}{
			"status": "error",
			"result": "用户名或密码错误",
			"error":  "",
		}
		jso := jsonapiobj.JsResult{}
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
	}
	if emailerr == nil {
		resAccount = emailout
	}else {
		resAccount = phoneout
	}
	return &BmResource.Response{Res: resAccount}
}

func (h AccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h AccountHandler) GetHandlerMethod() string {
	return h.Method
}
