/*
 * Author  : xzf
 * Time    : 2020-04-28 13:17:40
 * Email   : xpoony@163.com
 */

package gweb

import (
	"net/http"
	"io/ioutil"
)

//type samplePara struct {
//	QSingleMap map[string]string
//	QMultiMap map[string][]string
//	PostSingleMap map[string]string
//	PostMultiMap map[string][]string
//	Body []byte
//	IsBodyJson bool
//	originRequest *http.Request
//	writer http.ResponseWriter
//	httpMethod string
//}
func httpRequestToPara(writer http.ResponseWriter, request *http.Request) (para samplePara) {
	para = samplePara{
		QSingleMap:    map[string]string{},
		QMultiMap:     map[string][]string{},
		PostSingleMap: map[string]string{},
		PostMultiMap:  map[string][]string{},
		writer:        writer,
		httpMethod:    request.Method,
	}
	queryMap := request.URL.Query()
	if len(queryMap) != 0 {
		for key, slice:= range queryMap {
			switch len(slice) {
			case 0:
				para.QSingleMap[key] = ""
			case 1:
				para.QSingleMap[key] =slice[0]
			default:
				para.QMultiMap[key] = slice
			}
		}
	}
	body, err := ioutil.ReadAll(request.Body)
	if err == nil && len(body) != 0 {
		para.Body = body
	}
	err = request.ParseForm()
	postForm := request.PostForm
	if err == nil && len(postForm) != 0 {
		for key, slice:= range queryMap {
			switch len(slice) {
			case 0:
				para.PostSingleMap[key] = ""
			case 1:
				para.PostSingleMap[key] =slice[0]
			default:
				para.PostMultiMap[key] = slice
			}
		}
	}
	return
}
