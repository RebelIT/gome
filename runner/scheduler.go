package runner

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/devices/tuya"
	"log"
	"strconv"
	"strings"
	"time"
)

func GoGoScheduler() error {
	log.Println("[INFO] schedule runner, starting")
	for {
		doIt := true

		//get array of all devices in the database
		devs, err := devices.GetAllDevicesFromDb()
		if err != nil{
			log.Printf("[WARN] schedule runner, get all devices: %s", err)
			doIt = false
		}

		if len(devs) == 0{
			log.Printf("[WARN] schedule runner, no devices found in the database")
			doIt = false
		}

		if doIt {
			for _, dev := range devs {
				doItForReal := true
				d := devices.Devices{}

				//get device data from redis
				devData, err := devices.DbHashGet(dev)
				if err != nil {
					log.Printf("[WARN] schedule runner, unable to get dbData for %s: %s", dev, err)
					doItForReal = false
				}
				redis.ScanStruct(devData, &d)

				//get schedules data from redis
				hasSchedule, s, err := devices.ScheduleGet(d.Name)
				if err != nil{
					log.Printf("[WARN] schedule runner, unable to get schedule for %s: %s", dev, err)
					doItForReal = false
				}

				if !hasSchedule {
					log.Printf("[INFO] schedule runner, no schedule for %s", dev)
					doItForReal = false
				}

				if doItForReal {
					if s.Status != "enable" {
						log.Printf("[INFO] schedule runner, %s has schedule defined but not enabled", dev)
						doItForReal = false
					}
				}

				if doItForReal {
					go doSchedule(d, s.Schedules)
				}
			}
		}
		time.Sleep(time.Minute *common.SCHEDULE_MIN)
	}

	common.SendSlackAlert("[ERROR] scheduler, routine broke out of loop")
	return nil
}

func doSchedule(device devices.Devices, schedules []devices.Schedule) {
	_, iTime, day, _ := splitTime()  //custom parse date/time

	for _, schedule := range schedules {
		devStatusRtn, err := devices.StatusGet(device.Name)
		if err != nil {
			log.Printf("[ERROR] scheduler, %s : %s\n", device, err)
		}
		devOn, err := strconv.ParseBool(devStatusRtn)

		scheduleOutCol := []Validator{}
		for _, s := range schedules {
			if day == strings.ToLower(s.Day) {
				if s.Status == "enable" {
					scheduleOut := Validator{}

					onTime, _ := strconv.Atoi(s.On)   //time of day device is on
					offTime, _ := strconv.Atoi(s.Off) //time of day device is off

					doChange, isInScheduleBlock := whatDoIDo(devOn, iTime, onTime, offTime)

					scheduleOut.DoChange = doChange
					scheduleOut.InSchedule = isInScheduleBlock

					scheduleOutCol = append(scheduleOutCol, scheduleOut)
				}
			}
		}

		//if device is in any enabled schedule it must be on
		for _, s := range scheduleOutCol {
			if s.InSchedule && s.DoChange {
				if err := doDeviceSpecificAction(device.Device, device.Name, schedule.Action, "on"); err != nil { //change it to true
					log.Printf("[ERROR] schedule runner, %s failed to change powerstate: %s\n", device.Name, err)
					common.SendSlackAlert("[ERROR] schedule runner failed to change powerstate for " + device.Name)
				}
			}
		}

		//if device is not in  any enabled schedule it must be off
		if noSchedules(scheduleOutCol) {
			if devOn {
				if err := doDeviceSpecificAction(device.Device, device.Name, schedule.Action, "off"); err != nil { //change it to true
					log.Printf("[ERROR] schedule runner, %s failed to change powerstate: %s\n", device.Name, err)
					common.SendSlackAlert("[ERROR] schedule runner failed to change powerstate for " + device.Name)
				}
			}
		}
	}
}

func doDeviceSpecificAction(deviceType string, deviceName string, deviceAction string, deviceStatus string) error{
	switch deviceType {
	case "tuya":
		newStatus := false
		if deviceStatus == "on"{
			newStatus = true
		}
		if err := tuya.PowerControl(deviceName, newStatus); err != nil {
			log.Printf("[ERROR] scheduler, %s failed to change powerstate: %s\n", deviceName, err)
			common.SendSlackAlert("[ERROR] schedule runner failed to change powerstate for "+deviceName+" to "+deviceStatus)
		}
		return nil

	case "pi":
		return nil

	default:
		log.Printf("[WARN] schedule runner, %s no device types match", deviceName)
		return errors.New("no device types match for "+deviceName)

	}
}

func noSchedules(v []Validator) bool {
	for _, s := range v {
		if s.InSchedule {
			return false
		}
	}
	return true
}

func whatDoIDo(devOn bool, currentHour int, devOnTime int, devOffTime int) (changeState bool, inScheduleBlock bool){
	reverseCheck := false
	changeState = false
	inScheduleBlock = false

	if devOffTime <= devOnTime {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck{
		//does not span PM to AM
		inScheduleBlock = inBetween(currentHour, devOnTime, devOffTime)
	} else {
		//spans a day PM to AM reverse check the schedule
		inScheduleBlock = inBetweenReverse(currentHour, devOffTime, devOnTime)
	}

	if devOn{
		if inScheduleBlock{
			//leave it be change state:false
			changeState = false
			return changeState, inScheduleBlock
		} else {
			//change state:true. change the power control to false
			changeState = true
			return changeState, inScheduleBlock
		}
	} else {
		if inScheduleBlock{
			//change state:true. change the power control to true
			changeState = true
			return changeState, inScheduleBlock
		}else {
			//leave it be change state:false
			changeState = false
			return changeState, inScheduleBlock
		}
	}
}