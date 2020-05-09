/*
 * Author : xzf
 * Time    : 2020-05-05 01:38:54
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"io/ioutil"
)

var webApiMethodMap = map[string]struct{}{
	"WriteBody":    {},
	"GetGoRequest": {},
	"SetWriter":    {},
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

func (api *WebApi) GetFileSlice() (fileSlice []fileInfo) {
	req, ok := api.GetGoRequest()
	if !ok {
		logFunc("sjw6f4qbm can't get http.Request")
		return
	}
	debugLog("u7z5h5z74",req.MultipartForm)
	debugLog("tz3lrxdnh",req.MultipartForm.File)
	for _, fhSlice := range req.MultipartForm.File {
		for _, fh := range fhSlice {
			f, err := fh.Open()
			if err != nil {
				logFunc("wos0ysbkb file open err : " + err.Error())
				return
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				logFunc("gsvffvedz file read err : " + err.Error())
				return
			}
			fileSlice = append(fileSlice, fileInfo{
				Content:  content,
				FileName: fh.Filename,
			})
		}
	}
	return
}