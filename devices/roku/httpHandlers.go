package roku

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/notify"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

func HandleDetails(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	deviceName := vars["roku"]

	details, err := detailsGet(deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : details %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	notify.MetricHttpIn(r.URL.Path, http.StatusOK, r.Method)
	json.NewEncoder(w).Encode(details)
	return
}

func HandleStatus(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	status, err := StatusGet(deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : status %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	notify.MetricHttpIn(r.URL.Path, http.StatusOK, r.Method)
	json.NewEncoder(w).Encode(status)
	return
}

func HandleControl(w http.ResponseWriter, r *http.Request) {
	action := DeviceAction{}
	vars := mux.Vars(r)
	deviceName := vars["device"]
	appName := vars["app"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusBadRequest, r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &action); err != nil {
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := launchApp(deviceName,appName); err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}