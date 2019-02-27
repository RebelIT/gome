package devices

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func RpIotDeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	rpIotAliveStatus(deviceName)
	rpIotDisplayStatus(deviceName)

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

func rpIotAliveStatus(deviceName string){
	uriPart := "/api/alive"
	alive := false

	resp, err := PiGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
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
	return
}

func rpIotDisplayStatus(deviceName string){
	uriPart := "/api/display"
	state := false

	piBody := PiResponse{}

	resp, err := PiGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : hdmi status, %s\n", deviceName, err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &piBody)

	if piBody.Message == "1"{
		state = true
	}

	if err := database.DbSet(deviceName+"_display"+"_state", []byte(strconv.FormatBool(state))); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}
	return
}

func rpIotDisplayToggle(deviceName string, toggle string) error{
	uriPart := "/api/display/"+toggle

	resp, err := PiPost(deviceName, uriPart)
	if err != nil{
		return err
	}

	if resp.StatusCode != 200{
		return errors.Errorf("%s returned %d for %s", deviceName, resp.StatusCode, uriPart)
	}

	state := false
	if toggle == "on"{
		state = true
	}

	if err := database.DbSet(deviceName+"_display"+"_state", []byte(strconv.FormatBool(state))); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return errors.Errorf("%s failed to set device status after issuing display %s", deviceName, toggle)
	}


	return nil
}

// http wrappers
func PiGet(uriPart string, deviceName string) (response http.Response, err error) {
	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		return http.Response{}, err
	}

	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpGet(url, rpIotHeaders())
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}

func PiPost(deviceName string, uriPart string) (response http.Response, err error) {
	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		return http.Response{}, err
	}

	url := "http://"+d.Addr+":"+d.NetPort+uriPart

	resp, err := common.HttpPost(url,nil, rpIotHeaders())
	if err != nil{
		return http.Response{}, err
	}

	return resp, nil
}


// helpers
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

func rpIotHeaders()(headers map[string]string){
	s, err := common.GetSecrets()
	if err != nil{
		log.Printf("[ERROR] unable to set RaspberryPi rpIoT headers")
		return
	}

	h := map[string]string{
		"X-API-User": s.RpiotUser,
		"X-API-Token": s.RpiotToken,
	}
	return h
}