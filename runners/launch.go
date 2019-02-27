package runners

import (
	"github.com/rebelit/gome/runners/scheduler"
	"github.com/rebelit/gome/runners/status"
	"time"
)

func Launch(){
	go status.GoGoDeviceStatus()
	time.Sleep(time.Second *2)

	go scheduler.GoGoScheduler()
	time.Sleep(time.Second *5)

	//go aws.GoGoSQS()
	//time.Sleep(time.Second *2)
}