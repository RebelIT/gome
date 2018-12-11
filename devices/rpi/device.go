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

func Connect(address string) (Pi, error) {
	return Pi{address: address, client: &http.Client{}}, nil
}

func (p Pi) Do(uriPart string) (http.Response, error) {
	url := "http://"+p.address+":6660"+uriPart

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.Response{}, err
	}
	resp, err := p.client.Do(req)
	return *resp, err
}

func doAction(deviceName string, action string) error {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return err
	}

	p, err := Connect(d.Addr)
	if err != nil {
		return err
	}

	uri := "/action/"+action
	resp, err := p.Do(uri)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func DeviceStatus(addr string, deviceName string) {
	data := devices.Status{}
	uriPart := "/"

	p, err := Connect(addr)
	if err != nil {
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	resp, err := p.Do(uriPart)
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
