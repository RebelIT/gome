package roku

import (
	"github.com/pkg/errors"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"strings"
)

func rokuPost(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := http.Post(url, "", strings.NewReader(""))
	if err != nil{
		return *resp, err
	}
	return *resp, nil
}

func rokuGet(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":8060"+uriPart

	resp, err := http.Get(url)
	if err != nil{
		return *resp, err
	}
	return *resp, nil
}

func launchApp(deviceName string, app string) error {
	id, err := getAppId(app)
	if err != nil{
		return err
	}
	uri := "/launch/"+id
	resp, err := rokuPost(uri, deviceName)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return errors.New("non-200 status code return")
	}
	return nil
}

func DeviceStatus(deviceName string){
	data := devices.Status{}
	uriPart := "/"

	resp, err := rokuGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
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

	if err := devices.DbHashSet(deviceName+"_"+"status", data); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s :  device status done\n", deviceName)
	return
}