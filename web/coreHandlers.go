package listener

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/notify"
	"github.com/rebelit/gome/runner"
	"io/ioutil"
	"net/http"
	"os"
)

//core handlers to update the device inventory.  devices.json used on gome startup to load last state into redis
//
func getDevices(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] "+ r.Method + " " + r.RequestURI)
	var i Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] Loaded json")

	json.Unmarshal(deviceFile, &i)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	json.NewEncoder(w).Encode(i)

	return
}

func addDevice(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] "+ r.Method + " " + r.RequestURI)
	var i Devices
	fullDevs := &Inputs{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusBadRequest, r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("[DEBUG] read input")
	defer r.Body.Close()

	if err := json.Unmarshal(body, &i); err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] unmarshal input")

	//Read devices.json and unmarshal into struct
	deviceFile, err := os.OpenFile(common.FILE, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer deviceFile.Close()

	bytes, err := ioutil.ReadAll(deviceFile)
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] Loaded json")

	if err := json.Unmarshal(bytes, &fullDevs); err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Append new device to devices array struct
	fullDevs.Devices = append(fullDevs.Devices,i)

	newBytes, err := json.MarshalIndent(fullDevs, "", "    ")
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Write new file with the new device appended to struct
	_, err = deviceFile.WriteAt(newBytes, 0)
	if err != nil {
		fmt.Println(err)
		notify.MetricHttpIn(r.RequestURI, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Re-run device loader to add to DB cache
	runner.GoGODeviceLoader()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	json.NewEncoder(w).Encode(i)

}
