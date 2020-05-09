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
	"path/filepath"
	"io/ioutil"
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
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//http method not important,just support get and post
		//post api call by get won't get 500, but para value might set wrong
		if http.MethodGet != request.Method && http.MethodPost != request.Method {
			debugLog("pnj27n9nf call api "+request.Method)
			writer.WriteHeader(404)
			return
		}
		path := strings.Trim(request.URL.Path,"/")
		debugLog("pnj27n9nf call api "+path)
		req.Obj.SetWriter(writer,request) //interface call
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
	if req.FileRootPath != "" && req.FilePathPrefix != "" {
		filePathPrefix := "/" + req.FilePathPrefix + "/"
		mux.HandleFunc(filePathPrefix, func(writer http.ResponseWriter, request *http.Request) {
			tmpPath := strings.Replace(request.URL.Path, filePathPrefix, "", -1)
			absPath := filepath.Join(req.FileRootPath, tmpPath)
			//todo might have absPath err, if pop out ,fix it.
			file, err := os.Open(absPath)
			if err != nil {
				if !os.IsExist(err) {
					writer.WriteHeader(404)
					logFunc("jp6pngn65 file [ " + absPath + " ] no exist")
					return
				}
				writer.WriteHeader(500)
				logFunc("frsvv8h0m file [ " + absPath + " ] open error : " + err.Error())
				return
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				writer.WriteHeader(500)
				logFunc("ryn6e6kgz file [ " + absPath + " ] read error : " + err.Error())
				return
			}
			_, err = writer.Write(content)
			if err != nil {
				writer.WriteHeader(500)
				logFunc("1bumlw8m0 file [ " + absPath + " ] write error : " + err.Error())
				return
			}
		})
	}
	if req.Crt == "" && req.Key == "" {
		go func() {
			err := http.ListenAndServe(req.Addr, mux)
			if err!=nil{
				logFunc("zw2otslp6 ListenAndServe failed error : "+err.Error())
				return
			}
			fmt.Println("start http server success")
		}()
		return
	}
	if req.Crt != "" && req.Key != "" {
		go func() {
			err := http.ListenAndServeTLS(req.Addr, req.Crt, req.Key, mux)
			if err!=nil{
				logFunc("zw2otslp6 ListenAndServe failed error : "+err.Error())
				return
			}
			fmt.Println("start http server success")
		}()
		return
	}
	panic("Crt and Key must be both empty or both not empty")
}

func WaitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
