package scheduler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleScheduleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	hasSchedule, schedule, err := devices.ScheduleGet(deviceName)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}
	if !hasSchedule{
		devices.ReturnBad(w,r)
		return
	}

	devices.ReturnOk(w,r, schedule)
	return
}

func HandleScheduleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	in := devices.Schedules{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnBad(w,r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}

	if err := devices.ScheduleSet(&in, deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r,nil)
	return
}

func HandleScheduleDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	if err := devices.ScheduleDel(deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r,nil)
	return
}

func HandleScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	status := vars["status"]

	if err := devices.ScheduleUpdate(deviceName, status); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r,nil)
	return
}