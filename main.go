package main

import (
	"fmt"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmFactory"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmMaxDefine"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmSingleton"
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
	// 本机测试，添加上
	os.Setenv(BmSingleton.EnvHome, ".")

	phHome := os.Getenv(BmSingleton.EnvHome)
	fac := BmFactory.BmTable{}
	var pod = BmMaxDefine.Pod{ Name: BmSingleton.ProjectName, Factory:fac }
	pod.RegisterSerFromYAML(phHome + "/resource/def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(BmSingleton.EnvHome)
	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(BmSingleton.Version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)
	pod.RegisterAllFunctions(BmSingleton.Version, api)
	pod.RegisterAllMiddleware(api)
	handler := api.Handler().(*httprouter.Router)
	http.ListenAndServe(":"+bmRouter.Port, handler)

	fmt.Println("pod archi ends")
}
// package main
// import (
//     "crypto/tls"
//     "fmt"
//     "log"
//     "net"
//     "net/smtp"
// )
// func main() {
//     host := "smtpdm.aliyun.com"
//     port := 465
//     email := "uploadtooss@uploadtooss.pharbers.com"
//     password := "Black84244216Mirror"
//     toEmail := "543187000@qq.com"
//     header := make(map[string]string)
//     header["From"] = "法伯科技" + "<" + email + ">"
//     header["To"] = toEmail
//     header["Subject"] = "重置密码通知"
//     header["Content-Type"] = "text/html; charset=UTF-8"
//     body := "www.bing.com"
//     message := ""
//     for k, v := range header {
//         message += fmt.Sprintf("%s: %s\r\n", k, v)
//     }
//     message += "\r\n" + body
//     auth := smtp.PlainAuth(
//         "",
//         email,
//         password,
//         host,
//     )
//     err := SendMailUsingTLS(
//         fmt.Sprintf("%s:%d", host, port),
//         auth,
//         email,
//         []string{toEmail},
//         []byte(message),
//     )
//     if err != nil {
//         panic(err)
//     } else {
//         fmt.Println("Send mail success!")
//     }
// }
// //return a smtp client
// func Dial(addr string) (*smtp.Client, error) {
//     conn, err := tls.Dial("tcp", addr, nil)
//     if err != nil {
//         log.Println("Dialing Error:", err)
//         return nil, err
//     }
//     //分解主机端口字符串
//     host, _, _ := net.SplitHostPort(addr)
//     return smtp.NewClient(conn, host)
// }
// //参考net/smtp的func SendMail()
// //使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// //len(to)>1时,to[1]开始提示是密送
// func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
//     to []string, msg []byte) (err error) {
//     //create smtp client
//     c, err := Dial(addr)
//     if err != nil {
//         log.Println("Create smpt client error:", err)
//         return err
//     }
//     defer c.Close()
//     if auth != nil {
//         if ok, _ := c.Extension("AUTH"); ok {
//             if err = c.Auth(auth); err != nil {
//                 log.Println("Error during AUTH", err)
//                 return err
//             }
//         }
//     }
//     if err = c.Mail(from); err != nil {
//         return err
//     }
//     for _, addr := range to {
//         if err = c.Rcpt(addr); err != nil {
//             return err
//         }
//     }
//     w, err := c.Data()
//     if err != nil {
//         return err
//     }
//     _, err = w.Write(msg)
//     if err != nil {
//         return err
//     }
//     err = w.Close()
//     if err != nil {
//         return err
//     }
//     return c.Quit()
// }