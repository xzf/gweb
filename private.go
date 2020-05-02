package gweb

import (
	"reflect"
	"fmt"
)

func ParseWebApiObj(obj interface{}) {
	//todo if obk not a ptr ,define a var to get ptr for obj
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	var objPtr = obj
	if objValue.Kind() != reflect.Ptr {
		objPtr = &obj
		objValue = reflect.ValueOf(objPtr)
	}
	methodNum := objType.NumMethod()
	methodMap := map[string]func(samplePara){}

	//samplePara to func in
	for i := 0; i < objValue.NumMethod(); i++ {
		valueMethod := objValue.Method(i)
		typeMethod := objType.Method(i)
		var in []reflect.Value
		methodType := valueMethod.Type()
		switch valueMethod.Type().NumIn() {
		case 0:
			methodMap[typeMethod.Name] = func(para samplePara) {
				valueMethod.Call([]reflect.Value{})
			}
			continue
		case 1:
			inType := methodType.In(0)
			panicIfTypeUnsupport(inType)
			switch inType.Kind() {
			case reflect.String:
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			case reflect.Float32,reflect.Float64:
			case reflect.Bool:
			case reflect.Array,reflect.Slice:
			case reflect.Struct:
				for iti := 0; iti < inType.NumField(); iti++ {
					itf := inType.Field(iti)
					fmt.Println(inType.Name(), itf.Name, itf.Type)
				}
			default:
				unsupportKindPanic(inType.Kind())
			}

		default:
			for ii := 0; ii < methodType.NumIn(); ii++ {
				inType := methodType.In(ii)
				panicIfTypeUnsupport(inType)

				//p:=reflect.New(tt)

				in = append(in, )
				fmt.Println(inType.Name())
			}
		}

		//method.Call([]reflect.Value{
		//
		//})
	}
	return
	for i := 0; i < methodNum; i++ {
		methodName := objType.Method(i).Name
		//method := objValue.Method(i)
		//for ii:=0;i<objType.NumIn();ii++{
		//	inType:=objType.In(ii)
		//}
		//para := reflect.New(inPara).Type()
		//fmt.Println(methodName, " | ", method.Type.In(1))

		methodMap[methodName] = func(para samplePara) {
			//reflect.New(objType)
			//method.Call([]reflect.Value{})
		}
		fmt.Println(objType.Method(i).Type.NumIn())
		inPara := objType.Method(i).Type.In(1)
		for j := 0; j < inPara.NumField(); j++ {
			fmt.Println(methodName, inPara.Name(), inPara.Field(j).Name, inPara.Field(j).Type)
		}
		//remove  methods for WebApi
		//if methodName == "" {
		//	continue
		//}
		//method := objValue.Method(i)

		//fmt.Println(objValue.Method(i).String())
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
	for i := 0; i < 3; i++ {
		fmt.Println("v408h4t8d unsupport type", kind)
	}
	panic("v408h4t8d unsupport type " + kind.String())

}
