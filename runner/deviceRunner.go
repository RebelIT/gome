package runner

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"time"
)

func GoGoRunners() error {
	fmt.Println("[DEBUG] Starting Wrapper")
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
	for {
		fmt.Println("[DEBUG] Starting Device Loader")
		for _, d := range in.Devices {
			fmt.Println("[DEBUG] Adding: " + d.Name)
			CacheSetHash(db, redis.Args{d.Name}.AddFlat(d))
		}
		fmt.Println("[DEBUG] Resting device loader")
		time.Sleep(time.Second * 20)
	}
	return
}

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