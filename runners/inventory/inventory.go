package inventory

import (
	"github.com/rebelit/gome/devices"
	"log"
)

func GomeDevices() {
	var fName = "gome devices runner"
	log.Printf("[INFO] %s, starting", fName)

	LoadDevices() //optional file to bulk load devices on boot
	go devices.GetDeviceStatus()
}

