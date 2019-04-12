package BmMiddleware

import (
	"fmt"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"encoding/json"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
)

var BmCheckToken BmCheckTokenMiddleware

type BmCheckTokenMiddleware struct {
	Args []string
	rd   *BmRedis.BmRedis
}

type result struct {
	AllScope         string  `json:"all_scope"`
	AuthScope        string  `json:"auth_scope"`
	UserID           string  `json:"user_id"`
	ClientID         string  `json:"client_id"`
	Expires          float64 `json:"expires_in"`
	RefreshExpires   float64 `json:"refresh_expires_in"`
	Error            string  `json:"error"`
	ErrorDescription string  `json:"error_description"`
}

func (ctm BmCheckTokenMiddleware) NewCheckTokenMiddleware(args ...interface{}) BmCheckTokenMiddleware {
	var r *BmRedis.BmRedis
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	BmCheckToken = BmCheckTokenMiddleware{Args: ag, rd: r}
	return BmCheckToken
}

func (ctm BmCheckTokenMiddleware) DoMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request) {
	if _, err := ctm.CheckTokenFormFunction(w, r); err != nil {
		panic(err.Error())
	}
}

// TODO @Alex这块需要重构
func (ctm BmCheckTokenMiddleware) CheckTokenFormFunction(w http.ResponseWriter, r *http.Request) (rst *result, err error) {
	w.Header().Add("Content-Type", "application/json")

	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	version := "v0"
	resource := fmt.Sprint(ctm.Args[0], "/"+version+"/", "TokenValidation")
	mergeURL := strings.Join([]string{scheme, resource}, "")

	// 转发
	client := &http.Client{}
	req, _ := http.NewRequest("POST", mergeURL, nil)
	for k, v := range r.Header {
		req.Header.Add(k, v[0])
	}
	response, err := client.Do(req)
	if err != nil {
		return nil,err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil,err
	}

	temp := result{}
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil,err
	}

	if temp.Error != "" {
		err = errors.New(temp.ErrorDescription)
		return nil,err
	}

	return &temp, err
}