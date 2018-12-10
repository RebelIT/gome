package tuya

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
	deviceName := vars["device"]

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
	vars := mux.Vars(r)
	deviceName := vars["device"]
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

func HandleScheduleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	schedule, err := ScheduleGet(deviceName)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	json.NewEncoder(w).Encode(schedule)
	return
}

func HandleScheduleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	in := Schedules{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusBadRequest, r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := scheduleSet(&in, deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}

func HandleScheduleDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	if err := scheduleDel(deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}

func HandleScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	status := vars["status"]

	if err := scheduleUpdate(deviceName, status); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}