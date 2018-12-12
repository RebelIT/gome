package roku

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
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
	fmt.Printf("dbData:  %+v\n", d)
	url := "http://"+d.Addr+":"+d.NetPort+uriPart
	fmt.Printf("roku doing post %s\n")
	resp, err := http.Post(url, "", strings.NewReader(""))
	if err != nil{
		fmt.Println(err)
		return *resp, err
	}
	fmt.Printf("roku post done app %s\n")

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
		fmt.Println(err)

		return *resp, err
	}

	return *resp, nil
}

func launchApp(deviceName string, app string) error {
	id, err := getAppId(app)
	if err != nil{
		return err
	}
	fmt.Printf("roku details app %s device %s\n", id, deviceName)

	uri := "/launch/"+id
	resp, err := rokuPost(uri, deviceName)
	fmt.Printf("roku launched app %s\n", id)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func DeviceStatus(deviceName string){
	data := devices.Status{}
	uriPart := "/"

	resp, err := rokuGet(uriPart, deviceName)
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

	if _, err := c.Do("HMSET", redis.Args{deviceName+"_"+"status"}.AddFlat(data)...); err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s : status done\n", deviceName)
	return
}