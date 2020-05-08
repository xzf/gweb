/*
 * Author  : xzf
 * Time    : 2020-04-28 13:19:18
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
)

type NewHttpRequest struct {
	Addr           string
	Obj            webApiInterface
	Prefix         string
	FileRootPath   string //can be empty
	FilePathPrefix string //can be empty
	Crt            string //if empty and Key also empty mean http server
	Key            string //if empty and Crt also empty mean http server
}

func NewHttpServer(req NewHttpRequest) {
	if logFunc == nil {
		logFunc = func(logStr string) {
			fmt.Println(logStr)
		}
	}
	mux := http.NewServeMux()
	methodMap := parseWebApiObjToMethodMap(req.Obj)
	debugLog("9bb8941zd api map count:", len(methodMap), methodMap)
	for path, logicFunc := range methodMap {
		apiPath := "/" + path
		if req.Prefix != "" {
			apiPath = "/" + req.Prefix + "/" + path
		}
		mux.HandleFunc(apiPath, func(writer http.ResponseWriter, request *http.Request) {
			//http method not important,just support get and post
			//post api call by get won't get 500, but para value might set wrong
			if http.MethodGet != request.Method && http.MethodPost != request.Method {
				writer.WriteHeader(404)
				return
			}
			path := strings.Trim(request.URL.Path, "/")
			req.Obj.SetWriter(writer, request) //interface call
			debugLog("o09lp1lc4 HandleFunc id", getGoroutineId())
			debugLog("9bb8941zd call api:", path)
			logicFunc(writer, request)
		})
	}
	if req.FileRootPath != "" && req.FilePathPrefix != "" {
		mux.HandleFunc("/"+req.FilePathPrefix+"/", func(writer http.ResponseWriter, request *http.Request) {
			//todo file handler
			fmt.Println("file ", request.URL.Path)
		})
	}
	if req.Crt == "" && req.Key == "" {
		go http.ListenAndServe(req.Addr, mux)
		fmt.Println("start http server success")
		return
	}
	if req.Crt != "" && req.Key != "" {
		go http.ListenAndServeTLS(req.Addr, req.Crt, req.Key, mux)
		fmt.Println("start https server success")
		return
	}
	panic("Crt and Key must be both empty or both not empty")
}

func WaitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
