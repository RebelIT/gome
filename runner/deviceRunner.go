package runner

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"log"
	"time"
)

func GoGoRunners() error {
	log.Println("[INFO] runner, starting")
	for {
		devs, err := devices.LoadDevices()
		if err != nil{
			log.Printf("[WARN] runner, unable to load devices from file. skipping this round")
		}else{
			for _, d := range (devs.Devices) {
				switch d.Device {
				case "pi":
					go rpi.DeviceStatus(d.Name)

				case "roku":
					go roku.DeviceStatus(d.Name)

				case "tuya":
					go tuya.DeviceStatus(d.Addr, d.Id, d.Key, d.Name)

				default:
					log.Printf("[WARN] runner, %s no device types match", d.Name)
				}
			}
		}
		time.Sleep(time.Second *30)
	}

	notify.SendSlackAlert("[ERROR] runner, routine broke out of loop")
	return nil
}

func GoGODeviceLoader() error {
	log.Println("[INFO] loader, starting")
	devs, err := devices.LoadDevices()
	if err != nil{
		log.Printf("[WARN] runner, unable to load devices from file. skipping this round")
	}
	go DeviceLoader(devs)

	return nil
}

func DeviceLoader(in devices.Inputs) {
	//Load Devices into database from startup json
	db := in.Database
	c, err := redis.Dial("tcp", db)
	if err != nil {
		log.Println("[ERROR] Error writing to redis, catch it next time around")
	}

	defer c.Close()
	for _, d := range in.Devices {
		log.Printf("[INFO] loader, %s working", d.Name)

		if _, err := c.Do("HMSET", redis.Args{d.Name}.AddFlat(d)...);err != nil {
			log.Printf("[ERROR] loader, %s : %s\n ", d.Name, err)
		}
	}
	log.Println("[INFO] loader, all done")
	return
}