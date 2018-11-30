package main

import (
	"github.com/rebelit/gome/runner"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
)

func main(){
	listenOn := "6661"
	runner.GoGODeviceLoader()
	//go runner.GoGoRunners()
	go runner.GoGoScheduler()
	start(listenOn)

	return
}

func start(listenOn string){
	log.Printf("Starting Web Listener on :%v", listenOn)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+listenOn, router))
}