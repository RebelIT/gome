package rpi

import (
	"github.com/pkg/errors"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"strings"
)

func PiGet(uriPart string, deviceName string) (response http.Response, err error) {
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

func PiPost(deviceName string, uri string) (response http.Response, err error) {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return http.Response{}, err
	}

	url := "http://"+d.Addr+":"+d.NetPort+"/api/"+uri

	resp, err := http.Post(url, "", strings.NewReader(""))
	if err != nil{
		notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")
		return http.Response{}, err
	}
	notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")

	return *resp, nil
}

func DeviceStatus(deviceName string) {
	data := devices.Status{}
	uriPart := "/api/alive"

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

func compileUrl(uriPart string, d PiControl) (uri string, err error){
	switch uriPart {
	case "power":
		return  uriPart+"/"+d.Action, nil

	case "apt":
		if d.Package == ""{
			return uriPart+"/"+d.Action, nil
		} else{
			return uriPart+"/"+d.Package+"/"+d.Action, nil
		}

	case "service":
		return uriPart+"/"+d.Service+"/"+d.Action, nil

	case "display":
		return  uriPart+"/"+d.Action, nil

	case "gpio":
		return uriPart+"/"+d.PinNumber+"/"+d.Action, nil

	default:
		return "", errors.New("no pi component "+uriPart+" action found" )
	}

	return "",nil
}