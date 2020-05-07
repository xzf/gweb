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
	"os"
	"os/signal"
	"syscall"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
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
		method(writer,request)
	})
	fmt.Println("http://127.0.0.1:2333")
	go http.ListenAndServe(addr, mux)
	ch := make(chan os.Signal)
	obj.SetKillFunc(func() {
		ch <- os.Kill
		logFunc("tk7khmjfwn normal kill")
	})
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
