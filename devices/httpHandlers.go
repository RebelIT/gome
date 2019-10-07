package devices

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/runners/inventory"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HandleDetails(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	deviceName := vars["device"]

	details, err := GetDevice(deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : details %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,details)
	return
}

func HandleStatus(w http.ResponseWriter,r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]

	status, err := GetDeviceAliveState(deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : status %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,status)
	return
}

func GetDevices(w http.ResponseWriter,r *http.Request){
	log.Println("[DEBUG] "+ r.Method + " " + r.RequestURI)
	var i Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	json.Unmarshal(deviceFile, &i)

	common.ReturnOk(w,r,i)

	return
}

func AddDevice(w http.ResponseWriter,r *http.Request){
	var i DevicesOld
	fullDevs := &Inputs{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		common.ReturnBad(w,r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &i); err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	//Read devices.json and unmarshal into struct
	deviceFile, err := os.OpenFile(common.FILE, os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}
	defer deviceFile.Close()

	bytes, err := ioutil.ReadAll(deviceFile)
	if err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	if err := json.Unmarshal(bytes, &fullDevs); err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	//Append new device to devices array struct
	fullDevs.Devices = append(fullDevs.Devices,i)

	newBytes, err := json.MarshalIndent(fullDevs, "", "    ")
	if err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	//Write new file with the new device appended to struct
	_, err = deviceFile.WriteAt(newBytes, 0)
	if err != nil {
		log.Println(err)
		common.ReturnInternalError(w,r)
		return
	}

	//Re-run device loader to add to DB cache
	inventory.LoadDevices()

	common.ReturnOk(w,r,i)
	return
}

//**********************************************************************
// tuya device endpoints
func TuyaControl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["name"]
	state := vars["state"]

	action := false
	if state == "on"{
		action = true
	} else if state == "off"{
		action = false
	} else{
		common.ReturnBad(w,r)
		return
	}

	if err := TuyaPowerControl(deviceName, action); err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r, http.Response{})
	return
}

//**********************************************************************
// raspberryPi IoT device endpoints
func RpIotControl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	component := vars["component"]

	i := PiControl{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.ReturnBad(w, r)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &i); err != nil {
		common.ReturnInternalError(w, r)
		return
	}

	uri, err := compileUrl(component, i)
	if err != nil{
		common.ReturnBad(w, r)
		return
	}

	resp, err := PiPost(deviceName,uri)
	if err != nil{
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		common.ReturnInternalError(w, r)
		return
	}

	common.ReturnOk(w, r, resp)
	return
}

//**********************************************************************
// roku device endpoints
func RokuLaunchApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	appName := vars["app"]

	if err := launchApp(deviceName,appName); err != nil{
		common.ReturnInternalError(w,r)
		return
	}

	common.ReturnOk(w,r,http.Response{})
	return
}