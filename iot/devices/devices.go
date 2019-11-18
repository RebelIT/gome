package devices

import (
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

		for _, p := range devices.Profiles {
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
		time.Sleep(time.Minute * 5)
	}
}

func (p *Profile) SetDeviceStatus(action string, set bool) {
	doAction := Action{}
	isToggle := false

	if action == "" {
		isToggle = true
	} else {
		validAction, err := validateAction(*p, action)
		if err != nil {
			return
		}
		doAction = validAction
	}

	switch p.Make {
	case "pi":
		go p.ActionRpIot(doAction)

	case "roku":
		if isToggle {
			go p.ToggleRoku(set)
		} else {
			go p.ActionRoku(doAction)
		}

	case "tuya":
		go p.ToggleTuya(set)

	//case "newDeviceHere":
	//	go ToggleNewDevice(profile,set)

	default:
		log.Printf("[WARN] SetDeviceStatus, no device types match %s", p.Name)
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
