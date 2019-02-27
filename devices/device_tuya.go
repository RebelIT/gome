package devices

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"log"
	"strconv"
	"strings"
	"time"
)

func TuyaDeviceStatus (deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	doStatus := true
	alive := false
	powerState := false

	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	cliArgs, err := generateGetCliArgs(d)
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	if doStatus{
		cmdOut, err := tuyaCliWrapper("tuya-cli", cliArgs)
		if err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		}else{
			alive = true
		}

		if cmdOut == "true"{
			powerState = true
		}

		//status = is on the network and accessible
		if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}

		//state = the current switch power state on/true | off/false
		if err := database.DbSet(deviceName+"_power"+"_state", []byte(strconv.FormatBool(powerState))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

// device wrappers
func TuyaPowerControl(deviceName string, value bool) error {
	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		return err
	}

	cliArgs, err := generateSetCliArgs(d, value)
	if err != nil{
		return err
	}

	log.Printf("[INFO] issuing power control for %s\n", deviceName)
	cmdOut, err := tuyaCliWrapper("tuya-cli", cliArgs)
	if err != nil{
		return err
	}
	if cmdOut != "ok"{
		common.SendSlackAlert("Tuya PowerControl failed for "+d.Name+"("+d.NameFriendly+") to "+strconv.FormatBool(value)+"")
		return err
	}
	if err := UpdateStatus(deviceName, value); err != nil{
		log.Printf("[ERROR] Update Device Status, %s : %s", deviceName, err)
	}

	common.SendSlackAlert("Tuya PowerControl initiated for "+d.Name+"("+d.NameFriendly+") to "+strconv.FormatBool(value)+"")
	return nil
}

func tuyaCliWrapper(cmdName string, args []string) (cmdReturn string, error error){
	cmdOut, err := common.TryCommand(cmdName, args,5)
	if err != nil{
		return "",err
	}

	fmtOut := strings.Replace(string(cmdOut), "\n", "", -1)
	if fmtOut == "Set succeeded." {
		return "ok", nil

	}else if fmtOut == "false" {
		return fmtOut, nil

	} else if fmtOut == "true"{
		return fmtOut, nil

	} else{
		return "",errors.Errorf("cmd returnes could not be formatted: returned %s", fmtOut)

	}
}

func generateSetCliArgs(deviceDetails Devices, pwrState bool)(cliArg []string, err error){
	args := []string{}
	fmt.Printf("%+v\n",deviceDetails)
	switch deviceDetails.Type{
		case "outlet":
			args = []string{"set", "--set", strconv.FormatBool(pwrState), "--ip", deviceDetails.Addr,
							"--id", deviceDetails.Id, "--key", deviceDetails.Key}

		case "switch":
			args = []string{"set", "--set", strconv.FormatBool(pwrState), "--ip", deviceDetails.Addr,
				"--id", deviceDetails.Id, "--key", deviceDetails.Key, "--dps", deviceDetails.Dps}

		default:
			return args, errors.New("no device type "+deviceDetails.Type+" found in cli args switch")
	}
	return args, nil
}

func generateGetCliArgs(deviceDetails Devices)(cliArg []string, err error){
	args := []string{}

	switch deviceDetails.Type{
	case "outlet":
		args = []string{"get", "--ip", deviceDetails.Addr, "--id", deviceDetails.Id, "--key", deviceDetails.Key}

	case "switch":
		args = []string{"get", "--ip", deviceDetails.Addr, "--id", deviceDetails.Id, "--key", deviceDetails.Key,
						"--dps", deviceDetails.Dps}

	default:
		return args, errors.New("no device type "+deviceDetails.Type+" found in cli args switch")
	}
	return args, nil
}