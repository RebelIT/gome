package status

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
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
				doItForReal := true
				devData, err := database.DbHashGet(dev)
				if err != nil {
					log.Printf("[WARN] device status runner, unable to get dbData for %s: %s", dev, err)
					doItForReal = false
				}

				d := devices.Devices{}
				if err := redis.ScanStruct(devData, &d); err != nil{
					log.Printf("[WARN] device status runner, unable to parse dbData for %s: %s", dev, err)
					doItForReal = false
				}

				if doItForReal {
					go devices.GetDeviceStatus(d)
				}
			}
		}
		time.Sleep(time.Minute *common.STATUS_MIN)
	}

	log.Printf("[ERROR] device status runner broke out of loop")
	return
}