package BmHandler

import (
	//"fmt"
	"net/http"
	// "crypto/tls"  
	// "log"
	"encoding/json"
	"io/ioutil"
    // "net"
	// "net/smtp"
	"github.com/alfredyang1986/blackmirror/jsonapi/jsonapiobj"
	"github.com/go-gomail/gomail"
	"reflect"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
)
type SendemailHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}
func (h SendemailHandler) NewSendemailHandler(args ...interface{}) SendemailHandler {
	var m *BmMongodb.BmMongodb
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
	return SendemailHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}
type ToEmail struct {
	Email string `json:"email" bson:"email"`
}
func (h SendemailHandler) Sendemail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	rbody, _ := ioutil.ReadAll(r.Body)
	ToEmail := ToEmail{}
	json.Unmarshal(rbody, &ToEmail)
	toEmail := ToEmail.Email
	
    m := gomail.NewMessage()
    m.SetAddressHeader("From", h.Args[2], h.Args[4])  // 发件人
    m.SetHeader("To", m.FormatAddress(toEmail, ""))   // 收件人
    m.SetHeader("Subject", h.Args[5])  // 主题
    m.SetBody("text/html", h.Args[0])  // 正文

    d := gomail.NewPlainDialer(h.Args[3],465,h.Args[2], h.Args[1])  // 发送邮件服务器、端口、发件人账号、发件人密码
    if err := d.DialAndSend(m); err != nil {
        panic(err)
	}
	response := map[string]interface{}{
		"status": "ok",
		"result": "success",
		"error":  "",
	}
	jso := jsonapiobj.JsResult{}
	jso.Obj = response
	enc := json.NewEncoder(w)
	enc.Encode(jso.Obj)
	return 1
}

func (h SendemailHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h SendemailHandler) GetHandlerMethod() string {
	return h.Method
}
