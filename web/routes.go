package listener

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/runners/scheduler"
)
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
	//IoT Device Endpoints
	Route{"device", "GET", "/device", devices.GetDevices},
	Route{"device", "GET", "/device/{name}", devices.GetDeviceByName},
	Route{"device", "DELETE", "/device/{name}", devices.RemoveDevice},
	Route{"device", "POST", "/device", devices.AddUpdateDevice},
	Route{"device", "POST", "/device/{name}/{bool}", devices.ToggleDevice},

	Route{"device", "POST", "/api/tuya/{name}/{state}", devices.TuyaControl},
	Route{"device", "POST", "/api/roku/{name}/app/{app}", devices.RokuLaunchApp},
	Route{"device", "POST", "/api/pi/{name}/{component}", devices.RpIotControl},
	//Schedule Endpoints
	Route{"schedule", "GET", "/api/schedule/{device}", scheduler.HandleScheduleGet},
	Route{"schedule", "POST", "/api/schedule/{device}", scheduler.HandleScheduleSet},
	Route{"schedule", "DELETE", "/api/schedule/{device}", scheduler.HandleScheduleDel},
	Route{"schedule", "POST", "/api/schedule/{device}/{status}", scheduler.HandleScheduleUpdate},
	//Details Endpoints
	Route{"details", "GET", "/api/details/{device}", devices.HandleDetails},
	Route{"status", "GET", "/api/status/{device}", devices.HandleStatus},
}