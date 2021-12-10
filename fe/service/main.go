package main

import (
	"fmt"
	"net/http"
)

type FeService struct {
	config map[string]string
}

func (feService *FeService) Build() {
	feService.config = map[string]string{}
	feService.config["httpPort"] = "10001"
}

func (feService *FeService) Run() {
	// 管理静态文件目录
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Http service listened at : " + feService.config["httpPort"])
	// 启动HTTP服务
	err := http.ListenAndServe("0.0.0.0:"+feService.config["httpPort"], nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	feService := FeService{}
	feService.Build()
	go feService.Run()

	c := make(chan bool, 1)
	<-c
}
