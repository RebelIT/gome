package listener

import (
	"github.com/gorilla/mux"
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
	//Core
	Route{"coreDevice", "GET", "/api/device/all", getDevices},
	Route{"coreDevice", "POST", "/api/device", addDevice},
	//RaspberryPi
	Route{"rpi", "GET", "/api/rpi/{device}/details", rpi.HandleDetails},
	Route{"rpi", "GET", "/api/rpi/{device}/status", rpi.HandleStatus},
	Route{"rpi", "POST", "/api/rpi/{device}", rpi.DeviceControl},
	//Roku
	Route{"roku", "GET", "/api/roku/{roku}/details", roku.HandleDetails},
	Route{"roku", "GET", "/api/roku/{roku}/status", roku.HandleStatus},
	Route{"roku", "POST", "/api/roku/{roku}", roku.DeviceControl},
	//Tuya
	Route{"tuya", "GET", "/api/tuya/{device}/details", tuya.HandleDetails},
	Route{"tuya", "GET", "/api/tuya/{device}/status", tuya.HandleStatus},
	Route{"tuya", "POST", "/api/tuya/{device}", tuya.DeviceControl},
}