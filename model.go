/*
 * Author : xzf
 * Time    : 2020-04-26 00:50:22
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"reflect"
	"fmt"
)

func NewHttpsServer(addr string, obj interface{}) {
	//todo tls config, I think this thing just need call a check tls func before call http method, study this later.
	//ParseWebApiObj(obj)
}

type samplePara struct {
	QSingleMap map[string]string
	QMultiMap map[string][]string
	PostBody []byte
	IsBodyJson bool
}

func NewHttpServer(addr string, obj interface{}) {
	methodMap := map[string]func(writer *http.ResponseWriter, request *http.Request){}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		method, ok := methodMap[path]
		if !ok {
			writer.WriteHeader(404)
			//todo set 404 page
			return
		}
		//need ptr for http.ResponseWriter
		method(&writer, request)
	})
	http.ListenAndServe(addr, nil)
	waitForKill()
}

func ParseWebApiObj(obj interface{}) {
	//todo if obk not a ptr ,define a var to get ptr for obj
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	var objPtr = obj
	if objValue.Kind() != reflect.Ptr {
		objPtr = &obj
		objValue = reflect.ValueOf(objPtr)
	}
	methodNum := objType.NumMethod()
	methodMap := map[string]func(in ...interface{}){}
	for i := 0; i < objValue.NumMethod(); i++ {
		method := objValue.Method(i)
		var in []reflect.Value
		for ii := 0; ii < method.Type().NumIn(); ii++ {
			tt := method.Type().In(ii)
			switch tt.Kind() {
			case reflect.Struct:
				for iti := 0; iti < tt.NumField(); iti++ {
					itf := tt.Field(iti)
					fmt.Println(tt.Name(), itf.Name, itf.Type)

				}
			default:
				fmt.Println("unsupport method para type", tt.Kind().String())
				return
			}
			//p:=reflect.New(tt)

			in = append(in, )
			fmt.Println(tt.Name())
		}
		//method.Call([]reflect.Value{
		//
		//})
	}
	for i := 0; i < methodNum; i++ {
		methodName := objType.Method(i).Name
		method := objValue.Method(i)
		//for ii:=0;i<objType.NumIn();ii++{
		//	inType:=objType.In(ii)
		//}
		//para := reflect.New(inPara).Type()
		//fmt.Println(methodName, " | ", method.Type.In(1))

		methodMap[methodName] = func(in ...interface{}) {
			reflect.New(objType)
			method.Call([]reflect.Value{})
		}
		fmt.Println(objType.Method(i).Type.NumIn())
		inPara := objType.Method(i).Type.In(1)
		for j := 0; j < inPara.NumField(); j++ {
			fmt.Println(methodName, inPara.Name(), inPara.Field(j).Name, inPara.Field(j).Type)
		}
		//remove  methods for WebApi
		//if methodName == "" {
		//	continue
		//}
		//method := objValue.Method(i)

		//fmt.Println(objValue.Method(i).String())
	}
}

type WebApi struct {
	methodMap map[string]func(req http.Request, respResp http.ResponseWriter)
}

func waitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
