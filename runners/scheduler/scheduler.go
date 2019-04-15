package scheduler

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"github.com/rebelit/gome/devices"
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
				devData, err := database.DbHashGet(dev)
				if err != nil {
					log.Printf("[WARN] schedule runner, unable to get dbData for %s: %s", dev, err)
					doItForReal = false
				}
				redis.ScanStruct(devData, &d)

				//get schedules data from redis
				hasSchedule, s, err := scheduleGet(d.Name)
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

func doSchedule(device devices.Devices, schedules []Schedule) {
	_, _, day, _ := splitTime()  //custom parse date/time

	devStatus, err := devices.StatusGet(device.Name)
	if err != nil {
		log.Printf("[ERROR] scheduler, %s : %s\n", device, err)
		return
	}
	devAlive, _ := strconv.ParseBool(devStatus)

	if devAlive {
		for _, schedule := range schedules {
			today := []Schedule{}

			scheduleOutCol := []Validator{}

			devComponentState, err := devices.StateGet(device.Name, schedule.Component)
			if err != nil {
				log.Printf("[ERROR] scheduler, %s : %s\n", device, err)
				return
			}
			componentState, _ := strconv.ParseBool(devComponentState)

			if schedule.Status == "enable" {
				if day == strings.ToLower(schedule.Day) {
					today = append(today, schedule)
					///add an array of schedules to today, 

					//scheduleOut := Validator{}
					//
					//onTime, _ := strconv.Atoi(schedule.On)   //time of day device is on
					//offTime, _ := strconv.Atoi(schedule.Off) //time of day device is off
					//
					//doChange, isInScheduleBlock := whatDoIDo(componentState, iTime, onTime, offTime)
					//
					//scheduleOut.DoChange = doChange
					//scheduleOut.InSchedule = isInScheduleBlock
					//
					//scheduleOutCol = append(scheduleOutCol, scheduleOut)
				}
			}

			//if device is in any enabled schedule it must be on
			for _, s := range scheduleOutCol {
				if s.InSchedule && s.DoChange {
					log.Printf("[ANDY] %s turn on\n", device.Name)
					devices.DoScheduledAction(device.Device, device.Name, schedule.Component, "on") //turn it on
				}
			}

			//if device is not in  any enabled schedule it must be off
			if noSchedules(scheduleOutCol) {
				if componentState {
					log.Printf("[ANDY] %s turn off\n", device.Name)
					devices.DoScheduledAction(device.Device, device.Name, schedule.Component, "off") //turn it off
				}
			}
		}
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