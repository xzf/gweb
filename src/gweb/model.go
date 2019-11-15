package gweb

import (
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func NewHttpServer(addr string, obj interface{}) {
	t := reflect.TypeOf(obj)
	methodNum := t.NumMethod()
	for i := 0; i < methodNum; i++ {
		//去掉WebApi的method
	}
}

type WebApi struct {
}

func (w WebApi) ServeHTTP(resp http.ResponseWriter, req *http.Request) {}

func WaitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}