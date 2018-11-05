package rpi

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
//TODO:
//TODO: function to read devices.json and only return the datbase.  probably in core handlers for all devices to use.
//TODO:

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
	dev := uri[len(uri)-2]
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

	if err := cache.CacheSetHash(db, redis.Args{name+"_"+"status"}.AddFlat(data)); err != nil {
		fmt.Println("[ERROR] Error in adding "+name+" to cache will retry")
		return
	}
	fmt.Println("[DEBUG] Done with Device Status for "+name)
	return
}
