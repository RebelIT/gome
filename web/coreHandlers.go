package listener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const FILE  = "./devices.json"

//core handlers to update the device inventory.  devices.json used on gome startup to load last state into redis
//
func getDevices(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] "+ r.Method + " " + r.RequestURI)
	var i Inputs

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] Loaded json")

	json.Unmarshal(deviceFile, &i)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("[DEBUG] read input")
	defer r.Body.Close()

	if err := json.Unmarshal(body, &i); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] unmarshal input")

	deviceFile, err := os.OpenFile(FILE, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer deviceFile.Close()

	bytes, err := ioutil.ReadAll(deviceFile)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("[DEBUG] Loaded json")

	if err := json.Unmarshal(bytes, &fullDevs); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fullDevs.Devices = append(fullDevs.Devices,i)

	newBytes, err := json.MarshalIndent(fullDevs, "", "    ")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = deviceFile.WriteAt(newBytes, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO Add to cache

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)

}
