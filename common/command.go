package common

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func TryCommand(cmdName string, args []string) (string, error){
	maxRetry := 5
	retrySleep := time.Second * 2

	for i := 0; ;i++ {
		if i >= maxRetry{
			break
		}
		cmdOut, err := Command(cmdName, args)
		if err == nil{
			return cmdOut, err
		}
		log.Printf("[WARN] cmd %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %i reached for %s\n", maxRetry, cmdName)
}

func Command(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
