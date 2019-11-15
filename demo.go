package main

import (
	"fmt"
	"gweb"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8001", gweb.WebApi{})
	if err != nil {
		fmt.Println("81aiwu5pox", "http.ListenAndServe", err)
	}
	gweb.WaitForKill()
}
