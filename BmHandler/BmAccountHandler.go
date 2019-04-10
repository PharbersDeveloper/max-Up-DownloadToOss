package BmHandler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"

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

func (h AccountHandler) AccountValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	res := BmModel.Account{}
	json.Unmarshal(body, &res)
	var out BmModel.Account
	cond := bson.M{"account": res.Account, "password": res.Password}
	err = h.db.FindOneByCondition(&res, &out, cond)

	response := map[string]interface{}{
		"status": "",
		"result": nil,
		"error":  nil,
	}

	if err == nil && out.ID != "" {
		hex := md5.New()
		io.WriteString(hex, out.ID)
		out.Password = ""
		token := fmt.Sprintf("%x", hex.Sum(nil))
		err = h.rd.PushToken(token, time.Hour*24*365)
		out.Token = token

		response["status"] = "ok"
		response["result"] = out
		response["error"] = err

		//reval, _ := json.Marshal(response)
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 0
	} else {
		response["status"] = "error"
		response["error"] = "账户或密码错误！"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
}

func (h AccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h AccountHandler) GetHandlerMethod() string {
	return h.Method
}
