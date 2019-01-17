package rpi

import (
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
)

func PiGet(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := http.Get(url)
	if err != nil{
		notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")
		return *resp, err
	}
	notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")
	return *resp, nil
}

func PiPost(deviceName string, action string) error {
	uriPart := "/action/"+action
	resp, err := PiGet(uriPart, deviceName)
	if err != nil{
		notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")
		return err
	}
	notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")
	return nil
}

func DeviceStatus(deviceName string) {
	data := devices.Status{}
	uriPart := "/"

	resp, err := PiGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data.Alive = false
	} else {
		data.Alive = true
	}
	data.Device = deviceName

	if err := devices.DbHashSet(deviceName+"_"+"status", data); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s : device status done\n", deviceName)
	return
}
