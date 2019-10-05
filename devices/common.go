package devices

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// *****************************************************************
// General device functions
func GetAllDevicesFromDb() (devices []string, err error){  //Gets a full inventory of devices from redis
	keySearch := "*_device"
	keys, err := database.DbGetKeys(keySearch)
	if err != nil{
		return nil, err
	}
	return keys, nil
}

func GetDevice (deviceName string) (deviceDetails Devices, error error){  //Gets the device status from redis
	key := ""
	if strings.Contains(deviceName, "_device"){
		key = deviceName
	}else{
		key = deviceName+"_device"
	}

	d := Devices{}

	value, err := database.DbGet(key)
	if err != nil{
		return d, err
	}

	json.Unmarshal([]byte(value), &d)

	return d, nil
}

func UpdateDevice (d* Devices) error{  //Gets the device status from redis
	key := d.Name+"_device"

	value, err := json.Marshal(d)
	if err != nil{
		log.Println(err)
	}

	if err := database.DbSet(key, value); err != nil{
		return err
	}
	return nil
}

func GetDeviceAliveState (deviceName string) (status string, error error){  //Gets the device status from redis
	key := deviceName+"_alive"

	value, err := database.DbGet(key)
	if err != nil{
		return "", err
	}

	return value, nil
}

func UpdateDeviceAliveState(deviceName string, status bool) error{  //Update the device status in redis
	key := deviceName+"_alive"
	value := []byte(strconv.FormatBool(status))

	if err := database.DbSet(key,value); err != nil{
		return err
	}

	return nil
}

func GetDeviceComponentState(deviceName string, component string) (status string, error error){  //Gets the device component status from redis
	key := deviceName+"_"+component+"_state"

	value, err := database.DbGet(key)
	if err != nil{
		return "", err
	}

	return value, nil
}

func UpdateDeviceComponentState(deviceName string, component string, state bool) error{  //Update the device component state in redis
	key := deviceName+"_"+component+"_state"
	value := []byte(strconv.FormatBool(state))
	if err := database.DbSet(key,value); err != nil{
		return err
	}

	return nil
}

func LoadDevices() error{  //Load Devices into redis from devices.json file
	log.Printf("[INFO] device loader, starting")
	i, err := ReadDeviceFile()
	if err != nil{
		return err
	}

	if len(i.Devices) == 0{
		log.Printf("[WARN] device loader, no devices to load, skipping")
		return nil
	}

	for _, d := range i.Devices {
		log.Printf("[INFO] device loader, loading %s under '%s_device'", d.Name, d.Name)

		if err := UpdateDevice(&d); err != nil{
			return err
		}
	}
	log.Println("[INFO] device loader, all done")
	return nil
}

func ReadDeviceFile()(Inputs, error){  //Read the devices.json
	var in Inputs
	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		return in, err
	}
	if err := json.Unmarshal(deviceFile, &in); err != nil{
		return in, err
	}

	return in, nil
}


// *****************************************************************
// Runner device functions
func GetDeviceStatus(d* Devices) {
	delay := randomizeCollection()

	switch d.Device {
	case "pi":
		go RpIotDeviceStatus(d.Name, delay)

	case "roku":
		go RokuDeviceStatus(d.Name, delay)

	case "tuya":
		go TuyaDeviceStatus(d.Name, delay)

	default:
		log.Printf("[WARN] GetDeviceStatus, no device types match %s", d.Name)
	}
}

func DoScheduledAction(device string, deviceName string, deviceComponent string, deviceStatus string){
	switch device {
	case "tuya":
		if deviceComponent == "power"{
			newStatus := false
			if deviceStatus == "on"{
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
		log.Printf("[WARN] DoScheduledAction, %s no action %s found: %s\n", deviceName, deviceComponent)
		return

	case "pi":
		if deviceComponent == "display"{  //scheduler switch for rPioT display controls
			if err := rpIotDisplayToggle(deviceName, deviceStatus); err != nil{
				doScheduleEror(deviceName, deviceComponent, deviceStatus)
				return
			}
		} else{
			log.Printf("[WARN] DoScheduledAction, %s no action %s found: %s\n", deviceName, deviceComponent)
			return
		}
		doScheduleOk(deviceName, deviceComponent, deviceStatus)
		return

	default:
		log.Printf("[WARN] DoScheduledAction, no device types match %s", deviceName)
		return
	}
}

func DoWhatAlexaSays(deviceType string, deviceName string, deviceAction string) error{
	action := false

	common.MetricAws("alexa", "doAction", "nil",deviceName, deviceAction)

	switch deviceType{
	case "tuya":
		if deviceAction == "on"{
			action = true
		}
		if err := TuyaPowerControl(deviceName, action); err != nil{
			return err
		}
		return nil

	case "pi":
		_, err := PiPost(deviceName, deviceAction)
		if err != nil{
			return err
		}

	default:
		//no match
		return errors.New("no message in queue to parse")
	}

	return nil
}

func doScheduleEror(deviceName string, deviceComponent string, deviceState string){
	common.SendSlackAlert("Scheduler failed to do a scheduled action\n" +
		""+deviceName+" - "+deviceComponent+" - "+deviceState)
}

func doScheduleOk(deviceName string, deviceComponent string, deviceState string){
	common.SendSlackAlert("Scheduler did it's thing!\n" +
		""+deviceName+" - "+deviceComponent+" - "+deviceState)
}