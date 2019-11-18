package cron

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleScheduleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	hasSchedule, schedule, err := scheduleGet(deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !hasSchedule {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
	return
}

func HandleScheduleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	in := Schedules{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := scheduleSet(&in, deviceName); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleScheduleDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	if err := scheduleDel(deviceName); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	status := vars["status"]

	if err := scheduleUpdate(deviceName, status); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
