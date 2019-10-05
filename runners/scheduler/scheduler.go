package scheduler

import (
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"log"
	"strconv"
	"time"
)

func GoGoScheduler() error {
	log.Println("[INFO] schedule runner, starting")
	for {
		processSchedules := true

		devs, err := devices.GetAllDevicesFromDb()
		if err != nil{
			log.Printf("[WARN] schedule runner, get all devices: %s", err)
			processSchedules = false
		}

		if len(devs) == 0{
			log.Printf("[WARN] schedule runner, no devices found in the database")
			processSchedules = false
		}

		if processSchedules {
			for _, dev := range devs {
				doItForReal := true

				d, err := devices.GetDevice(dev)
				if err != nil {
					log.Printf("[WARN] schedule runner, unable to get dbData for %s: %s", dev, err)
				}

				//check if device is online 'alive'
				devStatus, err := devices.GetDeviceAliveState(d.Name)
				if err != nil {
					log.Printf("[ERROR] scheduler runner, unable to get %s alive status: %s\n", d.Name, err)
				}
				devAlive, _ := strconv.ParseBool(devStatus)

				if devAlive {
					//get schedules data from redis
					hasSchedule, s, err := scheduleGet(d.Name)
					if err != nil {
						log.Printf("[WARN] schedule runner, unable to get %s schedule: %s", dev, err)
						doItForReal = false
					}

					if !hasSchedule {
						log.Printf("[DEBUG] schedule runner, no schedule for %s", dev)
						doItForReal = false
					}

					if doItForReal {
						if s.Status != "enable" {
							log.Printf("[DEBUG] schedule runner, %s has schedule defined but not enabled", dev)
							doItForReal = false
						}
					}

					if doItForReal {
						go doSchedule(d, s.Schedules)
					}
				}
			}
		}
		time.Sleep(time.Minute *common.SCHEDULE_MIN)
	}
}

func doSchedule(device devices.Devices, schedules []Schedule) {
	_, iTime, day, _ := splitTime()  //custom parse date/time

	for _, schedule := range schedules {
		if schedule.Day == day {
			if schedule.Status == "enable" {

				devComponentState, err := devices.GetDeviceComponentState(device.Name, schedule.Component)
				if err != nil {
					log.Printf("[ERROR] doSchedule, get %s %s state: %s\n", device, schedule.Component, err)
					return
				}

				componentState, _ := strconv.ParseBool(devComponentState) //state of the device component in the schedule
				onTime, _ := strconv.Atoi(schedule.On)                    //time of day device is on
				offTime, _ := strconv.Atoi(schedule.Off)                  //time of day device is off

				doChange, inSchedule := changeComponentState(componentState, iTime, onTime, offTime)

				if doChange {
					if inSchedule {
						log.Printf("[DEBUG] %s turn on\n", device.Name)
						devices.DoScheduledAction(device.Device, device.Name, schedule.Component, "on") //turn it on
					}
					if !inSchedule {
						log.Printf("[ANDY] %s turn off\n", device.Name)
						devices.DoScheduledAction(device.Device, device.Name, schedule.Component, "off") //turn it off
					}
				}
			}
		}
	}
}

func changeComponentState(componentState bool, currentHour int, onTime int, offTime int) (changeState bool, inScheduleBlock bool){
	reverseCheck := false
	changeState = false
	inScheduleBlock = false

	if offTime <= onTime {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck{
		//does not span PM to AM
		inScheduleBlock = inBetween(currentHour, onTime, offTime)
	} else {
		//spans a day PM to AM reverse check the schedule
		inScheduleBlock = inBetweenReverse(currentHour, offTime, onTime)
	}

	if componentState {
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