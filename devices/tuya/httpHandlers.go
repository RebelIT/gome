package tuya

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/cache"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)
const FILE  = "./devices.json"

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

//http handler request to return specific device details
func HandleDetails(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] getting details for: " + r.URL.Path)
	vars := mux.Vars(r)
	dev := vars["device"]

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

	s, err := cache.GetStatus(db, dev+"_status")

	if s.Alive == a.Action {
		fmt.Println("##########Action match:: fail it")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else{
		fmt.Println("##########Action not match: set it")
		args := []string{"set","--id", dbRet.Id, "--key", dbRet.Key, "--set", strconv.FormatBool(a.Action)}
		cmdOut, err := command(string("tuya-cli"), args)
		if err != nil{
			fmt.Println("[ERROR] Error in tyua Cli")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		} else {
			fmtOut := strings.Replace(cmdOut, "\n", "", -1)
			if fmtOut == "Set succeeded."{
				w.WriteHeader(http.StatusOK)
				return
			} else{
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}