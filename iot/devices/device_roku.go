package devices

import (
	"fmt"
	db "github.com/rebelit/gome/database"
	"github.com/rebelit/gome/util/http"
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

func (p *Profile) ToggleRoku(powerstate bool) error {
	var control = ""
	if powerstate {
		control = "PowerOn"
	} else {
		control = "PowerOff"
	}

	url := fmt.Sprintf("http://%s:%s/%s/%s", p.Metadata.NetAddr, p.Metadata.Port, p.Metadata.UriPart, control)
	resp, err := http.HttpPost(url, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return err
	}

	return nil
}

func (p *Profile) ActionRoku(action Action) error {
	actionUri := action.constructAction()
	url := fmt.Sprintf("http://%s:%s/%s", p.Metadata.NetAddr, p.Metadata.Port, actionUri)

	resp, err := http.HttpPost(url, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return err
	}

	return nil
}

func StateRoku(name string, collectionDelayMin time.Duration) {
	log.Printf("INFO: StateRoku, %s device collection delayed +%d sec", name, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)
	alive := false

	value, err := db.Get(name)
	if err != nil {
		log.Printf("WARN: StateRoku, %s not found for status collection, is it orphaned: %s", name, err)
		return
	}

	p, err := stringToStruct(value)
	if err != nil {
		log.Printf("WARN: StateRoku, %s does not have a valid profile in the database: %s", name, err)
		return
	}

	url := fmt.Sprintf("http://%s:%s/", p.Metadata.NetAddr, p.Metadata.Port)

	resp, err := http.HttpGet(url, nil)
	if err != nil {
		//dont do anything, let return code handle it
	}

	if resp.StatusCode == 200 {
		alive = true
	}

	p.State.Alive = alive

	if err := db.Add(name, p.structToString()); err != nil {
		log.Printf("WARN: StateRoku, %s unable to update the status with: %s", name, err)
		return
	}

	log.Printf("INFO: StateRoku, %s status: %s done & done", name, strconv.FormatBool(alive))
	return
}
