package devices

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"log"
	"strconv"
	"strings"
	"time"
)

///New functions
func StateControlTuya(profile Profile, value bool) error {
	args := []string{
		"set",
		"--set", strconv.FormatBool(value), "--ip", profile.Metadata.NetAddr,
		"--id", profile.Metadata.Id,
		"--key", profile.Metadata.Key,
		"--dps", profile.Metadata.Dps,
	}

	log.Printf("[INFO] issuing power control for %s\n", profile.Name)
	cmdOut, err := tuyaCliWrapper("tuya-cli", args)
	if err != nil{
		return err
	}
	if cmdOut != "ok"{
		common.SendSlackAlert(fmt.Sprintf("Tuya PowerControl failed to set %s to %s", profile.Name,strconv.FormatBool(value)))
		return err
	}

	common.SendSlackAlert(fmt.Sprintf("PowerControlTuya changed %s to %s",profile.Name, strconv.FormatBool(value)))
	return nil
}


func TuyaDeviceStatus (deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	doStatus := true
	alive := false
	powerState := false

	d, err := GetDevice(deviceName)
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
		if err := UpdateDeviceAliveState(deviceName, alive); err != nil{
			log.Printf("[ERROR] %s : update device status, %s\n", deviceName, err)
		}

		if err := UpdateDeviceComponentState(deviceName, "power", powerState); err != nil{
			log.Printf("[ERROR] %s : update device status, %s\n", deviceName, err)
		}
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

// device wrappers


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

func generateSetCliArgs(profile Profile, pwrState bool)(cliArg []string, err error){
	args := []string{}
	args = []string{"set", "--set", strconv.FormatBool(pwrState), "--ip", deviceDetails.Addr,
				"--id", deviceDetails.Id, "--key", deviceDetails.Key, "--dps", deviceDetails.Dps}


	return args, nil
}

func generateGetCliArgs(deviceDetails DevicesOld)(cliArg []string, err error){
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