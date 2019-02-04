package rpi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleControl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	component := vars["component"]

	i := PiControl{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		devices.ReturnBad(w, r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &i); err != nil {
		devices.ReturnInternalError(w, r)
		return
	}

	uri, err := compileUrl(component, i)
	if err != nil{
		devices.ReturnBad(w, r)
		return
	}

	resp, err := PiPost(deviceName,uri)
	if err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w, r)
		return
	}

	devices.ReturnOk(w, r, resp)
	return
}