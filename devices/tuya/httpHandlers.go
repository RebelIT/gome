package tuya

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"log"
	"net/http"
)

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
		devices.ReturnBad(w,r)
		return
	}

	if err := PowerControl(deviceName, action); err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r, http.Response{})
	return
}