package devices

import (
	"fmt"
	"github.com/pkg/errors"
	db "github.com/rebelit/gome/database"
	"github.com/rebelit/gome/util/cmd"
	"github.com/rebelit/gome/util/notify"
	"log"
	"strconv"
	"strings"
	"time"
)

func (p *Profile) ToggleTuya(value bool) error {
	args := []string{
		"set",
		"--set", strconv.FormatBool(value),
		"--ip", p.Metadata.NetAddr,
		"--id", p.Metadata.Id,
		"--key", p.Metadata.Key,
		"--dps", p.Metadata.Dps,
	}

	log.Printf("INFO: ToggleTuya, issuing power control for %s", p.Name)
	cmdOut, err := tuyaCliWrapper("tuya-cli", args)
	if err != nil {
		return err
	}
	if cmdOut != "ok" {
		notify.Slack(fmt.Sprintf("Tuya PowerControl failed to set %s to %s", p.Name, strconv.FormatBool(value)))
		return err
	}

	notify.Slack(fmt.Sprintf("PowerControlTuya changed %s to %s", p.Name, strconv.FormatBool(value)))
	return nil
}

func StateTuya(name string, collectionDelayMin time.Duration) {
	log.Printf("INFO: StateTuya, %s device collection delayed +%d sec", name, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	alive := false
	status := false

	value, err := db.Get(name)
	if err != nil {
		log.Printf("WARN: StateTuya, %s not found for status collection, is it orphaned: %s", name, err)
		return
	}

	p, err := stringToStruct(value)
	if err != nil {
		log.Printf("WARN: StateTuya, %s does not have a valid profile in the database: %s", name, err)
		return
	}

	args := []string{
		"get",
		"--ip", p.Metadata.NetAddr,
		"--id", p.Metadata.Id,
		"--key", p.Metadata.Key,
		"--dps", p.Metadata.Dps,
	}

	cmdOut, err := tuyaCliWrapper("tuya-cli", args)
	if err != nil {
		log.Printf("ERROR: StateTuya, %s : device status, %s", name, err)
	} else {
		alive = true
	}

	if cmdOut == "true" {
		status = true
	}

	p.State.Alive = alive
	p.State.Status = status

	if err := db.Add(name, p.structToString()); err != nil {
		log.Printf("ERROR: StateTuya, %s : update device state, %s", name, err)
		return
	}

	log.Printf("INFO: StateTuya, %s status: %s done & done", name, strconv.FormatBool(alive))
	return
}

//tuya wrapper
func tuyaCliWrapper(cmdName string, args []string) (cmdReturn string, error error) {
	cmdOut, err := cmd.TryCommand(cmdName, args, 5)
	if err != nil {
		return "", err
	}

	fmtOut := strings.Replace(string(cmdOut), "\n", "", -1)
	if fmtOut == "Set succeeded." {
		return "ok", nil

	} else if fmtOut == "false" {
		return fmtOut, nil

	} else if fmtOut == "true" {
		return fmtOut, nil

	} else {
		return "", errors.Errorf("cmd returnes could not be formatted: returned %s", fmtOut)

	}
}
