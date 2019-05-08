package BmHandler

import (
	"encoding/json"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type VerifiedSmsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	s          *BmSms.BmSms
	r          *BmRedis.BmRedis
}

func (h VerifiedSmsHandler) NewVerifiedSmsHandler(args ...interface{}) VerifiedSmsHandler {
	var s *BmSms.BmSms
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
				if tm.Name() == "BmSms" {
					s = dm.(*BmSms.BmSms)
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

	return VerifiedSmsHandler{Method: md, HttpMethod: hm, Args: ag, s: s, r: r}
}

func (h VerifiedSmsHandler) VerifiedSmsCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sr := SmsRecord{}
	json.Unmarshal(body, &sr)
	rcode, err := h.r.GetPhoneCode(sr.Phone)
	response := map[string]interface{}{
		"status": "",
		"error":  nil,
	}
	if err==nil && rcode == sr.Code {
		response["status"] = "ok"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 0
	}

	response["status"] = "error"
	if err!=nil && err.Error() == "phoneCode expired" {
		response["error"] = "验证码过期"
	} else {
		response["error"] = "验证码错误！"
	}
	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 1
}

func (h VerifiedSmsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h VerifiedSmsHandler) GetHandlerMethod() string {
	return h.Method
}
