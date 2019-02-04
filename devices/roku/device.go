package roku

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"net/http"
	"time"
)

func rokuPost(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet("device_"+deviceName)
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*1)
	defer cncl()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil{
		return http.Response{}, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")
	return *resp, nil
}

func rokuGet(uriPart string, deviceName string) (http.Response, error) {
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

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil{
		return http.Response{}, err
	}

	notify.MetricHttpOut(deviceName, resp.StatusCode, "POST")
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

func DeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	fmt.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)
	data := devices.Status{}
	uriPart := "/"

	resp, err := rokuGet(uriPart, deviceName)
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