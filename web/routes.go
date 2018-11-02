package listener

import (
	"github.com/gorilla/mux"
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
	//Route{"RokuControl", "POST", "/roku", roku.DeviceControl},
	//Route{"RokuControl", "GET", "/roku", roku.GetOptRoot},
	//Route{"RokuControl", "GET", "/roku/launch", roku.GetOptLaunch},
}