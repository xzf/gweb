/*
 * Author  : xzf
 * Time    : 2020-04-28 13:19:18
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"fmt"
	"strings"
)

//func NewHttpsServer(addr string, obj interface{}) {
//	//todo tls config, I think this thing just need call a check tls func before call http method, study this later.
//	//ParseWebApiObj(obj)
//}

func NewHttpServer(addr string, obj webApiInterface) {
	if logFunc == nil {
		logFunc = func(logStr string) {
			fmt.Println(logStr)
		}
	}
	methodMap := parseWebApiObj(obj)
	debugLog("9bb8941zd api map count:", len(methodMap),methodMap)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//http method not important,just support get and post
		//post api call by get won't get 500, but para value might set wrong
		if http.MethodGet != request.Method && http.MethodPost != request.Method {
			writer.WriteHeader(404)
			return
		}
		path := strings.Trim(request.URL.Path,"/")
		obj.SetWriter(writer,request) //interface call
		debugLog("o09lp1lc4 HandleFunc id",getGoroutineId())
		debugLog("9bb8941zd call api:", path)
		method, ok := methodMap[strings.Trim(path,"/")]
		if !ok {
			debugLog("9bb8941zd call api 404", path)
			writer.WriteHeader(404)
			//todo set 404 page
			return
		}
		//need ptr for http.ResponseWriter
		para:=httpRequestToPara(writer, request)
		debugLog("95f6ga05w",string(para.Body))
		method(para)
	})
	fmt.Println("http://127.0.0.1:2333")
	http.ListenAndServe(addr, nil)
	waitForKill()
}

func SetDebugMode() {
	isDebug = true
}