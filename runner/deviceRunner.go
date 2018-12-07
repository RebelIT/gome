package runner

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices/roku"
	"github.com/rebelit/gome/devices/rpi"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"io/ioutil"
	"time"
)

func GoGoRunners() error {
	fmt.Println("[INFO] Starting runners")
	var in Inputs

	for {
		deviceFile, err := ioutil.ReadFile(common.FILE)
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
				// Assumption that all roku's use the same account and apps are in sync.
				go roku.DeviceApps(db, d.Addr, d.NetPort, d.Name)

			case "tuya":
				go tuya.DeviceStatus(db, d.Addr, d.Id, d.Key, d.Name)

			default:
				fmt.Println("[ERROR] No device typse match for "+ d.Name)
			}
		}
		time.Sleep(time.Second *60)
	}

	notify.SendSlackAlert("Device runner broke out of loop")
	return nil
}

func GoGODeviceLoader() error {
	fmt.Println("[INFO] Starting Device Loader")
	var in Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("[INFO] Loaded json")
	json.Unmarshal(deviceFile, &in)
	db := in.Database

	go DeviceLoader(db, in)

	return nil
}

func DeviceLoader(db string, in Inputs) {
	//Load Devices into database from startup json
	for _, d := range in.Devices {
		fmt.Println("[DEBUG] Adding: " + d.Name)
		if err := cache.SetHash(db, redis.Args{d.Name}.AddFlat(d));err != nil{
			fmt.Printf("error loading %s\n ",err)
		}
	}
	return
}