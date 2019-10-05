package status

import (
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"log"
	"time"
)

func GoGoDeviceStatus() {
	log.Println("[INFO] device status runner, starting")
	for {
		doIt := true

		devs, err := devices.GetAllDevicesFromDb()
		if err != nil {
			log.Printf("[WARN] device status runner, %s", err)
			doIt = false
		}
		if len(devs) == 0{
			log.Printf("[WARN] device status runner, no devices found in the database")
			doIt = false
		}

		if doIt {
			for _, dev := range devs {
				log.Printf("[DEBUG] device status runner, processing %s\n", dev)
				doItForReal := true
				d, err := devices.GetDevice(dev)
				if err != nil {
					log.Printf("[WARN] device status runner, unable to get dbData for %s: %s", dev, err)
					doItForReal = false
				}

				if doItForReal {
					go devices.GetDeviceStatus(&d)
				}
			}
		}
		time.Sleep(time.Minute *common.STATUS_MIN)
	}

	log.Printf("[ERROR] device status runner broke out of loop")
	return
}