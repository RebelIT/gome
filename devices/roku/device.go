package roku

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/cache"
	"github.com/rebelit/gome/common"
	"io/ioutil"
	"net/http"
	"strings"
)

//http handler request to return specific device details
func HandleDetails(w http.ResponseWriter,r *http.Request){
	fmt.Println("[DEBUG] getting details for: " + r.URL.Path)
	vars := mux.Vars(r)
	dev := vars["roku"]

	var in Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
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

//http handler request to check the status of a device that was found in the runner and stored in redis
func HandleStatus(w http.ResponseWriter,r *http.Request) {
	fmt.Println("[DEBUG] getting status for: " + r.URL.Path)
	uri := strings.Split(r.URL.Path, "/")
	vars := mux.Vars(r)
	dev := vars["roku"]
	action := uri[len(uri)-1]
	var in Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
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
	dev := vars["roku"]
	var in Inputs

	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Println(err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &a); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deviceFile, err := ioutil.ReadFile(common.FILE)
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
	//roku:port/launch/appid
	//roku:port/keypress/button
	url := "http://"+dbRet.Addr+":"+dbRet.NetPort+"/"+a.Action+"/"+a.ActionItem

	resp, err := http.Post(url, "", strings.NewReader(""))
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
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}

//Used for Device runner for running is alive inventory stored in redis
func DeviceStatus(db string, addr string, port string, name string){
	data := Status{}
	url := "http://"+addr+":"+port+"/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[ERROR] Error in Request to "+ url +" will Retry")
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

//Used for Device runner for roku app inventory stored in redis for launching apps
func DeviceApps(db string, addr string, port string, name string){
	apps := Apps{}
	url := "http://"+addr+":"+port+"/query/apps"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("[ERROR] Error in Request to "+ url +" will Retry")
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &apps)

	for _,i := range(apps.Apps){
		a := strings.Replace(i.App, " ", "", -1)
		cache.Set(db,a, i.Id)
	}
	return
}