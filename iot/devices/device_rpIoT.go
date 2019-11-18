package devices

import (
	"fmt"
	db "github.com/rebelit/gome/database"
	"github.com/rebelit/gome/util/http"
	"log"
	"strconv"
	"time"
)

func (p *Profile) ActionRpIot(action Action) error {
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

func StateRpIot(name string, collectionDelayMin time.Duration) {
	log.Printf("INFO: StateRpIot, %s device collection delayed +%d sec", name, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)
	alive := false

	value, err := db.Get(name)
	if err != nil {
		log.Printf("WARN: StateRpIot, %s not found for status collection, is it orphaned: %s", name, err)
		return
	}

	p, err := stringToStruct(value)
	if err != nil {
		log.Printf("WARN: StateRpIot, %s does not have a valid profile in the database: %s", name, err)
		return
	}

	url := fmt.Sprintf("http://%s:%s/api/alive", p.Metadata.NetAddr, p.Metadata.Port)

	headers := map[string]string{
		"X-API-User":  p.Metadata.Username,
		"X-API-Token": p.Metadata.Password,
		"accept":      "application/json",
		"contentType": "application/json",
	}

	resp, err := http.HttpGet(url, headers)
	if err != nil {
		//dont do anything, let return code handle it
	}

	if resp.StatusCode == 200 {
		alive = true
	}

	p.State.Alive = alive

	if err := db.Add(name, p.structToString()); err != nil {
		log.Printf("WARN: StateRpIot, %s unable to update the status with: %s", name, err)
		return
	}

	log.Printf("WARN: StateRpIot, %s status: %s done & done", name, strconv.FormatBool(alive))
	return
}
