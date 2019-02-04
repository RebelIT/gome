package devices

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

	hasSchedule, schedule, err := ScheduleGet(deviceName)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnInternalError(w,r)
		return
	}
	if !hasSchedule{
		ReturnBad(w,r)
		return
	}

	ReturnOk(w,r, schedule)
	return
}

func HandleScheduleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	in := Schedules{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnBad(w,r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnInternalError(w,r)
		return
	}

	if err := ScheduleSet(&in, deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnInternalError(w,r)
		return
	}

	ReturnOk(w,r,nil)
	return
}

func HandleScheduleDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	if err := ScheduleDel(deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnInternalError(w,r)
		return
	}

	ReturnOk(w,r,nil)
	return
}

func HandleScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	status := vars["status"]

	if err := ScheduleUpdate(deviceName, status); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		ReturnInternalError(w,r)
		return
	}

	ReturnOk(w,r,nil)
	return
}
