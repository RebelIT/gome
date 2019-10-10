package cron

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/common"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleScheduleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	hasSchedule, schedule, err := scheduleGet(deviceName)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}
	if !hasSchedule{
		common.ReturnBad(w,r)
		return
	}

	common.ReturnOk(w,r, schedule)
	return
}

func HandleScheduleSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	in := Schedules{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnBad(w,r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	if err := scheduleSet(&in, deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,nil)
	return
}

func HandleScheduleDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	if err := scheduleDel(deviceName); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,nil)
	return
}

func HandleScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	status := vars["status"]

	if err := scheduleUpdate(deviceName, status); err != nil{
		log.Printf("[ERROR] %s : schedule %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,nil)
	return
}