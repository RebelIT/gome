package scheduler

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/database"
	"log"
	"strconv"
	"strings"
	"time"
)

type Validator struct{
	DoChange 	bool
	InSchedule 	bool
}

type Schedules struct {
	Status 		string `json:"status"`
	Schedules		[]Schedule `json:"schedules"`
}

type Schedule struct{
	Day		string `json:"day"`
	Status 	string `json:"status"`
	Desc	string `json:"desc"`
	Action  string `json:"action"`
	On		string `json:"on"`
	Off 	string `json:"off"`
}

// *****************************************************************
// Scheduler wrapper functions
func scheduleSet(s* Schedules, device string) (error){
	data, err := json.Marshal(s)
	if err != nil{
		log.Println(err)
	}

	if err := database.DbSet(device+"_schedule", data); err != nil{
		return err
	}
	return nil
}

func scheduleGet(devName string) (hasSchedule bool, schedules Schedules, error error){
	s := Schedules{}

	value, err := database.DbGet(devName+"_schedule")
	if err != nil{
		return false, s, err
	}
	if value == ""{
		return false, s, nil
	}

	json.Unmarshal([]byte(value), &s)

	if len(s.Schedules) <= 1 {
		return false, s, errors.New("invalid schedule struct")
	}

	return true, s, nil
}

func scheduleDel(device string) (error){
	if err := database.DbDel(device+"_schedule"); err != nil{
		return err
	}
	return nil
}

func scheduleUpdate(device string, status string) (error){
	_, s, err := scheduleGet(device)
	if err != nil{
		return err
	}

	s.Status = status

	if err := scheduleSet(&s,device); err != nil{
		return err
	}

	return nil
}


// *****************************************************************
// misc helper functions
func splitTime()(strTime string, intTime int, weekday string, now time.Time){
	Now := time.Now()
	NowMinute := Now.Minute()
	NowHour := Now.Hour()
	NowDay := Now.Weekday()

	sTime := ""
	singleMinute := inBetween(NowMinute, 0,9)
	if singleMinute{
		sTime = strconv.Itoa(NowHour) + "0"+ strconv.Itoa(NowMinute)
	} else{
		sTime = strconv.Itoa(NowHour) + strconv.Itoa(NowMinute)
	}

	iTime, _ := strconv.Atoi(sTime)
	day := strings.ToLower(NowDay.String())

	return sTime, iTime, day, Now
}

func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func inBetweenReverse(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return false
	} else {
		return true
	}
}