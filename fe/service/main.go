package main

import (
	feService "github.com/stevewooo/testin/fe/service/feService"
)

func main() {
	feService := feService.FeService{}
	feService.Build()
	go feService.Run()

	c := make(chan bool, 1)
	<-c
}
