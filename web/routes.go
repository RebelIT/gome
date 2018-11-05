package listener

import (
	"github.com/gorilla/mux"
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
	Route{"getDevice", "GET", "/api/device/all", getDevices},
	Route{"addDevice", "POST", "/api/device", addDevice},
	Route{"getCalendar", "GET", "/api/calendar", rpi.HandleDetails},
	Route{"getCalendar", "GET", "/api/calendar/status", rpi.HandleStatus},
}