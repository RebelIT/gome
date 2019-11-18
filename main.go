package main

import (
	"flag"
	"github.com/rebelit/gome/iot/runners"
	"github.com/rebelit/gome/util/config"
	"github.com/rebelit/gome/web"
	"log"
	"net/http"
	"time"
)

func main() {
	flag.String("configFile", "/etc/apps/gome-server/configuration.json", "ConfigurationFile")
	flag.Parse()
	f := flag.Lookup("configFile")

	if err := config.LoadConfiguration(f.Value.String()); err != nil {
		log.Panic("PANIC: main, unable to load configuration\n")
	}

	time.Sleep(time.Second * 2)
	runners.Launch()
	start()
	return
}

func start() {
	log.Printf("INFO: main, starting gome http on :%v\n", config.App.HttpPort)
	router := listener.NewRouter()
	log.Fatal(http.ListenAndServe(":"+config.App.HttpPort, router))
}
