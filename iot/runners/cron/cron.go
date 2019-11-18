package cron

import (
	"github.com/rebelit/gome/iot/devices"
	"log"
	"strconv"
	"time"
)

func GomeSchedules() error {
	log.Println("[INFO] gome cron starting")
	for {
		devs, err := devices.GetAllProfiles()
		if err != nil {
			log.Printf("[ERROR] getting device profiles from the database, will retry: %s", err)
			continue
		}

		for _, profile := range devs.Profiles {
			schedules := profile.Schedules
			state := profile.State.Alive

			if len(schedules) <= 0 {
				log.Printf("[INFO] gomeCron, %s has no schedules", profile.Name)
				continue
			}
			if !state {
				log.Printf("[WARN] gomeCron, %s has schedules but is not on the network", profile.Name)
				continue
			}

			//get slice of enabled schedules for today's day of the week
			enabledSchedules := []devices.Schedule{}
			_, iTime, today, _ := splitTime()
			for _, s := range profile.Schedules {
				if s.Enabled { //schedule enabled
					if s.Day == today { //schedule is for today
						enabledSchedules = append(enabledSchedules, s)
					}
				}
			}

			go doSchedule(profile, enabledSchedules, iTime)
		}

		time.Sleep(time.Minute * 1)
	}
}

func doSchedule(profile devices.Profile, schedules []devices.Schedule, currentTime int) {
	currentStatus := profile.State.Status
	name := profile.Name

	for _, s := range schedules {
		action := s.Action
		on, _ := strconv.Atoi(s.ActionOn)   //time of day device is on in int type
		off, _ := strconv.Atoi(s.ActionOff) //time of day device is off in int type

		changeStatus, inSchedule := whatDoIDo(currentStatus, currentTime, on, off)

		if changeStatus {
			if inSchedule {
				log.Printf("[DEBUG] %s turn on\n", name)
				go profile.SetDeviceStatus(action, true) //toggle the action
			}
			if !inSchedule {
				log.Printf("[DEBUG] %s turn off\n", name)
				go profile.SetDeviceStatus(action, false) //toggle the action
			}
		}
	}
}

func whatDoIDo(componentState bool, currentHour int, onTime int, offTime int) (changeState bool, inScheduleBlock bool) {
	reverseCheck := false
	changeState = false
	inScheduleBlock = false

	if offTime <= onTime {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck {
		//does not span PM to AM
		inScheduleBlock = inBetween(currentHour, onTime, offTime)
	} else {
		//spans a day PM to AM reverse check the schedule
		inScheduleBlock = inBetweenReverse(currentHour, offTime, onTime)
	}

	if componentState {
		if inScheduleBlock {
			//leave it be change state:false
			changeState = false
			return changeState, inScheduleBlock
		} else {
			//change state:true. change the power control to false
			changeState = true
			return changeState, inScheduleBlock
		}
	} else {
		if inScheduleBlock {
			//change state:true. change the power control to true
			changeState = true
			return changeState, inScheduleBlock
		} else {
			//leave it be change state:false
			changeState = false
			return changeState, inScheduleBlock
		}
	}
}
