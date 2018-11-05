package main

import (
	"github.com/rebelit/gome/runner"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
)

func main(){
	listenOn := "6661"
	go runner.GoGODeviceLoader()
	go runner.GoGoRunners()
	start(listenOn)

	return
}

func start(listenOn string){
	log.Printf("Starting Web Listener on :%v", listenOn)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+listenOn, router))
}