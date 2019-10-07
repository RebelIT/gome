package inventory

import (
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"log"
	"time"
)

func GomeStatus() {
	var fName = "gome status runner"
	log.Printf("[INFO] %s, starting", fName)
	for {
		getStatus := true

		deviceNames, err := devices.GetAllDeviceNames()
		if err != nil {
			log.Printf("[WARN] %s, %s", fName, err)
			getStatus = false
		}
		if len(deviceNames) == 0 {
			log.Printf("[WARN] %s, nothing in the database", fName)
			getStatus = false
		}

		if getStatus {
			for _, name := range deviceNames {
				log.Printf("[DEBUG] %s, processing %s\n", fName, name)
				d, err := devices.GetProfile(name)
				if err != nil {
					log.Printf("[WARN] %s, unable to get profile for %s: %s", fName, name, err)
					continue
				}

				go devices.GetDeviceStatus(&d)
			}
		}
		time.Sleep(time.Minute * common.INVENTORY_MIN)
	}
}

