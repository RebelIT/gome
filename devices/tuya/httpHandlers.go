package tuya

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/cache"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const FILE  = "./devices.json"

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

func GetDetails(w http.ResponseWriter,r *http.Request){
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

func GetStatus(w http.ResponseWriter,r *http.Request) {
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

	fmt.Println("##########Action not match: set it")
	args := []string{"set","--id", dbRet.Id, "--key", dbRet.Key, "--set", strconv.FormatBool(a.Action)}
	cmdOut, err := tryTuyaCli(string("tuya-cli"), args)
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

func GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dev := vars["device"]

	resp, err := ScheduleGet(dev)
	if err != nil{
		fmt.Printf("unable to get schedule status:  %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}

func SetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dev := vars["device"]
	in := Schedule{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Printf("bad body %s\n", r.RequestURI)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &in); err != nil {
		fmt.Printf("cant unmarshal %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := scheduleSet(&in, dev); err != nil{
		fmt.Printf("unable to set schedule:  %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return
}

func DelSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dev := vars["device"]

	if err := scheduleDel(dev); err != nil{
		fmt.Printf("unable to delete schedule:  %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dev := vars["device"]
	status := vars["status"]

	if err := scheduleUpdate(dev, status); err != nil{
		fmt.Printf("unable to delete schedule:  %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}