package runners

import (
	"github.com/rebelit/gome/runners/aws"
	"github.com/rebelit/gome/runners/inventory"
	"github.com/rebelit/gome/runners/cron"
	"time"
)

func Launch(){
	go inventory.GomeDevices()
	time.Sleep(time.Second *2)

	go cron.GomeSchedules()
	time.Sleep(time.Second *5)

	go aws.GoGoSQS()
	time.Sleep(time.Second *2)
}