package rpi

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 5
}

func piGet(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := http.Get(url)
	if err != nil{
		return *resp, err
	}

	return *resp, nil
}

func doAction(deviceName string, action string) error {
	uriPart := "/action/"+action
	resp, err := piGet(uriPart, deviceName)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func DeviceStatus(deviceName string) {
	data := devices.Status{}
	uriPart := "/"

	resp, err := piGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")
		return
	}
	defer resp.Body.Close()

	notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")

	if resp.StatusCode != 200 {
		data.Alive = false
	} else {
		data.Alive = true
	}
	data.Device = deviceName

	c, err := devices.DbConnect()
	if err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{deviceName+"_"+"status"}.AddFlat(data)); err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s : status done\n", deviceName)
	return
}
