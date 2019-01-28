package roku

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/devices"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

func HandleLaunchApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceName := vars["device"]
	appName := vars["app"]

	log.Printf("roku app load %s %s\n", deviceName, appName)
	if err := launchApp(deviceName,appName); err != nil{
		devices.ReturnInternalError(w,r)
		return
	}

	devices.ReturnOk(w,r,http.Response{})
	return
}