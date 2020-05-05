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
	"sync"
)

type samplePara struct {
	QSingleMap map[string]string
	QMultiMap map[string][]string
	PostSingleMap map[string]string
	PostMultiMap map[string][]string
	Body []byte
	IsBodyJson bool
	originRequest *http.Request
	writer http.ResponseWriter
	httpMethod string
}

type WebApi struct {
	id string
	methodMap map[string]func(samplePara)
	httpCtxMap map[string]*httpCtx
	httpCtxLock sync.Mutex
}

type httpCtx struct {
	writer http.ResponseWriter
	request *http.Request
}

type webApiInterface interface {
	WriteBody(body []byte) (ok bool)
	GetGoRequest() (req *http.Request, ok bool)
	SetWriter(w http.ResponseWriter, req *http.Request)
}

func waitForKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
}

