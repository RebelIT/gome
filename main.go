package main

import (
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/runner"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
	"time"
)

func main(){

	listenOn := "6660"
	devices.LoadDevices()
	time.Sleep(time.Second *2)

	go runner.GoGoDeviceStatus()
	time.Sleep(time.Second *5)

	go runner.GoGoScheduler()
	time.Sleep(time.Second *5)

	go runner.GoGoSQS()

	start(listenOn)

	return
}

func start(listenOn string){
	log.Printf("Starting Web Listener on :%v\n", listenOn)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+listenOn, router))
}