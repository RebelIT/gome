package listener

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
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
	//Core
	Route{"coreDevice", "GET", "/api/device/all", getDevices},
	Route{"coreDevice", "POST", "/api/device", addDevice},
	//RaspberryPi
	Route{"rpi", "GET", "/api/{device}/details", rpi.HandleDetails},
	Route{"rpi", "GET", "/api/{device}/status", rpi.HandleStatus},
	Route{"rpi", "POST", "/api/{device}/action", rpi.DeviceControl},
	//Roku
	Route{"roku", "GET", "/api/{device}/details", roku.HandleDetails},
	Route{"roku", "GET", "/api/{device}/status", roku.HandleStatus},
	Route{"roku", "POST", "/api/{device}/action", roku.DeviceControl},
}