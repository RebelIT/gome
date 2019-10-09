package devices

import (
	"fmt"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"log"
	"strconv"
	"time"
)

//Roku App ID's
//const NETFLIX = 12
//const PLEX = 13535
//const SLING = 46041
//const PANDORA = 28
//const PRIME_VIDEO = 13
//const GOOGLE_PLAY = 50025
//const HBOGO = 8378
//const YOUTUBE = 837

/// new functions
func ToggleRoku(profile Profile, powerstate bool) error {
	var control = ""
	if powerstate{
		control = "PowerOn"
	} else{
		control = "PowerOff"
	}

	url := fmt.Sprintf("http://%s:%s/%s/%s", profile.Metadata.NetAddr, profile.Metadata.Port, profile.Metadata.UriPart, control)

	resp, err := common.HttpPost(url, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func ActionRoku(profile Profile, action Action) error {
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

func StateRoku(name string, collectionDelayMin time.Duration) {
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

	url := fmt.Sprintf("http://%s:%s/", p.Metadata.NetAddr, p.Metadata.Port,)

	resp, err := common.HttpGet(url, nil)
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