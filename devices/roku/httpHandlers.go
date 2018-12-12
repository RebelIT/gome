package roku

import (
	"github.com/gorilla/mux"
	"github.com/rebelit/gome/notify"
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
		log.Printf("[ERROR] %s : control %s, %s", deviceName, r.Method, err)
		notify.MetricHttpIn(r.URL.Path, http.StatusInternalServerError, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notify.MetricHttpIn(r.RequestURI, http.StatusOK, r.Method)
	return
}