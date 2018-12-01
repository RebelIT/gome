package main

import (
	"github.com/rebelit/gome/runner"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
	"time"
)

func main(){
	listenOn := "6661"
	runner.GoGODeviceLoader()
	time.Sleep(time.Second *3)

	go runner.GoGoRunners()
	time.Sleep(time.Second *5)

	go runner.GoGoScheduler()

	start(listenOn)

	return
}

func start(listenOn string){
	log.Printf("Starting Web Listener on :%v\n", listenOn)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+listenOn, router))
}