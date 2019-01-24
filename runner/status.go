package runner

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
	"github.com/rebelit/gome/devices/tuya"
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
			for _, dev := range (devs) {
				doItForReal := true
				devData, err := devices.DbHashGet(dev)
				if err != nil {
					log.Printf("[WARN] device status runner, unable to get dbData for %s: %s", dev, err)
					doItForReal = false
				}
				if doItForReal {
					d := devices.Devices{}
					redis.ScanStruct(devData, &d)

					switch d.Device {
					case "pi":
						go rpi.DeviceStatus(d.Name)

					case "roku":
						go roku.DeviceStatus(d.Name)

					case "tuya":
						go tuya.DeviceStatus(d.Addr, d.Id, d.Key, d.Name)

					default:
						log.Printf("[WARN] device status runner, %s no device types match", d.Name)
					}
				}
			}
		}
		time.Sleep(time.Second *10)
	}

	log.Printf("[ERROR] device status runner broke out of loop")
	return
}