package roku

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"net/http"
)

func HandleLaunchApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	appName := vars["app"]

	if err := launchApp(deviceName,appName); err != nil{
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r,http.Response{})
	return
}