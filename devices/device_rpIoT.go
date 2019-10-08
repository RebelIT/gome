package devices

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/// new functions
func StateControlRpIot(profile Profile, powerstate bool) error {
	var control = ""
	if powerstate{
		control = "on"
	} else{
		control = "off"
	}
	url := fmt.Sprintf("http://%s:%s/api/%s/%s",profile.Metadata.NetAddr, profile.Metadata.Port, profile.Action, control)

	resp, err := common.HttpPost(url, nil, rpIotHeaders())
	if err != nil {

		return err
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("%s returned %d for %s", profile.Name, resp.StatusCode, url)
	}

	return nil
}

func RpIotDeviceStatus(deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n", deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	rpIotAliveStatus(deviceName)
	rpIotDisplayStatus(deviceName)

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

func rpIotAliveStatus(deviceName string) {
	uriPart := "/api/alive"
	alive := false

	resp, err := PiGet(uriPart, deviceName)
	if err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		if err := db.Add(deviceName+"_"+"status", string([]byte(strconv.FormatBool(alive)))); err != nil {
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		alive = true
	}

	if err := db.Add(deviceName+"_"+"status", string([]byte(strconv.FormatBool(alive)))); err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}
	return
}

func rpIotDisplayStatus(deviceName string) {
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
	if err := json.Unmarshal(body, &piBody); err != nil {
		return
	}

	if piBody.Message == "1" {
		state = true
	}

	if err := db.Add(deviceName+"_display"+"_state", string([]byte(strconv.FormatBool(state)))); err != nil {
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return
	}
	return
}


// http wrappers
func PiGet(uriPart string, deviceName string) (response http.Response, err error) {
	d, err := GetDevice(deviceName)
	if err != nil {
		return http.Response{}, err
	}

	url := "http://" + d.Addr + ":" + d.NetPort + uriPart

	resp, err := common.HttpGet(url, rpIotHeaders())
	if err != nil {
		return http.Response{}, err
	}

	return resp, nil
}

func PiPost(deviceName string, uriPart string) (response http.Response, err error) {
	d, err := GetDevice(deviceName)
	if err != nil {
		return http.Response{}, err
	}

	url := "http://" + d.Addr + ":" + d.NetPort + uriPart

	resp, err := common.HttpPost(url, nil, rpIotHeaders())
	if err != nil {
		return http.Response{}, err
	}

	return resp, nil
}

// helpers
func compileUrl(uriPart string, d PiControl) (uri string, err error) {
	switch uriPart {
	case "power":
		return uriPart + "/" + d.Action, nil

	case "apt":
		if d.Package == "" {
			return uriPart + "/" + d.Action, nil
		} else {
			return uriPart + "/" + d.Package + "/" + d.Action, nil
		}

	case "service":
		return uriPart + "/" + d.Service + "/" + d.Action, nil

	case "display":
		return uriPart + "/" + d.Action, nil

	case "gpio":
		return uriPart + "/" + d.PinNumber + "/" + d.Action, nil

	default:
		return "", errors.New("no pi component " + uriPart + " action found")
	}
}

func rpIotHeaders() (headers map[string]string) {
	s, err := common.GetSecrets()
	if err != nil {
		log.Printf("[ERROR] unable to set RaspberryPi rpIoT headers")
		return
	}

	h := map[string]string{
		"X-API-User":  s.RpiotUser,
		"X-API-Token": s.RpiotToken,
	}
	return h
}
