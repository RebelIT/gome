package inventory

import (
	"encoding/json"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"io/ioutil"
	"log"
)

func LoadDevices() { //Load DevicesOld into redis from devices.json file
	var fName = "gome loader"
	log.Printf("[INFO] %s, starting", fName)
	i, err := ReadDeviceFile()
	if err != nil {
		return
	}

	if len(i.Profiles) == 0 {
		log.Printf("[WARN] %s, no devices to load, skipping",fName)
		return
	}

	for _, p := range i.Profiles {
		log.Printf("[INFO] %s, loading %s", fName, p.Name)

		if err := devices.UpdateProfile(p); err != nil{
			log.Printf("[ERROR] %s, loading %s",fName, p.Name)
			continue
		}
	}
	log.Printf("[INFO] %s, all done",fName)
	return
}

func ReadDeviceFile() (devices.Devices, error) { //Read the devices.json
	var list devices.Devices
	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		return list, err
	}
	if err := json.Unmarshal(deviceFile, &list); err != nil {
		return list, err
	}

	return list, nil
}