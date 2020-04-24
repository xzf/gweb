package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func NewHttpServer(addr string, obj interface{}) { //obj 必须用指针
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr {
		panic("obj need ptr")
	}
	methodNum := objType.NumMethod()
	//methodMap := map[string]func(req http.Request, respResp http.ResponseWriter){}
	methodMap := map[string]func(req http.Request, respResp http.ResponseWriter){}
	for i := 0; i < methodNum; i++ {
		methodName := objType.Method(i).Name
		method := objValue.Method(i)
		//para := reflect.New(inPara).Type()
		//fmt.Println(methodName, " | ", method.Type.In(1))

		methodMap[methodName] = func(req http.Request, respResp http.ResponseWriter) {
			reflect.New(objType)
			method.Call([]reflect.Value{})
		}
		fmt.Println(objType.Method(i).Type.NumIn())
		inPara := objType.Method(i).Type.In(1)
		for j := 0; j < inPara.NumField(); j++ {
			fmt.Println(methodName, inPara.Name(), inPara.Field(j).Name, inPara.Field(j).Type)
		}

		//去掉WebApi的method
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

func WaitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}
