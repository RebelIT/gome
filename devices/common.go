package devices

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"io/ioutil"
	"log"
)

// *****************************************************************
// General device functions
func StatusGet (device string) (status string, error error){  //Gets the device status from redis
	value, err := database.DbGet(device+"_status")
	if err != nil{
		return "", err
	}

	return value, nil
}

func UpdateStatus(deviceName string, status bool) error{  //Update the device status in redis
	statusData := Status{}
	statusData.Alive = status

	if err := database.DbHashSet(deviceName+"_"+"status", statusData); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return err
	}

	return nil
}

func DetailsGet (device string) (Devices, error){  //Gets the device details from redis
	d := Devices{}

	values, err := database.DbHashGet(device)
	if err != nil{
		return d, err
	}
	redis.ScanStruct(values, &d)
	return d, nil
}

func StateGet (device string, component string) (status string, error error){  //Gets the device status from redis
	/*component:
		power - tuya
		display - rpIoT HDMI status
	*/
	value, err := database.DbGet(device+"_"+component+"_state")
	if err != nil{
		return "", err
	}

	return value, nil
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
		if err := database.DbHashSet(d.Name+"_device",d); err != nil{
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

func GetAllDevicesFromDb() (devices []string, err error){  //Gets a full inventory of devices from redis
	keySearch := "*_device"
	keys, err := database.DbGetKeys(keySearch)
	if err != nil{
		return nil, err
	}
	return keys, nil
}


// *****************************************************************
// Runner device functions
func GetDeviceStatus(d Devices) {
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
				common.SendSlackAlert("[ERROR] DoScheduledAction failed to change "+deviceComponent+" for "+deviceName+" to "+deviceStatus)
			}
			return
		}
		log.Printf("[WARN] DoScheduledAction, %s no action %s found: %s\n", deviceName, deviceComponent)
		return

	case "pi":
		uriPart := "/api"
		doErr := false

		if deviceComponent == "display"{
			resp, err := PiPost(deviceName, uriPart+"/"+deviceComponent+"/"+deviceStatus)
			if err != nil{
				doErr = true
			}
			if resp.StatusCode != 200{
				doErr = true
			}
		} else{
			log.Printf("[WARN] DoScheduledAction, %s no action %s found: %s\n", deviceName, deviceComponent)
			return
		}

		if doErr{
			log.Printf("[ERROR] DoScheduledAction, %s failed to change %s to %s\n", deviceName, deviceComponent, deviceStatus)
			common.SendSlackAlert("[ERROR] DoScheduledAction failed to change "+deviceComponent+" for "+deviceName+" to "+deviceStatus)
			return
		}
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