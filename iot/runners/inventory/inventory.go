package inventory

import (
	"github.com/rebelit/gome/iot/devices"
	"log"
)

func GomeDevices() {
	var fName = "gome devices runner"
	log.Printf("[INFO] %s, starting", fName)

	go devices.GetDeviceStatus()
}
