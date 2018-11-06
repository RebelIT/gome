package rpi

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/cache"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
//TODO:
//TODO: function to read devices.json and only return the datbase.  probably in core handlers for all devices to use.
//TODO: a lot of this is junk code.. but it works. a lot to reactor later.

const FILE  = "./devices.json"

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

//http handler request to return specific device details
func HandleDetails(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] getting details for: " + r.URL.Path)
	uri := strings.Split(r.URL.Path, "/")
	dev := uri[len(uri)-1]
	var in Inputs

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.Unmarshal(deviceFile, &in)
	db := in.Database

	devDetail, err := cache.CacheGetHash(db, dev)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(devDetail)
	return
}

//http handler request to check the status of a device that was found i nthe runner and stored in redis
func HandleStatus(w http.ResponseWriter,r *http.Request) {
	fmt.Println("[DEBUG] getting status for: " + r.URL.Path)
	uri := strings.Split(r.URL.Path, "/")
	vars := mux.Vars(r)
	dev := vars["device"]
	action := uri[len(uri)-1]
	var in Inputs

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.Unmarshal(deviceFile, &in)
	db := in.Database

	s, err := cache.GetStatus(db, dev+"_"+action)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
	return
}

//http handler request to perform actions against the rpi device
func DeviceControl(w http.ResponseWriter, r *http.Request) {
	a := DeviceAction{}
	vars := mux.Vars(r)
	dev := vars["device"]
	var in Inputs

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.Unmarshal(deviceFile, &in)
	db := in.Database

	dbRet, err := cache.GetHashKey(db, redis.Args{dev})
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url := "http://"+dbRet.Addr+":"+dbRet.NetPort+"/action/"+a.Action

	resp, err := http.Get(url)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else{
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	return
}

//Used for Device runner for running is alive inventory stored in redis
func DeviceStatus(db string, addr string, port string, name string) {
	fmt.Println("[DEBUG] Starting Device Status for "+name)
	data := Status{}
	url := "http://"+addr+":"+port+"/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[ERROR] Error in Request to "+name+" will Retry")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data.Alive = false
	} else {
		data.Alive = true
	}
	data.Url = url
	data.Device = name

	if err := cache.SetHash(db, redis.Args{name+"_"+"status"}.AddFlat(data)); err != nil {
		fmt.Println("[ERROR] Error in adding "+name+" to cache will retry")
		return
	}
	fmt.Println("[DEBUG] Done with Device Status for "+name)
	return
}
