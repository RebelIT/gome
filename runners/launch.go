package runners

import (
	"github.com/rebelit/gome/runners/aws"
	"github.com/rebelit/gome/runners/inventory"
	"github.com/rebelit/gome/runners/scheduler"
	"time"
)

func Launch(){
	go inventory.GomeStatus()
	time.Sleep(time.Second *2)

	go scheduler.GoGoScheduler()
	time.Sleep(time.Second *5)

	go aws.GoGoSQS()
	time.Sleep(time.Second *2)
}