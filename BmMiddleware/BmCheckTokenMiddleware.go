package BmMiddleware

import (
	"fmt"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"encoding/json"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
	"time"
)

var BmCheckToken BmCheckTokenMiddleware

const PROJECT_NAME string = "FileUpAndDownLoad"

type BmCheckTokenMiddleware struct {
	Args []string
	rd   *BmRedis.BmRedis
}

type result struct {
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
	if err := ctm.CheckTokenFormFunction(w, r); err != nil {
		panic(err.Error())
	}
}

// TODO @Alex这块需要重构
func (ctm BmCheckTokenMiddleware) CheckTokenFormFunction(w http.ResponseWriter, r *http.Request) (err error) {
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
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	temp := result{}
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return
	}

	if temp.Error != "" {
		err = errors.New(temp.ErrorDescription)
		return
	}

	accessed, accessOpt := checkAccessScope(temp.AuthScope)	//允许访问(判断是否有授权)
	indate := checkIndateScope(accessOpt)	//有效期内(判断授权是否过期)
	// TODO @Jeorch这块需要重构
	if accessed && indate {
		r.URL.RawQuery = fmt.Sprint("accept=", getOperatedCompany(accessOpt))
	} else {
		err = errors.New("expired scope")
	}
	return

}

func checkAccessScope(userScope string) (accessed bool, accessOpt string) {

	accessed = false
	userAccessOpts := strings.Split(userScope, "/")[1]
	userAccessOptArr := strings.Split(userAccessOpts, ",")
	for _, userAccessOpt := range userAccessOptArr {
		userAccess := strings.Split(userAccessOpt, ":")[0]
		if userAccess == PROJECT_NAME {
			//TODO:目前只是检查有无访问项目的权限，还未进行具体操作权限的check => func checkOperationScope(operationCmd string) (allowed bool) {}
			accessed = true
			accessOpt = userAccessOpt
			return
		}
	}
	return
}

func checkIndateScope(accessOpt string) (indate bool) {

	indate = false
	if accessOpt == "" {
		return
	}
	accessOptArr := strings.Split(accessOpt, ":")
	if len(accessOptArr) < 2 { //项目名+操作
		return
	}
	operation := accessOptArr[1]
	optExpArr := strings.Split(operation, "#")
	if len(optExpArr) < 3 {	//操作范围+具体操作+过期时间
		return
	}
	expired, err := strconv.ParseFloat(optExpArr[2], 64)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	now := float64(time.Now().UnixNano() / 1e6)
	if expired > now {
		indate = true
		return
	}
	return
}

func checkOperationScope(operationCmd string) (allowed bool) {

	allowed = false
	operationCmdArr := strings.Split(operationCmd, "")
	if len(operationCmdArr) != 3 {
		panic("Scope OperationCmd Error!")
	}
	//TODO:针对不同情况验证权限[还需要再想想]
	if operationCmdArr[2] == "x" {
		allowed = true
	}
	return
}

func getOperatedCompany(scope string) (company string) {
	tempArr := strings.Split(scope, ":")
	if len(tempArr) < 2 {
		panic("Error getOperatedCompany")
	}
	companyOptArr := strings.Split(tempArr[1], "#")
	company = strings.ToLower(companyOptArr[0])
	return
}