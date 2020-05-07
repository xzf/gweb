/*
 * Author : xzf
 * Time    : 2020-05-05 01:38:54
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
)

var webApiMethodMap = map[string]struct{}{
	"WriteBody":    {},
	"GetGoRequest": {},
	"SetWriter":    {},
	"SetKillFunc":  {},
	"Kill":         {},
}

func (api *WebApi) WriteBody(body []byte) (ok bool) {
	writer, ok := api.getWriter()
	if !ok {
		return
	}
	_, err := writer.Write(body)
	ok = err == nil
	if !ok {
		logFunc("3rljce5cf write body failed " + err.Error())
	}
	return
}

func (api *WebApi) GetGoRequest() (req *http.Request, ok bool) {
	api.httpCtxLock.Lock()
	defer api.httpCtxLock.Unlock()
	if api.httpCtxMap == nil {
		api.httpCtxMap = map[string]*httpCtx{}
		return
	}
	ctx, ok := api.httpCtxMap[getGoroutineId()]
	if !ok {
		return
	}
	req = ctx.request
	return
}

func (api *WebApi) SetWriter(w http.ResponseWriter, req *http.Request) {
	if api.id == "" {
		api.id = "233"
	}
	api.httpCtxLock.Lock()
	defer api.httpCtxLock.Unlock()
	if api.httpCtxMap == nil {
		api.httpCtxMap = map[string]*httpCtx{}
	}
	//todo test goroutine id can repeat
	api.httpCtxMap[getGoroutineId()] = &httpCtx{
		writer:  w,
		request: req,
	}
}

func (api *WebApi) SetKillFunc(f func()) {
	api.killFunc = f
}

func (api *WebApi) Kill() {
	if api.killFunc != nil{
		api.killFunc()
	}
}

func (api *WebApi) getWriter() (w http.ResponseWriter, ok bool) {
	api.httpCtxLock.Lock()
	defer api.httpCtxLock.Unlock()
	if api.httpCtxMap == nil {
		api.httpCtxMap = map[string]*httpCtx{}
		debugLog("mso7435o5 api.httpCtxMap == nil")
		debugLog("mso7435o5 api.httpCtxMap == nil")
		debugLog("mso7435o5 api.httpCtxMap == nil")
		debugLog("mso7435o5 api.httpCtxMap == nil")
		return
	}
	ctx, ok := api.httpCtxMap[getGoroutineId()]
	if !ok {
		logFunc("khvexckic can not found httpCtx")
		return
	}
	w = ctx.writer
	return
}
