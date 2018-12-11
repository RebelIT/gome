package listener

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
	"github.com/rebelit/gome/devices/tuya"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	//Devices Endpoints
	Route{"device", "GET", "/api/device", getDevices},
	Route{"device", "POST", "/api/device", addDevice},
	Route{"device", "POST", "/api/device/tuya/{name}/{state}", tuya.HandleControl},
	Route{"device", "POST", "/api/device/roku/{name}/app/{app}", roku.HandleLaunchApp},
	Route{"device", "POST", "/api/device/pi/{name}/action", rpi.HandleControl},
	//Schedule Endpoints
	Route{"schedule", "GET", "/api/schedule/{device}", devices.HandleScheduleGet},
	Route{"schedule", "POST", "/api/schedule/{device}", devices.HandleScheduleSet},
	Route{"schedule", "DELETE", "/api/schedule/{device}", devices.HandleScheduleDel},
	Route{"schedule", "POST", "/api/schedule/{device}/{status}", devices.HandleScheduleUpdate},
	//Details Endpoints
	Route{"details", "GET", "/api/details/{device}", devices.HandleDetails},
	Route{"status", "GET", "/api/status/{device}", devices.HandleStatus},
}