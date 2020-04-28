/*
 * Author  : xzf
 * Time    : 2020-04-28 13:19:18
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
)

func NewHttpsServer(addr string, obj interface{}) {
	//todo tls config, I think this thing just need call a check tls func before call http method, study this later.
	//ParseWebApiObj(obj)
}

func NewHttpServer(addr string, obj interface{}) {
	methodMap := map[string]func(writer http.ResponseWriter, request *http.Request){}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		method, ok := methodMap[path]
		if !ok {
			writer.WriteHeader(404)
			//todo set 404 page
			return
		}
		//need ptr for http.ResponseWriter
		method(writer, request)
	})
	http.ListenAndServe(addr, nil)
	waitForKill()
}

