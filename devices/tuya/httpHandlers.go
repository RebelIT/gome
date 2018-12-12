package tuya

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

func HandleControl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["name"]
	state := vars["state"]

	action := false
	if state == "on"{
		action = true
	} else if state == "off"{
		action = false
	} else{
		log.Printf("[ERROR] %s : control %s, state %s not found", deviceName, r.Method, state)
		notify.MetricHttpIn(r.URL.Path, http.StatusBadRequest, r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := PowerControl(deviceName, action); err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}