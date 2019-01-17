package rpi

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
)

func HandleControl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	action := vars["action"]

	if err := piPost(deviceName,action); err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}