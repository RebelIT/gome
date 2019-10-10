package devices

import (
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	db "github.com/rebelit/gome/database"
	"log"
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

func SetDeviceStatus(profile Profile, action string, set bool) {
	doAction := Action{}
	isToggle := false

	if action == ""{
		isToggle = true
	} else {
		validAction, err := validateAction(profile,action)
		if err != nil {
			return
		}
		doAction = validAction
	}

	switch profile.Make {
	case "pi":
		go ActionRpIot(profile, doAction)

	case "roku":
		if isToggle{
			go ToggleRoku(profile, set)
		} else {
			go ActionRoku(profile, doAction)
		}

	case "tuya":
		go ToggleTuya(profile,set)

	//case "newDeviceHere":
	//	go ToggleNewDevice(profile,set)

	default:
		log.Printf("[WARN] SetDeviceStatus, no device types match %s", profile.Name)
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
// Runner device functions
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