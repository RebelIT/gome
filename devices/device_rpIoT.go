package devices

import (
	"fmt"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

/// new functions
func ActionRpIot(profile Profile, action Action) error {
	actionUri := action.constructAction()
	url := fmt.Sprintf("http://%s:%s/%s", profile.Metadata.NetAddr, profile.Metadata.Port, actionUri)

	resp, err := common.HttpPost(url, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func StateRpIot(name string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n", name, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)
	alive := false

	value, err := db.Get(name)
	if err != nil {
		log.Printf("[WARN] %s not found for status collection, is it orphaned: %s", name, err)
		return
	}

	p, err := stringToStruct(value)
	if err != nil {
		log.Printf("[WARN] %s does not have a valid profile in the database: %s", name, err)
		return
	}

	url := fmt.Sprintf("http://%s:%s/api/alive", p.Metadata.NetAddr, p.Metadata.Port,)

	headers := map[string]string{
		"X-API-User":  p.Metadata.Username,
		"X-API-Token": p.Metadata.Password,
		"accept": "application/json",
		"contentType": "application/json",
	}

	resp, err := common.HttpGet(url, headers)
	if err != nil {
		//dont do anything, let return code handle it
	}

	if resp.StatusCode == 200 {
		alive = true
	}

	p.State.Alive = alive

	if err := db.Add(name, p.structToString()); err != nil{
		log.Printf("[WARN] %s unable to update the status with: %s", name, err)
		return
	}

	log.Printf("[INFO] %s status: %s done & done", name, strconv.FormatBool(alive))
	return
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
