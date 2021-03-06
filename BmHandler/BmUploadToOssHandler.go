package BmHandler

import (
	"encoding/json"
	"fmt"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmSingleton"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/blackmirror/jsonapi/jsonapiobj"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/go-uuid"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type UploadToOssHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h UploadToOssHandler) NewUploadToOssHandler(args ...interface{}) UploadToOssHandler {
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
	return UploadToOssHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

func (h UploadToOssHandler) UploadToOss(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	fmt.Println("method:", "UploadToOssHandler", r.Method)
	w.Header().Add("Content-Type", "application/json")
	if r.Method == "GET" {
		errMsg := "upload request method error, please use POST."
		panic(errMsg)
		return 0
	} else {
		r.ParseMultipartForm(32 << 20)
		//file, handler, err := r.FormFile("file")
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			errMsg := "upload file key error, please use key 'file'."
			panic(errMsg)
			return 0
		}
		defer file.Close()

		var bmRouter BmConfig.BmRouterConfig
		bmRouter.GenerateConfig(BmSingleton.EnvHome)

		fn, err := uuid.GenerateUUID()
		if err != nil {
			fmt.Println(err)
			errMsg := "upload file key error, please use key 'file'."
			panic(errMsg)
			return 0
		}
		lsttmp := strings.Split(handler.Filename, ".")
		exname := lsttmp[len(lsttmp)-1]

		localDir := bmRouter.TmpDir +"/"+fn + "." + exname // handler.Filename
		f, err := os.OpenFile(localDir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("OpenFile error")
			fmt.Println(err)
			errMsg := "upload local file open error."
			panic(errMsg)
			return 0
		}
		
		defer os.Remove(localDir)
		io.Copy(f, file)
		result := map[string]string{
			//"file": handler.Filename,
			"file": fn,
		}
		accept:=r.Form["accept"]
		fnm:=accept[0]+"/"+fn
		des:=r.Form["des"]
		desc := ""
		if len(des)>0{
			desc = des[0]
		}
		client, err := oss.New(h.Args[1], h.Args[2], h.Args[3])
		if err != nil {
			// HandleError(err)
			panic("密钥出错")
		}
	
		bucket, err := client.Bucket("pharbers-max-bi")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		err = bucket.PutObjectFromFile(fnm, localDir)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		Size,_:=f.Seek(0, os.SEEK_END)
		sizestr:=strconv.FormatInt(Size/1024,10)+"MB"

		t := time.Now()
		tmp := t.Format("2006-01")
		
		filename:=lsttmp[0]
		bmfile := BmModel.Files{
			Name : filename,
			UploadTime : tmp,
			Describe : desc,
			Accept : accept[0],
			Uuid  : fn ,
			Size :  sizestr,
			Type : exname,
		}
		h.db.InsertBmObject(&bmfile)
		response := map[string]interface{}{
			"status": "ok",
			"result": result,
			"error":  "",
		}
		jso := jsonapiobj.JsResult{}
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		f.Close()
		return 1
	}
}

func (h UploadToOssHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h UploadToOssHandler) GetHandlerMethod() string {
	return h.Method
}
