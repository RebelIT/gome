package runner

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
)

const FILE  = "./devices.json"

func GoGoRunners() error {
	fmt.Println("[DEBUG] Starting runners")

	return nil
}

func GoGODeviceLoader() error {
	fmt.Println("[DEBUG] Starting Device Loader")
	var in Inputs

	deviceFile, err := ioutil.ReadFile("./devices.json")
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
		CacheSetHash(db, redis.Args{d.Name}.AddFlat(d))
	}
	return
}


//func inventoryDevices() (*Inputs, error) {
//	i := Inputs{}
//
//	deviceFile, err := ioutil.ReadFile(FILE)
//	if err != nil {
//		fmt.Println(err)
//		return &i, err
//	}
//	fmt.Println("[DEBUG] Loaded json")
//
//	json.Unmarshal(deviceFile, &i)
//	return &i, nil
//}
//
//func StatusRunner(db string, in Inputs) {
//	fmt.Println("[DEBUG] Starting Status Runners")
//	for _, d := range in.Devices {
//		info, err := CacheGetHash(db, "device_"+d.Name)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		d := net.Dialer{Timeout: time.Millisecond * 50}
//		conn, err := d.Dial("tcp", info.Addr+":"+info.NetPort)
//		if err != nil {
//			fmt.Println(info.Name + " is not available on " + info.NetPort)
//			CacheSet(db, "status_"+info.Name, "false")
//		} else {
//			defer conn.Close()
//			fmt.Println(info.Name + " is available on " + info.NetPort)
//			CacheSet(db, "status_"+info.Name, "true")
//		}
//	}
//}