package devices

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetDevices(w http.ResponseWriter,r *http.Request){
	type list struct{
		Devices []string `json:"devices"`
	}
	response := list{}

	values, err := db.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Devices = values

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func GetDeviceByName(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]

	value, err := db.Get(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := stringToStruct(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func RemoveDevice(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]

	if err := db.Del(name); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return
}

func AddUpdateDevice(w http.ResponseWriter,r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newProfile, err := stringToStruct(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.Add(newProfile.Name, newProfile.structToString());err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newProfile); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func ToggleDevice(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]
	state := vars["bool"]

	toggle, err := strconv.ParseBool(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	value, err := db.Get(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profile, err := stringToStruct(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if profile.State.Alive == false{
		w.WriteHeader(http.StatusInternalServerError)
	}

	if profile.State.Status == toggle{
		w.WriteHeader(http.StatusBadRequest)
	}

	switch profile.Make {
	case "tuya":
		if err := StateControlTuya(profile, toggle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		profile.State.Status = true

		if err := db.Add(profile.Name, profile.structToString()); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "roku":
		if err := StateControlRoku(profile,toggle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		profile.State.Status = true

		if err := db.Add(profile.Name, profile.structToString()); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "rpiot":
		if err := StateControlRpIot(profile,toggle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		profile.State.Status = true

		if err := db.Add(profile.Name, profile.structToString()); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}





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