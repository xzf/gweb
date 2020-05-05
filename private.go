package gweb

import (
	"reflect"
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"
	"runtime"
	"strings"
	"io/ioutil"
)

var logFunc func(string)

func debugLog(in ...interface{}) {
	if isDebug {
		fmt.Println(in...)
	}
}
func parseWebApiObj(obj interface{}) map[string]func(http.ResponseWriter, *http.Request) {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	debugLog("xv58t36dw num of method", objType.NumMethod())
	methodMap := map[string]func(http.ResponseWriter, *http.Request){}
	for i := 0; i < objType.NumMethod(); i++ {
		valueMethod := objValue.Method(i)
		if valueMethod.Type().PkgPath() != "" {
			continue
		}
		typeMethod := objType.Method(i)
		_, ok := webApiMethodMap[typeMethod.Name]
		if ok { // method of  WebApi
			continue
		}
		var in []reflect.Value
		methodType := valueMethod.Type()
		switch valueMethod.Type().NumIn() {
		case 0:
			methodMap[typeMethod.Name] = func(http.ResponseWriter, *http.Request) {
				valueMethod.Call(in)
			}
			continue
		case 1:
			inType := methodType.In(0)
			//because golang reflect can get func para name, golang base type unsupport
			if inType.Kind() != reflect.Struct {
				for i := 0; i < 3; i++ {
					logFunc("v408h4t8d unsupport type " + inType.Kind().String())
				}
				panic("v408h4t8d unsupport type " + inType.Kind().String())
			}
			//check is multi level struct or not, panic if type unsupport
			structFlag := false
			for i := 0; i < inType.NumField(); i++ {
				field := inType.Field(i)
				if field.PkgPath != "" {
					continue
				}
				panicIfTypeUnsupport(field.Type)
				if field.Type.Kind() == reflect.Struct {
					structFlag = true
				}
			}
			newObjValue := reflect.New(inType).Elem()
			formValueSetToInSlice := func(formMap map[string][]string) (success bool) {
				in = []reflect.Value{}
				for key, strValueSlice := range formMap {
					if len(strValueSlice) == 0 {
						continue
					}
					field, ok := inType.FieldByName(key)
					if !ok {
						continue
					}
					if field.PkgPath != "" { // none public field
						continue
					}
					fieldType := field.Type
					kind := fieldType.Kind()
					if kind != reflect.Array && kind != reflect.Slice {
						if len(strValueSlice) != 1 {
							debugLog("kdmfudu5m warning path : [" + typeMethod.Name + "] para's field  [" + field.Name + "] has more than one value [ " + strings.Join(strValueSlice, ",") + " ]")
						}
						var strValue string
						strValue = strValueSlice[0]
						getKindErrStr := func(code string) string {
							return code + " " + kind.String() + " field [" + field.Name + "] value wrong [ " + strValue + " ]"
						}
						switch kind {
						case reflect.String:
							debugLog("4pd3r2k42",strValue,field.Name,field.Type,field.Index)
							newObjValue.Field(field.Index[0]).SetString(strValue)
						case reflect.Bool:
							boolValue, err := strconv.ParseBool(strValue)
							if err != nil {
								logFunc(getKindErrStr("hh53qwl08"))
								return
							}
							newObjValue.Field(field.Index[0]).SetBool(boolValue)
						case reflect.Float64, reflect.Float32:
							fv := 64
							if kind == reflect.Float32 {
								fv = 32
							}
							floatValue, err := strconv.ParseFloat(strValue, fv)
							if err != nil {
								logFunc(getKindErrStr("via1vh6ky"))
								return
							}
							newObjValue.Field(field.Index[0]).SetFloat(floatValue)
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							tmpIntValue, err := strconv.ParseInt(strValue, 10, 64)
							if err != nil {
								logFunc(getKindErrStr("8914tszqn"))
								return
							}
							newObjValue.Field(field.Index[0]).SetInt(tmpIntValue)
							//switch kind {
							//	case reflect.Int:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int(tmpIntValue)))
							//	case reflect.Int8:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int8(tmpIntValue)))
							//	case reflect.Int16:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int16(tmpIntValue)))
							//	case reflect.Int32:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int32(tmpIntValue)))
							//	case reflect.Int64:
							//		newObjValue.Field(field.Index[0]).SetInt(tmpIntValue)
							//}
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							tmpUintValue, err := strconv.ParseUint(strValue, 10, 64)
							if err != nil {
								logFunc(getKindErrStr("8914tszqn"))
								return
							}
							newObjValue.Field(field.Index[0]).SetUint(tmpUintValue)
							//switch kind {
							//	case reflect.Uint:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int(tmpUintValue)))
							//	case reflect.Uint8:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int8(tmpUintValue)))
							//	case reflect.Uint16:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int16(tmpUintValue)))
							//	case reflect.Uint32:
							//		newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int32(tmpUintValue)))
							//	case reflect.Uint64:
							//		newObjValue.Field(field.Index[0]).SetUint(tmpUintValue)
							//}
						}
						continue
					}
					//slice and array
					isSlice := kind == reflect.Slice
					switch fieldType.Elem().Kind() {
					case reflect.String:
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(strValueSlice))
						} else {
							for index, oneStrValue := range strValueSlice {
								newObjValue.Field(field.Index[0]).Index(index).SetString(oneStrValue)
							}
						}
					case reflect.Bool:
						var boolSlice []bool
						if isSlice {
							boolSlice = make([]bool, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							boolValue, err := strconv.ParseBool(oneStrValue)
							if err != nil {
								logFunc("lewqa17bi" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								boolSlice = append(boolSlice, boolValue)
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetBool(boolValue)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(boolSlice))
						}
					case reflect.Float64:
						var float64Slice []float64
						if isSlice {
							float64Slice = make([]float64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							float64Value, err := strconv.ParseFloat(oneStrValue, 64)
							if err != nil {
								logFunc("ssex5xsb7" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								float64Slice = append(float64Slice, float64Value)
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetFloat(float64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(float64Slice))
						}
					case reflect.Float32:
						var float32Slice []float32
						if isSlice {
							float32Slice = make([]float32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							float32Value, err := strconv.ParseFloat(oneStrValue, 32)
							if err != nil {
								logFunc("boydstipm" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								float32Slice = append(float32Slice, float32(float32Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetFloat(float32Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(float32Slice))
						}
					case reflect.Int:
						var intSlice []int
						if isSlice {
							intSlice = make([]int, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								logFunc("xsd96lzpo" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								intSlice = append(intSlice, int(int64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(intSlice))
						}
					case reflect.Int8:
						var int8Slice []int8
						if isSlice {
							int8Slice = make([]int8, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								logFunc("v8jd6ehdu" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int8Slice = append(int8Slice, int8(int64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int8Slice))
						}
					case reflect.Int16:
						var int16Slice []int16
						if isSlice {
							int16Slice = make([]int16, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								logFunc("v6962e3rf" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int16Slice = append(int16Slice, int16(int64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int16Slice))
						}
					case reflect.Int32:
						var int32Slice []int32
						if isSlice {
							int32Slice = make([]int32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								logFunc("z2uqrori7" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int32Slice = append(int32Slice, int32(int64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int32Slice))
						}
					case reflect.Int64:
						var int64Slice []int64
						if isSlice {
							int64Slice = make([]int64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								logFunc("np7sytoo8" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int64Slice = append(int64Slice, int64Value)
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(int64Slice))
						}
					case reflect.Uint:
						var uintSlice []uint
						if isSlice {
							uintSlice = make([]uint, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								logFunc("9ud8qgwdg" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uintSlice = append(uintSlice, uint(uint64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(uintSlice))
						}
					case reflect.Uint8:
						var uint8Slice []uint8
						if isSlice {
							uint8Slice = make([]uint8, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								logFunc("tzmn1urh9" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint8Slice = append(uint8Slice, uint8(uint64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(uint8Slice))
						}
					case reflect.Uint16:
						var uint16Slice []uint16
						if isSlice {
							uint16Slice = make([]uint16, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								logFunc("i2wwwykyz" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint16Slice = append(uint16Slice, uint16(uint64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(uint16Slice))
						}
					case reflect.Uint32:
						var uint32Slice []uint32
						if isSlice {
							uint32Slice = make([]uint32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								logFunc("4320zcfy6" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint32Slice = append(uint32Slice, uint32(uint64Value))
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(uint32Slice))
						}
					case reflect.Uint64:
						var uint64Slice []uint64
						if isSlice {
							uint64Slice = make([]uint64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								logFunc("uuz5pmohk" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint64Slice = append(uint64Slice, uint64Value)
							} else {
								newObjValue.Field(field.Index[0]).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(field.Index[0]).Set(reflect.ValueOf(uint64Slice))
						}
					}
				}
				newObjValuePtr := newObjValue.Addr().Interface()
				in = append(in, reflect.ValueOf(newObjValuePtr).Elem())
				success = true
				return
			}
			methodMap[typeMethod.Name] = func(writer http.ResponseWriter, request *http.Request) {
				contentType := request.Header.Get("Content-Type")
				switch {
				case contentType == "application/json":
					in = []reflect.Value{}
					body, err := ioutil.ReadAll(request.Body)
					if err != nil {
						logFunc("u5suypqlv request body read error : " + err.Error())
						writer.WriteHeader(400) //bad request
						return
					}
					newObjValuePtr := newObjValue.Addr().Interface()
					err = json.Unmarshal(body, newObjValuePtr)
					//can be replace to a high qps json decoder
					if err != nil {
						writer.WriteHeader(400)
						logFunc("fgl7vg91q [ " + typeMethod.Name + " ] json Unmarshal error " + err.Error())
						return
					}
					in = append(in, reflect.ValueOf(newObjValuePtr).Elem())
					valueMethod.Call(in)
					return
				//has todo
				case strings.HasPrefix(contentType, "multipart/form-data"):
					if structFlag {
						debugLog("zzjvz6l3y warning path : [" + typeMethod.Name + "] para has unsupport struct field,if set it's value will return  bad request(http code 400)")
					}
					err := request.ParseMultipartForm(4096) //todo need user define
					if err != nil {
						logFunc("2ag9oo1lt request ParseMultipartForm error : " + err.Error())
						writer.WriteHeader(400) //bad request
						return
					}
					ok := formValueSetToInSlice(request.MultipartForm.Value)
					if !ok {
						writer.WriteHeader(400) //bad request
						return
					}
					valueMethod.Call(in)
				//has todo
				//and para field can not be struct
				case contentType == "application/x-www-form-urlencoded":
					if structFlag {
						debugLog("zzjvz6l3y warning path : [" + typeMethod.Name + "] para has unsupport struct field,if set it's value will return  bad request(http code 400)")
					}
					if request.Method == http.MethodPost {
						err := request.ParseForm()
						if err != nil {
							logFunc("11tp4qrno request ParseForm error : " + err.Error())
							writer.WriteHeader(400) //bad request
							return
						}
						ok := formValueSetToInSlice(request.Form)
						if !ok {
							writer.WriteHeader(400) //bad request
							return
						}
						debugLog("puyk0y72b", len(in),in)
						valueMethod.Call(in)
						return
					}
					//todo get  need to parse body
					writer.WriteHeader(400) //bad request
					logFunc("4alafpbfr request header content type is [ application/x-www-form-urlencoded ] but http method is get, can not parse form")
					return
				default:
					logFunc("sfqn3u5to ContentType [ " + contentType + " ] wrong")
					writer.WriteHeader(400) //bad request
					return
				}
				valueMethod.Call(in)
			}
		default:
			logInfo := "jsd4cq0w2 do not support multi para"
			for i := 0; i < 3; i++ {
				logFunc(logInfo)
			}
			panic(logInfo)
		}
	}
	return methodMap
}

//unsupport type :Func,Interface,byte、complex、chan、ptr and their slice type
var unSupportKindSlice = []reflect.Kind{ //todo maybe map faster
	reflect.Chan,
	reflect.Ptr,
	reflect.Complex64,
	reflect.Complex64,
	reflect.Interface, //must assign real type
	reflect.Func,
	reflect.Map,
}

func panicIfTypeUnsupport(typ reflect.Type) {
	kind := typ.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		kind = typ.Elem().Kind()
	}
	for _, item := range unSupportKindSlice {
		if item == kind {
			unsupportKindPanic(kind)
		}
	}
}

func unsupportKindPanic(kind reflect.Kind) {
	str := "v408h4t8d unsupport type " + kind.String()
	for i := 0; i < 3; i++ {
		logFunc(str)
	}
	panic(str)
}

func getGoroutineId() string {
	var buf [128]byte
	runtime.Stack(buf[:], false)
	allInfo := string(buf[10:])
	n := strings.Index(allInfo, " ")
	return allInfo[:n]
}
