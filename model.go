/*
 * Author : xzf
 * Time    : 2020-04-26 00:50:22
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"sync"
)

type WebApi struct {
	id string
	httpCtxMap map[string]*httpCtx
	httpCtxLock sync.Mutex
	killFunc func()
}

type httpCtx struct {
	writer http.ResponseWriter
	request *http.Request
}

type webApiInterface interface {
	WriteBody(body []byte) (ok bool)
	GetGoRequest() (req *http.Request, ok bool)
	SetWriter(w http.ResponseWriter, req *http.Request)
	SetKillFunc(func())
	Kill()
}


