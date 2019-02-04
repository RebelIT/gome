package rpi

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Timeout = time.Second * 2
}

func PiGet(uriPart string, deviceName string) (response http.Response, err error) {
	d, err := devices.DetailsGet("device_"+deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
	defer cncl()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	s, err := common.GetSecrets()
	if err != nil{
		return http.Response{}, err
	}
	req.Header.Set("X-API-User", s.RpiotUser)
	req.Header.Set("X-API-Token", s.RpiotToken)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")
	return *resp, nil
}

func PiPost(deviceName string, uri string) (response http.Response, err error) {
	d, err := devices.DetailsGet("device_"+deviceName)
	if err != nil{
		return http.Response{}, err
	}

	url := "http://"+d.Addr+":"+d.NetPort+"/api/"+uri

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	s, err := common.GetSecrets()
	if err != nil{
		return http.Response{}, err
	}
	req.Header.Set("X-API-User", s.RpiotUser)
	req.Header.Set("X-API-Token", s.RpiotToken)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")

	return *resp, nil
}

func DeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	fmt.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)
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

	log.Printf("[INFO] %s device status : done\n", deviceName)
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