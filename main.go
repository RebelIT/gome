package main

import (
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/runners"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
	"time"
)

func main(){

	listenOn := "6660"
	devices.LoadDevices()
	time.Sleep(time.Second *2)

	runners.Launch()

	start(listenOn)

	return
}

func start(listenOn string){
	log.Printf("Starting Web Listener on :%v\n", listenOn)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+listenOn, router))
}