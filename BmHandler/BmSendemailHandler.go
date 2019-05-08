package BmHandler

import (
	"fmt"
	"net/http"
	"crypto/tls"  
	"log"
	"encoding/json"
	"io/ioutil"
    "net"
    "net/smtp"
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
	response := map[string]interface{}{
		"status": "",
		"result": nil,
		"error":  nil,
	}
	json.Unmarshal(rbody, &ToEmail)
	toEmail := ToEmail.Email
	
    header := make(map[string]string)
    header["From"] = h.Args[4] + "<" + h.Args[2] + ">"
    header["To"] = toEmail
    header["Subject"] = h.Args[5]
    header["Content-Type"] = h.Args[6]
    body := "www.bing.com"
    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body
    auth := smtp.PlainAuth(
        "",
        h.Args[2],
        h.Args[3],
        h.Args[0],
    )
    err := h.SendMailUsingTLS(
		h.Args[0] + ":" + h.Args[1],
        auth,
        h.Args[2],
        []string{toEmail},
        []byte(message),
    )
    if err != nil {
        panic(err)
    } else {
        fmt.Println("Send mail success!")
	}
	response["status"] = "success"
	response["error"] = ""
	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 1
}

func (h SendemailHandler)Dial(addr string) (*smtp.Client, error) {
    conn, err := tls.Dial("tcp", addr, nil)
    if err != nil {
        log.Println("Dialing Error:", err)
        return nil, err
    }
    //分解主机端口字符串
    host, _, _ := net.SplitHostPort(addr)
    return smtp.NewClient(conn, host)
}

func (h SendemailHandler)SendMailUsingTLS(addr string, auth smtp.Auth, from string,
    to []string, msg []byte) (err error) {
    //create smtp client
    c, err := h.Dial(addr)
    if err != nil {
        log.Println("Create smpt client error:", err)
        return err
    }
    defer c.Close()
    if auth != nil {
        if ok, _ := c.Extension("AUTH"); ok {
            if err = c.Auth(auth); err != nil {
                log.Println("Error during AUTH", err)
                return err
            }
        }
    }
    if err = c.Mail(from); err != nil {
        return err
    }
    for _, addr := range to {
        if err = c.Rcpt(addr); err != nil {
            return err
        }
    }
    w, err := c.Data()
    if err != nil {
        return err
    }
    _, err = w.Write(msg)
    if err != nil {
        return err
    }
    err = w.Close()
    if err != nil {
        return err
    }
    return c.Quit()
}
func (h SendemailHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h SendemailHandler) GetHandlerMethod() string {
	return h.Method
}
