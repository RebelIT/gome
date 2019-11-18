package devices

import (
	"encoding/json"
	"github.com/gorilla/mux"
	db "github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetDevices(w http.ResponseWriter, r *http.Request) {
	type list struct {
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

func GetDeviceByName(w http.ResponseWriter, r *http.Request) {
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

func RemoveDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := db.Del(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return
}

func AddUpdateDevice(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	p := Profile{}
	if err := json.Unmarshal(body, &p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.Add(p.Name, p.structToString()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ToggleDevice(w http.ResponseWriter, r *http.Request) { //Toggle for simple on/off (true/false) IoT devices
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

	if profile.State.Alive == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	if profile.State.Status == toggle {
		w.WriteHeader(http.StatusBadRequest)
	}

	switch profile.Make {
	case "tuya":
		if err := profile.ToggleTuya(toggle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "roku":
		if err := profile.ToggleRoku(toggle); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profile.State.Status = toggle

	if err := db.Add(profile.Name, profile.structToString()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}

func ActionDevice(w http.ResponseWriter, r *http.Request) { //Toggle for simple on/off (true/false) IoT devices
	vars := mux.Vars(r)
	name := vars["name"]
	action := vars["action"]

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

	if profile.State.Alive == false {
		w.WriteHeader(http.StatusBadRequest)
	}

	doAction, err := validateAction(profile, action)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	switch profile.Make {
	case "roku":
		if err := profile.ActionRoku(doAction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "rpiot":
		if err := profile.ActionRpIot(doAction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}
