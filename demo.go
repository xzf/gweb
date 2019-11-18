package main

import (
	"fmt"
	"gweb"
)

type web struct {
	gweb.WebApi
}

//func(a web)FuncA(){
//	fmt.Println(111)
//}

func main() {
	w := web{}
	fmt.Println("------------")
	gweb.NewHttpServer("", &w)
	fmt.Println("------------")
}

type AReq struct {
	As string
	Ai int
}

func (w web) A(req AReq) {
	fmt.Println(req)
}

type BReq struct {
	Bs string
	Bi int
}

func (w web) B(req BReq) {
	fmt.Println(req)
}
