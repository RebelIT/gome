package roku

import (
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"github.com/rebelit/gome/devices"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	uriPart := "/"
	alive := false

	resp, err := rokuGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		alive = true
	}

	if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
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


// http wrappers
func rokuPost(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName+"_device")
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpPost(url,nil,nil)
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}

func rokuGet(uriPart string, deviceName string) (http.Response, error) {
	d, err := devices.DetailsGet(deviceName+"_device")
	if err != nil{
		return http.Response{}, err
	}
	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpGet(url, nil)
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}