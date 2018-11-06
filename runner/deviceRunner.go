package runner

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
	"io/ioutil"
	"time"
)

const FILE  = "./devices.json"

func GoGoRunners() error {
	fmt.Println("[DEBUG] Starting runners")
	var in Inputs

	for {
		deviceFile, err := ioutil.ReadFile(FILE)
		if err != nil {
			fmt.Println(err)
			return err
		}
		json.Unmarshal(deviceFile, &in)
		db := in.Database

		for _, d := range (in.Devices) {
			switch d.Device {
			case "pi":
				go rpi.DeviceStatus(db, d.Addr, d.NetPort, d.Name)

			case "roku":
				go roku.DeviceStatus(db, d.Addr, d.NetPort, d.Name)

			default:
				fmt.Println("[ERROR] No device type match for "+ d.Name)
			}
		}
		time.Sleep(time.Second *10)
	}

	return nil
}

func GoGODeviceLoader() error {
	fmt.Println("[DEBUG] Starting Device Loader")
	var in Inputs

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("[DEBUG] Loaded json")
	json.Unmarshal(deviceFile, &in)
	db := in.Database

	go DeviceLoader(db, in)

	return nil
}

func DeviceLoader(db string, in Inputs) {
	//Load Devices into database from startup json
	for _, d := range in.Devices {
		fmt.Println("[DEBUG] Adding: " + d.Name)
		cache.SetHash(db, redis.Args{d.Name}.AddFlat(d))
	}
	return
}