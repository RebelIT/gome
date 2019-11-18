package cron

type Validator struct {
	DoChange   bool
	InSchedule bool
}

type Schedules struct {
	Status    string     `json:"status"`
	Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
	Day       string `json:"day"`
	Status    string `json:"status"`
	Desc      string `json:"desc"`
	Component string `json:"component"`
	Action    string `json:"action"`
	On        string `json:"on"`
	Off       string `json:"off"`
}
