package gweb

import (
	"reflect"
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"
	"runtime"
	"strings"
)

var logFunc func(string)

func debugLog(in ...interface{}) {
	if isDebug {
		fmt.Println(in...)
	}
}
func parseWebApiObj(obj interface{}) map[string]func(samplePara) {
	//todo if obk not a ptr ,define a var to get ptr for obj
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	debugLog("xv58t36dw num of method", objType.NumMethod())
	methodMap := map[string]func(samplePara){}
	for i := 0; i < objType.NumMethod(); i++ {
		valueMethod := objValue.Method(i)
		typeMethod := objType.Method(i)
		var in []reflect.Value
		methodType := valueMethod.Type()
		switch valueMethod.Type().NumIn() {
		case 0:
			methodMap[typeMethod.Name] = func(para samplePara) {
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
			jsonApiFunc := func(code string, para samplePara) {
				if len(para.Body) == 0 {
					para.writer.WriteHeader(400)
					//do not need header, or you can all json header check
					logFunc(code + " need request json body [ " + typeMethod.Name + " ]")
					return
				}
				newObjValue := reflect.New(inType).Elem()
				newObjValuePtr := newObjValue.Addr().Interface()
				err := json.Unmarshal(para.Body, newObjValuePtr)
				//can be replace to a high qps json decoder
				if err != nil {
					para.writer.WriteHeader(400)
					logFunc(code + " [ " + typeMethod.Name + " ] json Unmarshal error " + err.Error())
					return
				}
				in = append(in, reflect.ValueOf(newObjValuePtr).Elem())
				valueMethod.Call(in)
				return
			}
			if structFlag {
				//read from body, is body empty,return para error and code 400
				methodMap[typeMethod.Name] = func(para samplePara) {
					jsonApiFunc("e60c6nbnw", para)
				}
				continue
			}
			// on level struct
			// also can set by json
			methodMap[typeMethod.Name] = func(para samplePara) {
				newObjValue := reflect.New(inType).Elem()
				//newObjValuePtr := newObjValue.Addr().Interface()
				setFlag := false
				for i := 0; i < inType.NumField(); i++ {
					field := inType.Field(i)
					if field.PkgPath != "" {
						continue
					}
					fieldType := inType.Field(i).Type
					kind := fieldType.Kind()
					//can continue
					if kind != reflect.Array && kind != reflect.Slice {
						var strValue string
						if para.httpMethod == http.MethodGet {
							strValue = para.QSingleMap[field.Name]
							if strValue == "" {
								strValue = para.PostSingleMap[field.Name]
							}
						}
						if para.httpMethod == http.MethodPost {
							strValue = para.PostSingleMap[field.Name]
							if strValue == "" {
								strValue = para.QSingleMap[field.Name]
							}
						}
						if strValue == "" {
							continue
						}
						getKindErrStr := func(code string) string {
							return code + " " + kind.String() + " field [" + field.Name + "] value wrong [ " + strValue + " ]"
						}
						switch kind {
						case reflect.String:
							newObjValue.Field(i).SetString(strValue)
						case reflect.Bool:
							boolValue, err := strconv.ParseBool(strValue)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc(getKindErrStr("hh53qwl08"))
								return
							}
							newObjValue.Field(i).SetBool(boolValue)
						case reflect.Float64, reflect.Float32:
							fv := 64
							if kind == reflect.Float32 {
								fv = 32
							}
							floatValue, err := strconv.ParseFloat(strValue, fv)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc(getKindErrStr("via1vh6ky"))
								return
							}
							newObjValue.Field(i).SetFloat(floatValue)
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							tmpIntValue, err := strconv.ParseInt(strValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc(getKindErrStr("8914tszqn"))
								return
							}
							newObjValue.Field(i).SetInt(tmpIntValue)
							//switch kind {
							//	case reflect.Int:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int(tmpIntValue)))
							//	case reflect.Int8:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int8(tmpIntValue)))
							//	case reflect.Int16:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int16(tmpIntValue)))
							//	case reflect.Int32:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int32(tmpIntValue)))
							//	case reflect.Int64:
							//		newObjValue.Field(i).SetInt(tmpIntValue)
							//}
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							tmpUintValue, err := strconv.ParseUint(strValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc(getKindErrStr("8914tszqn"))
								return
							}
							newObjValue.Field(i).SetUint(tmpUintValue)
							//switch kind {
							//	case reflect.Uint:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int(tmpUintValue)))
							//	case reflect.Uint8:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int8(tmpUintValue)))
							//	case reflect.Uint16:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int16(tmpUintValue)))
							//	case reflect.Uint32:
							//		newObjValue.Field(i).Set(reflect.ValueOf(int32(tmpUintValue)))
							//	case reflect.Uint64:
							//		newObjValue.Field(i).SetUint(tmpUintValue)
							//}
						}
						setFlag = true
						continue
					}
					var strValueSlice []string
					if para.httpMethod == http.MethodGet {
						strValueSlice = para.QMultiMap[field.Name]
						if len(strValueSlice) == 0 {
							strValueSlice = para.PostMultiMap[field.Name]
						}
					}
					if para.httpMethod == http.MethodPost {
						strValueSlice = para.PostMultiMap[field.Name]
						if len(strValueSlice) == 0 {
							strValueSlice = para.QMultiMap[field.Name]
						}
					}
					if len(strValueSlice) == 0 {
						continue
					}
					isSlice := kind == reflect.Slice
					switch fieldType.Elem().Kind() {
					case reflect.String:
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(strValueSlice))
						} else {
							for index, oneStrValue := range strValueSlice {
								newObjValue.Field(i).Index(index).SetString(oneStrValue)
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
								para.writer.WriteHeader(400)
								logFunc("lewqa17bi" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								boolSlice = append(boolSlice, boolValue)
							} else {
								newObjValue.Field(i).Index(index).SetBool(boolValue)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(boolSlice))
						}
					case reflect.Float64:
						var float64Slice []float64
						if isSlice {
							float64Slice = make([]float64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							float64Value, err := strconv.ParseFloat(oneStrValue, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("ssex5xsb7" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								float64Slice = append(float64Slice, float64Value)
							} else {
								newObjValue.Field(i).Index(index).SetFloat(float64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(float64Slice))
						}
					case reflect.Float32:
						var float32Slice []float32
						if isSlice {
							float32Slice = make([]float32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							float32Value, err := strconv.ParseFloat(oneStrValue, 32)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("boydstipm" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								float32Slice = append(float32Slice, float32(float32Value))
							} else {
								newObjValue.Field(i).Index(index).SetFloat(float32Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(float32Slice))
						}
					case reflect.Int:
						var intSlice []int
						if isSlice {
							intSlice = make([]int, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("xsd96lzpo" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								intSlice = append(intSlice, int(int64Value))
							} else {
								newObjValue.Field(i).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(intSlice))
						}
					case reflect.Int8:
						var int8Slice []int8
						if isSlice {
							int8Slice = make([]int8, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("v8jd6ehdu" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int8Slice = append(int8Slice, int8(int64Value))
							} else {
								newObjValue.Field(i).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(int8Slice))
						}
					case reflect.Int16:
						var int16Slice []int16
						if isSlice {
							int16Slice = make([]int16, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("v6962e3rf" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int16Slice = append(int16Slice, int16(int64Value))
							} else {
								newObjValue.Field(i).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(int16Slice))
						}
					case reflect.Int32:
						var int32Slice []int32
						if isSlice {
							int32Slice = make([]int32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("z2uqrori7" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int32Slice = append(int32Slice, int32(int64Value))
							} else {
								newObjValue.Field(i).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(int32Slice))
						}
					case reflect.Int64:
						var int64Slice []int64
						if isSlice {
							int64Slice = make([]int64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							int64Value, err := strconv.ParseInt(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("np7sytoo8" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								int64Slice = append(int64Slice, int64Value)
							} else {
								newObjValue.Field(i).Index(index).SetInt(int64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(int64Slice))
						}
					case reflect.Uint:
						var uintSlice []uint
						if isSlice {
							uintSlice = make([]uint, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("9ud8qgwdg" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uintSlice = append(uintSlice, uint(uint64Value))
							} else {
								newObjValue.Field(i).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(uintSlice))
						}
					case reflect.Uint8:
						var uint8Slice []uint8
						if isSlice {
							uint8Slice = make([]uint8, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("tzmn1urh9" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint8Slice = append(uint8Slice, uint8(uint64Value))
							} else {
								newObjValue.Field(i).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(uint8Slice))
						}
					case reflect.Uint16:
						var uint16Slice []uint16
						if isSlice {
							uint16Slice = make([]uint16, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("i2wwwykyz" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint16Slice = append(uint16Slice, uint16(uint64Value))
							} else {
								newObjValue.Field(i).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(uint16Slice))
						}
					case reflect.Uint32:
						var uint32Slice []uint32
						if isSlice {
							uint32Slice = make([]uint32, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("4320zcfy6" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint32Slice = append(uint32Slice, uint32(uint64Value))
							} else {
								newObjValue.Field(i).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(uint32Slice))
						}
					case reflect.Uint64:
						var uint64Slice []uint64
						if isSlice {
							uint64Slice = make([]uint64, len(strValueSlice))
						}
						for index, oneStrValue := range strValueSlice {
							uint64Value, err := strconv.ParseUint(oneStrValue, 10, 64)
							if err != nil {
								para.writer.WriteHeader(400)
								logFunc("uuz5pmohk" + " " + kind.String() + " field [" + field.Name + "] one value wrong [ " + oneStrValue + " ]")
								return
							}
							if isSlice {
								uint64Slice = append(uint64Slice, uint64Value)
							} else {
								newObjValue.Field(i).Index(index).SetUint(uint64Value)
							}
						}
						if isSlice {
							newObjValue.Field(i).Set(reflect.ValueOf(uint64Slice))
						}
					}
					setFlag = true
				}
				if !setFlag {
					jsonApiFunc("i4wqyxplp", para)
				}
				return
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

func (para samplePara) GetOneString(method string) {
	if method == http.MethodGet {
		if len(para.QMultiMap) != 0 {

			return
		}
		return
	}
}

//unsupport type :Func,Interface,byte、complex、chan、ptr and their slice type
var unSupportKindSlice = []reflect.Kind{ //todo maybe map faster
	reflect.Chan,
	//reflect.Slice, //todo check item type
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

/*
for i := 0; i < methodNum; i++ {
		methodName := objType.Method(i).Name
		//method := objValue.Method(i)
		//for ii:=0;i<objType.NumIn();ii++{
		//	inType:=objType.In(ii)
		//}
		//para := reflect.New(inPara).Type()
		//logFunc(methodName, " | ", method.Type.In(1))

		methodMap[methodName] = func(para samplePara) {
			//reflect.New(objType)
			//method.Call([]reflect.Value{})
		}
		//logFunc(objType.Method(i).Type.NumIn())
		//inPara := objType.Method(i).Type.In(1)
		//for j := 0; j < inPara.NumField(); j++ {
		//	logFunc(methodName, inPara.Name(), inPara.Field(j).Name, inPara.Field(j).Type)
		//}
		//remove  methods for WebApi
		//if methodName == "" {
		//	continue
		//}
		//method := objValue.Method(i)

		//logFunc(objValue.Method(i).String())
	}
*/



func getGoroutineId() string {
	var buf [128]byte
	runtime.Stack(buf[:], false)
	allInfo := string(buf[10:])
	n := strings.Index(allInfo, " ")
	return allInfo[:n]
}