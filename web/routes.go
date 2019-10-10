package listener

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/runners/cron"
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
	Route{"device", "POST", "/device/{name}/toggle/{bool}", devices.ToggleDevice},
	Route{"device", "POST", "/device/{name}/action/{action}", devices.ActionDevice},

	//Schedule Endpoints
	Route{"schedule", "GET", "/api/schedule/{device}", cron.HandleScheduleGet},
	Route{"schedule", "POST", "/api/schedule/{device}", cron.HandleScheduleSet},
	Route{"schedule", "DELETE", "/api/schedule/{device}", cron.HandleScheduleDel},
	Route{"schedule", "POST", "/api/schedule/{device}/{status}", cron.HandleScheduleUpdate},
}

