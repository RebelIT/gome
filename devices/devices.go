package devices

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"log"
	"strings"
	"time"
)

func GetDeviceStatus() {
	//function checks each device if it is online and the toggle state, updates accordingly
	for {
		delay := randomizeCollection()
		devices, err := GetAllProfiles()
		if err != nil {
			log.Printf("unable to get all devices for runner; %s", err)
			continue
		}

		for _, p := range devices.Profiles{
			switch p.Make {
			case "pi":
				go StateRpIot(p.Name, delay)

			case "roku":
				go StateRoku(p.Name, delay)

			case "tuya":
				go StateTuya(p.Name, delay)

			//case "newDeviceHere":
			//	go StateNewDevice(p.Name, delay)

			default:
				log.Printf("[WARN] GetDeviceStatus, no device types match %s", p.Name)
			}
		}
		time.Sleep(time.Minute * common.INVENTORY_MIN)
	}
}

func GetProfile(name string) (profile Profile, error error) { //Gets the device profile from the database
	value, err := db.Get(name)
	if err != nil {
		return Profile{}, err
	}

	profile, err = stringToStruct(value)
	if error != nil {
		return Profile{}, nil
	}

	return profile, nil
}

func GetAllProfiles() (profiles Devices, err error) { //Gets the full inventory of device profiles from the database
	profiles = Devices{}
	list := []Profile{}

	keys, err := db.GetAll()
	if err != nil {
		return Devices{}, err
	}

	for _, key := range keys {
		p, err := GetProfile(key)
		if err != nil {
			continue
		}
		list = append(list, p)
	}
	profiles.Profiles = list

	return profiles, nil
}










// *****************************************************************
// General device functions
func GetDevice(deviceName string) (deviceDetails DevicesOld, error error) { //Gets the device status from the database
	key := ""
	if strings.Contains(deviceName, "_device") {
		key = deviceName
	} else {
		key = deviceName + "_device"
	}

	d := DevicesOld{}

	value, err := db.Get(key)
	if err != nil {
		return d, err
	}

	if err := json.Unmarshal([]byte(value), &d); err != nil {
		return d, err
	}

	return d, nil
}

func GetDeviceAliveState(deviceName string) (status string, error error) { //Gets the device status from the database
	key := deviceName + "_alive"

	value, err := db.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}


func GetDeviceComponentState(deviceName string, component string) (status string, error error) { //Gets the device component status from the database
	key := deviceName + "_" + component + "_state"

	value, err := db.Get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}


// *****************************************************************
// Runner device functions


func DoScheduledAction(device string, deviceName string, deviceComponent string, deviceStatus string) {
	switch device {
	case "tuya":
		if deviceComponent == "power" {
			newStatus := false
			if deviceStatus == "on" {
				newStatus = true
			}
			if err := TuyaPowerControl(deviceName, newStatus); err != nil {
				log.Printf("[ERROR] DoScheduledAction, %s failed to change %s to %s\n", deviceName, deviceComponent, deviceStatus)
				doScheduleEror(deviceName, deviceComponent, deviceStatus)
				return
			}
			doScheduleOk(deviceName, deviceComponent, deviceStatus)
			return
		}
		log.Printf("[WARN] DoScheduledAction, %s no action %s found\n", deviceName, deviceComponent)
		return

	case "pi":
		if deviceComponent == "display" { //scheduler switch for rPioT display controls
			if err := rpIotDisplayToggle(deviceName, deviceStatus); err != nil {
				doScheduleEror(deviceName, deviceComponent, deviceStatus)
				return
			}
		} else {
			log.Printf("[WARN] DoScheduledAction, %s no action %s found\n", deviceName, deviceComponent)
			return
		}
		doScheduleOk(deviceName, deviceComponent, deviceStatus)
		return

	default:
		log.Printf("[WARN] DoScheduledAction, no device types match %s", deviceName)
		return
	}
}

func DoWhatAlexaSays(deviceType string, deviceName string, deviceAction string) error {
	action := false

	common.MetricAws("alexa", "doAction", "nil", deviceName, deviceAction)

	switch deviceType {
	case "tuya":
		if deviceAction == "on" {
			action = true
		}
		if err := TuyaPowerControl(deviceName, action); err != nil {
			return err
		}
		return nil

	case "pi":
		_, err := PiPost(deviceName, deviceAction)
		if err != nil {
			return err
		}

	default:
		//no match
		return errors.New("no message in queue to parse")
	}

	return nil
}

func doScheduleEror(deviceName string, deviceComponent string, deviceState string) {
	common.SendSlackAlert("Scheduler failed to do a scheduled action\n" +
		"" + deviceName + " - " + deviceComponent + " - " + deviceState)
}

func doScheduleOk(deviceName string, deviceComponent string, deviceState string) {
	common.SendSlackAlert("Scheduler did it's thing!\n" +
		"" + deviceName + " - " + deviceComponent + " - " + deviceState)
}
